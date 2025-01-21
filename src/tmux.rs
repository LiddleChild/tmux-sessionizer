use std::{process::Command, str::from_utf8};

use chrono::{DateTime, Local};

use crate::utils::str_to_datetime;

#[derive(Debug)]
pub struct Session {
    pub name: String,
    pub created_at: DateTime<Local>,
    pub is_attached: bool,
}

impl Session {
    fn new_from_string(s: &str) -> Result<Self, &'static str> {
        let elements: Vec<&str> = s.split(':').map(|s| s.trim()).collect();

        if elements.len() < 3 {
            return Err("invalid string format");
        }

        let name = String::from(elements[0]);

        let created_at: DateTime<Local> = str_to_datetime(elements[1])?;

        let is_attached = match elements[2] {
            "1" => true,
            _ => false,
        };

        Ok(Self {
            name,
            created_at,
            is_attached,
        })
    }
}

pub fn list_sessions() -> Result<Vec<Session>, &'static str> {
    let output = Command::new("tmux")
        .args([
            "list-sessions",
            "-F",
            "#{session_name}:#{session_created}:#{session_attached}",
        ])
        .output();

    let output = match output {
        Ok(output) => output,
        Err(_) => return Err("cannot list tmux sessions"),
    };

    let stdout = match from_utf8(&output.stdout) {
        Ok(stdout) => stdout,
        Err(_) => return Err("cannot read stdout"),
    };

    let lines: Result<Vec<Session>, _> = stdout
        .lines()
        .map(|line| Session::new_from_string(line))
        .collect();

    lines
}
