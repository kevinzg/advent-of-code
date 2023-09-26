use std::io;

fn main() {
    let data = {
        let mut line = String::new();
        io::stdin().read_line(&mut line).unwrap();
        line.split(" ")
            .map(|s| s.trim().parse().unwrap())
            .collect::<Vec<i32>>()
    };
    println!("{}", sum_metadata(&data).1);
}

fn sum_metadata(data: &[i32]) -> (usize, i32) {
    let n = data[0] as usize;
    let m = data[1] as usize;
    let mut s = 0;
    let mut offset = 2;
    for _ in 0..n {
        let (l, z) = sum_metadata(&data[offset..]);
        offset += l;
        s += z;
    }
    for k in 0..m {
        s += data[offset + k as usize];
    }
    return (offset + m, s);
}
