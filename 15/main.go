package main

import (
	"flag"
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/martin-nyaga/aoc-2022/util"
)

type Point [2]int
type Area [4]Point
type Line [2]Point

type SensorBeaconPair struct {
	sensor Point
	beacon Point
	area   *Area
}

func abs(x int) int {
	return int(math.Abs(float64(x)))
}

func orientation(a, b, c Point) int {
	b0 := Point{b[0] - a[0], b[1] - a[1]}
	c0 := Point{c[0] - a[0], c[1] - a[1]}
	return int(math.Copysign(1, float64(b0[0]*c0[0]+b0[1]*c0[1])))
}

func (a Area) includes(p Point) bool {
	o1 := orientation(a[0], a[1], p)
	o2 := orientation(a[1], a[2], p)
	o3 := orientation(a[2], a[3], p)
	o4 := orientation(a[3], a[0], p)

	return o1 == o2 && o2 == o3 && o3 == o4
}

func (a Area) size() int {
	xSpan := (a[1][0] - a[0][0]) + 1
	ySpan := (a[2][1] - a[0][1]) + 1
	return xSpan * ySpan
}

func (sb *SensorBeaconPair) Area() Area {
	if sb.area != nil {
		return *sb.area
	}
	dbx := int(math.Abs(float64(sb.sensor[0] - sb.beacon[0])))
	dby := int(math.Abs(float64(sb.sensor[1] - sb.beacon[1])))
	left := Point{sb.sensor[0] - (dby + dbx), sb.sensor[1]}
	right := Point{sb.sensor[0] + (dby + dbx), sb.sensor[1]}
	top := Point{sb.sensor[0], sb.sensor[1] - (dbx + dby)}
	bottom := Point{sb.sensor[0], sb.sensor[1] + (dbx + dby)}

	sb.area = &Area{top, right, bottom, left}

	return *sb.area
}

func (sb *SensorBeaconPair) fringeAreas() []Area {
	areas := make([]Area, 0)
	xDiff := abs(sb.Area()[3][0] - sb.sensor[0])
	yDiff := abs(sb.Area()[0][1] - sb.sensor[1])

	topLeftX := sb.Area()[3][0]
	topLeftY := sb.Area()[0][1]

	for j := 0; j < 4; j++ {
		for i := 0; i < 4; i++ {
			if j > 0 && j < 3 && i > 0 && i < 3 {
				continue
			}
			topLeft := Point{topLeftX + i*xDiff, topLeftY + j*yDiff}
			topRight := Point{topLeft[0] + xDiff, topLeft[1]}
			bottomRight := Point{topRight[0], topRight[1] + yDiff}
			bottomLeft := Point{topLeft[0], topLeft[1] + yDiff}
			area := Area{topLeft, topRight, bottomRight, bottomLeft}
			areas = append(areas, area)
		}
	}

	return areas
}

func (sb *SensorBeaconPair) scansArea(area Area) bool {
	return sb.scansPoint(area[0]) &&
		sb.scansPoint(area[1]) &&
		sb.scansPoint(area[2]) &&
		sb.scansPoint(area[3])
}

func (sb *SensorBeaconPair) scansPoint(point Point) bool {
	return sb.Area().includes(point)
}

func parseInput() []SensorBeaconPair {
	sensorBeaconPairs := make([]SensorBeaconPair, 0)
	file := util.NewInputFile("15")
	for _, line := range file.ReadLines() {
		var sx, sy, bx, by int
		var err error
		arr := strings.Split(line, " ")
		sx, err = strconv.Atoi(arr[2][2 : len(arr[2])-1])
		util.HandleError(err)
		sy, err = strconv.Atoi(arr[3][2 : len(arr[3])-1])
		util.HandleError(err)
		bx, err = strconv.Atoi(arr[8][2 : len(arr[8])-1])
		util.HandleError(err)
		by, err = strconv.Atoi(arr[9][2:len(arr[9])])
		util.HandleError(err)
		sensorBeaconPairs = append(sensorBeaconPairs, SensorBeaconPair{
			sensor: Point{sx, sy},
			beacon: Point{bx, by},
		})
	}

	return sensorBeaconPairs
}

func part1(sensorBeaconPairs []SensorBeaconPair) {
	minX := math.MaxInt
	maxX := math.MinInt
	for _, sb := range sensorBeaconPairs {
		if sb.Area()[3][0] < minX {
			minX = sb.Area()[3][0]
		}
		if sb.Area()[1][0] > maxX {
			maxX = sb.Area()[1][0]
		}
	}

	var targetY int
	if *util.UseSampleInput {
		targetY = 10
	} else {
		targetY = 2000000
	}
	result := 0
	for x := minX; x <= maxX; x++ {
		point := Point{x, targetY}
		cantContain := false
		// Check the point is scanned by at least one sensor
		for _, sb := range sensorBeaconPairs {
			if sb.scansPoint(point) {
				cantContain = true
				break
			}
		}

		// Make sure there's no actual beacons
		for _, sb := range sensorBeaconPairs {
			if sb.beacon == point {
				cantContain = false
				break
			}
		}

		if cantContain {
			result += 1
		}
	}
	fmt.Println("Part 1:", result)
}

