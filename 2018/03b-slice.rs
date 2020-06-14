use std::collections::HashMap;
use std::io;
use std::num::ParseIntError;
use std::str::FromStr;

#[derive(Debug, Clone)]
struct Cell {
    row: i32,
    col: i32,
}

#[derive(Debug, Clone)]
struct SubGrid {
    id: i32,
    start: Cell,
    width: i32,
    height: i32,
}

impl Cell {
    fn as_tuple(self) -> (i32, i32) {
        (self.row, self.col)
    }
}

impl SubGrid {
    fn size(&self) -> i32 {
        return self.width * self.height;
    }
}

impl IntoIterator for SubGrid {
    type Item = Cell;
    type IntoIter = SubGridIterator;

    fn into_iter(self) -> Self::IntoIter {
        SubGridIterator {
            sub_grid: self,
            index: 0,
        }
    }
}

#[derive(Debug)]
struct SubGridIterator {
    sub_grid: SubGrid,
    index: i32,
}

impl Iterator for SubGridIterator {
    type Item = Cell;

    fn next(&mut self) -> Option<Self::Item> {
        if self.index >= self.sub_grid.size() {
            return None;
        }

        let drow = self.index / self.sub_grid.width;
        let dcol = self.index % self.sub_grid.width;

        let row = self.sub_grid.start.row + drow;
        let col = self.sub_grid.start.col + dcol;

        self.index += 1;

        Some(Cell { row, col })
    }
}

impl FromStr for SubGrid {
    type Err = ParseIntError;

    // #123 @ 3,2: 5x4
    fn from_str(s: &str) -> Result<Self, Self::Err> {
        let num = s.find('#').unwrap();
        let at = s.find('@').unwrap();
        let comma = s.find(',').unwrap();
        let colon = s.find(':').unwrap();
        let ex = s.find('x').unwrap();

        Ok(SubGrid {
            id: s[num + 1..at].trim().parse()?,
            start: Cell {
                col: s[at + 1..comma].trim().parse()?,
                row: s[comma + 1..colon].trim().parse()?,
            },
            width: s[colon + 1..ex].trim().parse()?,
            height: s[ex + 1..].trim().parse()?,
        })
    }
}

fn main() {
    let mut map: HashMap<(i32, i32), i32> = HashMap::new();
    let mut claims: Vec<SubGrid> = Vec::new();

    loop {
        let mut rect = String::new();

        match io::stdin().read_line(&mut rect) {
            Ok(0) => break,
            Ok(_) => (),
            Err(err) => panic!(err),
        }

        let sub_grid = SubGrid::from_str(&rect).expect("Error parsing rect");
        claims.push(sub_grid.clone());

        for c in sub_grid.into_iter() {
            let value = map.entry(c.as_tuple()).or_insert(0);
            *value += 1;
        }
    }

    for g in claims {
        let mut ok = true;
        let id = g.id;

        for c in g.into_iter() {
            if map[&c.as_tuple()] != 1 {
                ok = false;
                break;
            }
        }

        if ok {
            println!("{}", id);
        }
    }
}
