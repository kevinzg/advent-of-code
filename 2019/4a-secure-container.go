package main

import (
	"fmt"
	"strconv"
)

func check(a int) bool {
	number := []rune(strconv.Itoa(a))

	var last rune = '0'
	digits := make(map[rune]int)

	for _, c := range number {
		if c < last {
			return false
		}

		digits[c] += 1
		last = c
	}

	for _, count := range digits {
		if count >= 2 {
			return true
		}
	}

	return false
}

func main() {
	var n, m int

	fmt.Scanf("%d-%d", &n, &m)

	s := 0
	for i := n; i <= m; i++ {
		if check(i) {
			s++
		}
	}

	fmt.Println(s)
}
