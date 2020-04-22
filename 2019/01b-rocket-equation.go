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

		i = i/3 - 2
		for i > 0 {
			s += i
			i = i/3 - 2
		}
	}

	fmt.Println(s)
}
