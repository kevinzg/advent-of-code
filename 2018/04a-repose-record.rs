use std::collections::HashMap;
use std::io;

type GuardId = String;

#[derive(Clone)]
struct Interval(Vec<char>);

impl Interval {
    fn new() -> Self {
        Interval(
            ['.'].repeat(60)
        )
    }

    fn score(&self) -> i32 {
        self.0.iter().map(|c| if *c == '#' { 1 } else { 0 }).fold(0, |acc, x| acc + x)
    }
}

fn main() {
    let mut map: HashMap<GuardId, Vec<Interval>> = HashMap::new();
    let mut events: Vec<String> = vec![];

    loop {
        let mut line = String::new();

        match io::stdin().read_line(&mut line) {
            Ok(0) => break,
            Ok(_) => (),
            Err(err) => panic!(err),
        }

        events.push(line);
    }

    events.sort_unstable();

    let mut current_guard = String::new();
    let mut guard_status = Interval::new();

    for event in events {
        if event.contains("Guard") {
            if current_guard != "" {
                let v = map.entry(current_guard.to_owned()).or_insert(vec![]);
                v.push(guard_status.clone());

                guard_status = Interval::new();
            }

            current_guard = event
                .split(" ")
                .skip_while(|t| *t != "Guard")
                .skip(1)
                .next()
                .unwrap()
                .to_owned();

        } else {
            let event_type = if event.contains("falls asleep") {
                '#'
            } else {
                '.'
            };

            let mut time = event.split(' ').skip(1).next().unwrap()
                .split(|c: char| c.is_ascii_punctuation());
            let hour: i32 = time.next().unwrap().parse().unwrap();
            let minute: i32 = time.next().unwrap().parse().unwrap();

            assert_eq!(hour, 0);

            for i in minute..60 {
                guard_status.0[i as usize] = event_type;
            }
        }
    }

    let mut best_guard = String::new();
    let mut best_score = 0;

    for (gid, v) in map.iter() {
        let score = v.iter().map(|i| i.score()).fold(0, |acc, x| acc + x);
        if score > best_score {
            best_guard = gid.clone();
            best_score = score;
        }
    }

    let mut minutes = vec![0; 60];
    let mut best_minute = 0;
    let v = &map[&best_guard];

    for i in v {
        for (k, c) in i.0.iter().enumerate() {
            if *c == '#' {
                minutes[k] += 1;
            }

            if minutes[k] > minutes[best_minute] {
                best_minute = k
            }
        }
    }

    let gid: usize = best_guard[1..].parse().unwrap();
    println!("{}", gid * best_minute);
}
