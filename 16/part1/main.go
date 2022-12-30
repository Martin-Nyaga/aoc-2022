package main

import (
	"flag"
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/martin-nyaga/aoc-2022/util"
	"github.com/martin-nyaga/aoc-2022/util/set"
	"github.com/martin-nyaga/aoc-2022/util/slices"
)

type Valve struct {
	name        string
	rate        int
	connections []string
}

type State struct {
	currentValve       string
	openSet            set.Set[string]
	valves             map[string]*Valve
	currentMinute      int
	accumulatedRelease int
	path               []string
}

type Path struct {
	path  []string
	valve string
}

func (p Path) Dist() int { return len(p.path) }

func (s *State) pathBetween(a, b string) []string {
	queue := make([]Path, 0)
	queue = append(queue, Path{make([]string, 0), a})

	visited := set.NewSet[string]()
	var bestPath []string
	minDist := math.MaxInt

	for len(queue) > 0 {
		curr, err := slices.Shift(&queue)
		util.HandleError(err)

		if visited.Has(curr.valve) {
			continue
		}
		visited.Add(curr.valve)

		if curr.valve == b && curr.Dist() < minDist {
			minDist = curr.Dist()
			bestPath = curr.path
		}

		for _, next := range s.valves[curr.valve].connections {
			nextPath := make([]string, 0)
			nextPath = append(nextPath, curr.path...)
			nextPath = append(nextPath, next)
			queue = append(queue, Path{nextPath, s.valves[next].name})
		}
	}

	return bestPath
}

func (state *State) releasePressure() {
	for _, valve := range state.valves {
		if state.openSet.Has(valve.name) {
			state.accumulatedRelease += valve.rate
		}
	}
}

func (state *State) nextStates() []State {
	next := make([]State, 0)

	if state.currentMinute == 30 {
		return next
	}

	// If all non zero valves are open, just complete the simulation
	allOpen := true
	for _, valve := range state.valves {
		if valve.rate > 0 && !state.openSet.Has(valve.name) {
			allOpen = false
		}
	}
	if allOpen {
		minutesLeft := 30 - state.currentMinute
		accumulatedRelease := state.accumulatedRelease
		state.openSet.Each(func(valve string) {
			v := state.valves[valve]
			accumulatedRelease += v.rate * minutesLeft
		})
		next = append(next, State{
			currentValve:       state.currentValve,
			openSet:            state.openSet,
			valves:             state.valves,
			currentMinute:      30,
			accumulatedRelease: accumulatedRelease,
			path:               state.path,
		})
		return next
	}

	// Try get to and open all valves I haven't opened
	for _, valve := range state.valves {
		if state.openSet.Has(valve.name) {
			continue
		}

		// Don't bother opening jammed valves
		if valve.rate == 0 {
			continue
		}

		path := state.pathBetween(state.valves[state.currentValve].name, valve.name)
		// Don't bother opening valves which are too far to make a difference
		if 30-state.currentMinute-len(path) < 2 {
			continue
		}

		nextSet := set.NewSet(state.openSet.ToSlice()...)
		nextSet.Add(valve.name)
		accumulatedRelease := state.accumulatedRelease
		state.openSet.Each(func(valve string) {
			v := state.valves[valve]
			accumulatedRelease += v.rate * (len(path) + 1)
		})
		nextPath := make([]string, 0)
		nextPath = append(nextPath, state.path...)
		nextPath = append(nextPath, path...)
		next = append(next, State{
			currentValve:       valve.name,
			openSet:            nextSet,
			valves:             state.valves,
			currentMinute:      state.currentMinute + len(path) + 1,
			accumulatedRelease: accumulatedRelease,
			path:               nextPath,
		})
	}
	if len(next) > 0 {
		return next
	}

	// I haven't opened all valves, but can't get anywhere in reasonable time, so
	// just simulate the remaining time
	minutesLeft := 30 - state.currentMinute
	accumulatedRelease := state.accumulatedRelease
	state.openSet.Each(func(valve string) {
		v := state.valves[valve]
		accumulatedRelease += v.rate * minutesLeft
	})
	next = append(next, State{
		currentValve:       state.currentValve,
		openSet:            state.openSet,
		valves:             state.valves,
		currentMinute:      30,
		accumulatedRelease: accumulatedRelease,
		path:               state.path,
	})
	return next
}

func parseInput() map[string]*Valve {
	file := util.NewInputFile("16")

	valves := make(map[string]*Valve)
	for _, line := range file.ReadLines() {
		var valveName, rateStr string
		fmt.Sscanf(line, "Valve %s has flow %s tunnels lead to valves", &valveName, &rateStr)
		rateStr = strings.Split(rateStr, "=")[1]
		intRate, err := strconv.Atoi(rateStr[:len(rateStr)-1])
		util.HandleError(err)
		connectedValves := strings.Split(strings.TrimSpace(strings.Split(line, "valve")[1][1:]), ", ")
		valves[valveName] = &Valve{
			name:        valveName,
			rate:        intRate,
			connections: connectedValves,
		}
	}

	return valves
}

func main() {
	flag.Parse()

	valves := parseInput()
	state := State{currentValve: "AA", valves: valves}
	queue := make([]State, 0)
	queue = append(queue, state)

	var bestState *State
	for len(queue) > 0 {
		state, err := slices.Shift(&queue)
		util.HandleError(err)
		fmt.Println(state.currentMinute)
		if state.currentMinute == 30 && (bestState == nil || (state.accumulatedRelease > bestState.accumulatedRelease)) {
			bestState = &state
			fmt.Println("new best")
			fmt.Println("minute:", state.currentMinute)
			fmt.Println("path:", state.path)
			fmt.Println("released:", state.accumulatedRelease)
			opened := ""
			for _, valve := range state.valves {
				if state.openSet.Has(valve.name) {
					opened += valve.name + ", "
				}
			}
			fmt.Println("opened:", opened)
			fmt.Println("---")

			continue
		}

		nextStates := state.nextStates()
		queue = append(queue, nextStates...)
	}

	fmt.Println("Part 1:", bestState.accumulatedRelease)
}
