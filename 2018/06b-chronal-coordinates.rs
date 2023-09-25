use std::cmp;
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

    let mut counter = 0;

    for x in min_x..=max_x {
        for y in min_y..=max_y {
            let sum: i32 = coordinates
                .iter()
                .map(|p| (x - p.0).abs() + (y - p.1).abs())
                .fold(0, |acc, d| acc + d);

            if sum < 10000 {
                counter += 1
            }
        }
    }

    println!("{}", counter);
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
