package main

import (
	"fmt"
	"strings"
)

func count(bodies map[string][]string, current string, acc int) int {
	sats := bodies[current]

	sum := acc

	for _, name := range sats {
		sum += count(bodies, name, acc+1)
	}

	return sum
}

func main() {
	var line string
	var bodies = make(map[string][]string)

	for {
		_, err := fmt.Scanln(&line)

		if err != nil {
			break
		}

		values := strings.Split(line, ")")

		left := values[0]
		right := values[1]

		bodies[left] = append(bodies[left], right)
	}

	ans := count(bodies, "COM", 0)
	fmt.Println(ans)
}
