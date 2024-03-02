use std::{collections::BTreeMap, fs};

use super::testcase::Testcase;

pub struct TestcasesGetter {
    dir_path: String,
    in_file_paths: BTreeMap<String, String>,
    out_file_paths: BTreeMap<String, String>,
    testcases: Option<Vec<Testcase>>,
}

impl TestcasesGetter {
    pub fn new(dir_path: String) -> Self {
        Self {
            dir_path,
            in_file_paths: BTreeMap::new(),
            out_file_paths: BTreeMap::new(),
            testcases: None,
        }
    }

    fn fetch_testcases(&mut self) {
        let dir = match fs::read_dir(&self.dir_path) {
            Ok(dir) => dir,
            Err(_) => return,
        };

        for file in dir {
            let file = file.unwrap();
            if file.file_type().unwrap().is_file() {
                let file_path = file.file_name().into_string().unwrap();
                if file_path.ends_with(".in") {
                    let (prefix, _) = file_path.split_at(file_path.len() - ".in".len());
                    self.in_file_paths.insert(prefix.into(), file_path.clone());
                } else if file_path.ends_with(".out") {
                    let (prefix, _) = file_path.split_at(file_path.len() - ".out".len());
                    self.out_file_paths.insert(prefix.into(), file_path.clone());
                }
            }
        }

        let mut testcases = vec![];

        for (key, in_file_path) in self.in_file_paths.iter() {
            match self.out_file_paths.get(key) {
                Some(out_file_path) => {
                    testcases.push(Testcase::new(in_file_path.clone(), out_file_path.clone()));
                }
                None => {}
            }
        }

        self.testcases = Some(testcases);
    }

    pub fn get_testcases(&mut self) -> &Vec<Testcase> {
        if self.testcases.is_none() {
            self.fetch_testcases();
        }
        self.testcases
            .as_ref()
            .unwrap_or_else(|| panic!("can't fetch testcases in {}", self.dir_path))
    }
}
