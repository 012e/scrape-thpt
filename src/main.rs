use regex::Regex;
use std::collections::HashMap;
use rusqlite::{Connection, Result};
// use rusqlite::NO_PARAMS;

fn parse_input_string(input_string: &str, re: &Regex) -> HashMap<String, f32> {
    let mut subjects = HashMap::new();
    for capture in re.captures_iter(input_string) {
        let subject_name = capture.get(1).unwrap().as_str().to_string();
        let grade = capture.get(2).unwrap().as_str().parse::<f32>().unwrap();
        subjects.insert(subject_name, grade);
    }

    subjects
}

fn main() -> Result<()> {
    let conn = rusqlite::Connection::open("student.db")?;

    let pattern = r"(\w+(?: \w+)*)\s*:\s*(\d+\.\d+)";
    let re = Regex::new(pattern).unwrap();
    Ok(())
}
