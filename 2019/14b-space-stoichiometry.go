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
	in  map[string]int64
}

type factory struct {
	n         int
	state     []int64
	reactions [][]int64
}

func (f *factory) solve(target, ore int, desired int64) int64 {
	f.state[target] = -desired
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
	r := make([][]int64, n)
	for k, v := range reactions {
		idx[v.out] = k
	}
	idx["ORE"] = n - 1
	for k, v := range reactions {
		r[k] = mapToSlice(n, v.in, idx)
		r[k][k] = int64(-v.n)
	}

	fuel := binarySearch(1, 10000000, 1000000000000, func(t int) int64 {
		f := factory{
			n:         n,
			state:     make([]int64, n),
			reactions: r,
		}
		return f.solve(idx["FUEL"], idx["ORE"], int64(t))
	})
	fmt.Println(fuel)
}

func binarySearch(l, h int, target int64, fn func(int) int64) int {
	if l >= h-1 {
		return l
	}

	m := (l + h) / 2
	res := fn(m)

	if res > target {
		return binarySearch(l, m, target, fn)
	} else {
		return binarySearch(m, h, target, fn)
	}
}

func add(a []int64, b []int64) []int64 {
	c := make([]int64, len(a))
	for i, v := range a {
		c[i] = v + b[i]
	}
	return c
}

func mul(a []int64, k int64) []int64 {
	c := make([]int64, len(a))
	for i, v := range a {
		c[i] = v * k
	}
	return c
}

func mapToSlice(n int, in map[string]int64, idx map[string]int) []int64 {
	s := make([]int64, n)
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

	chemicals := make(map[string]int64)
	for next := nextToken(); next != "=>" && next != ""; next = nextToken() {
		c, n := readChemical()
		chemicals[c] = int64(n)
	}
	readToken() // "=>"

	c, n := readChemical()

	return reaction{
		out: c,
		n:   n,
		in:  chemicals,
	}
}
