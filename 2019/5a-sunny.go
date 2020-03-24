package main

import (
	"fmt"
	"strconv"
	"strings"
)

type intcode struct {
	memory             []int
	input              []int
	instructionPointer int
	inputPointer       int
	modeMask           int
}

func (c *intcode) getOperands(n int, pos int) []int {
	operands := make([]int, n)
	modeMask := c.modeMask

	for i, p := range c.memory[pos : pos+n] {
		mode := modeMask % 10
		modeMask = modeMask / 10

		if mode == 0 {
			operands[i] = c.memory[p]
		} else {
			operands[i] = p
		}
	}

	return operands
}

var operations = map[int]func(c *intcode) int{
	1: func(c *intcode) int {
		operands := c.getOperands(2, c.instructionPointer+1)
		dest := c.memory[c.instructionPointer+3]

		c.memory[dest] = operands[0] + operands[1]

		c.instructionPointer += 4

		return 0
	},
	2: func(c *intcode) int {
		operands := c.getOperands(2, c.instructionPointer+1)
		dest := c.memory[c.instructionPointer+3]

		c.memory[dest] = operands[0] * operands[1]

		c.instructionPointer += 4

		return 0
	},
	3: func(c *intcode) int {
		dest := c.memory[c.instructionPointer+1]
		c.read(dest)

		c.instructionPointer += 2

		return 0
	},
	4: func(c *intcode) int {
		operands := c.getOperands(1, c.instructionPointer+1)
		c.write(operands[0])

		c.instructionPointer += 2

		return 0
	},
	99: func(c *intcode) int {
		return 1
	},
}

func (c *intcode) read(dest int) {
	value := c.input[c.inputPointer]
	c.inputPointer++
	c.memory[dest] = value
}

func (c *intcode) write(value int) {
	fmt.Println(value)
}

func (c *intcode) run() int {
	for {
		instruction := c.memory[c.instructionPointer]

		opcode := instruction % 100
		c.modeMask = instruction / 100

		terminate := operations[opcode](c)

		if terminate != 0 {
			break
		}
	}

	return 0
}

func main() {
	var line string
	var code []int

	fmt.Scanln(&line)

	for _, s := range strings.Split(line, ",") {
		value, _ := strconv.Atoi(s)
		code = append(code, value)
	}

	input := []int{1}

	computer := intcode{memory: code, input: input}
	computer.run()
}
