package main

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

type intcode struct {
	memory             []int
	input              []int
	instructionPointer int
	inputPointer       int
	modeMask           int
	outputCallback     func(int)
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
	5: func(c *intcode) int {
		operands := c.getOperands(2, c.instructionPointer+1)

		if operands[0] != 0 {
			c.instructionPointer = operands[1]
		} else {
			c.instructionPointer += 3
		}

		return 0
	},
	6: func(c *intcode) int {
		operands := c.getOperands(2, c.instructionPointer+1)

		if operands[0] == 0 {
			c.instructionPointer = operands[1]
		} else {
			c.instructionPointer += 3
		}

		return 0
	},
	7: func(c *intcode) int {
		operands := c.getOperands(2, c.instructionPointer+1)
		dest := c.memory[c.instructionPointer+3]

		if operands[0] < operands[1] {
			c.memory[dest] = 1
		} else {
			c.memory[dest] = 0
		}

		c.instructionPointer += 4

		return 0
	},
	8: func(c *intcode) int {
		operands := c.getOperands(2, c.instructionPointer+1)
		dest := c.memory[c.instructionPointer+3]

		if operands[0] == operands[1] {
			c.memory[dest] = 1
		} else {
			c.memory[dest] = 0
		}

		c.instructionPointer += 4

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
	c.outputCallback(value)
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

const n int = 5

func main() {
	var line string
	var code []int

	fmt.Scanln(&line)

	for _, s := range strings.Split(line, ",") {
		value, _ := strconv.Atoi(s)
		code = append(code, value)
	}

	phaseSettings := [n]int{0, 1, 2, 3, 4}
	idx := [n]int{-1}

	max := math.MinInt32

	callback := func(arr [n]int) {
		previousOutput := 0

		for _, phase := range arr {
			input := []int{phase, previousOutput}

			localMemory := make([]int, len(code))
			copy(localMemory, code)

			computer := intcode{
				memory: code,
				input:  input,
				outputCallback: func(output int) {
					previousOutput = output
				},
			}

			computer.run()
		}

		if previousOutput > max {
			max = previousOutput
		}
	}

	permute(phaseSettings, 0, idx, 0, callback)

	fmt.Println(max)
}

func permute(arr [n]int, mask int, idx [n]int, m int, callback func([n]int)) {
	if m == n {
		permutation := [n]int{}

		for i := 0; i < n; i++ {
			permutation[i] = arr[idx[i]]
		}

		callback(permutation)
		return
	}

	for k := 0; k < n; k++ {
		bit := (1 << k)
		if ^mask&bit != 0 {
			idx[k] = arr[m]
			permute(arr, mask|bit, idx, m+1, callback)
			idx[k] = -1
		}
	}
}
