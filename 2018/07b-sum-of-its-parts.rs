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
        deps[d.1 as usize - b'A' as usize] |= 4 << (d.0 as u32 - b'A' as u32) | 1;
        deps[d.0 as usize - b'A' as usize] |= 1;
    }

    let mut t: i32 = 0;
    let mut free_at: Vec<i32> = vec![0, 0, 0, 0, 0];
    let mut current_task: Vec<i32> = vec![-1, -1, -1, -1, -1];

    loop {
        let free_workers: Vec<usize> = free_at
            .iter()
            .enumerate()
            .filter(|(_, c)| **c <= t)
            .map(|(idx, _)| idx)
            .collect();

        for w in free_workers.iter() {
            let idx = current_task[*w];
            if idx < 0 || deps[idx as usize] == 0 {
                continue;
            }
            deps[idx as usize] = 0;
            for i in deps.iter_mut() {
                *i &= !(4 << idx);
            }
        }

        let next_tasks: Vec<usize> = deps
            .iter()
            .enumerate()
            .filter(|(_, c)| **c == 1)
            .map(|(idx, _)| idx)
            .take(free_workers.len())
            .collect();

        for (n, w) in next_tasks.into_iter().zip(free_workers.into_iter()) {
            free_at[w] = t + 60 + 1 + n as i32;
            current_task[w] = n as i32;
            deps[n] |= 2;
        }

        if let Some(x) = free_at.iter().filter(|c| **c > t).min() {
            t = *x;
        } else {
            break;
        }
    }

    println!("{}", t);
}

#[derive(Copy, Clone, Debug, Hash, Eq, PartialEq)]
struct Dep(char, char);
