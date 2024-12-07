package main

import (
	_ "embed"
	"fmt"
	"strconv"
	"strings"
)

//go:embed example.txt
var example string

//go:embed input.txt
var input string

func main() {
	star1(example)
	star1(input)
	star2(example)
	star2(input)
}

func parseInput(input string) []Equation {
	eqStrs := strings.Split(input, "\n")
	eqs := make([]Equation, 0, len(eqStrs))
	for _, eqStr := range eqStrs {
		eqs = append(eqs, NewEquation(eqStr))
	}
	return eqs
}

func star1(input string) {
	eqs := parseInput(input)
	acc := 0
	for _, eq := range eqs {
		if star1IsValid(eq) {
			acc += eq.Val
		}
	}
	fmt.Println(acc)
}

func star1IsValid(eq Equation) bool {
	n := pow(2, len(eq.Nums)-1)

	for i := 0; i < n; i++ {
		ops := make([]Operator, len(eq.Nums)-1)
		for j := range ops {
			if i&(1<<uint(j)) != 0 {
				ops[j] = OpMul
			} else {
				ops[j] = OpAdd
			}
		}
		if eq.Solve(ops) {
			return true
		}
	}
	return false
}

func star2(input string) {
	eqs := parseInput(input)
	acc := 0
	for _, eq := range eqs {
		if star2IsValid(eq) {
			acc += eq.Val
		}
	}
	fmt.Println(acc)
}

func star2IsValid(eq Equation) bool {
	n := pow(3, len(eq.Nums)-1)

	for i := 0; i < n; i++ {
		ops := make([]Operator, len(eq.Nums)-1)
		b3 := strconv.FormatInt(int64(i), 3)
		b3 = strings.Repeat("0", len(eq.Nums)-1-len(b3)) + b3
		for j := range ops {
			switch b3[j] {
			case '0':
				ops[j] = OpAdd
			case '1':
				ops[j] = OpMul
			case '2':
				ops[j] = OpCat
			}
		}
		if eq.Solve(ops) {
			return true
		}
	}
	return false
}

type Operator int

const (
	OpAdd Operator = iota
	OpMul
	OpCat
)

type Equation struct {
	Val  int
	Nums []int
}

func NewEquation(s string) Equation {
	parts := strings.Split(s, ":")

	val, err := strconv.Atoi(parts[0])
	panicErr(err)

	numStrs := strings.Split(strings.TrimSpace(parts[1]), " ")
	nums := make([]int, 0, len(numStrs))
	for _, numStr := range numStrs {
		num, err := strconv.Atoi(numStr)
		panicErr(err)
		nums = append(nums, num)
	}

	return Equation{
		Val:  val,
		Nums: nums,
	}
}

func (e Equation) Solve(ops []Operator) bool {
	acc := e.Nums[0]
	nums := e.Nums[1:]

	for _, op := range ops {
		switch op {
		case OpAdd:
			acc += nums[0]
		case OpMul:
			acc *= nums[0]
		case OpCat:
			var err error
			acc, err = strconv.Atoi(fmt.Sprintf("%d%d", acc, nums[0]))
			panicErr(err)
		}
		nums = nums[1:]
	}

	return acc == e.Val
}

func pow(x, y int) int {
	if y == 0 {
		return 1
	}
	return x * pow(x, y-1)
}

func panicErr(err error) {
	if err != nil {
		panic(err)
	}
}
