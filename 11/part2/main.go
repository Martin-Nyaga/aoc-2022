package main

import (
	"flag"
	"fmt"
	"sort"

	"github.com/martin-nyaga/aoc-2022/util"
)

type Item map[int]int

func makeItems(rawItems []int, divisors []int) []*Item {
	items := make([]*Item, 0, len(rawItems))
	for _, rawItem := range rawItems {
		item := make(Item)
		for _, divisor := range divisors {
			item[divisor] = rawItem
		}
		items = append(items, &item)
	}
	return items
}

func (i *Item) Add(x int) {
	for base, item := range *i {
		(*i)[base] = (item % base) + (x % base)
	}
}

func (i *Item) Mul(x int) {
	for base, item := range *i {
		(*i)[base] = ((item % base) * (x % base)) % base
	}
}

func (i *Item) Square() {
	for base, item := range *i {
		(*i)[base] = ((item % base) * (item % base)) % base
	}
}

func (i *Item) DivisibleBy(x int) bool {
	return (*i)[x]%x == 0
}

type Monkey struct {
	items       []*Item
	operation   func(*Item)
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
		m.Inspect(item)
		var target *Monkey
		if item.DivisibleBy(m.divisor) {
			target = troop[m.trueTarget]
		} else {
			target = troop[m.falseTarget]
		}
		target.items = append(target.items, item)
	}
}

func (m *Monkey) Inspect(item *Item) {
	m.inspections += 1
	m.operation(item)
}

func parseInput() []*Monkey {
	if *util.UseSampleInput {
		divisors := []int{23, 19, 13, 17}
		return []*Monkey{
			{
				items:       makeItems([]int{79, 98}, divisors),
				operation:   func(old *Item) { old.Mul(19) },
				divisor:     23,
				trueTarget:  2,
				falseTarget: 3,
			},
			{
				items:       makeItems([]int{54, 65, 75, 74}, divisors),
				operation:   func(old *Item) { old.Add(6) },
				divisor:     19,
				trueTarget:  2,
				falseTarget: 0,
			},
			{
				items:       makeItems([]int{79, 60, 97}, divisors),
				operation:   func(old *Item) { old.Square() },
				divisor:     13,
				trueTarget:  1,
				falseTarget: 3,
			},
			{
				items:       makeItems([]int{74}, divisors),
				operation:   func(old *Item) { old.Add(3) },
				divisor:     17,
				trueTarget:  0,
				falseTarget: 1,
			},
		}
	}
	divisors := []int{11, 3, 5, 7, 19, 2, 13, 17}
	return []*Monkey{
		{
			items:       makeItems([]int{56, 52, 58, 96, 70, 75, 72}, divisors),
			operation:   func(old *Item) { old.Mul(17) },
			divisor:     11,
			trueTarget:  2,
			falseTarget: 3,
		},
		{
			items:       makeItems([]int{75, 58, 86, 80, 55, 81}, divisors),
			operation:   func(old *Item) { old.Add(7) },
			divisor:     3,
			trueTarget:  6,
			falseTarget: 5,
		},
		{
			items:       makeItems([]int{73, 68, 73, 90}, divisors),
			operation:   func(old *Item) { old.Square() },
			divisor:     5,
			trueTarget:  1,
			falseTarget: 7,
		},
		{
			items:       makeItems([]int{72, 89, 55, 51, 59}, divisors),
			operation:   func(old *Item) { old.Add(1) },
			divisor:     7,
			trueTarget:  2,
			falseTarget: 7,
		},
		{
			items:       makeItems([]int{76, 76, 91}, divisors),
			operation:   func(old *Item) { old.Mul(3) },
			divisor:     19,
			trueTarget:  0,
			falseTarget: 3,
		},
		{
			items:       makeItems([]int{88}, divisors),
			operation:   func(old *Item) { old.Add(4) },
			divisor:     2,
			trueTarget:  6,
			falseTarget: 4,
		},
		{
			items:       makeItems([]int{64, 63, 56, 50, 77, 55, 55, 86}, divisors),
			operation:   func(old *Item) { old.Add(8) },
			divisor:     13,
			trueTarget:  4,
			falseTarget: 0,
		},
		{
			items:       makeItems([]int{79, 58}, divisors),
			operation:   func(old *Item) { old.Add(6) },
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
	for round < 10000 {
		for _, monkey := range troop {
			monkey.PlayTurn(troop)
		}
		round += 1
	}

	sort.Sort(ByInspections(troop))
	monkeyBusiness := troop[0].inspections * troop[1].inspections

	fmt.Println("Part 2:", monkeyBusiness)
}
