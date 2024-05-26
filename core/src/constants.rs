use once_cell::sync::Lazy;

use crate::util::config::GLOB_CONFIG;

pub static CONFIG_PATH: &str = "../config.toml";

pub static SERVER_ADDR: Lazy<String> = Lazy::new(|| format!("127.0.0.1:{}", GLOB_CONFIG.lock().unwrap().get_port()));
