package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

type intcode struct {
	memory             []int
	instructionPointer int
	relativeBase       int
	inputChannel       <-chan int
	outputChannel      chan<- int
	exitChannel        chan<- int
}

type instruction struct {
	opcode int
	params int
	args   []memref
}

type memref struct {
	mode  int
	value int
}

func newInstruction(memory []int, ip int) instruction {
	rawInstruction := memory[ip]

	opcode := rawInstruction % 100
	modeMask := rawInstruction / 100

	inst := instructions[opcode]
	inst.args = make([]memref, inst.params)

	for i := 0; i < inst.params; i++ {
		inst.args[i].mode = modeMask % 10
		inst.args[i].value = memory[ip+1+i]

		modeMask /= 10
	}

	return inst
}

func (inst *instruction) execute(c *intcode) int {
	return 0
}

// opcodes
const (
	SUM = iota + 1
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

var instructions = map[int]instruction{
	SUM: instruction{
		opcode: SUM,
		params: 3,
	},
	MUL: instruction{
		opcode: MUL,
		params: 3,
	},
	READ: instruction{
		opcode: READ,
		params: 1,
	},
	WRITE: instruction{
		opcode: WRITE,
		params: 1,
	},
	IFNOTZERO: instruction{
		opcode: IFNOTZERO,
		params: 2,
	},
	IFZERO: instruction{
		opcode: IFZERO,
		params: 2,
	},
	LESS: instruction{
		opcode: LESS,
		params: 3,
	},
	EQUALS: instruction{
		opcode: EQUALS,
		params: 3,
	},
	RBADD: instruction{
		opcode: RBADD,
		params: 1,
	},
	TERMINATE: instruction{
		opcode: TERMINATE,
		params: 0,
	},
}

var executableCode = map[int]func(c *intcode, inst *instruction) int{
	SUM: func(c *intcode, inst *instruction) int {
		c.memrefwrite(
			inst.args[2],
			c.memrefread(inst.args[0])+c.memrefread(inst.args[1]),
		)
		c.instructionPointer += 4
		return 0
	},
	MUL: func(c *intcode, inst *instruction) int {
		c.memrefwrite(
			inst.args[2],
			c.memrefread(inst.args[0])*c.memrefread(inst.args[1]),
		)
		c.instructionPointer += 4
		return 0
	},
	READ: func(c *intcode, inst *instruction) int {
		c.read(
			c.memrefptr(inst.args[0]),
		)
		c.instructionPointer += 2
		return 0
	},
	WRITE: func(c *intcode, inst *instruction) int {
		c.write(
			c.memrefread(inst.args[0]),
		)
		c.instructionPointer += 2
		return 0
	},
	IFNOTZERO: func(c *intcode, inst *instruction) int {
		a := c.memrefread(inst.args[0])
		b := c.memrefread(inst.args[1])

		if a != 0 {
			c.instructionPointer = b
		} else {
			c.instructionPointer += 3
		}
		return 0
	},
	IFZERO: func(c *intcode, inst *instruction) int {
		a := c.memrefread(inst.args[0])
		b := c.memrefread(inst.args[1])

		if a == 0 {
			c.instructionPointer = b
		} else {
			c.instructionPointer += 3
		}
		return 0
	},
	LESS: func(c *intcode, inst *instruction) int {
		a := c.memrefread(inst.args[0])
		b := c.memrefread(inst.args[1])
		var value int

		if a < b {
			value = 1
		} else {
			value = 0
		}

		c.memrefwrite(inst.args[2], value)

		c.instructionPointer += 4

		return 0
	},
	EQUALS: func(c *intcode, inst *instruction) int {
		a := c.memrefread(inst.args[0])
		b := c.memrefread(inst.args[1])
		var value int

		if a == b {
			value = 1
		} else {
			value = 0
		}

		c.memrefwrite(inst.args[2], value)

		c.instructionPointer += 4

		return 0
	},
	RBADD: func(c *intcode, inst *instruction) int {
		a := c.memrefread(inst.args[0])
		c.relativeBase += a

		c.instructionPointer += 2

		return 0
	},
	TERMINATE: func(c *intcode, inst *instruction) int {
		close(c.outputChannel)
		c.exitChannel <- 1
		return 1
	},
}

func (c *intcode) read(dest int) {
	value, ok := <-c.inputChannel
	if !ok {
		log.Fatalf("FATAL: Reading from closed channel: %v", c.inputChannel)
	}
	c.memwrite(dest, value)
}

func (c *intcode) write(value int) {
	c.outputChannel <- value
}

func (c *intcode) run() {
	for {
		inst := newInstruction(c.memory, c.instructionPointer)

		terminate := executableCode[inst.opcode](c, &inst)

		if terminate != 0 {
			break
		}
	}
}

func (c *intcode) memread(dir int) int {
	if n := len(c.memory); dir >= n {
		c.memexpand(dir - n + 1)
	}
	return c.memory[dir]
}

func (c *intcode) memwrite(dir int, value int) {
	if n := len(c.memory); dir >= n {
		c.memexpand(dir - n + 1)
	}
	c.memory[dir] = value
}

func (c *intcode) memrefptr(r memref) int {
	if r.mode == 0 {
		return r.value

	} else if r.mode == 1 {
		return -1
	} else if r.mode == 2 {
		return c.relativeBase + r.value
	}
	log.Panic("Bad mode received ", r.mode)
	return -1
}

func (c *intcode) memrefread(r memref) int {
	ptr := c.memrefptr(r)
	if ptr < 0 {
		return r.value
	}
	return c.memread(ptr)
}

func (c *intcode) memrefwrite(r memref, value int) {
	ptr := c.memrefptr(r)
	if ptr < 0 {
		log.Panic("Bad write ptr received, mode is ", r.mode)
	}
	c.memwrite(ptr, value)
}

func (c *intcode) memexpand(size int) {
	c.memory = append(c.memory, make([]int, size)...)
}

func (c *intcode) memdump() {
	var memStr strings.Builder
	var insStr strings.Builder

	memStr.WriteString("| ")
	insStr.WriteString("  ")

	for i, value := range c.memory {
		memStr.WriteString(fmt.Sprintf("%04d | ", value))
		if i == c.instructionPointer {
			insStr.WriteString("INS^   ")
		} else if i%4 == 0 {
			insStr.WriteString(fmt.Sprintf("  %02d   ", i))
		} else {
			insStr.WriteString("       ")
		}
	}

	log.Println(memStr.String())
	log.Println(insStr.String())
}

func main() {
	var line string
	var code []int

	fmt.Scanln(&line)

	for _, s := range strings.Split(line, ",") {
		value, _ := strconv.Atoi(s)
		code = append(code, value)
	}

	inputChannel := make(chan int)
	outputChannel := make(chan int)
	exitChannel := make(chan int)

	computer := intcode{
		memory:        code,
		inputChannel:  inputChannel,
		outputChannel: outputChannel,
		exitChannel:   exitChannel,
	}

	go func() {
		for x := range outputChannel {
			fmt.Println(x)
		}
		exitChannel <- 1
	}()

	go func() {
		inputChannel <- 1
		exitChannel <- 1
	}()

	go computer.run()

	<-exitChannel
	<-exitChannel
	<-exitChannel
}
