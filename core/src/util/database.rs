use std::sync::Mutex;

use once_cell::sync::Lazy;
use postgres::{Client, NoTls};

use super::config::GLOB_CONFIG;

pub static GLOB_DATABASE: Lazy<Mutex<Database>> = Lazy::new(|| Mutex::new(Database::default()));
static MAIN_TBL_NAME: &str = "submissions";

pub struct Database {
    host: String,
    port: String,
    username: String,
    password: String,
    database: String,
    client: Client,
}

impl Default for Database {
    fn default() -> Self {
        let sql_config = GLOB_CONFIG.lock().unwrap().sql_config.clone();
        let host = sql_config.get("host").unwrap().clone();
        let port = sql_config.get("port").unwrap().clone();
        let username = sql_config.get("username").unwrap().clone();
        let password = sql_config.get("password").unwrap().clone();
        let database = sql_config.get("database").unwrap().clone();
        let url = format!("postgresql://{username}:{password}@{host}:{port}/{database}");
        log::debug!("Postgres URL: {}", url);

        Self {
            client: Client::connect(url.as_str(), NoTls).unwrap(),
            host,
            port,
            username,
            password,
            database,
        }
    }
}

impl Database {
    pub fn update_record(&mut self, id: u64, time: u32, mem: u32, status: i32) {
        let query = format!("UPDATE {MAIN_TBL_NAME} SET (time_used, memo_used, status) = ($1, $2, $3) WHERE id = $4");
        let id = id as i64;
        let time = time as i64;
        let mem = mem as i64;
        let status = status as i64;
        self.client.execute(&query, &[&time, &mem, &status, &id]).unwrap();
    }
}
