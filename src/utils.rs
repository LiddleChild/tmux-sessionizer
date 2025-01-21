use chrono::{DateTime, Local, TimeZone};

pub fn str_to_datetime(str: &str) -> Result<DateTime<Local>, &'static str> {
    let timestamp: i64 = match str.parse() {
        Ok(timestamp) => timestamp,
        Err(_) => return Err("cannot convert timestamp"),
    };

    let timestamp = match Local.timestamp_opt(timestamp, 0) {
        chrono::offset::LocalResult::Single(timestamp) => timestamp,
        _ => return Err("cannot convert timestamp"),
    };

    Ok(timestamp)
}
