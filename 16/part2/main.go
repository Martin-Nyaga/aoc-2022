package main

import (
	"errors"
	"flag"
	"fmt"
	"math"
	"os"
	"runtime/pprof"
	"sort"
	"strconv"
	"strings"

	"github.com/martin-nyaga/aoc-2022/util"
	"github.com/martin-nyaga/aoc-2022/util/pqueue"
	"github.com/martin-nyaga/aoc-2022/util/set"
	"github.com/martin-nyaga/aoc-2022/util/slices"
)

type Valve struct {
	name        string
	rate        int
	connections []string
}

type State struct {
	actors                        [2]Actor
	openSet                       set.Set[string]
	valves                        map[string]*Valve
	valvesToOpen                  []string
	currentMinute                 int
	accumulatedRelease            int
	_key                          string
	_cumulativeReleasablePressure int
}

type Actor struct {
	name         string
	path         []string
	nextSteps    []string
	currentValve string
}

type Path struct {
	path  []string
	valve string
}

type PathCalculator struct {
	cache map[[2]string][]string
}

func (p *PathCalculator) pathBetween(s *State, a, b string) []string {
	pair := [2]string{a, b}
	if path, exists := p.cache[pair]; exists {
		return path
	}
	path := s.pathBetween(a, b)
	p.cache[pair] = path
	return path
}

func (p Path) Dist() int { return len(p.path) }

