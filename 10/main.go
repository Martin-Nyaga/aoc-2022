package main

import (
	"flag"
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/martin-nyaga/aoc-2022/util"
	"github.com/martin-nyaga/aoc-2022/util/slices"
)

const (
	Addx = "addx"
	Noop = "noop"
)

type Insn struct {
	kind string
	arg  int
}

type Cpu struct {
	x       int
	cycles  int
	samples []int
	crt     Crt
}

func newCpu() Cpu {
	return Cpu{
		x:       1,
		cycles:  0,
		samples: make([]int, 0),
		crt:     Crt{screen: [6][40]byte{}},
	}
}

func (c *Cpu) Process(insn *Insn) {
	switch insn.kind {
	case Addx:
		c.nextCycle()
		c.nextCycle()
		c.x += insn.arg
	case Noop:
		c.nextCycle()
	}
}

func (c *Cpu) nextCycle() {
	c.cycles += 1
	if c.cycles == 20 || c.cycles <= 220 && ((c.cycles-20)%40 == 0) {
		c.samples = append(c.samples, c.cycles*c.x)
	}
	c.crt.Draw(c.x)
}

type Crt struct {
	currentPixel int
	screen       [6][40]byte
}

func (c *Crt) Draw(x int) {
	row := c.currentPixel / 40
	col := c.currentPixel % 40
	if math.Abs(float64(x-col)) <= 1 {
		c.screen[row][col] = '#'
	} else {
		c.screen[row][col] = '.'
	}
	c.currentPixel += 1
}

func (c *Crt) Print() {
	for _, row := range c.screen {
		for _, px := range row {
			fmt.Print(string(px))
		}
		fmt.Println()
	}
}

func parseInput() []Insn {
	lines := util.NewInputFile("10").ReadLines()
	insns := make([]Insn, 0, len(lines))
	for _, line := range lines {
		arr := strings.Split(line, " ")
		var insn Insn
		insn.kind = arr[0]
		switch arr[0] {
		case Addx:
			arg, err := strconv.Atoi(arr[1])
			util.HandleError(err)
			insn.arg = arg
		}
		insns = append(insns, insn)
	}

	return insns
}

func main() {
	flag.Parse()
	insns := parseInput()

	cpu := newCpu()
	for _, insn := range insns {
		cpu.Process(&insn)
	}

	fmt.Println("Part 1:", slices.Sum(cpu.samples))
	fmt.Println("Part 2:")
	cpu.crt.Print()
}
