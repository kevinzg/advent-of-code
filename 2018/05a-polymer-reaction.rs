use std::collections::VecDeque;
use std::io;

enum ReactionResult {
    Boom,
    Nope,
}

fn react(pol: &Vec<char>, i: i32, j: i32) -> ReactionResult {
    if i < 0 {
        return ReactionResult::Nope;
    }

    let p = pol[i as usize];
    let q = pol[j as usize];

    if p.to_ascii_lowercase() == q.to_ascii_lowercase() {
        if p.is_ascii_lowercase() == !q.is_ascii_lowercase() {
            ReactionResult::Boom
        } else {
            ReactionResult::Nope
        }
    } else {
        ReactionResult::Nope
    }
}

fn main() {
    let mut pol = String::new();
    io::stdin().read_line(&mut pol).unwrap();

    let pol: Vec<char> = pol.trim().chars().collect();

    let mut st: VecDeque<i32> = VecDeque::new();
    let mut curr: i32 = 0;

    while curr < pol.len() as i32 {
        let prev: i32 = *st.back().unwrap_or(&-1);

        match react(&pol, prev, curr) {
            ReactionResult::Boom => {
                st.pop_back();
            }
            ReactionResult::Nope => {
                st.push_back(curr);
            }
        }
        curr += 1;
    }

    println!("{}", st.len())
}
