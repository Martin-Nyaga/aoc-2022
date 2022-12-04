package main

import (
	"flag"
	"fmt"

	"github.com/martin-nyaga/aoc-2022/util"
	set "k8s.io/apimachinery/pkg/util/sets"
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
		first := make(map[byte]bool)
		for i := 0; i < len(bytes)/2; i++ {
			first[bytes[i]] = true
		}

		second := make(map[byte]bool)
		for i := len(bytes) / 2; i < len(bytes); i++ {
			if first[bytes[i]] {
				dups = append(dups, bytes[i])
			}
			second[bytes[i]] = true
		}
	}

	total := 0
	for _, item := range dups {
		total += Priority(item)
	}

	labels := 0
	for i := 0; i < (len(input) - 2); i += 3 {
		group := input[i : i+3]
		sets := make([]set.Byte, 3)
		for i, elf := range group {
			bytes := []byte(elf)
			sets[i] = set.NewByte(bytes...)
		}
		common, _ := sets[0].Intersection(sets[1]).Intersection(sets[2]).PopAny()
		labels += Priority(common)
	}

	fmt.Println("Part 1:", total)
	fmt.Println("Part 2:", labels)
}
