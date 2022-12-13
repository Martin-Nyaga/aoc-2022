package main

import (
	"flag"
	"fmt"
	"math"

	"github.com/martin-nyaga/aoc-2022/util"
	"github.com/martin-nyaga/aoc-2022/util/pqueue"
	"github.com/martin-nyaga/aoc-2022/util/set"
)

type Point [2]int

type HeightMap struct {
	start Point
	grid  [][]byte
	goal  Point
}

func (h *HeightMap) AccessibleNeighbours(point Point) []Point {
	x := point[0]
	y := point[1]
	result := make([]Point, 0)
	if (x + 1) < len(h.grid[0]) {
		nextPoint := Point{x + 1, y}
		if h.canReach(point, nextPoint) {
			result = append(result, nextPoint)
		}
	}
	if (x - 1) >= 0 {
		nextPoint := Point{x - 1, y}
		if h.canReach(point, nextPoint) {
			result = append(result, nextPoint)
		}
	}
	if (y + 1) < len(h.grid) {
		nextPoint := Point{x, y + 1}
		if h.canReach(point, nextPoint) {
			result = append(result, nextPoint)
		}
	}
	if (y - 1) >= 0 {
		nextPoint := Point{x, y - 1}
		if h.canReach(point, nextPoint) {
			result = append(result, nextPoint)
		}
	}
	return result
}

func (h *HeightMap) At(p Point) byte {
	return h.grid[p[1]][p[0]]
}

func (h *HeightMap) canReach(source, dest Point) bool {
	return (int(h.At(dest)) - int(h.At(source))) <= 1
}

func (h *HeightMap) Heuristic(point PointWithSteps) int {
	return Manhattan(point.Point, h.goal) + point.steps
}

func Manhattan(p1, p2 Point) int {
	return int(math.Abs(float64(p2[0]-p1[0])) + math.Abs(float64(p2[1]-p1[1])))
}

func (h *HeightMap) Print() {
	for _, row := range h.grid {
		for _, val := range row {
			fmt.Print(string(val))
		}
		fmt.Println()
	}
}

func parseInput() (HeightMap, []Point) {
	file := util.NewInputFile("12")
	grid := make([][]byte, 0)
	potentialStarts := make([]Point, 0)
	var start, goal [2]int

	for j, line := range file.ReadLines() {
		row := make([]byte, 0)
		for i, char := range []byte(line) {
			if char == 'E' {
				goal = [2]int{i, j}
				char = 'z'
			}
			if char == 'S' {
				start = [2]int{i, j}
				char = 'a'
			}
			if char == 'a' {
				potentialStarts = append(potentialStarts, Point{i, j})
			}
			row = append(row, char)
		}

		grid = append(grid, row)
	}

	return HeightMap{
		start: start,
		grid:  grid,
		goal:  goal,
	}, potentialStarts
}

type PointWithSteps struct {
	Point
	steps int
	path  []Point
}

func main() {
	flag.Parse()

	heightMap, starts := parseInput()
	var minDist int
	for _, start := range starts {
		heightMap.start = start
		pq := pqueue.NewPqueue[int, PointWithSteps](pqueue.MinQueue)
		pq.Push(0, PointWithSteps{heightMap.start, 0, []Point{heightMap.start}})
		visited := set.NewSet[Point]()
		var winningPoint PointWithSteps
		for pq.Len() > 0 {
			nextPoint, err := pq.Pop()
			util.HandleError(err)

			if visited.Has(nextPoint.Point) {
				continue
			}
			visited.Add(nextPoint.Point)

			if nextPoint.Point == heightMap.goal {
				winningPoint = nextPoint
				if (minDist == 0) || (minDist > winningPoint.steps) {
					minDist = winningPoint.steps
				}
				break
			}

			neighbours := heightMap.AccessibleNeighbours(nextPoint.Point)
			for _, point := range neighbours {
				nextPath := make([]Point, 0)
				for _, p := range nextPoint.path {
					nextPath = append(nextPath, p)
				}
				nextPath = append(nextPath, point)

				pointToAdd := PointWithSteps{point, nextPoint.steps + 1, nextPath}
				pq.Push(heightMap.Heuristic(pointToAdd), pointToAdd)
			}
		}
	}

	fmt.Println("Part 2:", minDist)
}
