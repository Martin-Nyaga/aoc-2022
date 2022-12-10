package main

import (
	"flag"
	"fmt"
	"math"

	"github.com/martin-nyaga/aoc-2022/util"
	"github.com/martin-nyaga/aoc-2022/util/set"
)

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

type Point [2]int

type TrackedPoint struct {
	point   Point
	tracker set.Set[Point]
	next    *TrackedPoint
}

func newTrackedPoint(point Point) TrackedPoint {
	var t TrackedPoint
	t.point = point
	t.tracker = set.NewSet(point)
	t.next = nil
	return t
}

func (t *TrackedPoint) MoveAndPropagate(m *Move) {
	switch m.direction {
	case Up:
		for i := 0; i < m.steps; i++ {
			t.point[1] -= 1
			t.Propagate()
		}
	case Down:
		for i := 0; i < m.steps; i++ {
			t.point[1] += 1
			t.Propagate()
		}
	case Left:
		for i := 0; i < m.steps; i++ {
			t.point[0] -= 1
			t.Propagate()
		}
	case Right:
		for i := 0; i < m.steps; i++ {
			t.point[0] += 1
			t.Propagate()
		}
	}

	t.tracker.Add(t.point)
}

func (t *TrackedPoint) Propagate() {
	if t.next != nil {
		t.next.FollowAndPropagate(&t.point)
	}
}

func (t *TrackedPoint) FollowAndPropagate(prev *Point) {
	dx := prev[0] - t.point[0]
	dy := prev[1] - t.point[1]
	if math.Abs(float64(dx)) <= 1 && math.Abs(float64(dy)) <= 1 {
		// Touching, nothing to do
		return
	}

	if math.Abs(float64(dx)) == 2 && dy == 0 {
		// 2 steps left or right
		t.point[0] += dx / 2
	} else if math.Abs(float64(dy)) == 2 && dx == 0 {
		// 2 steps up or down
		t.point[1] += dy / 2
	} else {
		// diagonal
		t.point[0] += int(math.Copysign(1, float64(dx)))
		t.point[1] += int(math.Copysign(1, float64(dy)))
	}

	t.tracker.Add(t.point)
	t.Propagate()
}

func parseInput() []Move {
	file := util.NewInputFile("9")
	moves := make([]Move, 0)
	for _, line := range file.ReadLines() {
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

	head := newTrackedPoint(Point{0, 0})
	end := newTrackedPoint(Point{0, 0})
	tail := &end
	head.next = tail

	for _, move := range moves {
		head.MoveAndPropagate(&move)
	}

	fmt.Println("Part 1:", tail.tracker.Len())

	head = newTrackedPoint(Point{0, 0})
	curr := &head
	for i := 0; i < 9; i++ {
		next := newTrackedPoint(Point{0, 0})
		tail = &next
		curr.next = tail
		curr = tail
	}

	for _, move := range moves {
		head.MoveAndPropagate(&move)
	}

	fmt.Println("Part 2:", tail.tracker.Len())
}
