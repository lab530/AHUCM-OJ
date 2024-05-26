use std::{
    ffi::CString,
    fs::{self, File},
    os::fd::FromRawFd,
};

use log::debug;
use nix::{
    libc::{alarm, exit, fdopen, freopen, ptrace, rlimit, rusage, setrlimit, setrlimit64, wait4, PTRACE_KILL, PTRACE_O_TRACEEXIT, PTRACE_O_TRACESYSGOOD, PTRACE_SETOPTIONS, RLIMIT_AS, RLIMIT_CPU, RLIMIT_FSIZE, RLIMIT_NPROC, RLIMIT_STACK, SIGALRM, SIGCHLD, SIGKILL, SIGXCPU, SIGXFSZ, STDERR_FILENO, STDIN_FILENO, STDOUT_FILENO, WEXITSTATUS, WIFEXITED, WIFSIGNALED, WTERMSIG, __WALL},
    sys::wait::waitpid,
    unistd::{execvp, fork, write, ForkResult},
};

use crate::util::{comparer::Comparer, config::GLOB_CONFIG, database::GLOB_DATABASE, unix::{get_file_size, get_proc_status}};

use super::testcases_getter::TestcasesGetter;

macro_rules! c_string {
    ($string:expr) => {
        CString::new($string).unwrap()
    };
}

const UNIT_MB: u64 = 1048576;
const FSIZE_LIM: u64 = UNIT_MB << 7;
const MEM_LIM: u64 = UNIT_MB << 7;
const PROC_LIM: u64 = 750;

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

#[derive(Debug, PartialEq, Eq)]
pub enum ExecutionResult {
    Halt,
    Compiling,
    Running,
    Accpected(usize, u32, u32),
    WrongAnswer(usize, usize),
    PresentationError,
    CompilationError(String),
    RuntimeError(String),
    TimeLimitExceeded(u32, u32),
    MemoLimitExceeded(u32, u32),
	OutputLimitExceeded,
    UnknownError(String),
}

/*
	Pending 0
	PendingRejudge 1
	Compiling 2
	Running 3
	Accepted 4
	PresentationError 5
	WrongAnswer 6
	TimeLimitExceeded 7
	MemoryLimitExceeded 8
	OutputLimitExceeded 9
	RuntimeError 10
	CompileError 11
	UnknownError 12
* */

impl From<&ExecutionResult> for i32 {
    fn from(value: &ExecutionResult) -> Self {
        match value {
            ExecutionResult::Halt => 12,
            ExecutionResult::Compiling => 2,
            ExecutionResult::Running => 3,
            ExecutionResult::Accpected(_, _, _) => 4,
            ExecutionResult::WrongAnswer(_, _) => 6,
            ExecutionResult::PresentationError => 5,
            ExecutionResult::CompilationError(_) => 11,
            ExecutionResult::RuntimeError(_) => 10,
            ExecutionResult::TimeLimitExceeded(_, _) => 7,
            ExecutionResult::MemoLimitExceeded(_, _) => 8,
            ExecutionResult::OutputLimitExceeded => 9,
            ExecutionResult::UnknownError(_) => 12,
        }
    }
}

#[derive(Default)]
struct RunContext {
    pub max_testcase_cnt: usize,
    pub run_testcase_cnt: usize,
    pub wrong_answer_cnt: usize,
    pub mem_exceeded_cnt: usize,
    pub runtime_error_cnt: usize,
    pub output_exceeded_cnt: usize,
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

    pub fn wrong_answer_cnt(&self) -> usize {
        // self.run_testcase_cnt.load(Ordering::SeqCst)
        self.wrong_answer_cnt
    }

    pub fn mem_exceeded_cnt(&self) -> usize {
        self.mem_exceeded_cnt
    }

    pub fn runtime_error_cnt(&self) -> usize {
        self.runtime_error_cnt
    }

