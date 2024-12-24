package main

import (
	_ "embed"
	"fmt"
	"strconv"
	"strings"
)

//go:embed example1.txt
var example1 string

//go:embed example2.txt
var example2 string

//go:embed input.txt
var input string

type Computer struct {
	A, B, C int
	IP      int
	Program []Op
	Output  []int
}

func (c *Computer) Run() {
	c.Output = []int{}
	c.IP = 0

	for c.IP < len(c.Program) {
		var (
			op      = c.Program[c.IP]
			operand = int(c.Program[c.IP+1])
		)
		//fmt.Printf("%v %v %-48b %-48b %-48b %v\n", op, operand, c.A, c.B, c.C, c.Output)

		switch op {
		case OpADV:
			rop := c.ComboOperand(operand)
			c.A = c.A / pow(2, rop)
		case OpBXL:
			c.B = c.B ^ operand
		case OpBST:
			c.B = c.ComboOperand(operand) % 8
		case OpJNZ:
			if c.A == 0 {
				break
			}
			c.IP = operand
			goto next
		case OpBXC:
			c.B = c.B ^ c.C
		case OpOUT:
			rop := c.ComboOperand(operand)
			out := rop % 8
			c.Output = append(c.Output, out)
		case OpBDV:
			rop := c.ComboOperand(operand)
			c.B = c.A / pow(2, rop)
		case OpCDV:
			rop := c.ComboOperand(operand)
			c.C = c.A / pow(2, rop)
		}
		c.IP += 2
	next:
	}
}

func (c *Computer) ComboOperand(o int) int {
	switch o {
	case 0:
		return 0
	case 1:
		return 1
	case 2:
		return 2
	case 3:
		return 3
	case 4:
		return c.A
	case 5:
		return c.B
	case 6:
		return c.C
	default:
		panic("catching fire")
	}
}

func (c *Computer) FmtOutput() string {
	buf := strings.Builder{}
	for i, o := range c.Output {
		buf.WriteString(strconv.Itoa(o))
		if i < len(c.Output)-1 {
			buf.WriteString(",")
		}
	}
	return buf.String()
}

func (c *Computer) FmtProgram() string {
	buf := strings.Builder{}
	for i, o := range c.Program {
		buf.WriteString(strconv.Itoa(int(o)))
		if i < len(c.Program)-1 {
			buf.WriteString(",")
		}
	}
	return buf.String()
}

func pow(x, y int) int {
	if y == 0 {
		return 1
	}
	acc := x
	for i := 1; i < y; i++ {
		acc *= x
	}
	return acc
}

type Op int

const (
	OpADV Op = iota
	OpBXL
	OpBST
	OpJNZ
	OpBXC
	OpOUT
	OpBDV
	OpCDV
)

func (o Op) String() string {
	switch o {
	case OpADV:
		return "ADV"
	case OpBXL:
		return "BXL"
	case OpBST:
		return "BST"
	case OpJNZ:
		return "JNZ"
	case OpBXC:
		return "BXC"
	case OpOUT:
		return "OUT"
	case OpBDV:
		return "BDV"
	case OpCDV:
		return "CDV"
	default:
		return "???"
	}
}

func main() {
	star1(example1)
	star1(input)
	//star2(example2)
	star2(input)
}

func star1(s string) {
	c := parseInput(s)
	c.Run()
	fmt.Println(c.FmtOutput())
}

type State struct {
	Found string
	N     int
}

func (s State) A() int {
	a, _ := strconv.ParseUint(s.Found+strings.Repeat("0", 3*(15-s.N)), 2, 64)
	return int(a)
}

func genDigit() []string {
	return []string{"000", "001", "010", "011", "100", "101", "110", "111"}
}

func star2(s string) {
	c := parseInput(s)

	_, ib, ic := c.A, c.B, c.C

	try := []State{}

	for i, d := range genDigit() {
		if i == 0 {
			continue
		}
		try = append(try, State{d, 0})
	}

	found := -1

	for ; len(try) > 0; try = try[1:] {
		cur := try[0]

		c.A, c.B, c.C = cur.A(), ib, ic
		c.Run()
		fmt.Println(cur, cur.A(), c.Output, c.Program)
		if c.Output[15-cur.N] == int(c.Program[15-cur.N]) {
			if cur.N == 15 {
				found = cur.A()
				break
			}
			for _, d := range genDigit() {
				try = append(try, State{cur.Found + d, cur.N + 1})
			}
		}
	}

	fmt.Println(found)
	for i := 1; i < 10; i++ {
		c.A, c.B, c.C = found-i, ib, ic
		c.Run()
		fmt.Println(found-i, c.FmtOutput())
	}
}

func parseInput(s string) *Computer {
	c := &Computer{}
	parts := strings.Split(s, "\n\n")
	fmt.Sscanf(parts[0],
		"Register A: %d\nRegister B: %d\nRegister C: %d",
		&c.A, &c.B, &c.C)
	programStr := strings.TrimPrefix(parts[1], "Program: ")
	opStrs := strings.Split(programStr, ",")
	for _, opStr := range opStrs {
		op, err := strconv.Atoi(opStr)
		panicErr(err)
		c.Program = append(c.Program, Op(op))
	}

	return c
}

func panicErr(err error) {
	if err != nil {
		panic(err)
	}
}
