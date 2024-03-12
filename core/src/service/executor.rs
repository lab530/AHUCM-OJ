use std::{
    ffi::CString,
    fs::{self, File},
    os::fd::FromRawFd,
    sync::{
        atomic::{AtomicU32, AtomicUsize, Ordering},
        Arc,
    },
};

use nix::{
    libc::{exit, fdopen, freopen, STDERR_FILENO},
    sys::wait::waitpid,
    unistd::{execvp, fork, write, ForkResult},
};

use crate::util::config::GLOB_CONFIG;

use super::testcases_getter::TestcasesGetter;

macro_rules! c_string {
    ($string:expr) => {
        CString::new($string).unwrap()
    };
}

#[derive(Debug)]
pub enum CompilationError {
    Error(String),
    MissingLang(String),
    ForkFailed,
    FileSystemError,
}

#[derive(Debug)]
pub enum RunningError {
    ForkFailed,
    NonEmptyStderr(String),
}

pub enum ExecutionResult {
    Halt,
    Accpected(usize, u32, u32),
    WrongAnswer(usize, usize),
    PresentationError,
    CompilationError(String),
    RuntimeError(String),
    TimeLimitExceeded(u32, u32),
    MemoLimitExceeded(u32, u32),
    UnknownError(String),
}

#[derive(Default)]
struct RunContext {
    pub max_testcase_cnt: AtomicUsize,
    pub run_testcase_cnt: AtomicUsize,
    pub passed_testcase_cnt: AtomicUsize,
    pub elapsed_time: AtomicU32,
    pub used_memory: AtomicU32,
}

impl RunContext {
    pub fn max_testcase_cnt(&self) -> usize {
        self.max_testcase_cnt.load(Ordering::Relaxed)
    }

    pub fn run_testcase_cnt(&self) -> usize {
        self.run_testcase_cnt.load(Ordering::SeqCst)
    }

    pub fn passed_testcase_cnt(&self) -> usize {
        self.run_testcase_cnt.load(Ordering::SeqCst)
    }

    pub fn elapsed_time(&self) -> u32 {
        self.elapsed_time.load(Ordering::SeqCst)
    }

    pub fn used_memory(&self) -> u32 {
        self.used_memory.load(Ordering::SeqCst)
    }
}

pub struct Executor {
    source_path: String,
    lang: String,
    target_path: String,
    log_path: String,

    problem_id: u64,
    mem_limit: u32,
    time_limit: u32,
    testcases_path: String,

    submission_id: u64,

    is_compiling_done: bool,
    is_running_done: bool,

    run_ctx: RunContext,
}

impl Executor {
    pub fn new(
        source_path: String,
        lang: String,
        problem_id: u64,
        mem_limit: u32,
        time_limit: u32,
        testcases_path: String,
        submission_id: u64,
    ) -> Self {
        Self {
            source_path,
            lang,
            target_path: format!("{}.exe", source_path),
            log_path: format!("{}.log", source_path),
            problem_id,
            mem_limit,
            time_limit,
            testcases_path,
            submission_id,
            is_compiling_done: false,
            is_running_done: false,
            run_ctx: RunContext::default(),
        }
    }

