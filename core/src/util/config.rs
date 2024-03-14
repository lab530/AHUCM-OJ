use log::{debug, error, warn};
use std::{
    collections::{BTreeMap, BTreeSet},
    fs, sync::Mutex,
};
use once_cell::sync::Lazy;

use crate::constants::CONFIG_PATH;

pub static GLOB_CONFIG: Lazy<Mutex<Config>> = Lazy::new(|| {
    Mutex::new(Config::new())
});

#[derive(Debug, Default)]
pub struct Config {
    pub sql_config: BTreeMap<String, String>,
    support_langs: BTreeSet<String>,
    compile_commands: BTreeMap<String, Vec<String>>,
    run_commands: BTreeMap<String, Vec<String>>,
    compiler_output_channels: BTreeMap<String, &'static str>,
}

impl Config {
    pub fn new() -> Self {
        let mut config = Self::default();
        if let Some(err) = config.reload() {
            error!("{}", err);
        }
        config
    }

    pub fn reload(&mut self) -> Option<String> {
        let mut new_self = Self::default();

        let config_raw_text = match fs::read_to_string(CONFIG_PATH) {
            Ok(s) => s,
            Err(e) => return Some(format!("config can't be read or missing: {e}")),
        };

        let root: toml::Table = match toml::from_str(&config_raw_text) {
            Ok(root) => root,
            Err(e) => return Some(format!("wrong config format: {e}")),
        };

        let sql_list = match root.get("sql") {
            Some(v) => v,
            _ => return Some("wrong config format, missing sql".into()),
        };
        if let toml::Value::Table(sql_list) = sql_list {
            for (key, value) in sql_list {
                if let toml::Value::String(value) = value {
                    self.sql_config.insert(key.clone(), value.clone());
                } else {
                    return Some(format!("wrong config format, sql.{key} should have a string as the value"));
                }
            }
        } else {
            return Some("wrong config format, [sql] should be a table".into());
        }

        let lang_list = match root.get("languages") {
            Some(v) => v,
            _ => return Some("wrong config format, missing languages".into()),
        };
        if let toml::Value::Array(lang_list) = lang_list {
            for lang in lang_list {
                if let toml::Value::String(lang) = lang {
                    new_self.support_langs.insert(lang.clone());
                } else {
                    return Some(
                        "wrong config format, there are not all strings in languages".into(),
                    );
                }
            }
        }

        let compile_list = match root.get("compile") {
            Some(v) => v,
            _ => return Some("wrong config format, missing [compile]".into()),
        };
        if let toml::Value::Table(compile_list) = compile_list {
            for (lang, command) in compile_list {
                if !new_self.support_langs.contains(lang) {
                    continue;
                }
                if new_self.compile_commands.contains_key(lang) {
                    warn!("duplicated key {lang} in [compile]");
                    continue;
                }
                let command = match command {
                    toml::Value::String(command) => command,
                    _ => return Some(format!("wrong config format, table value should be a string in [compile], key = {lang}")),
                };
                new_self.compile_commands.insert(
                    lang.clone(),
                    command
                        .split_ascii_whitespace()
                        .map(|e| e.to_string())
                        .collect(),
                );
            }
        } else {
            return Some("wrong config format, [compile] should be a table".into());
        }
        for lang in new_self.support_langs.iter() {
            if !new_self.compile_commands.contains_key(lang) {
                new_self.compile_commands.insert(lang.clone(), vec![]);
            }
        }

        let run_list = match root.get("run") {
            Some(v) => v,
            _ => return Some("wrong config format, missing [run]".into()),
        };
        if let toml::Value::Table(run_list) = run_list {
            for (lang, command) in run_list {
                if !new_self.support_langs.contains(lang) {
                    continue;
                }
                if new_self.run_commands.contains_key(lang) {
                    warn!("duplicated key {lang} in [run]");
                    continue;
                }
                let command = match command {
                    toml::Value::String(command) => command,
                    _ => return Some(format!("wrong config format, table value should be a string in [run], key = {lang}")),
                };
                new_self.run_commands.insert(
                    lang.clone(),
                    command
                        .split_ascii_whitespace()
                        .map(|e| e.to_string())
                        .collect(),
                );
            }
        }
        for lang in new_self.support_langs.iter() {
            if !new_self.run_commands.contains_key(lang) {
                return Some(format!(
                    "wrong config format, no run recipe for language {lang}"
                ));
            }
        }

        let channel_list = match root.get("channel") {
            Some(v) => v,
            _ => return Some("wrong config format, missing [channel]".into()),
        };
        if let toml::Value::Table(channel_list) = channel_list {
            for (lang, channel) in channel_list {
                if !new_self.support_langs.contains(lang) {
                    continue;
                }
                if new_self.compiler_output_channels.contains_key(lang) {
                    warn!("duplicated key {lang} in [channel]");
                    continue;
                }
                let channel = match channel {
                     toml::Value::String(channel) => channel,
                    _ => return Some(format!("wrong config format, table value should be a string in [channel], key = {lang}")),
                };

                static SUPPORT_CHANNELS: [&str; 2] = ["stdout", "stderr"];
                let mut inserted = false;
                for support_channel in SUPPORT_CHANNELS {
                    if channel == support_channel {
                        new_self
                            .compiler_output_channels
                            .insert(lang.clone(), support_channel);
                        inserted = true;
                        break;
                    }
                }
                if !inserted {
                    return Some(format!("wrong config format, key = {lang} in [channel] got a value = {channel}, should be in {:?}", SUPPORT_CHANNELS));
                }
            }
        }
        for lang in new_self.support_langs.iter() {
            if !new_self.compiler_output_channels.contains_key(lang) {
                new_self
                    .compiler_output_channels
                    .insert(lang.clone(), "stderr");
            }
        }

        self.support_langs = new_self.support_langs;
        self.compile_commands = new_self.compile_commands;
        self.run_commands = new_self.run_commands;
        self.compiler_output_channels = new_self.compiler_output_channels;

        debug!("got a new config: {:?}", self);

        None
    }

    fn get_command_aux(
        &self,
        tree: &BTreeMap<String, Vec<String>>,
        lang: &str,
        source_path: &str,
        target_path: &str,
    ) -> Option<Vec<String>> {
        if let Some(command) = tree.get(lang) {
            let mut ret: Vec<String> = vec![];
            for s in command {
                if s == "$source" {
                    ret.push(source_path.to_string());
                } else if s == "$target" {
                    ret.push(target_path.to_string());
                } else {
                    ret.push(s.to_string());
                }
            }
            Some(ret)
        } else {
            None
        }
    }

    pub fn get_compile_command(
        &self,
        lang: &str,
        source_path: &str,
        target_path: &str,
    ) -> Option<Vec<String>> {
        self.get_command_aux(&self.compile_commands, lang, source_path, target_path)
    }

    pub fn get_run_command(
        &self,
        lang: &str,
        source_path: &str,
        target_path: &str,
    ) -> Option<Vec<String>> {
        self.get_command_aux(&self.run_commands, lang, source_path, target_path)
    }

    pub fn get_compiler_output_channel(&self, lang: &str) -> Option<&str> {
        self.compiler_output_channels.get(lang).map(|s| *s)
    }
}
