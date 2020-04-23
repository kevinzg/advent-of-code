package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"strconv"
	"strings"
	"sync"
)

type intcode struct {
	name               string
	pid                int
	memory             []int
	instructionPointer int
	modeMask           int
	ioChannels         []chan int
	exitChannel        chan int
	inputChannel       int
	outputChannel      int
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
		log.Printf("[%s] Closing channel %d", c.name, c.outputChannel)
		close(c.ioChannels[c.outputChannel])
		log.Printf("[%s] Executing instruction 99", c.name)
		c.exitChannel <- 1
		return 1
	},
}

func (c *intcode) read(dest int) {
	log.Printf("[%s] Reading from channel: %v", c.name, c.inputChannel)
	value, ok := <-c.ioChannels[c.inputChannel]
	if !ok {
		log.Fatalf("[%s] FATAL: Reading from closed channel: %d", c.name, c.inputChannel)
	}
	log.Printf("[%s] Got value (%d) from channel: %v", c.name, value, c.inputChannel)
	c.memory[dest] = value
}

func (c *intcode) write(value int) {
	log.Printf("[%s] Writing value (%d) to channel: %v\n", c.name, value, c.outputChannel)
	c.ioChannels[c.outputChannel] <- value
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

func (c *intcode) memdump() {
	var memStr strings.Builder
	var insStr strings.Builder

	memStr.WriteString("| ")
	insStr.WriteString("  ")

	for i, value := range c.memory {
		memStr.WriteString(fmt.Sprintf("%04d | ", value))
		if i == c.instructionPointer {
			insStr.WriteString("INS^   ")
		} else {
			insStr.WriteString("       ")
		}
	}

	log.Println(memStr.String())
	log.Println(insStr.String())
}

const n int = 5

func main() {
	log.SetOutput(ioutil.Discard)

	var line string
	var code []int

	fmt.Scanln(&line)

	for _, s := range strings.Split(line, ",") {
		value, _ := strconv.Atoi(s)
		code = append(code, value)
	}

	phaseSettings := [n]int{5, 6, 7, 8, 9}
	idx := [n]int{-1}

	max := math.MinInt32

	callback := func(arr [n]int) {
		output := runAmplifiers(code, arr)
		if output > max {
			max = output
		}
	}

	permute(phaseSettings, 0, idx, 0, callback)

	fmt.Println(max)
}

func permute(arr [n]int, mask int, idx [n]int, m int, callback func([n]int)) {
	if m == n {
		callback(idx)
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

func runAmplifiers(code []int, phaseSettings [n]int) int {
	initialInput := 0

	ioChannels := make([]chan int, n)
	exitChannels := make([]chan int, n)
	for i := range ioChannels {
		ioChannels[i] = make(chan int, 2)
		exitChannels[i] = make(chan int)
	}

	names := []string{"Alfa", "Bravo", "Charlie", "Delta", "Echo"}

	var wg sync.WaitGroup

	for i, phase := range phaseSettings {
		wg.Add(1)

		go func(amplifier int, phase int, inputChannel chan int) {
			defer wg.Done()

			log.Printf("[main] Writing input phase %d to channel %d:", phase, amplifier)
			inputChannel <- phase

			if amplifier == 0 {
				log.Printf("[main] Writing initial input to channel %d\n", amplifier)
				inputChannel <- initialInput
			}
		}(i, phase, ioChannels[i])
	}

	wg.Wait()

	log.Printf("Initial input has been written, starting intcode computers")

	for amplifier := range phaseSettings {
		nextAmplifier := (amplifier + 1) % n

		localMemory := make([]int, len(code))
		copy(localMemory, code)

		computer := intcode{
			name:          names[amplifier],
			pid:           amplifier,
			memory:        localMemory,
			ioChannels:    ioChannels,
			inputChannel:  amplifier,
			outputChannel: nextAmplifier,
			exitChannel:   exitChannels[amplifier],
		}

		go computer.run()
	}

	for i := range exitChannels {
		if i != 0 {
			go func(i int, exitChannel <-chan int) {
				<-exitChannel
				log.Printf("Got exit from amplifier %d", i)
			}(i, exitChannels[i])
		}
	}

	<-exitChannels[0]
	log.Printf("Got exit from amplifier %d", 0)

	log.Printf("[main] Reading final output from channel %v", 0)
	return <-ioChannels[0]
}
