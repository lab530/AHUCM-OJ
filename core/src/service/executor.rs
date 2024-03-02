use std::sync::{
    atomic::{AtomicU32, AtomicUsize, Ordering},
    Arc,
};

use super::testcases_getter::TestcasesGetter;

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

    problem_id: u64,
    mem_limit: u32,
    time_limit: u32,
    testcases_path: String,

    submission_id: u64,

    is_compiling_done: bool,
    is_running_done: bool,

    run_ctx: Arc<RunContext>,
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
            problem_id,
            mem_limit,
            time_limit,
            testcases_path,
            submission_id,
            is_compiling_done: false,
            is_running_done: false,
            run_ctx: Arc::new(RunContext::default()),
        }
    }

    pub fn execute(&mut self) -> ExecutionResult {
        if let Some(err_msg) = self.compile() {
            return ExecutionResult::CompilationError(err_msg);
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
                self.run_ctx.clone(),
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

    fn compile(&self) -> Option<String> {
        // regular POSIX workflow

        None
    }

    fn run(&self, input_path: &str, output_path: &str, ctx: Arc<RunContext>) {
        // regular POSIX workflow
    }

    fn update_db(&self, execute_result: &ExecutionResult) {
        // TODO: logic for updating db
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