func (s *State) pathBetween(a, b string) []string {
	queue := make([]Path, 0, 10)
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

func (state *State) accumulateRelease(minutes int) int {
	accumulatedRelease := state.accumulatedRelease
	state.openSet.Each(func(valve string) {
		v := state.valves[valve]
		accumulatedRelease += v.rate * minutes
	})
	return accumulatedRelease
}

func (state *State) addNextSteps(actorIndex int, nextSteps []string) State {
	nextSteps = append(nextSteps, "open")
	actor := state.actors[actorIndex]
	nextActor := Actor{
		name:         actor.name,
		currentValve: actor.currentValve,
		path:         actor.path,
		nextSteps:    nextSteps,
	}
	nextActors := [2]Actor{}
	if actorIndex == 0 {
		nextActors[0] = nextActor
		nextActors[1] = state.actors[1]
	} else {
		nextActors[0] = state.actors[0]
		nextActors[1] = nextActor
	}

	return State{
		actors:             nextActors,
		openSet:            state.openSet,
		valves:             state.valves,
		valvesToOpen:       state.valvesToOpen,
		currentMinute:      state.currentMinute,
		accumulatedRelease: state.accumulatedRelease,
	}
}

func (state *State) simulateStepsForActorWithLeastStepsToTake() State {
	actorIndex := state.actorWithLeastStepsToTake()
	actor := state.actors[actorIndex]
	if len(actor.nextSteps) == 0 {
		panic("How did this happen?")
	}
	nextSet := set.NewSet(state.openSet.ToSlice()...)

	nextActor := Actor{
		name:         actor.name,
		path:         make([]string, 0, len(actor.path)),
		nextSteps:    make([]string, 0, len(actor.nextSteps)),
		currentValve: actor.currentValve,
	}
	nextActor.path = append(nextActor.path, actor.path...)
	nextActor.nextSteps = append(nextActor.nextSteps, actor.nextSteps...)

	otherActorIndex := (actorIndex + 1) % 2
	otherActor := state.actors[otherActorIndex]
	nextOtherActor := Actor{
		name:         otherActor.name,
		path:         make([]string, 0, len(otherActor.path)),
		nextSteps:    make([]string, 0, len(otherActor.nextSteps)),
		currentValve: otherActor.currentValve,
	}
	nextOtherActor.path = append(nextOtherActor.path, otherActor.path...)
	nextOtherActor.nextSteps = append(nextOtherActor.nextSteps, otherActor.nextSteps...)

	steps := len(nextActor.nextSteps)
	for len(nextActor.nextSteps) > 0 {
		nextStep, err := slices.Shift(&nextActor.nextSteps)
		util.HandleError(err)
		if nextStep == "open" {
			nextSet.Add(nextActor.currentValve)
		} else {
			nextActor.currentValve = nextStep
		}
		nextActor.path = append(nextActor.path, nextStep)

		nextOtherActorStep, err := slices.Shift(&nextOtherActor.nextSteps)
		if nextOtherActorStep == "open" {
			nextSet.Add(nextOtherActor.currentValve)
		} else {
			nextOtherActor.currentValve = nextOtherActorStep
		}
		nextOtherActor.path = append(nextOtherActor.path, nextOtherActorStep)
	}
	accumulatedRelease := state.accumulateRelease(steps)

	nextActors := [2]Actor{}
	if actorIndex == 0 {
		nextActors[0] = nextActor
		nextActors[1] = nextOtherActor
	} else {
		nextActors[0] = nextOtherActor
		nextActors[1] = nextActor
	}

	return State{
		actors:             nextActors,
		openSet:            nextSet,
		valves:             state.valves,
		valvesToOpen:       state.valvesToOpen,
		currentMinute:      state.currentMinute + steps,
		accumulatedRelease: accumulatedRelease,
	}
}

func (state *State) allActorsHaveStepsToTake() bool {
	allActorsHaveStepsToTake := true
	for _, actor := range state.actors {
		if len(actor.nextSteps) == 0 && !state.allValvesAreOpenOrAssigned() {
			allActorsHaveStepsToTake = false
			break
		}
	}
	return allActorsHaveStepsToTake
}

func (state *State) actorWithLeastStepsToTake() int {
	firstActorSteps := len(state.actors[0].nextSteps)
	secondActorSteps := len(state.actors[1].nextSteps)
	if firstActorSteps > 0 && secondActorSteps > 0 {
		if firstActorSteps < secondActorSteps {
			return 0
		} else {
			return 1
		}
	} else {
		if firstActorSteps > 0 {
			return 0
		}
		if secondActorSteps > 0 {
			return 1
		}
	}

	panic("How did I get here?")
}

func (state *State) DebugPrint() {
	if !*Debug {
		return
	}

	fmt.Println("---State---")
	fmt.Println("CurrentMinute", state.currentMinute)
	fmt.Println("accumulatedRelease", state.accumulatedRelease)
	fmt.Println("My valve", state.actors[0].currentValve)
	fmt.Println("My path", state.actors[0].path)
	fmt.Println("My next steps", state.actors[0].nextSteps)
	fmt.Println("elephant valve", state.actors[1].currentValve)
	fmt.Println("elephant path", state.actors[1].path)
	fmt.Println("elephant next steps", state.actors[1].nextSteps)
	opened := ""
	for _, valve := range state.valves {
		if state.openSet.Has(valve.name) {
			opened += valve.name + ", "
		}
	}
	fmt.Println("opened", opened)
	fmt.Println("------")
}

func (state *State) ForcePrint() {
	old := *Debug
	defer func() { *Debug = old }()
	*Debug = true
	state.DebugPrint()
}

func (state *State) allValvesAreOpen() bool {
	allOpen := true
	for _, valve := range state.valvesToOpen {
		if !state.openSet.Has(valve) {
			allOpen = false
			break
		}
	}
	return allOpen
}

func (state *State) valveIsOpenOrClaimed(valveName string) bool {
	if state.openSet.Has(valveName) {
		return true
	}
	someActorHasClaimed := false
	for _, actor := range state.actors {
		if len(actor.nextSteps) > 0 && ((actor.currentValve == valveName && actor.nextSteps[0] == "open") || (len(actor.nextSteps) > 1 && actor.nextSteps[len(actor.nextSteps)-2] == valveName)) {
			someActorHasClaimed = true
		}
	}
	return someActorHasClaimed
}

func (state *State) allValvesAreOpenOrAssigned() bool {
	allOpen := true
	for _, valve := range state.valves {
		if valve.rate > 0 && !state.valveIsOpenOrClaimed(valve.name) {
			allOpen = false
		}
	}
	return allOpen
}

func (state *State) bestNextStepsForActor(pathCalculator *PathCalculator, actor Actor) ([][]string, error) {
	Debugln("Searching for best next steps for", actor.name)
	state.DebugPrint()
	paths := make([][]string, 0)
	for _, valveName := range state.valvesToOpen {
		if state.valveIsOpenOrClaimed(valveName) {
			Debugln("skipped", valveName, "as it was open or claimed")
			continue
		}

		path := pathCalculator.pathBetween(state, state.valves[actor.currentValve].name, valveName)

		// Don't bother opening valves which are too far to make a difference
		if 26-state.currentMinute-len(path) < 2 {
			Debugln("skipped", valveName, "as", actor.name, "won't get there in time")
			continue
		}

		// If this path has a higher rate valve along the way, don't bother going
		// all the way
		releasable := (26 - state.currentMinute - len(path) - 1) * state.valves[valveName].rate
		hasBetterOption := false
		for i, waypoint := range path {
			if i == len(path)-1 {
				continue
			}
			if state.valveIsOpenOrClaimed(waypoint) {
				continue
			}
			waypointReleasable := (26 - state.currentMinute - i - 2) * state.valves[waypoint].rate
			if waypointReleasable >= releasable {
				hasBetterOption = true
				Debugln("skipped", valveName, "as", waypoint, "is a better option along the way")
				break
			}
		}

		// if this valve can be opened more efficiently by the other actor, don't
		// bother adding it
		var otherActor *Actor
		for _, a := range state.actors {
			if a.name != actor.name {
				otherActor = &a
				break
			}
		}
		var otherActorValve string
		if len(otherActor.nextSteps) > 1 {
			otherActorValve = otherActor.nextSteps[len(otherActor.nextSteps)-2]
		} else {
			otherActorValve = otherActor.currentValve
		}
		otherPath := pathCalculator.pathBetween(state, otherActorValve, valveName)
		if len(otherPath)+len(otherActor.nextSteps) < len(path) {
			Debugln("skipped", valveName, "as", otherActor.name, "can get there much quicker")
			hasBetterOption = true
		}

		if !hasBetterOption {
			paths = append(paths, path)
		}
	}

	if len(paths) == 0 {
		return nil, errors.New("no paths to go to")
	}

	return paths, nil
}

func (state *State) nextStates(pathCalculator *PathCalculator) []State {
	Debugln("At state")
	state.DebugPrint()

	next := make([]State, 0)

	if state.currentMinute == 26 {
		Debugln("Ran out of time")
		return next
	}

	// If all non zero valves are open, just complete the simulation and don't
	// bother moving any more
	if state.allValvesAreOpen() {
		Debugln("Already opened all valves")
		minutesLeft := 26 - state.currentMinute
		accumulatedRelease := state.accumulateRelease(minutesLeft)
		next = append(next, State{
			actors:             state.actors,
			openSet:            state.openSet,
			valves:             state.valves,
			valvesToOpen:       state.valvesToOpen,
			currentMinute:      26,
			accumulatedRelease: accumulatedRelease,
		})
		return next
	}

	// Let actors move until at least one has no more steps to take
	if state.allActorsHaveStepsToTake() {
		Debugln("Can take steps first")
		nextState := state.simulateStepsForActorWithLeastStepsToTake()
		Debugln("After taking required steps")
		nextState.DebugPrint()
		next = append(next, nextState)
		return next
	}

	Debugln("No steps for at least one actor, finding somewhere for them to go")

	// Give any available moves to any actors without moves
	for i, actor := range state.actors {
		// Don't scan reachable states for actors with moves
		if len(actor.nextSteps) > 0 {
			continue
		}

		Debugln("Scanning reachable paths for", actor.name)
		bestPaths, err := state.bestNextStepsForActor(pathCalculator, actor)
		if err != nil {
			Debugln(actor.name, "Couldn't go anywhere")
			// No where to go for this guy
			continue
		}
		for _, bestPath := range bestPaths {
			nextState := state.addNextSteps(i, bestPath)
			Debugln(actor.name, "can take path", bestPath)

			for j, otherActor := range state.actors {
				if j == i {
					continue
				}

				// Don't scan reachable states for actors with moves
				if len(otherActor.nextSteps) > 0 {
					next = append(next, nextState)
					continue
				}

				bestPaths, err := nextState.bestNextStepsForActor(pathCalculator, otherActor)
				if err != nil {
					Debugln(otherActor.name, "Couldn't go anywhere")
					// No where to go for this guy
					next = append(next, nextState)
					continue
				}

				for _, bestPath := range bestPaths {
					Debugln(otherActor.name, "can take path", bestPath)

					otherActorNextState := nextState.addNextSteps(j, bestPath)
					next = append(next, otherActorNextState)
				}
			}
		}
	}

	if len(next) > 0 {
		Debugln("Can open more")
		for _, n := range next {
			n.DebugPrint()
		}
		Debugln("============")
		return next
	}

	Debugln("Can't get anywhere quick enough")
	// I haven't opened all valves, but can't get anywhere in reasonable time, so
	// just simulate the remaining time
	minutesLeft := 26 - state.currentMinute
	accumulatedRelease := state.accumulateRelease(minutesLeft)
	next = append(next, State{
		actors:             state.actors,
		openSet:            state.openSet,
		valves:             state.valves,
		valvesToOpen:       state.valvesToOpen,
		currentMinute:      26,
		accumulatedRelease: accumulatedRelease,
	})

	return uniqueByKey(next)
}

func (state *State) cumulativeReleasablePressure() int {
	if state._cumulativeReleasablePressure != 0 {
		return state._cumulativeReleasablePressure
	}
	result := state.accumulatedRelease

	for _, valveName := range state.valvesToOpen {
		if state.openSet.Has(valveName) {
			result += (26 - state.currentMinute) * state.valves[valveName].rate
		}
	}

	for _, actor := range state.actors {
		if len(actor.nextSteps) > 0 {
			var actorValveName string
			if len(actor.nextSteps) > 1 {
				actorValveName = actor.nextSteps[len(actor.nextSteps)-2]
			} else {
				actorValveName = actor.currentValve
			}
			actorReleasable := (26 - state.currentMinute - len(actor.nextSteps)) * state.valves[actorValveName].rate
			result += actorReleasable
		}
	}

	state._cumulativeReleasablePressure = result
	return result
}

func uniqueByKey(states []State) []State {
	set := map[string]*State{}
	for _, state := range states {
		set[state.key()] = &state
	}
	result := make([]State, 0)
	for _, state := range set {
		result = append(result, *state)
	}
	return result
}

func (state *State) key() string {
	if len(state._key) > 0 {
		return state._key
	}
	actor1 := strings.Join(state.actors[0].path, "") + "|" + strings.Join(state.actors[0].nextSteps, "")
	actor2 := strings.Join(state.actors[1].path, "") + "|" + strings.Join(state.actors[1].nextSteps, "")
	arr := []string{actor1, actor2}
	sort.Strings([]string(arr))
	key := strconv.Itoa(state.accumulatedRelease) + ";" + strconv.Itoa(state.currentMinute) + ";" + strings.Join(arr, ",")
	state._key = key
	return key
}

func parseInput() (map[string]*Valve, []string) {
	file := util.NewInputFile("16")

	valves := make(map[string]*Valve)
	valvesToOpen := make([]string, 0)
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
		if intRate > 0 {
			valvesToOpen = append(valvesToOpen, valveName)
		}
	}

	return valves, valvesToOpen
}

