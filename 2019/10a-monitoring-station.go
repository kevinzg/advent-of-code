package main

import (
	"fmt"
)

func main() {
	var line string
	grid := make([]string, 0)

	for {
		_, err := fmt.Scanln(&line)

		if err != nil {
			break
		}

		grid = append(grid, line)
	}

	n := len(grid)
	m := len(grid[0])
	best := 0

	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if grid[i][j] == '#' {
				visible := scan(grid, i, j, n, m)
				if visible > best {
					best = visible
				}
			}
		}
	}

	fmt.Println(best)
}

func scan(grid []string, x, y, n, m int) int {
	count := 0
	asteroid := make(map[string]bool)

	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if i == x && j == y {
				continue
			}

			dx := i - x
			dy := j - y
			d := gcd(abs(dx), abs(dy))
			dx = dx / d
			dy = dy / d
			tan := fmt.Sprintf("%d/%d", dx, dy)

			if _, checked := asteroid[tan]; !checked {
				ch := raycast(grid, x, y, dx, dy, n, m)
				if ch == '#' {
					count++
				}
				asteroid[tan] = true
			}
		}
	}
	return count
}

func raycast(grid []string, x, y, dx, dy, n, m int) rune {
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
