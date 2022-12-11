package main

import (
	"flag"
	"fmt"
	"sort"

	"github.com/martin-nyaga/aoc-2022/util"
)

type Monkey struct {
	items       []int
	operation   func(int) int
	divisor     int
	trueTarget  int
	falseTarget int
	inspections int
}

func (m *Monkey) PlayTurn(troop []*Monkey) {
	itemCount := len(m.items)
	for i := 0; i < itemCount; i++ {
		item := m.items[0]
		m.items = m.items[1:]
		item = m.Inspect(item)
		item = CalmDown(item)
		var target *Monkey
		if item%m.divisor == 0 {
			target = troop[m.trueTarget]
		} else {
			target = troop[m.falseTarget]
		}
		target.items = append(target.items, item)
	}
}

func (m *Monkey) Inspect(item int) int {
	m.inspections += 1
	return m.operation(item)
}

func CalmDown(item int) int {
	return item / 3
}

func parseInput() []*Monkey {
	if *util.UseSampleInput {
		return []*Monkey{
			{
				items:       []int{79, 98},
				operation:   func(old int) int { return old * 19 },
				divisor:     23,
				trueTarget:  2,
				falseTarget: 3,
			},
			{
				items:       []int{54, 65, 75, 74},
				operation:   func(old int) int { return old + 6 },
				divisor:     19,
				trueTarget:  2,
				falseTarget: 0,
			},
			{
				items:       []int{79, 60, 97},
				operation:   func(old int) int { return old * old },
				divisor:     13,
				trueTarget:  1,
				falseTarget: 3,
			},
			{
				items:       []int{74},
				operation:   func(old int) int { return old + 3 },
				divisor:     17,
				trueTarget:  0,
				falseTarget: 1,
			},
		}
	}
	return []*Monkey{
		{
			items:       []int{56, 52, 58, 96, 70, 75, 72},
			operation:   func(old int) int { return old * 17 },
			divisor:     11,
			trueTarget:  2,
			falseTarget: 3,
		},
		{
			items:       []int{75, 58, 86, 80, 55, 81},
			operation:   func(old int) int { return old + 7 },
			divisor:     3,
			trueTarget:  6,
			falseTarget: 5,
		},
		{
			items:       []int{73, 68, 73, 90},
			operation:   func(old int) int { return old * old },
			divisor:     5,
			trueTarget:  1,
			falseTarget: 7,
		},
		{
			items:       []int{72, 89, 55, 51, 59},
			operation:   func(old int) int { return old + 1 },
			divisor:     7,
			trueTarget:  2,
			falseTarget: 7,
		},
		{
			items:       []int{76, 76, 91},
			operation:   func(old int) int { return old * 3 },
			divisor:     19,
			trueTarget:  0,
			falseTarget: 3,
		},
		{
			items:       []int{88},
			operation:   func(old int) int { return old + 4 },
			divisor:     2,
			trueTarget:  6,
			falseTarget: 4,
		},
		{
			items:       []int{64, 63, 56, 50, 77, 55, 55, 86},
			operation:   func(old int) int { return old + 8 },
			divisor:     13,
			trueTarget:  4,
			falseTarget: 0,
		},
		{
			items:       []int{79, 58},
			operation:   func(old int) int { return old + 6 },
			divisor:     17,
			trueTarget:  1,
			falseTarget: 5,
		},
	}
}

type ByInspections []*Monkey

func (a ByInspections) Len() int           { return len(a) }
func (a ByInspections) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByInspections) Less(i, j int) bool { return a[i].inspections > a[j].inspections }

func main() {
	flag.Parse()

	troop := parseInput()

	round := 0
	for round < 20 {
		for _, monkey := range troop {
			monkey.PlayTurn(troop)
		}
		round += 1
	}

	sort.Sort(ByInspections(troop))
	monkeyBusiness := troop[0].inspections * troop[1].inspections

	fmt.Println("Part 1:", monkeyBusiness)
}
