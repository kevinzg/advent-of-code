package main

import (
	"fmt"
	"strconv"
	"strings"
)

func run(noun int, verb int, memory []int) int {
	code := make([]int, len(memory))
	copy(code, memory)

	code[1] = noun
	code[2] = verb

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

	return code[0]
}

func main() {
	var line string
	var code []int

	fmt.Scanln(&line)

	for _, s := range strings.Split(line, ",") {
		value, _ := strconv.Atoi(s)
		code = append(code, value)
	}

	var noun, verb int
	var desired = 19690720

	loop:
	for noun = 0; noun < 100; noun++ {
		for verb = 0; verb < 100; verb++ {
			output := run(noun, verb, code)
			if output == desired {
				break loop
			}
		}
	}

	fmt.Println(noun * 100 + verb)

}
