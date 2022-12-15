package main

import (
	"fmt"
	"math"
	"sort"
)

type direction [2]int
type point [2]int

func main() {
	grid, n, m := readGrid()

	best := 0
	x := 0
	y := 0

	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if grid[i][j] == '#' {
				visible := countVisible(grid, i, j, n, m)
				if visible > best {
					best = visible
					x = i
					y = j
				}
			}
		}
	}

	the200th := vaporize(grid, x, y, n, m, 200)
	fmt.Println(the200th[1]*100 + the200th[0])
}

// vaporize returns the position of the k-th asteroid being vaporized
// when the laser is in position [x, y]
func vaporize(grid []string, x, y, n, m, k int) [2]int {
	// list of asteroids found at each direction
	asteroids := make(map[direction][]point)

	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if i == x && j == y {
				continue
			}
			d := dir(x, y, i, j)
			if _, casted := asteroids[d]; casted {
				continue
			}
			asteroids[d] = lasercast(grid, x, y, d, n, m)
		}
	}

	// directions where there are asteroids
	directions := make([]direction, 0)
	for k, v := range asteroids {
		if len(v) > 0 {
			directions = append(directions, k)
		}
	}

	sort.Slice(directions, func(aa, bb int) bool {
		a := directions[aa]
		b := directions[bb]
		angleA := math.Atan2(float64(-a[0]), float64(a[1]))
		angleB := math.Atan2(float64(-b[0]), float64(b[1]))

		angleA -= math.Pi / 2
		angleB -= math.Pi / 2

		if angleA <= 0 {
			angleA += 2 * math.Pi
		}
		if angleB <= 0 {
			angleB += 2 * math.Pi
		}

		return -angleA < -angleB
	})

	vaporized := 0
	for a := 0; a < 10000; a++ {
		rotation := 1 + a/len(directions)
		d := directions[a%len(directions)]
		if len(asteroids[d]) >= rotation {
			vaporized++
			if vaporized == k {
				return asteroids[d][rotation-1]
			}
		}
	}

	panic(fmt.Errorf("Could not find asteroid %dth after 10k iterations", k))
}

// lasercast casts a ray from position [x, y] with direction d
// and returns a list of the asteroids that can be reached at that direction
func lasercast(grid []string, x, y int, d direction, n, m int) []point {
	dx := d[0]
	dy := d[1]
	x += dx
	y += dy
	list := make([]point, 0)
	for x >= 0 && y >= 0 && x < n && y < m {
		if grid[x][y] == '#' {
			list = append(list, [2]int{x, y})
		}
		x += dx
		y += dy
	}
	return list
}

// countVisible returns how many asteroids are visible from position [x, y]
func countVisible(grid []string, x, y, n, m int) int {
	count := 0
	asteroids := make(map[direction]bool)

	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if i == x && j == y {
				continue
			}

			d := dir(x, y, i, j)

			if _, checked := asteroids[d]; checked {
				continue
			}
			ch := raycast(grid, x, y, d, n, m)
			if ch == '#' {
				count++
			}
			asteroids[d] = true
		}
	}
	return count
}

// raycast casts a ray from position [x, y] with direction d
// and returns '#' if it finds an asteroid or '.' if there's only space
func raycast(grid []string, x, y int, d direction, n, m int) rune {
	dx := d[0]
	dy := d[1]
	x += dx
	y += dy
	for x >= 0 && y >= 0 && x < n && y < m {
		if grid[x][y] == '#' {
			return '#'
		}
		x += dx
		y += dy
	}
	return '.'
}

// readGrid reads the grid from stdin and returns it as a slice of strings along the size
// n is the number of rows
// m is the number of columns
func readGrid() (grid []string, n, m int) {
	var line string
	grid = make([]string, 0)

	for {
		_, err := fmt.Scanln(&line)
		if err != nil {
			break
		}
		grid = append(grid, line)
	}

	n = len(grid)
	m = len(grid[0])
	return
}

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

func dir(x, y, i, j int) direction {
	dx := i - x
	dy := j - y
	d := gcd(abs(dx), abs(dy))
	dx = dx / d
	dy = dy / d
	return [2]int{dx, dy}
}
