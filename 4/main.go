package main

import (
	"flag"
	"fmt"
	"strconv"
	"strings"

	"github.com/martin-nyaga/aoc-2022/util"
)

type Range [2]int

func NewRange(str string) Range {
	sections := strings.Split(str, "-")
	first, err := strconv.Atoi(sections[0])
	util.HandleError(err)
	last, err := strconv.Atoi(sections[1])
	return Range{first, last}
}

func (r Range) Covers(o Range) bool {
	return r[0] <= o[0] && r[1] >= o[1]
}

func (r Range) Intersects(o Range) bool {
	return (r[0] <= o[0] && r[1] >= o[0]) || (r[0] <= o[1] && r[1] >= o[1])
}

type RangePair [2]Range

func (r RangePair) HasFullContainment() bool {
	return r[0].Covers(r[1]) || r[1].Covers(r[0])
}

func (r RangePair) HasIntersection() bool {
	if r.HasFullContainment() {
		return true
	}
	return r[0].Intersects(r[1])
}

func parseInput() []RangePair {
	file := util.NewInputFile("4")

	pairs := make([]RangePair, 0)
	for _, line := range file.ReadLines() {
		strs := strings.Split(line, ",")
		first := NewRange(strs[0])
		last := NewRange(strs[1])
		pairs = append(pairs, RangePair{first, last})
	}
	return pairs
}

func main() {
	flag.Parse()

	containedPairs := 0
	overlappingPairs := 0
	for _, pair := range parseInput() {
		if pair.HasFullContainment() {
			containedPairs += 1
			overlappingPairs += 1
			continue
		}

		if pair.HasIntersection() {
			overlappingPairs += 1
		}
	}

	fmt.Println("Part 1:", containedPairs)
	fmt.Println("Part 2:", overlappingPairs)
}
