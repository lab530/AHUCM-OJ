use crate::util::comparer::{Comparer, ComparerErrorWrapper};

pub struct Testcase {
    input_path: String,
    output_path: String,
}

pub enum TestcaseResult {
    Accpected,
    WrongAnswer(usize, String, String),
    PresentationError,
}

impl Testcase {
    pub fn new(input_path: String, output_path: String) -> Self {
        Self {
            input_path,
            output_path,
        }
    }

    pub fn get_input_path(&self) -> &str {
        &self.input_path
    }

    pub fn get_output_path(&self) -> &str {
        &self.output_path
    }

    pub fn compare_output_with(&self, test_output_path: &str) -> TestcaseResult {
        match Comparer::compare_two_files(&self.output_path, test_output_path) {
            Ok(_) => TestcaseResult::Accpected,
            Err(ComparerErrorWrapper {
                error,
                line_cnt,
                lhs_line,
                rhs_line,
            }) => {
                if ComparerError::LineNotEqual == error {
                    TestcaseResult::WrongAnswer(line_cnt, lhs_line, rhs_line)
                } else {
                    TestcaseResult::PresentationError
                }
            }
        }
    }
}
