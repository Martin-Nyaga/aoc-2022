package main

import (
	"flag"
	"fmt"

	"github.com/martin-nyaga/aoc-2022/utils"
)

func readInput() []byte {
	return utils.NewInputFile("2").ReadBytes()
}

func main() {
	flag.Parse()

	var z1, z2 int
	var t, m1, m2 byte
	for i, x := range readInput() {
		switch i % 4 {
		case 0:
			t = x - 65
		case 2:
			m1 = x - 88
			m2 = (m1 + 2) % 4
			if m2 == 0 {
				m2 += 1
			}
			m2 = (t + m2) % 3
		case 3:
			d1 := (int(t) - int(m1)%3 + 3) % 3
			d2 := (int(t) - int(m2)%3 + 3) % 3
			z1 += int(m1) + 1
			z2 += int(m2) + 1
			if d1 != 1 {
				z1 += 3 + (d1/2)*3
			}
			if d2 != 1 {
				z2 += 3 + (d2/2)*3
			}
		}
	}

	fmt.Println("Part 1:", z1)
	fmt.Println("Part 2:", z2)
}
