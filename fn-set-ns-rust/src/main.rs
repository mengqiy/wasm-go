use std::io::{self, BufRead};
use set_ns::set_ns_transformer;

fn main() {
   let stdin = io::stdin();
   for line in stdin.lock().lines() {
       println!("{}", set_ns_transformer(&line.unwrap()));
   }
}
