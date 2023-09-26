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
    let mut offset = 2;
    let mut values = Vec::with_capacity(n);
    values.resize(n, 0);
    for i in 0..n {
        let (l, z) = sum_metadata(&data[offset..]);
        offset += l;
        values[i] = z;
    }
    let s = if n == 0 {
        data[offset..offset + m].iter().sum()
    } else {
        data[offset..offset + m]
            .iter()
            .map(|c| {
                if *c <= 0 || *c > n as i32 {
                    0
                } else {
                    values[*c as usize - 1]
                }
            })
            .sum()
    };
    return (offset + m, s);
}
