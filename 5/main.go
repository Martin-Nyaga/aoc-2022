package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/martin-nyaga/aoc-2022/util"
)

type Stack []byte

func (s *Stack) Push(val ...byte) {
	*s = append(*s, val...)
}

func (s *Stack) Pop(n int) []byte {
	val := (*s)[len(*s)-n : len(*s)]
	*s = (*s)[:len(*s)-(n)]
	return val
}

type Move [3]int

func (m Move) Execute(stacks []Stack) {
	count := m[0]
	source := m[1]
	target := m[2]

	for i := 0; i < count; i++ {
		stacks[target].Push(stacks[source].Pop(1)...)
	}
}

func (m Move) ExecuteInOrder(stacks []Stack) {
	count := m[0]
	source := m[1]
	target := m[2]
	stacks[target].Push(stacks[source].Pop(count)...)
}

type Stacks []Stack

func (s Stacks) Tops() []byte {
	result := make([]byte, 0)
	for _, stack := range s {
		if len(stack) > 0 {
			result = append(result, stack[len(stack)-1])
		}
	}
	return result
}

func (s Stacks) Print() {
	for _, stack := range s {
		fmt.Println(string(stack))
	}
}

func parseInput() (Stacks, []Move) {
	str := util.NewInputFile("5").ReadToString()
	arr := strings.Split(str, "\n\n")
	rawStacks := strings.Split(arr[0], "\n")
	rawMoves := strings.Split(strings.TrimSpace(arr[1]), "\n")

	stacksCount := len(strings.Fields(rawStacks[len(rawStacks)-1]))
	stacks := make([]Stack, stacksCount, stacksCount)
	for i := len(rawStacks) - 2; i >= 0; i-- {
		slice := rawStacks[i]
		for j, ch := range []byte(slice) {
			if j%4 == 1 && ch != ' ' {
				stacks[j/4].Push(ch)
			}
		}
	}

	movesCount := len(rawMoves)
	moves := make([]Move, movesCount, movesCount)
	for i, line := range rawMoves {
		var count, source, target int
		_, err := fmt.Fscanf(strings.NewReader(line), "move %d from %d to %d", &count, &source, &target)
		util.HandleError(err)
		moves[i] = Move{count, source - 1, target - 1}
	}

	return stacks, moves
}

func main() {
	flag.Parse()
	stacks, moves := parseInput()
	for _, move := range moves {
		move.Execute(stacks)
	}

	fmt.Println("Part 1:", string(stacks.Tops()))

	stacks, moves = parseInput()
	for _, move := range moves {
		move.ExecuteInOrder(stacks)
	}

	fmt.Println("Part 2:", string(stacks.Tops()))
}
