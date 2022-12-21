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
		if c.ip < 0 {
			break
		}
		op, args := c.readInstruction()
		terminate := instructions[op].exec(c, args)
		if terminate > 0 {
			break
		}
	}
}

func (c *intcode) forceStop() {
	c.ip = -1000
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

func run() map[cell]int {
	dx := [5]int{0, 0, 0, -1, 1}
	dy := [5]int{0, 1, -1, 0, 0}

	move := func(c cell, d int) cell {
		return cell{c[0] + dx[d], c[1] + dy[d]}
	}
	invert := func(d int) int {
		if d%2 == 1 {
			return d + 1 // 1 and 3
		}
		return d - 1 // 2 and 4
	}

	var computer intcode
	var droid cell
	var command int
	grid := make(map[cell]int)
	tb := make([]int, 0, 100)
	isTraceback := false

	sendInput := func() int {
		for k := 1; k <= 4; k++ {
			if grid[move(droid, k)] == 0 {
				command = k
				return k
			}
		}
		if len(tb) == 0 {
			computer.forceStop()
			return 1
		}
		prev := tb[len(tb)-1]
		tb = tb[:len(tb)-1]
		command = invert(prev)
		isTraceback = true
		return command
	}

	processOutput := func(x int) {
		if x == 0 {
			grid[move(droid, command)] = '#'
		} else if x == 1 {
			droid = move(droid, command)
			if !isTraceback {
				tb = append(tb, command)
			}
			grid[droid] = '.'
		} else if x == 2 {
			droid = move(droid, command)
			if !isTraceback {
				tb = append(tb, command)
			}
			grid[droid] = 'O'
		}
		isTraceback = false
	}

	computer = intcode{
		mem: readProgram(),
		in:  sendInput,
		out: processOutput,
	}

	computer.run()

	return grid
}

func solve(grid map[cell]int, start cell) int {
	dx := [5]int{0, 0, 0, -1, 1}
	dy := [5]int{0, 1, -1, 0, 0}
	move := func(c cell, d int) cell {
		return cell{c[0] + dx[d], c[1] + dy[d]}
	}

	queue := make([]cell, 0, 100)
	qpush := func(c cell) {
		queue = append(queue, c)
	}
	qpop := func() cell {
		c := queue[len(queue)-1]
		queue = queue[:len(queue)-1]
		return c
	}
	qempty := func() bool {
		return len(queue) == 0
	}

	dist := make(map[cell]int)

	qpush(start)
	dist[start] = 0
	maxDistance := 0

	for !qempty() {
		c := qpop()
		d := dist[c]

		for k := 1; k <= 4; k++ {
			nc := move(c, k)
			if _, ok := dist[nc]; ok {
				continue
			}
			if grid[nc] == 'O' {
				log.Panic("Somehow I returned to the start")
			} else if grid[nc] == '#' {
				dist[nc] = -1
			} else if grid[nc] == '.' {
				dist[nc] = d + 1
				if d+1 > maxDistance {
					maxDistance = d + 1
				}
				qpush(nc)
			} else {
				log.Panic("Reached unreachable cell")
			}
		}
	}

	return maxDistance
}

func main() {
	grid := run()
	var o2 cell
	for k, v := range grid {
		if v == 'O' {
			o2 = k
			break
		}
	}
	fmt.Println(solve(grid, o2))
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

func printGrid(grid map[cell]int, d cell) {
	minFn := func(a, b int) int {
		if a < b {
			return a
		}
		return b
	}
	maxFn := func(a, b int) int {
		if a > b {
			return a
		}
		return b
	}
	var minX, maxX, minY, maxY int
	for k := range grid {
		minX = minFn(k[0], minX)
		maxX = maxFn(k[0], maxX)
		minY = minFn(k[1], minY)
		maxY = maxFn(k[1], maxY)
	}

	var str strings.Builder
	for y := maxY; y >= minY; y-- {
		for x := minX; x <= maxX; x++ {
			color := rune(grid[cell{x, y}])
			if color == 0 {
				color = ' '
			}
			if d[0] == x && d[1] == y {
				color = 'D'
			}
			str.WriteRune(color)
		}
		str.WriteRune('\n')
	}

	fmt.Print(str.String())
}
