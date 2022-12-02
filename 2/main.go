package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/martin-nyaga/aoc-2022/utils"
)

const (
	Rock int = iota
	Paper
	Scissors
)

const (
	_ int = iota
	Win
	Lose
	Draw
)

var Theirs = map[string]int{
	"A": Rock,
	"B": Paper,
	"C": Scissors,
}

var Mine = map[string]int{
	"X": Rock,
	"Y": Paper,
	"Z": Scissors,
}

var Endings = map[string]int{
	"X": Lose,
	"Y": Draw,
	"Z": Win,
}

type Round [2]int

func (r Round) Score() int {
	mine := r[1]
	theirs := r[0]

	diff := ((theirs-mine)%3 + 3) % 3

	score := mine + 1

	if diff == 2 {
		score += 6
	} else if diff == 0 {
		score += 3
	}

	return score
}

func (r Round) ScoreRigged() int {
	theirs := r[0]
	mine := (theirs + r[1]) % 3
	return Round{theirs, mine}.Score()
}

func parseInput(theirMap, myMap map[string]int) []Round {
	file := utils.NewInputFile("2")

	result := make([]Round, 0)
	for _, line := range file.ReadLines() {
		arr := strings.Split(line, " ")
		result = append(result, Round{theirMap[arr[0]], myMap[arr[1]]})
	}
	return result
}

func main() {
	flag.Parse()

	rounds := parseInput(Theirs, Mine)
	total := 0
	for _, round := range rounds {
		total += round.Score()
	}

	rounds = parseInput(Theirs, Endings)
	riggedTotal := 0
	for _, round := range rounds {
		riggedTotal += round.ScoreRigged()
	}

	fmt.Println("Part 1:", total)
	fmt.Println("Part 2:", riggedTotal)
}
