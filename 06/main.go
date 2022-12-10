package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/martin-nyaga/aoc-2022/util"
	"github.com/martin-nyaga/aoc-2022/util/slices"
)

type Window struct {
	size       int
	charCounts map[byte]int
	slice      []byte
}

func newWindow(bytes []byte) Window {
	newBytes := make([]byte, len(bytes))
	copy(newBytes, bytes)
	charCounts := map[byte]int{}
	window := Window{size: len(bytes), charCounts: charCounts, slice: newBytes}
	for _, ch := range bytes {
		window.incChar(ch)
	}
	return window
}

func (w *Window) IsUnique() bool {
	return len(w.charCounts) == w.size
}

func (w *Window) Add(ch byte) {
	w.ShiftChar()
	w.incChar(ch)
	w.slice = append(w.slice, ch)
}

func (w *Window) ShiftChar() {
	ch, err := slices.Shift(&w.slice)
	util.HandleError(err)
	w.decChar(ch)
}

func (w *Window) incChar(ch byte) {
	if _, exists := w.charCounts[ch]; exists {
		w.charCounts[ch] += 1
	} else {
		w.charCounts[ch] = 1
	}
}

func (w *Window) decChar(ch byte) {
	w.charCounts[ch] -= 1
	if w.charCounts[ch] == 0 {
		delete(w.charCounts, ch)
	}
}

func parseInput() []byte {
	return util.NewInputFile("6").ReadBytes()
}

func main() {
	flag.Parse()

	bytes := []byte(strings.TrimSpace(string(parseInput())))

	window := newWindow(bytes[0:4])
	var packetStartIndex int
	for i := 4; i < len(bytes); i++ {
		packetStartIndex = i
		window.Add(bytes[i])
		if window.IsUnique() {
			break
		}
	}

	window = newWindow(bytes[0:14])
	var messageStartIndex int
	for i := 4; i < len(bytes); i++ {
		messageStartIndex = i
		window.Add(bytes[i])
		if window.IsUnique() {
			break
		}
	}

	fmt.Println("Part 1:", packetStartIndex+1)
	fmt.Println("Part 2:", messageStartIndex+1)
}