    pub fn execute(&mut self) -> ExecutionResult {
        if let Err(e) = self.compile() {
            return ExecutionResult::CompilationError(format!("{:?}", e));
        }
        self.is_compiling_done = true;

        let mut testcase_getter = TestcasesGetter::new(self.testcases_path.clone());
        let testcases = testcase_getter.get_testcases();
        self.run_ctx
            .max_testcase_cnt
            .store(testcases.len(), Ordering::Relaxed);

        for testcase in testcases {
            self.run(
                testcase.get_input_path(),
                testcase.get_output_path(),
            );
        }

        let result: ExecutionResult;
        if self.run_ctx.max_testcase_cnt() as u32 * self.time_limit < self.run_ctx.elapsed_time() {
            result = ExecutionResult::TimeLimitExceeded(
                self.run_ctx.elapsed_time(),
                self.run_ctx.max_testcase_cnt() as u32 * self.time_limit,
            );
        } else if self.run_ctx.max_testcase_cnt() as u32 * self.mem_limit
            < self.run_ctx.used_memory()
        {
            result = ExecutionResult::MemoLimitExceeded(
                self.run_ctx.used_memory(),
                self.run_ctx.max_testcase_cnt() as u32 * self.mem_limit,
            );
        } else if self.run_ctx.max_testcase_cnt() == self.run_ctx.passed_testcase_cnt() {
            result = ExecutionResult::Accpected(
                self.run_ctx.max_testcase_cnt(),
                self.run_ctx.elapsed_time(),
                self.run_ctx.used_memory(),
            );
        } else {
            result = ExecutionResult::WrongAnswer(
                self.run_ctx.passed_testcase_cnt(),
                self.run_ctx.max_testcase_cnt(),
            );
        }

        self.update_db(&result);

        result
    }

    fn compile(&self) -> Result<(), CompilationError> {
        let command = GLOB_CONFIG
            .get_compile_command(&self.lang, &self.source_path, &self.target_path)
            .ok_or(CompilationError::MissingLang(format!(
                "missing lang `{}`",
                self.lang
            )))?
            .iter()
            .map(|s| c_string!(s.as_str()))
            .collect::<Vec<_>>();

        match unsafe { fork() } {
            Ok(ForkResult::Parent { child, .. }) => {
                waitpid(child, None).unwrap();
            }
            Ok(ForkResult::Child) => {
                let log_path = c_string!(self.log_path.as_str());
                let w_mode = c_string!("w");

                unsafe {
                    let log_output = fdopen(STDERR_FILENO, w_mode.as_ptr());
                    freopen(log_path.as_ptr(), w_mode.as_ptr(), log_output);
                }

                match execvp(&command[0], &command) {
                    Ok(_) => unreachable!(),
                    Err(errno) => write(
                        unsafe { File::from_raw_fd(STDERR_FILENO) },
                        format!("Execvp error, errno = {:?}\n", errno).as_bytes(),
                    )
                    .ok(),
                };

                unsafe { exit(0) };
            }
            _ => return Err(CompilationError::ForkFailed),
        };

        let log = fs::read_to_string(self.log_path).map_err(|_e| { CompilationError::FileSystemError })?;
        if log.is_empty() {
            Ok(())
        } else {
            Err(CompilationError::Error(log))
        }
    }

    fn run(&self, input_path: &str, output_path: &str) -> Result<(), RunningError>{
        let command = GLOB_CONFIG
            .get_compile_command(&self.lang, &self.source_path, &self.target_path)
            .ok_or(CompilationError::MissingLang(format!(
                "missing lang `{}`",
                self.lang
            )))?
            .iter()
            .map(|s| c_string!(s.as_str()))
            .collect::<Vec<_>>();

        let redirect_stdout_path = format!("{}.stdout", self.source_path);
        let redirect_stderr_path = format!("{}.stderr", self.source_path);

        match unsafe { fork() } {
            Ok(ForkResult::Parent{ child, .. }) => {
                waitpid(child, None).unwrap();
            },
            Ok(ForkResult::Child) => {
                let input_path = c_string!(input_path);
                let redirect_stdout_path = c_string!(redirect_stdout_path.as_str());
                let redirect_stderr_path = c_string!(redirect_stderr_path.as_str());

                // TODO:
            },
            _ => return Err(RunningError::ForkFailed),
        };

        Ok(())
    }

    fn update_db(&self, execute_result: &ExecutionResult) {
        // TODO: logic for updating db

    pub fn clean(&self) {
        self.after_run_clean();
        self.after_compile_clean();
    }

    fn after_compile_clean(&self) {
        if self.is_compiling_done {}
    }

    fn after_run_clean(&self) {
        if self.is_running_done {}
    }
}