func filterScannedOrOutOfBoundsAreas(sensorBeaconPairs *[]SensorBeaconPair, areasToScan *[]Area, minCoordinate, maxCoordinate int) []Area {
	filteredAreasToScan := make([]Area, 0)
	for _, area := range *areasToScan {
		if area[1][0] < minCoordinate ||
			area[0][0] > maxCoordinate ||
			area[2][1] < minCoordinate ||
			area[0][1] > maxCoordinate {
			continue
		}

		areaIsScanned := false
		for _, sb := range *sensorBeaconPairs {
			if sb.scansArea(area) {
				areaIsScanned = true
				break
			}
		}

		if !areaIsScanned {
			filteredAreasToScan = append(filteredAreasToScan, area)
		}
	}
	return filteredAreasToScan
}

func splitAreasToScan(areas *[]Area) []Area {
	nextAreas := make([]Area, 0)
	for _, area := range *areas {
		topLeft := area[0]
		topRight := area[1]
		bottomRight := area[2]
		bottomLeft := area[3]
		absXDiff := abs(topLeft[0] - bottomRight[0])
		xDiff := absXDiff / 2
		absYDiff := abs(topLeft[1] - bottomRight[1])
		yDiff := absYDiff / 2
		topMiddleA := Point{topLeft[0] + xDiff, topLeft[1]}
		topMiddleB := Point{topLeft[0] + xDiff + 1, topLeft[1]}
		leftMiddleA := Point{topLeft[0], topLeft[1] + yDiff}
		leftMiddleB := Point{topLeft[0], topLeft[1] + yDiff + 1}
		rightMiddleA := Point{topRight[0], topRight[1] + yDiff}
		rightMiddleB := Point{topRight[0], topRight[1] + yDiff + 1}
		bottomMiddleA := Point{topLeft[0] + xDiff, bottomLeft[1]}
		bottomMiddleB := Point{topLeft[0] + xDiff + 1, bottomLeft[1]}
		middleMiddleA := Point{topLeft[0] + xDiff, topLeft[1] + yDiff}
		middleMiddleB := Point{topLeft[0] + xDiff + 1, topLeft[1] + yDiff}
		middleMiddleC := Point{topLeft[0] + xDiff + 1, topLeft[1] + yDiff + 1}
		middleMiddleD := Point{topLeft[0] + xDiff, topLeft[1] + yDiff + 1}

		if xDiff != 0 && yDiff != 0 {
			q1 := Area{topLeft, topMiddleA, middleMiddleA, leftMiddleA}
			q2 := Area{topMiddleB, topRight, rightMiddleA, middleMiddleB}
			q3 := Area{middleMiddleC, rightMiddleB, bottomRight, bottomMiddleB}
			q4 := Area{leftMiddleB, middleMiddleD, bottomMiddleA, bottomLeft}
			nextAreas = append(nextAreas, q1, q2, q3, q4)
			continue
		}
		if xDiff == 0 && yDiff != 0 {
			top := Area{topLeft, topLeft, leftMiddleA, leftMiddleA}
			bottom := Area{leftMiddleB, leftMiddleB, bottomLeft, bottomLeft}
			nextAreas = append(nextAreas, top, bottom)
			continue
		}
		if xDiff != 0 && yDiff == 0 {
			left := Area{topLeft, topMiddleA, topMiddleA, topLeft}
			right := Area{topMiddleB, topRight, topRight, topMiddleB}
			nextAreas = append(nextAreas, left, right)
			continue
		}
		nextAreas = append(nextAreas, area)
	}
	return nextAreas
}

func maxSize(areas *[]Area) int {
	maxSize := math.MinInt
	for _, area := range *areas {
		size := area.size()
		if size > maxSize {
			maxSize = size
		}
	}
	return maxSize
}

func part2(sensorBeaconPairs []SensorBeaconPair) {
	var minCoordinate = 0
	var maxCoordinate int
	if *util.UseSampleInput {
		maxCoordinate = 20
	} else {
		maxCoordinate = 4000000
	}

	areasToScan := make([]Area, 0)
	for _, sb := range sensorBeaconPairs {
		areasToScan = append(areasToScan, sb.fringeAreas()...)
	}

	fmt.Println("Max size", maxSize(&areasToScan))
	prevLen := 0
	for len(areasToScan) != prevLen {
		prevLen = len(areasToScan)
		fmt.Println("Before filtering", len(areasToScan))
		filteredAreasToScan := filterScannedOrOutOfBoundsAreas(&sensorBeaconPairs, &areasToScan, minCoordinate, maxCoordinate)
		fmt.Println("After filtering", len(filteredAreasToScan))
		areasToScan = splitAreasToScan(&filteredAreasToScan)
		fmt.Println("After splitting", len(areasToScan))
		fmt.Println("Max size", maxSize(&areasToScan))
		fmt.Println("-")
	}

	fmt.Println("Split settled, scanning points in", len(areasToScan), "areas")

	var distressPoint *Point

outer:
	for _, area := range areasToScan {
		for x := area[0][0]; x <= area[2][0]; x++ {
			for y := area[0][1]; y <= area[2][1]; y++ {
				point := Point{x, y}
				scanned := false
				for _, sb := range sensorBeaconPairs {
					if sb.scansPoint(point) {
						scanned = true
						break
					}
				}

				if !scanned {
					distressPoint = &point
					break outer
				}
			}
		}
	}

	if distressPoint == nil {
		fmt.Println("Didn't find it, sorry")
		return
	}

	fmt.Println(distressPoint)
	fmt.Println("Part 2:", (distressPoint[0]*4000000)+distressPoint[1])
}

var runPart1 = flag.Bool("part1", true, "Run part 1")
var runPart2 = flag.Bool("part2", false, "Run part 2")

func main() {
	flag.Parse()
	sensorBeaconPairs := parseInput()

	if *runPart2 {
		part2(sensorBeaconPairs)
	} else {
		part1(sensorBeaconPairs)
	}
}
