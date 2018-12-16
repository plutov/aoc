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
	Name string
	N    int
	F    func(op, before, after []int) []int
}

var ops = map[string]op{
	"addr": op{
		Name: "addr",
		F: func(op, before, after []int) []int {
			after[op[3]] = before[op[1]] + before[op[2]]
			return after
		},
	},
	"addi": op{
		Name: "addi",
		F: func(op, before, after []int) []int {
			after[op[3]] = before[op[1]] + op[2]
			return after
		},
	},
	"mulr": op{
		Name: "mulr",
		F: func(op, before, after []int) []int {
			after[op[3]] = before[op[1]] * before[op[2]]
			return after
		},
	},
	"muli": op{
		Name: "muli",
		F: func(op, before, after []int) []int {
			after[op[3]] = before[op[1]] * op[2]
			return after
		},
	},
	"banr": op{
		Name: "banr",
		F: func(op, before, after []int) []int {
			after[op[3]] = before[op[1]] & before[op[2]]
			return after
		},
	},
	"bani": op{
		Name: "bani",
		F: func(op, before, after []int) []int {
			after[op[3]] = before[op[1]] & op[2]
			return after
		},
	},
	"borr": op{
		Name: "borr",
		F: func(op, before, after []int) []int {
			after[op[3]] = before[op[1]] | before[op[2]]
			return after
		},
	},
	"bori": op{
		Name: "bori",
		F: func(op, before, after []int) []int {
			after[op[3]] = before[op[1]] | op[2]
			return after
		},
	},
	"setr": op{
		Name: "setr",
		F: func(op, before, after []int) []int {
			after[op[3]] = before[op[1]]
			return after
		},
	},
	"seti": op{
		Name: "seti",
		F: func(op, before, after []int) []int {
			after[op[3]] = op[1]
			return after
		},
	},
	"gtir": op{
		Name: "gtir",
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
		Name: "gtri",
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
		Name: "gtrr",
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
		Name: "eqir",
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
		Name: "eqri",
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
		Name: "eqrr",
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

func main() {
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
	opAnNumber := make(map[string][]int)
	for _, s := range samples {
		var successfull int
		for _, op := range ops {
			var after []int
			for _, b := range s.before {
				after = append(after, b)
			}

			op.F(s.op, s.before, after)
			if reflect.DeepEqual(after, s.after) {
				successfull++

				_, ok := opAnNumber[op.Name]
				if !ok {
					opAnNumber[op.Name] = []int{}
				}

				if !intInSlice(opAnNumber[op.Name], s.op[0]) {
					opAnNumber[op.Name] = append(opAnNumber[op.Name], s.op[0])
				}
			}
		}

		if successfull >= 3 {
			moreThanThree++
		}
	}

	fmt.Printf("Answer 1: %d\n", moreThanThree)

	var foundOps []int
	for len(foundOps) < len(ops) {
		for name, o := range opAnNumber {
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
	scanner2 := bufio.NewScanner(program)
	for scanner2.Scan() {
		line := scanner2.Text()
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

	fmt.Printf("Answer 2: %d\n", start[0])
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
