package main

import (
	"fmt"
	"strconv"
	"strings"
)

func main() {
	var line string
	var code []int

	fmt.Scanln(&line)

	for _, s := range strings.Split(line, ",") {
		value, _ := strconv.Atoi(s)
		code = append(code, value)
	}

	code[1] = 12
	code[2] = 2

	i := 0

loop:
	for i < len(code) {
		opcode := code[i]

		var in0, in1, out int

		if opcode == 1 || opcode == 2 {
			in0 = code[i+1]
			in1 = code[i+2]
			out = code[i+3]
		}

		switch opcode {
		case 1:
			code[out] = code[in0] + code[in1]
		case 2:
			code[out] = code[in0] * code[in1]
		case 99:
			break loop
		}

		i += 4
	}

	fmt.Println(code[0])
}
