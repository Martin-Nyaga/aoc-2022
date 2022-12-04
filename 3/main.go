package main

import (
	"flag"
	"fmt"

	"github.com/martin-nyaga/aoc-2022/util"
)

func parseInput() []string {
	file := util.NewInputFile("3")
	return file.ReadLines()
}

func Priority(item byte) int {
	if item < 91 {
		return int(item) - 65 + 27
	} else {
		return int(item) - 96
	}
}

func main() {
	flag.Parse()
	input := parseInput()

	dups := make([]byte, 0)
	for _, line := range input {
		bytes := []byte(line)
		first := util.NewByteSet()
		for i := 0; i < len(bytes)/2; i++ {
			first.Add(bytes[i])
		}

		second := util.NewByteSet()
		for i := len(bytes) / 2; i < len(bytes); i++ {
			if first.Has(bytes[i]) && !second.Has(bytes[i]) {
				dups = append(dups, bytes[i])
			}
			second.Add(bytes[i])
		}
	}

	total := 0
	for _, item := range dups {
		total += Priority(item)
	}

	labels := 0
	for i := 0; i < (len(input) - 2); i += 3 {
		group := input[i : i+3]
		sets := make([]util.ByteSet, 3)
		for i, elf := range group {
			bytes := []byte(elf)
			sets[i] = util.NewByteSet(bytes...)
		}
		common, _ := sets[0].Intersection(&sets[1]).Intersection(&sets[2]).PopAny()
		labels += Priority(common)
	}

	fmt.Println("Part 1:", total)
	fmt.Println("Part 2:", labels)
}
