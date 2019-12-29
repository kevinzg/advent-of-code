package main

import (
	"fmt"
	"strconv"
	"strings"
)

type cell struct {
	x, y int
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func cells(grid map[cell]int, path string, mask int) {
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

	x, y := 0, 0

	for _, code := range codes {
		direction := code[0]
		length, _ := strconv.Atoi(code[1:])

		for i := 0; i < length; i++ {
			x += dx[direction]
			y += dy[direction]
			grid[cell{x, y}] |= mask
		}
	}
}

func main() {
	var path0, path1 string

	grid := make(map[cell]int)

	fmt.Scanln(&path0)
	fmt.Scanln(&path1)

	cells(grid, path0, 0b01)
	cells(grid, path1, 0b10)

	min := 1 << 30

	for cell, mask := range grid {
		distance := abs(cell.x) + abs(cell.y)
		if mask == 0b11 && distance < min {
			min = distance
		}
	}

	fmt.Println(min)
}
