use regex::Regex;

pub fn set_ns_transformer(input: &str) -> String {
    let re = Regex::new("namespace: .*$").unwrap();
    return re.replace(input, "namespace: test").to_string()
}
