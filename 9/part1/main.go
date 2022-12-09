package main

import (
	"flag"
	"fmt"
	"math"

	"github.com/martin-nyaga/aoc-2022/util"
	"github.com/martin-nyaga/aoc-2022/util/set"
)

type Point [2]int

func (head *Point) Move(m *Move, next *Point, tracker *set.Set[Point]) {
	switch m.direction {
	case Up:
		for i := 0; i < m.steps; i++ {
			head[1] -= 1
			next.Follow(head, tracker)
		}
	case Down:
		for i := 0; i < m.steps; i++ {
			head[1] += 1
			next.Follow(head, tracker)
		}
	case Left:
		for i := 0; i < m.steps; i++ {
			head[0] -= 1
			next.Follow(head, tracker)
		}
	case Right:
		for i := 0; i < m.steps; i++ {
			head[0] += 1
			next.Follow(head, tracker)
		}
	}
}

func (tail *Point) Follow(prev *Point, tracker *set.Set[Point]) {
	dx := prev[0] - tail[0]
	dy := prev[1] - tail[1]
	if math.Abs(float64(dx)) <= 1 && math.Abs(float64(dy)) <= 1 {
		// Touching, nothing to do
		return
	}

	if math.Abs(float64(dx)) == 2 && dy == 0 {
		// 2 steps left or right
		tail[0] += dx / 2
	} else if math.Abs(float64(dy)) == 2 && dx == 0 {
		// 2 steps up or down
		tail[1] += dy / 2
	} else {
		// diagonal
		tail[0] += int(math.Copysign(1, float64(dx)))
		tail[1] += int(math.Copysign(1, float64(dy)))
	}
	tracker.Add(*tail)
}

type Move struct {
	direction byte
	steps     int
}

const (
	Up    = 'U'
	Down  = 'D'
	Right = 'R'
	Left  = 'L'
)

func parseInput() []Move {
	file := util.NewInputFile("9").ReadLines()
	moves := make([]Move, 0)
	for _, line := range file {
		var direction string
		var steps int
		fmt.Sscanf(line, "%1s %d", &direction, &steps)
		moves = append(moves, Move{[]byte(direction)[0], steps})
	}
	return moves
}

func main() {
	flag.Parse()
	moves := parseInput()
	head := Point{0, 0}
	tail := Point{0, 0}
	tracker := set.NewSet[Point]()
	tracker.Add(tail)
	for _, move := range moves {
		head.Move(&move, &tail, &tracker)
	}

	fmt.Println("Part 1:", tracker.Len())
}
