package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

const (
	sand = iota
	water
	clay
	waterMove
)

const springX = 500
const springY = 0

type clayMap map[p]bool

type p struct {
	x, y int
}

var xMax, yMax int

var yMin = 1000000
var xMin = 1000000

func main() {
	inputFile, _ := os.Open("./input")
	scanner := bufio.NewScanner(inputFile)

	var cm = make(clayMap)

	for scanner.Scan() {
		line := scanner.Text()
		lineParts := strings.Split(line, ", ")

		// x can be on left or on right side of input line
		var xIndex, yIndex int
		if strings.Contains(lineParts[0], "x") {
			yIndex = 1
		} else {
			xIndex = 1
		}

		// x
		xs := strings.Split(strings.Replace(lineParts[xIndex], "x=", "", -1), "..")
		xFrom, _ := strconv.Atoi(xs[0])
		xTo := xFrom
		if len(xs) == 2 {
			xTo, _ = strconv.Atoi(xs[1])
		}
		if xTo > xMax {
			xMax = xTo
		}
		if xFrom < xMin {
			xMin = xFrom
		}

		// y
		ys := strings.Split(strings.Replace(lineParts[yIndex], "y=", "", -1), "..")
		yFrom, _ := strconv.Atoi(ys[0])
		yTo := yFrom
		if len(ys) == 2 {
			yTo, _ = strconv.Atoi(ys[1])
		}
		if yTo > yMax {
			yMax = yTo
		}
		if yFrom < yMin {
			yMin = yFrom
		}

		// clay map
		for y := yFrom; y <= yTo; y++ {
			for x := xFrom; x <= xTo; x++ {
				cm[p{x, y}] = true
			}
		}
	}

	leave := p{springX, springY}
	leaves := []p{leave}
	wet := map[p]bool{}

	for len(leaves) > 0 {
		last := len(leaves) - 1
		leave, leaves = leaves[last], leaves[:last]
		if cm[leave] {
			continue
		}

		f := leave
		for f.y+1 <= yMax && !cm[p{f.x, f.y + 1}] {
			f.y++
			wet[f] = true
		}

		if !cm[p{f.x, f.y + 1}] {
			continue
		}

		locked := true

		l := f
		for !cm[p{l.x - 1, l.y}] && cm[p{l.x - 1, l.y + 1}] {
			l.x--
			wet[l] = true
		}
		if !cm[p{l.x - 1, l.y}] {
			l.x--
			locked = false
			if !wet[l] {
				wet[l] = true
				leaves = append(leaves, leave, l)
			}
		}

		r := f
		for !cm[p{r.x + 1, r.y}] && cm[p{r.x + 1, r.y + 1}] {
			r.x++
			wet[r] = true
		}
		if !cm[p{r.x + 1, r.y}] {
			r.x++
			if locked {
				locked = false
				if !wet[r] {
					wet[r] = true
					leaves = append(leaves, leave, r)
				}
			} else {
				if !wet[r] {
					wet[r] = true
					leaves = append(leaves, r)
				}
			}
		}
		if locked {
			for x := l.x; x <= r.x; x++ {
				cm[p{x, f.y}] = true
			}
			leaves = append(leaves, leave)
		}
	}

	minX, maxX := math.MaxInt32, math.MinInt32
	for w := range wet {
		if w.x < minX {
			minX = w.x
		}
		if w.x > maxX {
			maxX = w.x
		}
	}

	var answer1, answer2 int
	for y := yMin; y <= yMax; y++ {
		for x := minX; x <= maxX; x++ {
			pp := p{x, y}
			if wet[pp] {
				answer1++
				if cm[pp] {
					answer2++
				}
			}
		}
	}

	fmt.Printf("Answer1 : %d\n", answer1)
	fmt.Printf("Answer2 : %d", answer2)
}
