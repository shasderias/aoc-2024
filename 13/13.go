package main

import (
	_ "embed"
	"fmt"
	"strings"
)

//go:embed example.txt
var example string

//go:embed input.txt
var input string

type Machine struct {
	AX, AY int
	BX, BY int
	PX, PY int
}

func (m Machine) Win(a, b int) bool {
	return a*m.AX+b*m.BX == m.PX && a*m.AY+b*m.BY == m.PY
}

func (m Machine) star1Min() int {
	minTokens := 999999
	for i := 0; i < 100; i++ {
		for j := 0; j < 100; j++ {
			if !m.Win(i, j) {
				continue
			}
			minTokens = min(minTokens, i*3+j)
		}
	}
	if minTokens == 999999 {
		return -1
	}
	return minTokens
}

func (m Machine) star2Min() int {
	bx := m.BX * m.AY
	px := m.PX * m.AY
	by := m.BY * m.AX
	py := m.PY * m.AX

	if (px-py)%(bx-by) != 0 {
		return -1
	}
	bPresses := (px - py) / (bx - by)
	if bPresses < 0 {
		return -1
	}

	if ((m.PX - (m.BX * bPresses)) % m.AX) != 0 {
		return -1
	}
	aPresses := (m.PX - (m.BX * bPresses)) / m.AX
	return aPresses*3 + bPresses
}

func main() {
	star1(example)
	star1(input)
	star2(example)
	star2(input)
}

func star1(s string) {
	machines := parseInput(s)
	acc := 0
	for _, m := range machines {
		if t := m.star1Min(); t != -1 {
			acc += t
		}
	}
	fmt.Println(acc)
}

func star2(s string) {
	machines := parseInput(s)
	acc := 0
	for _, m := range machines {
		m.PX += 10000000000000
		m.PY += 10000000000000
		if t := m.star2Min(); t != -1 {
			acc += t
		}
	}
	fmt.Println(acc)
}

func parseInput(s string) (machines []Machine) {
	inputs := strings.Split(s, "\n\n")
	for _, input := range inputs {
		machine := Machine{}
		_, err := fmt.Sscanf(input, "Button A: X+%d, Y+%d\nButton B: X+%d, Y+%d\nPrize: X=%d, Y=%d", &machine.AX, &machine.AY, &machine.BX, &machine.BY, &machine.PX, &machine.PY)
		panicErr(err)
		machines = append(machines, machine)
	}
	return machines
}

func panicErr(err error) {
	if err != nil {
		panic(err)
	}
}
