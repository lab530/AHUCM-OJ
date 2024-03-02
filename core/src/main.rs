mod api;
mod model;
mod service;
mod util;
mod constants;

use util::{config::Config, thread_pool::ThreadPool};

fn main() {
    // let thread_pool = ThreadPool::default();
    // let exceutor = Executor::new(1001, 100, 100, "./assets/1001".into());
    env_logger::init();
    let _config = Config::new();
}
