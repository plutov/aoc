package main

import (
	"bufio"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
)

type sample struct {
	before []int
	op     []int
	after  []int
}

type op struct {
	N int
	F func(op, before, after []int) []int
}

var ops = map[string]op{
	"addr": op{
		F: func(op, before, after []int) []int {
			after[op[3]] = before[op[1]] + before[op[2]]
			return after
		},
	},
	"addi": op{
		F: func(op, before, after []int) []int {
			after[op[3]] = before[op[1]] + op[2]
			return after
		},
	},
	"mulr": op{
		F: func(op, before, after []int) []int {
			after[op[3]] = before[op[1]] * before[op[2]]
			return after
		},
	},
	"muli": op{
		F: func(op, before, after []int) []int {
			after[op[3]] = before[op[1]] * op[2]
			return after
		},
	},
	"banr": op{
		F: func(op, before, after []int) []int {
			after[op[3]] = before[op[1]] & before[op[2]]
			return after
		},
	},
	"bani": op{
		F: func(op, before, after []int) []int {
			after[op[3]] = before[op[1]] & op[2]
			return after
		},
	},
	"borr": op{
		F: func(op, before, after []int) []int {
			after[op[3]] = before[op[1]] | before[op[2]]
			return after
		},
	},
	"bori": op{
		F: func(op, before, after []int) []int {
			after[op[3]] = before[op[1]] | op[2]
			return after
		},
	},
	"setr": op{
		F: func(op, before, after []int) []int {
			after[op[3]] = before[op[1]]
			return after
		},
	},
	"seti": op{
		F: func(op, before, after []int) []int {
			after[op[3]] = op[1]
			return after
		},
	},
	"gtir": op{
		F: func(op, before, after []int) []int {
			if op[1] > before[op[2]] {
				after[op[3]] = 1
			} else {
				after[op[3]] = 0
			}
			return after
		},
	},
	"gtri": op{
		F: func(op, before, after []int) []int {
			if before[op[1]] > op[2] {
				after[op[3]] = 1
			} else {
				after[op[3]] = 0
			}
			return after
		},
	},
	"gtrr": op{
		F: func(op, before, after []int) []int {
			if before[op[1]] > before[op[2]] {
				after[op[3]] = 1
			} else {
				after[op[3]] = 0
			}
			return after
		},
	},
	"eqir": op{
		F: func(op, before, after []int) []int {
			if op[1] == before[op[2]] {
				after[op[3]] = 1
			} else {
				after[op[3]] = 0
			}
			return after
		},
	},
	"eqri": op{
		F: func(op, before, after []int) []int {
			if before[op[1]] == op[2] {
				after[op[3]] = 1
			} else {
				after[op[3]] = 0
			}
			return after
		},
	},
	"eqrr": op{
		F: func(op, before, after []int) []int {
			if before[op[1]] == before[op[2]] {
				after[op[3]] = 1
			} else {
				after[op[3]] = 0
			}
			return after
		},
	},
}

var opNumbers map[string][]int

func main() {
	answer1 := getAnswerOne()
	fmt.Printf("Answer 1: %d\n", answer1)

	answer2 := getAnswerTwo()
	fmt.Printf("Answer 2: %d\n", answer2)
}

func getAnswerOne() int {
	var samples []sample

	samplesFile, _ := os.Open("./samples")
	scanner := bufio.NewScanner(samplesFile)
	var s sample
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "Before:") {
			s = sample{
				before: lineToIntSlice("Before", line, ", "),
			}
		} else if strings.Contains(line, "After:") {
			s.after = lineToIntSlice("After", line, ", ")
		} else if len(line) != 0 {
			s.op = lineToIntSlice("", line, " ")
		} else {
			samples = append(samples, s)
		}
	}

	var moreThanThree int
	opNumbers = make(map[string][]int)
	for _, s := range samples {
		var successfull int
		for opName, op := range ops {
			var after []int
			for _, b := range s.before {
				after = append(after, b)
			}

			op.F(s.op, s.before, after)
			if reflect.DeepEqual(after, s.after) {
				successfull++

				_, ok := opNumbers[opName]
				if !ok {
					opNumbers[opName] = []int{}
				}

				if !intInSlice(opNumbers[opName], s.op[0]) {
					opNumbers[opName] = append(opNumbers[opName], s.op[0])
				}
			}
		}

		if successfull >= 3 {
			moreThanThree++
		}
	}

	return moreThanThree
}

func getAnswerTwo() int {
	var foundOps []int
	for len(foundOps) < len(ops) {
		for name, o := range opNumbers {
			var uniqueN []int
			for _, n := range o {
				if !intInSlice(foundOps, n) {
					uniqueN = append(uniqueN, n)
				}
			}

			if len(uniqueN) == 1 {
				temp := ops[name]
				temp.N = uniqueN[0]
				ops[name] = temp

				foundOps = append(foundOps, uniqueN[0])
			}
		}
	}

	start := []int{0, 0, 0, 0}
	program, _ := os.Open("program")
	scanner := bufio.NewScanner(program)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) != 0 {
			op := lineToIntSlice("", line, " ")
			for _, o := range ops {
				if o.N == op[0] {
					var after []int
					for _, b := range start {
						after = append(after, b)
					}

					start = o.F(op, start, after)
				}
			}
		}
	}

	return start[0]
}
func intInSlice(s []int, n int) bool {
	for _, i := range s {
		if i == n {
			return true
		}
	}

	return false
}

func lineToIntSlice(prefix, line, del string) (res []int) {
	numbersStr := strings.Replace(line, prefix+":", "", -1)
	numbersStr = strings.Replace(numbersStr, "]", "", -1)
	numbersStr = strings.Replace(numbersStr, "[", "", -1)
	numbersStr = strings.TrimSpace(numbersStr)

	numbers := strings.Split(numbersStr, del)
	for _, n := range numbers {
		nInt, _ := strconv.Atoi(n)
		res = append(res, nInt)
	}

	return
}
