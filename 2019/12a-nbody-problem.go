package main

import (
	"fmt"
	"log"
)

type vec3 [3]int

func main() {
	pos := readInput()
	var vel [4]vec3

	steps := 1000

	for i := 0; i < steps; i++ {
		g := gravity(pos)
		vel = add(vel, g)
		pos = add(pos, vel)
	}

	fmt.Println(energy(pos, vel))
}

func energy(pos, vel [4]vec3) int {
	var s int
	for i := 0; i < 4; i++ {
		s += abssum(pos[i]) * abssum(vel[i])
	}
	return s
}

func gravity(pos [4]vec3) (g [4]vec3) {
	for i := 0; i < 4; i++ {
		for j := i + 1; j < 4; j++ {
			for k := 0; k < 3; k++ {
				d := 0
				if pos[i][k] < pos[j][k] {
					d = 1
				} else if pos[i][k] > pos[j][k] {
					d = -1
				}
				g[i][k] += d
				g[j][k] -= d
			}
		}
	}
	return
}

func add(a, b [4]vec3) (c [4]vec3) {
	for i := 0; i < 4; i++ {
		for k := 0; k < 3; k++ {
			c[i][k] = a[i][k] + b[i][k]
		}
	}
	return
}

func abssum(a vec3) int {
	return abs(a[0]) + abs(a[1]) + abs(a[2])
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

func readInput() [4]vec3 {
	var m [4]vec3
	for i := range m {
		_, err := fmt.Scanf("<x=%d, y=%d, z=%d>\n", &m[i][0], &m[i][1], &m[i][2])
		if err != nil {
			log.Panic(err)
		}
	}
	return m
}
