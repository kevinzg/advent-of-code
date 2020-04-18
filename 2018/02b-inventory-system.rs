use std::collections::HashMap;
use std::io;
use std::iter::FromIterator;

fn main() {
    let mut ids: Vec<String> = Vec::new();

    loop {
        let mut id = String::new();

        match io::stdin().read_line(&mut id) {
            Ok(0) => break,
            Ok(_) => (),
            Err(err) => panic!(err),
        }

        ids.push(id.trim().to_string());
    }

    // BK-Tree, not the best solution tho
    let mut tree = Node::new(&ids[0]);

    for i in 1..ids.len() {
        let id = &ids[i];

        match tree.search(&id) {
            Some(node) => {
                print_solution(&node.id, &id);
                break;
            }
            None => tree.insert(&id),
        }
    }
}

fn print_solution(a: &String, b: &String) {
    let solution = a
        .chars()
        .zip(b.chars())
        .filter(|(x, y)| x == y)
        .map(|(x, _)| x);

    println!("{}", String::from_iter(solution));
}

fn compare(a: &String, b: &String) -> usize {
    if a.len() != b.len() {
        return a.len() + b.len();
    }

    a.chars()
        .zip(b.chars())
        .map(|(x, y)| if x != y { 1 } else { 0 })
        .fold(0, |acc, k| acc + k)
}

#[derive(Debug)]
struct Node {
    id: String,
    nodes: HashMap<usize, Node>,
}

impl Node {
    fn new(id: &String) -> Self {
        Node {
            id: id.to_string(),
            nodes: HashMap::new(),
        }
    }

    fn insert(&mut self, other_id: &String) {
        let d = compare(&self.id, &other_id);

        match self.nodes.get_mut(&d) {
            Some(node) => node.insert(other_id),
            None => {
                self.nodes.insert(d, Node::new(other_id));
                ()
            }
        }
    }

    fn search(&self, other_id: &String) -> Option<&Self> {
        let d = compare(&self.id, &other_id);

        if d == 0 {
            panic!("Duplicate string found")
        } else if d == 1 {
            return Some(&self);
        }

        for i in (d - 1)..=(d + 1) {
            if let Some(node) = self.nodes.get(&i) {
                match node.search(other_id) {
                    Some(node) => return Some(node),
                    None => continue,
                };
            }
        }

        return None;
    }
}
