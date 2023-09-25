use std::cmp;
use std::collections::HashMap;
use std::io;
use std::num::ParseIntError;
use std::str::FromStr;

const INF: i32 = 1e30 as i32;

fn main() {
    let mut coordinates: Vec<Point> = vec![];
    let mut min_x = INF;
    let mut min_y = INF;
    let mut max_x = -INF;
    let mut max_y = -INF;

    loop {
        let mut line = String::new();

        match io::stdin().read_line(&mut line) {
            Ok(0) => break,
            Ok(_) => (),
            Err(err) => panic!("{}", err),
        }

        let coord = Point::from_str(&line).expect("Error parsing point");
        coordinates.push(coord);

        min_x = cmp::min(min_x, coord.0);
        max_x = cmp::max(max_x, coord.0);
        min_y = cmp::min(min_y, coord.1);
        max_y = cmp::max(max_y, coord.1);
    }

    let mut counter: HashMap<Point, i32> = HashMap::new();

    for x in min_x..=max_x {
        for y in min_y..=max_y {
            let closest: Closest = coordinates
                .iter()
                .map(|p| (p, (x - p.0).abs() + (y - p.1).abs()))
                .fold(Closest::None, |acc, (p, d)| match acc {
                    Closest::None => Closest::Single(d, *p),
                    Closest::Single(bd, _) => {
                        if d < bd {
                            Closest::Single(d, *p)
                        } else if d == bd {
                            Closest::Multiple(d)
                        } else {
                            acc
                        }
                    }
                    Closest::Multiple(bd) => {
                        if d < bd {
                            Closest::Single(d, *p)
                        } else {
                            acc
                        }
                    }
                });

            match closest {
                Closest::None => panic!("There should be a closest coordinate"),
                Closest::Single(_, bp) => {
                    if x == min_x || x == max_x || y == min_y || y == max_y {
                        counter.entry(bp).and_modify(|c| *c = -1).or_insert(-1);
                    } else {
                        counter
                            .entry(bp)
                            .and_modify(|c| *c = if *c >= 0 { *c + 1 } else { -1 })
                            .or_insert(1);
                    }
                }
                _ => {}
            }
        }
    }

    let max = counter.into_values().fold(0, |acc, s| cmp::max(acc, s));

    println!("{}", max);
}

enum Closest {
    Single(i32, Point),
    Multiple(i32),
    None,
}

#[derive(Copy, Clone, Debug, Hash, Eq, PartialEq)]
struct Point(i32, i32);

impl FromStr for Point {
    type Err = ParseIntError;

    // 1, 1
    fn from_str(s: &str) -> Result<Self, Self::Err> {
        let nums: Vec<_> = s.split(", ").collect();

        return Ok(Point(nums[0].trim().parse()?, nums[1].trim().parse()?));
    }
}
