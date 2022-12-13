package main

import (
	"errors"

	"github.com/martin-nyaga/aoc-2022/util"
)

func parseInput() [][2]Packet {
	lines := util.NewInputFile("13").ReadLines()
	packetPairs := make([][2]Packet, 0)
	var packetPair [2]Packet
	i := 0
	for i < len(lines) {
		line := lines[i]
		if len(line) == 0 {
			packetPairs = append(packetPairs, packetPair)
		} else {
			packetPair[0] = parsePacket(line)
			i += 1
			packetPair[1] = parsePacket(line)
			i += 1
		}
	}
	return packetPairs
}
