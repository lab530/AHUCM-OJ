use std::fs;

use nix::sys::stat::stat;

pub fn get_proc_status(pid: i32, mark: &str) -> Option<u32> {
    let read = fs::read_to_string(format!("/proc/{pid}/status")).unwrap();
    let lines = read.split('\n');

    for line in lines {
        if line.starts_with(mark) {
            let numeric = line
                .strip_prefix(mark)
                .unwrap()
                .trim()
                .split(' ')
                .collect::<Vec<_>>()[0]
                .parse::<u32>()
                .unwrap();
            return Some(numeric);
        }
    }

    None
}

pub fn get_file_size(path: &str) -> usize {
    let f_stat = stat(path).unwrap();
    f_stat.st_size as usize
}
