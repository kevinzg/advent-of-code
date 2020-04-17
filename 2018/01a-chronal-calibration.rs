use std::io;

fn main() {
    let mut frequency = 0;

    loop {
        let mut change = String::new();
        match io::stdin().read_line(&mut change) {
            Ok(0) => break,
            Ok(_) => (),
            Err(err) => panic!(err),
        }

        let change: i32 = change.trim().parse().expect("Wrong input");

        frequency += change;
    }

    println!("{}", frequency);
}
