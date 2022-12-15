package main

import (
	"flag"
	"fmt"
	"math"
	"os"
	"runtime/pprof"
	"strings"

	tm "github.com/buger/goterm"
	"github.com/eiannone/keyboard"
	"github.com/martin-nyaga/aoc-2022/util"
	"github.com/martin-nyaga/aoc-2022/util/set"
)

var logBox = tm.NewBox(30|tm.PCT, 20, 5)
var caveBox = tm.NewBox(50|tm.PCT, 50, 0)

var Draw = flag.Bool("draw", false, "Draw the simulation each frame")
var Step = flag.Bool("step", false, "Advance the simulation by pressing any key, esc exits")
var Prof = flag.String("prof", "", "Generate cpu profile")

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

func (p Path) eachSegment(fn func(Point, Point, Orientation) bool) {
	var result bool
	for i := 0; i < len(p)-1; i += 1 {
		if p[i][0] == p[i+1][0] {
			// vertical
			if p[i][1] > p[i+1][1] {
				result = fn(p[i+1], p[i], Vertical)
			} else {
				result = fn(p[i], p[i+1], Vertical)
			}
		} else {
			// horizontal
			if p[i][0] > p[i+1][0] {
				result = fn(p[i+1], p[i], Horizontal)
			} else {
				result = fn(p[i], p[i+1], Horizontal)
			}
		}
		if !result {
			break
		}
	}
}

func (p Path) blocks(point Point) bool {
	doesBlock := false

	p.eachSegment(func(start, end Point, orientation Orientation) bool {
		if orientation == Vertical {
			if start[0] == point[0] && start[1] <= point[1] && end[1] >= point[1] {
				doesBlock = true
				return false
			}
		} else {
			if start[1] == point[1] && start[0] <= point[0] && end[0] >= point[0] {
				doesBlock = true
				return false
			}
		}
		return true
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

		settledAtRoot := false
		settled := false

		for {
			// Pause and draw the current state for debuging
			if *Draw {
				pause()
				c.grains.Add(g)
				c.Draw()
				c.grains.Remove(g)
			}

			if settled {
				log(fmt.Sprintf("%#v settled", g))
				if g == c.nextGrain() {
					log(fmt.Sprintf("%#v blocked the entrance, the sand has stopped!", g))
					settledAtRoot = true
				}
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

		c.grains.Add(g)
		if settledAtRoot {
			log(fmt.Sprintf("%#v filled up the cave", g))
			break
		}
	}
}

func (c *Cave) minX() int {
	minX := 500
	for _, path := range c.rockPaths {
		for _, point := range path {
			if point[0] != math.MinInt && point[0] < minX {
				minX = point[0]
			}
		}
	}
	return minX - 20
}

func (c *Cave) maxY() int {
	maxY := 0
	for _, path := range c.rockPaths {
		for _, point := range path {
			if point[1] > maxY {
				maxY = point[1]
			}
		}
	}
	return maxY + 2
}

func (c *Cave) Draw() {
	tm.Clear()

	for _, path := range c.rockPaths {
		path.eachSegment(func(start, end Point, orientation Orientation) bool {
			if orientation == Vertical {
				for y := start[1]; y <= end[1]; y++ {
					tm.MoveCursor(start[0]-c.minX(), y)
					tm.Print("#")
				}
			}
			if orientation == Horizontal {
				startX := start[0]
				if startX == math.MinInt {
					startX = c.minX()
				}
				endX := end[0]
				if endX == math.MaxInt {
					endX = 1000
				}
				for x := startX; x <= endX; x++ {
					tm.MoveCursor(x-c.minX(), start[1])
					tm.Print("#")
				}
			}
			return true
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
	cave := Cave{rockPaths: paths, grains: set.NewSet[Point]()}
	maxY := cave.maxY()
	cave.rockPaths = append(cave.rockPaths,
		Path{Point{math.MinInt, maxY}, Point{math.MaxInt, maxY}},
	)
	return cave
}

func main() {
	flag.Parse()

	if *Prof != "" {
		f, err := os.Create(*Prof)
		util.HandleError(err)
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	cave := parseInput()
	cave.addSandUntilDone()

	fmt.Println("Part 2:", cave.grains.Len())
}