var MaxIter = flag.Int("maxiter", math.MaxInt, "Maximum number of iterations")
var Prof = flag.String("prof", "", "Generate cpu profile")
var Debug = flag.Bool("debug", false, "Print debug output")

func Debugln(args ...interface{}) {
	if *Debug {
		fmt.Println(args...)
	}
}

func main() {
	flag.Parse()

	if *Prof != "" {
		f, err := os.Create(*Prof)
		util.HandleError(err)
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	valves, valvesToOpen := parseInput()
	Debugln("Ordered valves", valvesToOpen)
	state := State{actors: [2]Actor{
		{name: "me", currentValve: "AA", path: []string{}},
		{name: "elephant", currentValve: "AA", path: []string{}},
	}, valves: valves, valvesToOpen: valvesToOpen}
	queue := pqueue.NewPqueue[int, State](pqueue.MaxQueue)
	queue.Push(0, state)
	visited := set.NewSet[string]()

	bestCumulativeReleasablePressure := math.MinInt
	var bestState *State
	pathCalculator := PathCalculator{cache: make(map[[2]string][]string)}
	i := 0
	for !queue.Empty() && i < *MaxIter {
		i += 1
		state, err := queue.Pop()
		util.HandleError(err)

		Debugln(state.key())
		if visited.Has(state.key()) {
			Debugln("Visited!")
			continue
		}
		visited.Add(state.key())

		if i%10000 == 0 {
			fmt.Println("Current state minute:", state.currentMinute)
			fmt.Println("Visited states:", len(visited))
			fmt.Println("Queue size:", queue.Len())
			if bestState != nil {
				fmt.Println("Best so far:", bestState.accumulatedRelease)
			}
			fmt.Println("Theoretical best:", bestCumulativeReleasablePressure)
			fmt.Println()
		}
		if state.currentMinute == 26 && (bestState == nil || (state.accumulatedRelease > bestState.accumulatedRelease)) {
			bestState = &state

			fmt.Println("new best:", bestState.accumulatedRelease)
			fmt.Println()

			if *Debug {
				Debugln("minute:", state.currentMinute)
				Debugln("my path:", state.actors[0].path)
				Debugln("elephant path:", state.actors[1].path)
				Debugln("released:", state.accumulatedRelease)
				opened := ""
				for _, valve := range state.valves {
					if state.openSet.Has(valve.name) {
						opened += valve.name + ", "
					}
				}
				Debugln("opened:", opened)
				Debugln("---")
			}

			continue
		}

		nextStates := state.nextStates(&pathCalculator)
		currentReleasable := state.cumulativeReleasablePressure()
		for _, nextState := range nextStates {
			releasable := nextState.cumulativeReleasablePressure()
			Debugln("Releasable:", releasable)

			if len(nextState.actors[0].nextSteps) > 0 && len(nextState.actors[1].nextSteps) > 0 {
				if releasable > currentReleasable {
					queue.Push(releasable, nextState)
				}
			} else {
				queue.Push(releasable, nextState)
			}

			if releasable > bestCumulativeReleasablePressure {
				bestCumulativeReleasablePressure = releasable
			}
		}
	}

	if bestState != nil {
		fmt.Println("Total visited", len(visited))
		fmt.Println("Part 2:", bestState.accumulatedRelease)
		bestState.ForcePrint()
	} else {
		fmt.Println("Couldn't find it, sorry")
	}
}
