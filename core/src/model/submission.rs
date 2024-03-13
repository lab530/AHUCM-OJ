use serde::{Deserialize, Serialize};

#[derive(Default, Debug, Serialize, Deserialize)]
pub struct Submission {
    pub source_path: String,
    pub lang: String,
    pub problem_id: u64,
    pub mem_limit: u32,
    pub time_limit: u32,
    pub testcases_path: String,
    pub submission_id: u64,
}
