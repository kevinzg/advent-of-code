package main

import (
	"fmt"
	"log"
)

type vec3 [3]int

func main() {
	pos := readInput()
	var cl [3]int

	for i := 0; i < 3; i++ {
		cl[i] = findCycleLength([4]int{
			pos[0][i],
			pos[1][i],
			pos[2][i],
			pos[3][i],
		}, [4]int{})
	}

	fmt.Println(lcm(lcm(cl[0], cl[1]), cl[2]))
}

func findCycleLength(pos, vel [4]int) int {
	targetPos := pos
	targetVel := vel
	c := 0
	for {
		gra := gravity(pos)
		vel = add(vel, gra)
		pos = add(pos, vel)
		c++
		if pos == targetPos && vel == targetVel {
			return c
		}
	}
}

func gravity(pos [4]int) (g [4]int) {
	for i := 0; i < 4; i++ {
		for j := i + 1; j < 4; j++ {
			d := 0
			if pos[i] < pos[j] {
				d = 1
			} else if pos[i] > pos[j] {
				d = -1
			}
			g[i] += d
			g[j] -= d
		}
	}
	return
}

func add(a, b [4]int) (c [4]int) {
	for k := 0; k < 4; k++ {
		c[k] = a[k] + b[k]
	}
	return
}

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func lcm(a, b int) int {
	g := gcd(a, b)
	return a * b / g
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
