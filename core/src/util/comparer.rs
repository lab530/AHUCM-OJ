use std::{
    fs::File,
    io::{BufRead, BufReader},
};

pub struct Comparer {}

#[derive(Debug, PartialEq, Eq)]
pub enum ComparerError {
    LineNotEqual,
    EOFNotCorresponding,
}

#[derive(Debug)]
pub struct ComparerErrorWrapper {
    pub error: ComparerError,
    pub line_cnt: usize,
    pub lhs_line: String,
    pub rhs_line: String,
}

impl ComparerErrorWrapper {
    pub fn new(error: ComparerError, line_cnt: usize, lhs_line: String, rhs_line: String) -> Self {
        Self {
            error,
            line_cnt,
            lhs_line,
            rhs_line,
        }
    }
}

impl Comparer {
    pub fn compare_two_files(
        lhs_file_path: &str,
        rhs_file_path: &str,
    ) -> Result<usize, ComparerErrorWrapper> {
        let lhs_file = File::open(lhs_file_path).unwrap();
        let rhs_file = File::open(rhs_file_path).unwrap();

        let mut lhs_reader = BufReader::new(lhs_file);
        let mut rhs_reader = BufReader::new(rhs_file);

        let mut lhs_line = String::new();
        let mut rhs_line = String::new();

        let mut line_cnt: i32 = 1;
        loop {
            let lhs_line_len = match lhs_reader.read_line(&mut lhs_line) {
                Ok(len) => len as i64,
                Err(_) => -1,
            };
            let rhs_line_len = match rhs_reader.read_line(&mut rhs_line) {
                Ok(len) => len as i64,
                Err(_) => -1,
            };

            // both of files got EOF
            if lhs_line_len < 0 && rhs_line_len < 0 {
                break;
            }

            // two valid line were not equal
            if lhs_line_len >= 0 && rhs_line_len >= 0 && lhs_line != rhs_line {
                return Err(ComparerErrorWrapper::new(
                    ComparerError::LineNotEqual,
                    line_cnt as usize,
                    lhs_line,
                    rhs_line,
                ));
            }

            if lhs_line_len < 0 || rhs_line_len < 0 {
                if lhs_line_len < 0 && rhs_line_len == 0 {
                    // lhs got EOF but rhs did not
                    let next_rhs_line_len = rhs_reader.read_line(&mut rhs_line);
                    if next_rhs_line_len.is_err() {
                        break;
                    } else {
                        return Err(ComparerErrorWrapper::new(
                            ComparerError::EOFNotCorresponding,
                            line_cnt as usize,
                            String::new(),
                            String::new(),
                        ));
                    }
                } else if rhs_line_len < 0 && lhs_line_len == 0 {
                    // rhs got EOF but lhs did not
                    let next_lhs_line_len = lhs_reader.read_line(&mut lhs_line);
                    if next_lhs_line_len.is_err() {
                        break;
                    } else {
                        return Err(ComparerErrorWrapper::new(
                            ComparerError::EOFNotCorresponding,
                            line_cnt as usize,
                            String::new(),
                            String::new(),
                        ));
                    }
                } else {
                    return Err(ComparerErrorWrapper::new(
                        ComparerError::LineNotEqual,
                        line_cnt as usize,
                        if lhs_line_len < 0 {
                            String::new()
                        } else {
                            lhs_line
                        },
                        if rhs_line_len < 0 {
                            String::new()
                        } else {
                            rhs_line
                        },
                    ));
                }
            }

            line_cnt += 1;
        }

        Ok(line_cnt as usize)
    }
}
