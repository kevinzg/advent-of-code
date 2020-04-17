use std::collections::HashSet;
use std::io;

fn main() {
    let mut frequency = 0;

    let mut changes: Vec<i32> = Vec::new();

    loop {
        let mut change = String::new();
        match io::stdin().read_line(&mut change) {
            Ok(0) => break,
            Ok(_) => (),
            Err(err) => panic!(err),
        }

        let change: i32 = change.trim().parse().expect("Wrong input");

        changes.push(change)
    }

    let mut counter = HashSet::new();
    counter.insert(0);

    let mut found = false;

    while !found {
        for change in &changes {
            frequency += change;

            if counter.contains(&frequency) {
                found = true;
                break;
            }

            counter.insert(frequency);
        }
    }

    println!("{}", frequency);
}
