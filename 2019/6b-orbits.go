package main

import (
	"fmt"
	"strings"
)

func path(parent map[string]string, target string) []string {
	targetPath := make([]string, 0)
	end := "COM"

	for target != end {
		targetPath = append(targetPath, target)
		target = parent[target]
	}

	targetPath = append(targetPath, end)

	return targetPath
}

func main() {
	var line string
	var parent = make(map[string]string)

	for {
		_, err := fmt.Scanln(&line)

		if err != nil {
			break
		}

		values := strings.Split(line, ")")

		left := values[0]
		right := values[1]

		parent[right] = left
	}

	me := "YOU"
	santa := "SAN"

	myPath := path(parent, me)
	santaPath := path(parent, santa)

	i := 1

	for {
		if myPath[len(myPath)-i] != santaPath[len(santaPath)-i] {
			ans := len(myPath) + len(santaPath) - 2*i
			fmt.Println(ans)
			break
		}
		i++
	}
}
