package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/martin-nyaga/aoc-2022/util"
)

type Grid [][]byte

// ScannerFunc scans the grid in a specific direction, invoking the callback for
// each encountered tree
type ScannerFunc func(g Grid, i, j int, callback ScannerCallback)

// ScanCallback is passed to the scan{Up,Down,Left,Right} functions. If it
// returns true, the scan is continued, otherwise, it stops
type ScannerCallback func(tree, next byte) bool

func (g Grid) scanUp(i, j int, callback ScannerCallback) {
	if j == 0 {
		return
	}
	tree := g[j][i]
	for row := j - 1; row >= 0; row-- {
		result := callback(tree, g[row][i])
		if !result {
			break
		}
	}
}

func (g Grid) scanDown(i, j int, callback ScannerCallback) {
	if j == (len(g) - 1) {
		return
	}
	tree := g[j][i]
	for row := j + 1; row < len(g); row++ {
		result := callback(tree, g[row][i])
		if !result {
			break
		}
	}
}

func (g Grid) scanLeft(i, j int, callback ScannerCallback) {
	if i == 0 {
		return
	}
	tree := g[j][i]
	for col := i - 1; col >= 0; col-- {
		result := callback(tree, g[j][col])
		if !result {
			break
		}
	}
}

func (g Grid) scanRight(i, j int, callback ScannerCallback) {
	if i == (len(g[0]) - 1) {
		return
	}
	tree := g[j][i]
	for col := i + 1; col < len(g[0]); col++ {
		result := callback(tree, g[j][col])
		if !result {
			break
		}
	}
}

func (g Grid) IsVisible(i, j int) bool {
	visible := g.visibleFromDirection(i, j, Grid.scanUp) ||
		g.visibleFromDirection(i, j, Grid.scanDown) ||
		g.visibleFromDirection(i, j, Grid.scanLeft) ||
		g.visibleFromDirection(i, j, Grid.scanRight)
	return visible
}

func (g Grid) visibleFromDirection(i, j int, scan ScannerFunc) bool {
	canBeSeen := true
	scan(g, i, j, func(tree, next byte) bool {
		if tree <= next {
			canBeSeen = false
			return false
		}
		return true
	})
	return canBeSeen
}

func (g Grid) ScenicScore(i, j int) int {
	top := g.countTreesSeenInDirection(i, j, Grid.scanUp)
	down := g.countTreesSeenInDirection(i, j, Grid.scanDown)
	left := g.countTreesSeenInDirection(i, j, Grid.scanLeft)
	right := g.countTreesSeenInDirection(i, j, Grid.scanRight)
	return top * down * left * right
}

func (g Grid) countTreesSeenInDirection(i, j int, scan ScannerFunc) int {
	seen := 0
	scan(g, i, j, func(tree, next byte) bool {
		seen += 1
		if next >= tree {
			return false
		}
		return true
	})
	return seen
}

func parseInput() Grid {
	file := util.NewInputFile("8").ReadLines()
	trees := make([][]byte, 0)
	for _, line := range file {
		row := []byte(strings.TrimSpace(line))
		trees = append(trees, row)
	}
	return trees
}

func main() {
	flag.Parse()

	grid := parseInput()
	visibleCount := 0
	for j, row := range grid {
		for i := range row {
			if grid.IsVisible(i, j) {
				visibleCount += 1
			}
		}
	}

	highestScore := 0
	for j, row := range grid {
		for i := range row {
			score := grid.ScenicScore(i, j)
			if score > highestScore {
				highestScore = score
			}
		}
	}

	fmt.Println("Part 1:", visibleCount)
	fmt.Println("Part 2:", highestScore)
}
