package main

import (
	"flag"
	"fmt"
	"sort"
	"strconv"

	"github.com/martin-nyaga/aoc-2022/util"
	"github.com/martin-nyaga/aoc-2022/util/slices"
)

func parseInput() [][]int {
	file := util.NewInputFile("1")

	deers := make([][]int, 0)
	deer := make([]int, 0)

	lines := file.ReadLines()
	for i, line := range lines {
		if len(line) > 0 {
			n, err := strconv.Atoi(line)
			util.HandleError(err)
			deer = append(deer, n)
		}

		if (len(line) == 0) || (i == (len(lines) - 1)) {
			deers = append(deers, deer)
			deer = make([]int, 0)
		}
	}

	return deers
}

func main() {
	flag.Parse()

	deers := parseInput()
	totals := make([]int, len(deers))
	for _, deer := range deers {
		totals = append(totals, slices.Sum(deer))
	}

	sort.Ints(totals)
	most := totals[len(totals)-1]
	topThree := slices.Sum(totals[len(totals)-3:])

	fmt.Println("Part 1:", most)
	fmt.Println("Part 2:", topThree)
}
