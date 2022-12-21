package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

type reaction struct {
	out string
	n   int
	in  map[string]int
}

type factory struct {
	n         int
	state     []int
	reactions [][]int
}

func (f *factory) solve(target int, ore int) int {
	f.state[target] = -1

	for {
		stop := true
		for chemical, amount := range f.state {
			if chemical == ore {
				continue
			}
			if amount < 0 {
				f.produce(chemical)
				stop = false
			}
		}
		if stop {
			break
		}
	}
	return -f.state[ore]
}

func (f *factory) produce(chemical int) {
	r := f.reactions[chemical]
	amount := -f.state[chemical]
	if amount <= 0 {
		log.Panic("Produce should only be called if the amount needed is positive")
	}
	times := amount / -r[chemical]
	if (amount % -r[chemical]) != 0 {
		times++
	}
	f.state = add(f.state, mul(r, -times))
}

func main() {
	reactions := readInput()
	n := len(reactions) + 1
	idx := make(map[string]int)
	r := make([][]int, n)
	for k, v := range reactions {
		idx[v.out] = k
	}
	idx["ORE"] = n - 1
	for k, v := range reactions {
		r[k] = mapToSlice(n, v.in, idx)
		r[k][k] = -v.n
	}
	f := factory{
		n:         n,
		state:     make([]int, n),
		reactions: r,
	}
	fmt.Println(f.solve(idx["FUEL"], idx["ORE"]))
}

func add(a []int, b []int) []int {
	c := make([]int, len(a))
	for i, v := range a {
		c[i] = v + b[i]
	}
	return c
}

func mul(a []int, k int) []int {
	c := make([]int, len(a))
	for i, v := range a {
		c[i] = v * k
	}
	return c
}

func mapToSlice(n int, in map[string]int, idx map[string]int) []int {
	s := make([]int, n)
	for k, v := range in {
		s[idx[k]] = v
	}
	return s
}

func readInput() []reaction {
	reactions := make([]reaction, 0)
	reader := bufio.NewReader(os.Stdin)

	for {
		line, err := reader.ReadString('\n')
		line = strings.TrimRight(line, "\n")
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Panic(err)
		}
		reactions = append(reactions, parseLine(line))
	}
	return reactions
}

func parseLine(line string) reaction {
	tokens := strings.Split(line, " ")
	i := 0

	readToken := func() string {
		i++
		return tokens[i-1]
	}
	nextToken := func() string {
		if i >= len(tokens) {
			return ""
		}
		return tokens[i]
	}
	readChemical := func() (c string, n int) {
		n, err := strconv.Atoi(readToken())
		if err != nil {
			log.Panic(err)
		}
		c = strings.TrimRight(readToken(), ",")
		return
	}

	chemicals := make(map[string]int)
	for next := nextToken(); next != "=>" && next != ""; next = nextToken() {
		c, n := readChemical()
		chemicals[c] = n
	}
	readToken() // "=>"

	c, n := readChemical()

	return reaction{
		out: c,
		n:   n,
		in:  chemicals,
	}
}