    pub fn output_exceeded_cnt(&self) -> usize {
        self.output_exceeded_cnt
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
        self.update_db(&ExecutionResult::Compiling);
        if let Err(e) = self.compile() {
            let compilation_error = ExecutionResult::CompilationError(format!("{:?}", e));
            self.update_db(&compilation_error);
            return compilation_error;
        }

        self.update_db(&ExecutionResult::Compiling);
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
        // TODO: update result branching logic
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
        } else if self.run_ctx.wrong_answer_cnt() != 0 {
            result = ExecutionResult::Accpected(
                self.run_ctx.max_testcase_cnt(),
                self.run_ctx.elapsed_time(),
                self.run_ctx.used_memory(),
            );
        } else {
            result = ExecutionResult::WrongAnswer(
                self.run_ctx.max_testcase_cnt() - self.run_ctx.wrong_answer_cnt(),
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
                let cpu = 50u32;
                let rlim_cpu = rlimit { rlim_cur: cpu as u64, rlim_max: cpu as u64 };
                unsafe {
                    setrlimit(RLIMIT_CPU, &rlim_cpu);
                    alarm(0);
                    if cpu > 0 {
                        alarm(cpu);
                    } else {
                        alarm(1);
                    }
                }

                let fsize = 500 * UNIT_MB;
                let rlim_fsize = rlimit { rlim_cur: fsize, rlim_max: fsize };
                unsafe { setrlimit(RLIMIT_FSIZE, &rlim_fsize) };

                let mem = UNIT_MB << 12;
                let rlim_as = rlimit { rlim_cur: mem, rlim_max: mem };
                unsafe { setrlimit(RLIMIT_AS, &rlim_as) };
 
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

        let mut res = ExecutionResult::Halt;

        match unsafe { fork() } {
            Ok(ForkResult::Parent { child, .. }) => {
                debug!("Judging, child pid: {}", child.as_raw());
                let mut status: i32 = 0;
                let mut ruse: rusage = unsafe { std::mem::zeroed() };
                let mut first = true;
                let mut tick: usize = 0;

                let std_output_size = get_file_size(output_path);
                loop {
                    tick += 1;
                    unsafe { wait4(child.as_raw(), &mut status, __WALL, &mut ruse) };
                    log::debug!("pid: {:?}, status: {status}", child.as_raw());
                    if first {
                        unsafe { ptrace(PTRACE_SETOPTIONS, child.as_raw(), 0, PTRACE_O_TRACESYSGOOD | PTRACE_O_TRACEEXIT) };
                        first = false;
                    }

                    let temp = get_proc_status(child.as_raw(), "VmPeak:").unwrap() as u64;
                    log::debug!("temp: {}", temp);
                    if temp > self.mem_limit as u64 * UNIT_MB {
                        unsafe { ptrace(PTRACE_KILL, child.as_raw(), 0, 0) };
                        res = ExecutionResult::MemoLimitExceeded(temp as u32, self.mem_limit);
                        break;
                    }
                    if temp > self.run_ctx.used_memory as u64 {
                        self.run_ctx.used_memory = temp as u32;
                    }

                    // runtime error: ret val is not 0
                    if WIFEXITED(status) && WEXITSTATUS(status) != 0 {
                        res = ExecutionResult::RuntimeError("111".into());
                        break;
                    }

                    if tick % 250 == 0 {
                        // output size is greater than double of standard output's or stderr
                        if get_file_size(&redirect_stdout_path) > std_output_size * 2 || get_file_size(&redirect_stderr_path) > 0 {
                            unsafe { ptrace(PTRACE_KILL, child.as_raw(), 0, 0)};
                            res = ExecutionResult::RuntimeError("222".into());
                            break;
                        }
                    }

                    let exitcode = WEXITSTATUS(status) % 256;
                    if exitcode != 0 && exitcode != 5 && exitcode != 17 && exitcode != 23 && exitcode != 133 {
                        if exitcode == SIGCHLD || exitcode == SIGALRM {
                            unsafe { alarm(0) };
                            res = ExecutionResult::RuntimeError("333".into());
                            self.run_ctx.elapsed_time = self.time_limit;
                        } else if exitcode == SIGKILL || exitcode == SIGXCPU {
                            res = ExecutionResult::RuntimeError("444".into());
                            self.run_ctx.elapsed_time = self.time_limit;
                        } else if exitcode == SIGXFSZ {
                            res = ExecutionResult::OutputLimitExceeded;
                            self.run_ctx.output_exceeded_cnt += 1;
                        } else {
                            res = ExecutionResult::RuntimeError("555".into());
                            self.run_ctx.elapsed_time = self.time_limit;
                        }
                        unsafe { ptrace(PTRACE_KILL, child.as_raw(), 0, 0)};
                        break;
                    }

                    if WIFSIGNALED(status) {
                        let sig = WTERMSIG(status);

                        if sig == SIGCHLD || sig == SIGALRM {
                            unsafe { alarm(0) };
                            res = ExecutionResult::RuntimeError("666".into());
                            self.run_ctx.runtime_error_cnt += 1;
                        } else if sig == SIGKILL || sig == SIGXCPU {
                            res = ExecutionResult::RuntimeError("777".into());
                            self.run_ctx.runtime_error_cnt += 1;
                        } else if sig == SIGXFSZ {
                            res = ExecutionResult::OutputLimitExceeded;
                            self.run_ctx.output_exceeded_cnt += 1;
                        } else {
                            res = ExecutionResult::RuntimeError("888".into());
                            self.run_ctx.runtime_error_cnt += 1;
                        }
                        break;
                    }
                }
                
                unsafe { ptrace(PTRACE_KILL, child.as_raw(), 0, 0) };

                if res == ExecutionResult::Halt {
                    self.run_ctx.elapsed_time += (ruse.ru_utime.tv_sec * 1000 + ruse.ru_utime.tv_usec / 1000 + ruse.ru_stime.tv_sec * 1000 + ruse.ru_stime.tv_usec / 1000) as u32;
                }
            }
            Ok(ForkResult::Child) => {
                let time = self.time_limit - self.run_ctx.elapsed_time() + 1;
                let rlim_cpu = rlimit { rlim_cur: time as u64, rlim_max: time as u64 + 1 };
                unsafe {
                    setrlimit(RLIMIT_CPU, &rlim_cpu);
                    alarm(0);
                    alarm(if self.time_limit > 1 { self.time_limit } else { 1 });
                }

                let rlim_fsize = rlimit { rlim_cur: FSIZE_LIM, rlim_max: FSIZE_LIM + UNIT_MB };
                unsafe { setrlimit(RLIMIT_FSIZE, &rlim_fsize) };

                let rlim_proc = rlimit { rlim_cur: PROC_LIM, rlim_max: PROC_LIM };
                unsafe { setrlimit(RLIMIT_NPROC, &rlim_proc) };

                let rlim_stack = rlimit { rlim_cur: UNIT_MB << 8, rlim_max: UNIT_MB << 8 };
                unsafe { setrlimit(RLIMIT_STACK, &rlim_stack) };

                let rlim_mem = rlimit { rlim_cur: self.mem_limit as u64 / 2 * 3 * UNIT_MB, rlim_max: self.mem_limit as u64 * 2 * UNIT_MB };
                unsafe { setrlimit(RLIMIT_AS, &rlim_mem) };

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
        if res == ExecutionResult::Halt {
            if let Err(e) = Comparer::compare_two_files(output_path, &redirect_stdout_path) {
                debug!("{:?}", e);
                res = ExecutionResult::WrongAnswer(1, 1);
            } else {
                res = ExecutionResult::Accpected(1, 1, 1);
            }
        }

        match res {
            ExecutionResult::Accpected(_, _, _) => {}
            ExecutionResult::WrongAnswer(_, _) => self.run_ctx.wrong_answer_cnt += 1,
            ExecutionResult::RuntimeError(_) => self.run_ctx.runtime_error_cnt += 1,
            ExecutionResult::OutputLimitExceeded => self.run_ctx.output_exceeded_cnt += 1,
            ExecutionResult::TimeLimitExceeded(_, _) => self.run_ctx.mem_exceeded_cnt += 1,
            _ => {}
        }

        Ok(())
    }

    fn update_db(&self, execute_result: &ExecutionResult) {
        GLOB_DATABASE.lock().unwrap().update_record(
            self.submission_id,
            self.run_ctx.elapsed_time(),
            self.run_ctx.used_memory(),
            execute_result.into(),
        );
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
