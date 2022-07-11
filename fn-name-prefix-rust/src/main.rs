use std::io::{self, BufRead};
use regex::Regex;

fn main() {
    let stdin = io::stdin();
    for line in stdin.lock().lines() {
        println!("{}", name_prefix_transformer(&line.unwrap()));
    }
}

pub fn name_prefix_transformer(input: &str) -> String {
    let re = Regex::new("name: (?P<name>.*)").unwrap();
    return re.replace(input, "name: pre-$name").to_string()
}

