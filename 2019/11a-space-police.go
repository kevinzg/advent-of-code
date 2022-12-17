package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

type intcode struct {
	// memory
	mem []int
	// instruction pointer
	ip int
	// relative base
	rb int
	// input channel
	in <-chan int
	// output channel
	out chan<- int
	// exit channel
	exit chan<- int
}

type opcode int

const (
	SUM opcode = iota + 1
	MUL
	READ
	WRITE
	IFNOTZERO
	IFZERO
	LESS
	EQUALS
	RBADD
	TERMINATE = 99
)

// param mode
type paramMode int

const (
	POSITION paramMode = iota
	IMMEDIATE
	RELATIVE
)

type instructionSpec struct {
	// How many params this instruction takes
	n int
	// Runs the instruction. Return value: 0 means continue, >0 means terminate
	exec func(c *intcode, args []memref) int
}

type memref struct {
	mode  paramMode
	value int
}

func (c *intcode) run() {
	for {
		op, args := c.readInstruction()
		terminate := instructions[op].exec(c, args)
		if terminate > 0 {
			break
		}
	}
}

var instructions = map[opcode]instructionSpec{
	SUM: {
		n: 3,
		exec: func(c *intcode, args []memref) int {
			value := c.memrefread(args[0]) + c.memrefread(args[1])
			c.memrefwrite(args[2], value)
			c.ip += 4
			return 0
		},
	},
	MUL: {
		n: 3,
		exec: func(c *intcode, args []memref) int {
			value := c.memrefread(args[0]) * c.memrefread(args[1])
			c.memrefwrite(args[2], value)
			c.ip += 4
			return 0
		},
	},
	READ: {
		n: 1,
		exec: func(c *intcode, args []memref) int {
			c.memrefwrite(args[0], c.read())
			c.ip += 2
			return 0
		},
	},
	WRITE: {
		n: 1,
		exec: func(c *intcode, args []memref) int {
			c.write(c.memrefread(args[0]))
			c.ip += 2
			return 0
		},
	},
	IFNOTZERO: {
		n: 2,
		exec: func(c *intcode, args []memref) int {
			a := c.memrefread(args[0])
			b := c.memrefread(args[1])
			if a != 0 {
				c.ip = b
			} else {
				c.ip += 3
			}
			return 0
		},
	},
	IFZERO: {
		n: 2,
		exec: func(c *intcode, args []memref) int {
			a := c.memrefread(args[0])
			b := c.memrefread(args[1])
			if a == 0 {
				c.ip = b
			} else {
				c.ip += 3
			}
			return 0
		},
	},
	LESS: {
		n: 3,
		exec: func(c *intcode, args []memref) int {
			a := c.memrefread(args[0])
			b := c.memrefread(args[1])
			var value int
			if a < b {
				value = 1
			} else {
				value = 0
			}
			c.memrefwrite(args[2], value)
			c.ip += 4
			return 0
		},
	},
	EQUALS: {
		n: 3,
		exec: func(c *intcode, args []memref) int {
			a := c.memrefread(args[0])
			b := c.memrefread(args[1])
			var value int
			if a == b {
				value = 1
			} else {
				value = 0
			}
			c.memrefwrite(args[2], value)
			c.ip += 4
			return 0
		},
	},
	RBADD: {
		n: 1,
		exec: func(c *intcode, args []memref) int {
			a := c.memrefread(args[0])
			c.rb += a
			c.ip += 2
			return 0
		},
	},
	TERMINATE: {
		n: 0,
		exec: func(c *intcode, args []memref) int {
			close(c.out)
			c.exit <- 1
			return 1
		},
	},
}

// reads the next instruction at the current ip
// returns its opcode and the parsed params for the operation
func (c *intcode) readInstruction() (opcode, []memref) {
	rawInstruction := c.mem[c.ip]

	op := opcode(rawInstruction % 100)
	modeMask := rawInstruction / 100

	inst := instructions[op]
	args := make([]memref, inst.n)

	for i := 0; i < inst.n; i++ {
		args[i].mode = paramMode(modeMask % 10)
		args[i].value = c.mem[c.ip+1+i]

		modeMask /= 10
	}

	return op, args
}

// returns a pointer (i.e. a position in memory)
func (c *intcode) memrefptr(r memref) int {
	if r.mode == POSITION {
		return r.value
	} else if r.mode == IMMEDIATE {
		log.Panic("Trying to reference memory using immediate mode")
	} else if r.mode == RELATIVE {
		return c.rb + r.value
	}
	log.Panic("Bad mode received ", r.mode)
	return -1
}

// returns the value for the memref
func (c *intcode) memrefread(r memref) int {
	if r.mode == IMMEDIATE {
		return r.value
	}
	ptr := c.memrefptr(r)
	return c.memread(ptr)
}

// writes the value for the memref
func (c *intcode) memrefwrite(r memref, value int) {
	if r.mode == IMMEDIATE {
		log.Panic("Trying to write using immediate mode reference", r.mode)
	}
	ptr := c.memrefptr(r)
	c.memwrite(ptr, value)
}

// returns the value for the given pointer
func (c *intcode) memread(ptr int) int {
	if n := len(c.mem); ptr >= n {
		return 0
	}
	return c.mem[ptr]
}

// writes the value to memory at the given pointer
func (c *intcode) memwrite(ptr int, value int) {
	// increases the memory size if the pointer is outside the current memory bounds
	if n := len(c.mem); ptr >= n {
		size := ptr - n + 1
		c.mem = append(c.mem, make([]int, size)...)
	}
	c.mem[ptr] = value
}

// reads from the input and returns the value
func (c *intcode) read() int {
	value, ok := <-c.in
	if !ok {
		log.Fatalf("FATAL: Reading from closed channel: %v", c.in)
	}
	return value
}

// writes to the output
func (c *intcode) write(value int) {
	c.out <- value
}

type cell [2]int

func main() {
	inputChannel := make(chan int)
	outputChannel := make(chan int)
	exitChannel := make(chan int)

	computer := intcode{
		mem:  readProgram(),
		in:   inputChannel,
		out:  outputChannel,
		exit: exitChannel,
	}

	grid := make(map[cell]int)
	dir := 0 // up: 0, right: 1, down: 2, left: 3
	pos := cell{}
	dx := [4]int{0, 1, 0, -1}
	dy := [4]int{1, 0, -1, 0}
	dr := [2]int{-1, 1}

	go func() {
		for {
			inputChannel <- grid[pos]
			color, ok := <-outputChannel
			if !ok {
				break
			}
			rotation, ok := <-outputChannel
			if !ok {
				break
			}
			grid[pos] = color
			dir = (dir + dr[rotation] + 4) % 4
			pos[0] += dx[dir]
			pos[1] += dy[dir]
		}
	}()

	go computer.run()

	<-exitChannel

	fmt.Println(len(grid))
}

func readProgram() []int {
	var line string
	var code []int

	fmt.Scanln(&line)

	for _, s := range strings.Split(line, ",") {
		value, _ := strconv.Atoi(s)
		code = append(code, value)
	}
	return code
}
