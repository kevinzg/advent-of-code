use std::io;

fn main() {
    let mut twos = 0;
    let mut threes = 0;

    loop {
        let mut id = String::new();

        match io::stdin().read_line(&mut id) {
            Ok(0) => break,
            Ok(_) => (),
            Err(err) => panic!(err),
        }

        id = id.trim().to_string();

        let c = count(&id);

        twos += c.0;
        threes += c.1;
    }

    println!("{}", twos * threes);
}

fn count(id: &String) -> (i32, i32) {
    let mut counter: [i32; 26] = [0; 26];

    for c in id.chars() {
        let key = (c as u8 - 'a' as u8) as usize;
        counter[key] += 1;
    }

    let mut twos = 0;
    let mut threes = 0;

    for c in counter.iter() {
        if *c == 2 {
            twos = 1;
        }

        if *c == 3 {
            threes = 1;
        }
    }

    return (twos, threes);
}
