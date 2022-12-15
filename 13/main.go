package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"sort"

	"github.com/martin-nyaga/aoc-2022/util"
)

type Ordering int

const (
	Correct Ordering = iota
	Incorrect
	Continue
)

func compareNumbers(left, right float64) Ordering {
	if left < right {
		return Correct
	} else if left == right {
		return Continue
	} else {
		return Incorrect
	}
}

func compareArrays(left, right []interface{}) Ordering {
	var ordering Ordering
	if len(left) < len(right) {
		ordering = Correct
	} else if len(left) > len(right) {
		ordering = Incorrect
	} else {
		ordering = Continue
	}

	for i, l := range left {
		if i >= len(right) {
			ordering = Incorrect
			break
		}

		elementOrdering := checkOrder(l, right[i])
		if elementOrdering != Continue {
			ordering = elementOrdering
			break
		}
	}

	return ordering
}

func compareNumberWithArray(left float64, right []interface{}) Ordering {
	return checkOrder([]interface{}{left}, right)
}

func compareArrayWithNumber(left []interface{}, right float64) Ordering {
	return checkOrder(left, []interface{}{right})
}

func checkOrder(left, right interface{}) Ordering {
	switch l := left.(type) {
	case float64:
		switch r := right.(type) {
		case float64:
			return compareNumbers(l, r)
		case []interface{}:
			return compareNumberWithArray(l, r)
		}
	case []interface{}:
		switch r := right.(type) {
		case float64:
			return compareArrayWithNumber(l, r)
		case []interface{}:
			return compareArrays(l, r)
		default:
			fmt.Printf("Right failed to match %#v, %T\n", r, r)
		}
	default:
		fmt.Printf("Left failed to match %#v, %T\n", l, l)
	}
	panic("Match failed!")
}

func parseInput() [][2][]interface{} {
	lines := util.NewInputFile("13").ReadLines()
	packetPairs := make([][2][]interface{}, 0)
	var packetPair [2][]interface{}
	i := 0
	for i < len(lines) {
		if len(lines[i]) == 0 {
			packetPairs = append(packetPairs, packetPair)
			packetPair = [2][]interface{}{}
			i += 1
		} else {
			json.Unmarshal([]byte(lines[i]), &packetPair[0])
			i += 1
			json.Unmarshal([]byte(lines[i]), &packetPair[1])
			i += 1
		}
	}
	// append final pair
	packetPairs = append(packetPairs, packetPair)

	return packetPairs
}

type ByOrder [][]interface{}

func (a ByOrder) Len() int           { return len(a) }
func (a ByOrder) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByOrder) Less(i, j int) bool { return checkOrder(a[i], a[j]) == Correct }

func main() {
	flag.Parse()
	packetPairs := parseInput()
	result := 0
	for i, pair := range packetPairs {
		if checkOrder(interface{}(pair[0]), interface{}(pair[1])) == Correct {
			result += i + 1
		}
	}
	fmt.Println("Part 1:", result)
	packets := make([][]interface{}, 0)
	for _, pair := range packetPairs {
		packets = append(packets, pair[0], pair[1])
	}

	var divider1, divider2 []interface{}
	json.Unmarshal([]byte("[[2]]"), &divider1)
	json.Unmarshal([]byte("[[6]]"), &divider2)
	packets = append(packets, divider1, divider2)
	sort.Sort(ByOrder(packets))
	var i, j int
	for index, packet := range packets {
		marshalled, err := json.Marshal(packet)
		util.HandleError(err)
		if string(marshalled) == "[[2]]" {
			i = index + 1
		}
		if string(marshalled) == "[[6]]" {
			j = index + 1
		}
	}
	fmt.Println(i, j)
	fmt.Println("Part 2:", i*j)
}
