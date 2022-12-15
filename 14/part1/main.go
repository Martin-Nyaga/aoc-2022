package main

import (
	"flag"
	"fmt"
	"strings"

	tm "github.com/buger/goterm"
	"github.com/eiannone/keyboard"
	"github.com/martin-nyaga/aoc-2022/util"
	"github.com/martin-nyaga/aoc-2022/util/set"
)

var logBox = tm.NewBox(30|tm.PCT, 20, 0)
var caveBox = tm.NewBox(50|tm.PCT, 50, 0)

var Draw = flag.Bool("draw", false, "Draw the simulation each frame")
var Step = flag.Bool("step", false, "Advance the simulation by pressing any key, esc exits")

func log(str string) {
	fmt.Fprint(logBox, str+"\n")
}

type Point [2]int
type Path []Point

type Orientation int

const (
	Vertical Orientation = iota
	Horizontal
)

func (p Path) eachSegment(fn func(Point, Point, Orientation)) {
	for i := 0; i < len(p)-1; i += 1 {
		if p[i][0] == p[i+1][0] {
			// vertical
			if p[i][1] > p[i+1][1] {
				fn(p[i+1], p[i], Vertical)
			} else {
				fn(p[i], p[i+1], Vertical)
			}
		} else {
			// horizontal
			if p[i][0] > p[i+1][0] {
				fn(p[i+1], p[i], Horizontal)
			} else {
				fn(p[i], p[i+1], Horizontal)
			}
		}
	}
}

func (p Path) blocks(point Point) bool {
	doesBlock := false

	p.eachSegment(func(start, end Point, orientation Orientation) {
		if orientation == Vertical {
			if start[0] == point[0] && start[1] <= point[1] && end[1] >= point[1] {
				doesBlock = true
			}
		} else {
			if start[1] == point[1] && start[0] <= point[0] && end[0] >= point[0] {
				doesBlock = true
			}
		}
	})

	return doesBlock
}

type Cave struct {
	rockPaths []Path
	grains    set.Set[Point]
}

func (c *Cave) nextGrain() Point {
	return Point{500, 0}
}

func (c *Cave) grainIsFallingIntoAbyss(g Point) bool {
	indeed := true
outer:
	for _, path := range c.rockPaths {
		for _, point := range path {
			if point[1] > g[1] {
				indeed = false
				break outer
			}
		}
	}
	return indeed
}

func (c *Cave) canMoveTo(p Point) bool {
	if c.grains.Has(p) {
		log("Blocked by grain!")
		return false
	}
	for _, path := range c.rockPaths {
		if path.blocks(p) {
			log("Blocked by rock!")
			return false
		}
	}

	return true
}

func pause() {
	if *Step {
		// Progress by pressing any key
		_, key, err := keyboard.GetSingleKey()
		util.HandleError(err)
		if key == keyboard.KeyEsc {
			panic("Escape pressed!")
		}
	}
}

func (c *Cave) addSandUntilDone() {
	for {
		g := c.nextGrain()

		fellIntoTheAbyss := false
		settled := false

		for {
			// Pause and draw the current state for debuging
			if *Draw {
				pause()
				c.grains.Add(g)
				c.Draw()
				c.grains.Remove(g)
			}
			if c.grainIsFallingIntoAbyss(g) {
				log(fmt.Sprintf("%#v fell into the abyss", g))
				fellIntoTheAbyss = true
				break
			}

			if settled {
				log(fmt.Sprintf("%#v settled", g))
				break
			}

			// Try move down
			down := Point{g[0], g[1] + 1}
			if c.canMoveTo(down) {
				log(fmt.Sprintf("%#v can move down", g))
				g[1] += 1
				continue
			}
			// Try move diagonally left
			downLeft := Point{g[0] - 1, g[1] + 1}
			if c.canMoveTo(downLeft) {
				log(fmt.Sprintf("%#v can move down left", g))
				g[0] -= 1
				g[1] += 1
				continue
			}
			// Try move diagonally right
			downRight := Point{g[0] + 1, g[1] + 1}
			if c.canMoveTo(downRight) {
				log(fmt.Sprintf("%#v can move down right", g))
				g[0] += 1
				g[1] += 1
				continue
			}

			settled = true
		}

		if fellIntoTheAbyss {
			log(fmt.Sprintf("%#v fell into the abyss", g))
			break
		} else {
			c.grains.Add(g)
		}
	}
}

func (c *Cave) minX() int {
	minX := 500
	for _, path := range c.rockPaths {
		for _, point := range path {
			if point[0] < minX {
				minX = point[0]
			}
		}
	}
	return minX
}

func (c *Cave) Draw() {
	tm.Clear()

	for _, path := range c.rockPaths {
		path.eachSegment(func(start, end Point, orientation Orientation) {
			if orientation == Vertical {
				for y := start[1]; y <= end[1]; y++ {
					tm.MoveCursor(start[0]-c.minX(), y)
					tm.Print("#")
				}
			}
			if orientation == Horizontal {
				for x := start[0]; x <= end[0]; x++ {
					tm.MoveCursor(x-c.minX(), start[1])
					tm.Print("#")
				}
			}
		})
	}

	c.grains.Each(func(g Point) {
		tm.MoveCursor(g[0]-c.minX(), g[1])
		tm.Print("o")
	})

	// print log box
	tm.Print(tm.MoveTo(logBox.String(), 70|tm.PCT, 5|tm.PCT))
	tm.Flush()
	// Truncate logbox after a while
	if logBox.Buf.Len() >= 400 {
		logBox.Buf.Reset()
	}
	fmt.Println()
}

func parseInput() Cave {
	lines := util.NewInputFile("14").ReadLines()
	paths := make([]Path, 0)
	for _, line := range lines {
		rawPoints := strings.Split(line, " -> ")
		path := make([]Point, 0)
		for _, point := range rawPoints {
			var x, y int
			fmt.Sscanf(point, "%d,%d", &x, &y)
			path = append(path, Point{x, y})
		}
		paths = append(paths, path)
	}
	return Cave{rockPaths: paths, grains: set.NewSet[Point]()}
}

func main() {
	flag.Parse()
	cave := parseInput()
	cave.addSandUntilDone()

	fmt.Println("Part 1:", cave.grains.Len())
}
