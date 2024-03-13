use std::{
    ffi::CString,
    fs::{self, File},
    os::fd::FromRawFd,
    sync::{
        atomic::{AtomicU32, AtomicUsize, Ordering},
        Arc,
    },
};

use log::{debug, info};
use nix::{
    libc::{exit, fdopen, freopen, STDERR_FILENO, STDIN_FILENO, STDOUT_FILENO},
    sys::wait::waitpid,
    unistd::{execvp, fork, write, ForkResult},
};

use crate::util::{comparer::Comparer, config::GLOB_CONFIG};

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
    MissingLang(String),
    NonEmptyStderr(String),
}

#[derive(Debug)]
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
    pub max_testcase_cnt: usize,
    pub run_testcase_cnt: usize,
    pub passed_testcase_cnt: usize,
    pub elapsed_time: u32,
    pub used_memory: u32,
}

impl RunContext {
    pub fn max_testcase_cnt(&self) -> usize {
        // self.max_testcase_cnt.load(Ordering::Relaxed)
        self.max_testcase_cnt
    }

    pub fn run_testcase_cnt(&self) -> usize {
        // self.run_testcase_cnt.load(Ordering::SeqCst)
        self.run_testcase_cnt
    }

    pub fn passed_testcase_cnt(&self) -> usize {
        // self.run_testcase_cnt.load(Ordering::SeqCst)
        self.passed_testcase_cnt
    }

    pub fn elapsed_time(&self) -> u32 {
        // self.elapsed_time.load(Ordering::SeqCst)
        self.elapsed_time
    }

    pub fn used_memory(&self) -> u32 {
        // self.used_memory.load(Ordering::SeqCst)
        self.used_memory
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
            source_path: source_path.clone(),
            lang,
            target_path: format!("{}.exe", source_path),
            log_path: format!("{}.log", source_path),
            problem_id,
            mem_limit,
            time_limit,
            testcases_path,
            submission_id,
            is_compiling_done: true,
            is_running_done: true,
            run_ctx: RunContext::default(),
        }
    }

    pub fn execute(&mut self) -> ExecutionResult {
        if let Err(e) = self.compile() {
            return ExecutionResult::CompilationError(format!("{:?}", e));
        }

        let mut testcase_getter = TestcasesGetter::new(self.testcases_path.clone());
        let testcases = testcase_getter.get_testcases();
        // self.run_ctx
        //     .max_testcase_cnt
        //     .store(testcases.len(), Ordering::Relaxed);
        self.run_ctx.max_testcase_cnt = testcases.len();

        for testcase in testcases {
            let result = self.run(testcase.get_input_path(), testcase.get_output_path());
            if let Err(e) = result {
                debug!("{:?}", e);
            }
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
            .lock()
            .unwrap()
            .get_compile_command(&self.lang, &self.source_path, &self.target_path)
            .ok_or(CompilationError::MissingLang(format!(
                "missing lang `{}`",
                self.lang
            )))?
            .iter()
            .map(|s| c_string!(s.as_str()))
            .collect::<Vec<_>>();
        debug!("{:?}", command);

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

        let log =
            fs::read_to_string(&self.log_path).map_err(|_e| CompilationError::FileSystemError)?;
        if log.is_empty() {
            Ok(())
        } else {
            Err(CompilationError::Error(log))
        }
    }

    fn run(&mut self, input_path: &str, output_path: &str) -> Result<(), RunningError> {
        let command = GLOB_CONFIG
            .lock()
            .unwrap()
            .get_run_command(&self.lang, &self.source_path, &self.target_path)
            .ok_or(RunningError::MissingLang(format!(
                "missing lang `{}`",
                self.lang
            )))?
            .iter()
            .map(|s| c_string!(s.as_str()))
            .collect::<Vec<_>>();
        debug!("{:?}", command);
        debug!("input_path: {}, output_path: {}", input_path, output_path);

        let redirect_stdout_path = format!("{}.stdout", self.source_path);
        let redirect_stderr_path = format!("{}.stderr", self.source_path);

        match unsafe { fork() } {
            Ok(ForkResult::Parent { child, .. }) => {
                waitpid(child, None).unwrap();
            }
            Ok(ForkResult::Child) => {
                let input_path = c_string!(input_path);
                let redirect_stdout_path = c_string!(redirect_stdout_path.as_str());
                let redirect_stderr_path = c_string!(redirect_stderr_path.as_str());
                let r_mode = c_string!("r");
                let w_mode = c_string!("w");

                unsafe {
                    let stdin = fdopen(STDIN_FILENO, r_mode.as_ptr());
                    freopen(input_path.as_ptr(), r_mode.as_ptr(), stdin);
                    let stdout = fdopen(STDOUT_FILENO, w_mode.as_ptr());
                    freopen(redirect_stdout_path.as_ptr(), w_mode.as_ptr(), stdout);
                    let stderr = fdopen(STDERR_FILENO, w_mode.as_ptr());
                    freopen(redirect_stderr_path.as_ptr(), w_mode.as_ptr(), stderr);
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
            _ => return Err(RunningError::ForkFailed),
        };

        self.run_ctx.run_testcase_cnt += 1;
        if let Err(e) = Comparer::compare_two_files(output_path, &redirect_stdout_path) {
            debug!("{:?}", e);
        } else {
            self.run_ctx.passed_testcase_cnt += 1;
            // TODO: update these two field
            self.run_ctx.used_memory += 100;
            self.run_ctx.elapsed_time += 100;
        }

        Ok(())
    }

    fn update_db(&self, execute_result: &ExecutionResult) {
        // TODO: logic for updating db
        info!("update_db")
    }

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
