package main

import (
	"fmt"
	"strconv"
	"strings"
)

type cell struct {
	x, y int
}

type value struct {
	mask int
	d    [2]int
}

func cells(grid map[cell]value, path string, p int) {
	dx := map[byte]int{
		'U': 0,
		'D': 0,
		'R': 1,
		'L': -1,
	}

	dy := map[byte]int{
		'U': 1,
		'D': -1,
		'R': 0,
		'L': 0,
	}

	codes := strings.Split(path, ",")

	d := 0
	x, y := 0, 0
	mask := 1 << p

	for _, code := range codes {
		direction := code[0]
		length, _ := strconv.Atoi(code[1:])

		for i := 0; i < length; i++ {
			x += dx[direction]
			y += dy[direction]
			d++

			val := grid[cell{x, y}]
			val.mask |= mask
			if val.d[p] == 0 {
				val.d[p] = d
			}

			grid[cell{x, y}] = val
		}
	}
}

func main() {
	var path0, path1 string

	grid := make(map[cell]value)

	fmt.Scanln(&path0)
	fmt.Scanln(&path1)

	cells(grid, path0, 0)
	cells(grid, path1, 1)

	min := 1 << 30

	for _, val := range grid {
		distance := val.d[0] + val.d[1]
		if val.mask == 0b11 && distance < min {
			min = distance
		}
	}

	fmt.Println(min)
}
