package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

type intcode struct {
	// memory
	mem []int
	// instruction pointer
	ip int
	// relative base
	rb int
	// input callback
	in func() int
	// output callback
	out func(int)
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
			c.memrefwrite(args[0], c.in())
			c.ip += 2
			return 0
		},
	},
	WRITE: {
		n: 1,
		exec: func(c *intcode, args []memref) int {
			c.out(c.memrefread(args[0]))
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

type cell [2]int

type tileid int

const (
	EMPTY tileid = iota
	WALL
	BLOCK
	PADDLE
	BALL
)

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

func run() int {
	var ball cell
	var paddle cell

	grid := make(map[cell]tileid)
	score := 0

	tidToChar := [5]rune{' ', '|', '◼', '-', 'o'}

	sendInput := func() int {
		time.Sleep(1 * time.Millisecond)
		if ball[0] < paddle[0] {
			return -1
		} else if ball[0] > paddle[0] {
			return 1
		}
		return 0
	}

	processOutput := func(x, y, z int) {
		if x == -1 && y == 0 {
			score = z
			fmt.Print("\033[1;1H") // Move cursor
			fmt.Printf("Score: %d", score)
		} else {
			tid := tileid(z)
			grid[cell{x, y}] = tid
			fmt.Printf("\033[%d;%dH", y+2, x+1)
			fmt.Printf(string(tidToChar[z]))

			if tid == PADDLE {
				paddle = cell{x, y}
			} else if tid == BALL {
				ball = cell{x, y}
			}
		}
	}

	groupOutput := func() func(int) {
		var t [3]int
		k := 0
		return func(c int) {
			t[k] = c
			k = (k + 1) % 3
			if k == 0 {
				processOutput(t[0], t[1], t[2])
			}
		}
	}()

	computer := intcode{
		mem: readProgram(),
		in:  sendInput,
		out: groupOutput,
	}

	computer.mem[0] = 2

	computer.run()

	// To see the score
	time.Sleep(1 * time.Second)

	return score
}

func main() {
	var score int

	// Only for Linux
	fmt.Print("\0337")     // save cursor position
	fmt.Print("\033[?47h") // switch to alternate screen
	fmt.Print("\033[?25l") // hide cursor

	defer func() {
		fmt.Print("\033[2J")   // clear screen
		fmt.Print("\033[?47l") // switch back to normal screen
		fmt.Print("\033[?25h") // show the cursor
		fmt.Print("\0338")     // restore cursor position
		fmt.Println(score)
	}()

	score = run()
}
