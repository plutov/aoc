package main

import (
	"bufio"
	"fmt"
	"os"
)

var surface = [][]byte{}

const mins = 1000000000

// I got these numbers after some manual debugging. It counts, don't blame me!
const uniqueCount = 485
const step = 28

type p struct {
	x int
	y int
}

func main() {
	inputFile, _ := os.Open("./input")
	scanner := bufio.NewScanner(inputFile)

	var i int
	for scanner.Scan() {
		surface = append(surface, []byte(scanner.Text()))
		i++
	}

	uniqueResults := make(map[int]int)
	minToRes := make(map[int]int)

	for m := 0; m < mins; m++ {
		updates := make(map[p]byte)
		for x, row := range surface {
			for y := range row {
				if becomeTree(surface, x, y) {
					updates[p{x, y}] = '|'
				} else if becomeLumberyard(surface, x, y) {
					updates[p{x, y}] = '#'
				} else if becomeOpen(surface, x, y) {
					updates[p{x, y}] = '.'
				}
			}
		}

		for point, kind := range updates {
			surface[point.x][point.y] = kind
		}

		res := result()
		uniqueResults[res] = m
		minToRes[m] = res

		if len(uniqueResults) == uniqueCount {
			for j := mins - 1; j > 0; j -= step {
				answer2, ok := minToRes[j]
				if ok {
					fmt.Printf("Answer 2: %d\n", answer2)
					os.Exit(0)
				}
			}
		}
	}

	fmt.Printf("Answer 1: %d\n", result())
}

func becomeTree(s [][]byte, x, y int) bool {
	if s[x][y] != '.' {
		return false
	}

	trees := countAdjacent(s, x, y, '|')
	return trees >= 3
}

func becomeLumberyard(s [][]byte, x, y int) bool {
	if s[x][y] != '|' {
		return false
	}

	lumberyards := countAdjacent(s, x, y, '#')
	return lumberyards >= 3
}

func becomeOpen(s [][]byte, x, y int) bool {
	if s[x][y] != '#' {
		return false
	}

	lumberyards := countAdjacent(s, x, y, '#')
	trees := countAdjacent(s, x, y, '|')
	return lumberyards < 1 || trees < 1
}

func countAdjacent(s [][]byte, x, y int, kind byte) int {
	var points = []p{
		p{x - 1, y},
		p{x - 1, y + 1},
		p{x, y + 1},
		p{x + 1, y + 1},
		p{x + 1, y},
		p{x + 1, y - 1},
		p{x, y - 1},
		p{x - 1, y - 1},
	}

	var n int
	for _, p := range points {
		if p.x >= 0 && p.y >= 0 && p.x < len(s) && p.y < len(s[0]) && s[p.x][p.y] == kind {
			n++
		}
	}

	return n
}

func result() int {
	var trees, lumberyards int
	for _, row := range surface {
		for _, col := range row {
			if col == '|' {
				trees++
			} else if col == '#' {
				lumberyards++
			}
		}
	}

	return trees * lumberyards
}
