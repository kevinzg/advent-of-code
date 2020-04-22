package main

import (
	"fmt"
)

func main() {
	var s, i int

	for {
		_, err := fmt.Scanf("%d", &i)

		if err != nil {
			break
		}

		s += i/3 - 2
	}

	fmt.Println(s)
}
