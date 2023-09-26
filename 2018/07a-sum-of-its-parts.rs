use std::io;

fn main() {
    let mut deps: Vec<u32> = Vec::with_capacity(26);
    deps.resize(26, 0);

    loop {
        let mut line = String::new();

        match io::stdin().read_line(&mut line) {
            Ok(0) => break,
            Ok(_) => (),
            Err(err) => panic!("{}", err),
        }

        let words: Vec<_> = line.split(" ").collect();

        let d = Dep(
            words[1].chars().next().unwrap(),
            words[7].chars().next().unwrap(),
        );
        deps[d.1 as usize - b'A' as usize] |= 2 << (d.0 as u32 - b'A' as u32) | 1;
        deps[d.0 as usize - b'A' as usize] |= 1;
    }

    let mut solution = String::new();

    loop {
        let idx = if let Some((idx, _)) = deps.iter().enumerate().find(|(_, c)| **c == 1) {
            idx
        } else {
            break;
        };
        deps[idx] = 0;
        solution.push(char::from_u32(b'A' as u32 + idx as u32).unwrap());
        for i in deps.iter_mut() {
            *i &= !(2 << idx as u32);
        }
    }

    println!("{}", solution);
}

#[derive(Copy, Clone, Debug, Hash, Eq, PartialEq)]
struct Dep(char, char);
