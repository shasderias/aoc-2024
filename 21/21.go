package main

import (
	_ "embed"
	"fmt"
	"slices"
	"strconv"
	"strings"

	"gonum.org/v1/gonum/stat/combin"
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

func star1(s string) {
	inputs := parseInput(s)

	np := NewNumPad()
	ap := NewArrowPad()

	calc := func(pad Pad, in Input) (o []Input) {
		lens := []int{}
		choices := [][]string{}

		for p := range in.Iter() {
			routes := pad.ShortestRoutes(p[0], p[1])
			choices = append(choices, routeSliceToStringSlice(routes))
			lens = append(lens, len(routes))
		}

		cp := combin.Cartesian(lens)

		for _, p := range cp {
			buf := strings.Builder{}
			for idx, pick := range p {
				buf.WriteString(choices[idx][pick])
			}
			o = append(o, Input{buf.String(), 0, nil})
		}
		return o
	}

	acc := 0
	for _, in1 := range inputs {
		minLen := 9999
		for _, in2 := range calc(np, in1) {
			for _, in3 := range calc(ap, in2) {
				for _, o := range calc(ap, in3) {
					minLen = min(minLen, o.Len())
				}
			}
		}
		complexity := in1.NumPart() * minLen
		fmt.Println(in1, in1.NumPart(), minLen, complexity)
		acc += complexity
	}
	fmt.Println(acc)
}

func star2(s string) {
	inputs := parseInput(s)

	np := NewNumPad()
	ap := NewArrowPad()

	precompute := func(s string) string {
		switch s {
		case "AA":
			return "A"
		case "^^":
			return "A"
		case "vv":
			return "A"
		case "<<":
			return "A"
		case ">>":
			return "A"
		}
		minLen := 999999999999999999
		var minInput Input

		calc := func(pad Pad, in Input) (o []Input) {
			if in.depth == 4 {
				if in.Len() < minLen {
					minInput = in
					minLen = in.Len()
				}
				return
			}
			lens := []int{}
			choices := [][]string{}

			for p := range in.Iter2() {
				routes := pad.ShortestRoutes(p[0], p[1])
				choices = append(choices, routeSliceToStringSlice(routes))
				lens = append(lens, len(routes))
			}

			cp := combin.Cartesian(lens)

			for _, p := range cp {
				buf := strings.Builder{}
				for idx, pick := range p {
					buf.WriteString(choices[idx][pick])
				}
				o = append(o, Input{buf.String(), in.depth + 1, append(append([]string{}, in.history...), in.s)})
			}
			return o
		}
		que := []Input{{s, 0, nil}}
		for ; len(que) > 0; que = que[1:] {
			results := calc(ap, que[0])
			que = append(que, results...)
		}

		return minInput.history[1]
	}

	numpadCalc := func(pad Pad, in Input) (o []Input) {
		lens := []int{}
		choices := [][]string{}

		for p := range in.Iter() {
			routes := pad.ShortestRoutes(p[0], p[1])
			choices = append(choices, routeSliceToStringSlice(routes))
			lens = append(lens, len(routes))
		}

		cp := combin.Cartesian(lens)

		for _, p := range cp {
			buf := strings.Builder{}
			for idx, pick := range p {
				buf.WriteString(choices[idx][pick])
			}
			o = append(o, Input{buf.String(), 0, nil})
		}
		return o
	}

	symbols := []string{"A", "^", "v", "<", ">"}
	combi := combin.Cartesian([]int{5, 5})

	breedMap := map[string]string{}

	for _, c := range combi {
		target := symbols[c[0]] + symbols[c[1]]
		postCycle := precompute(target)
		if postCycle != "A" {
			postCycle = "A" + postCycle
		}
		breedMap[target] = postCycle
		fmt.Println("computed", target, "->", postCycle)
	}

	cycleMap := map[string]func(c int) map[string]int{}

	for sf, ef := range breedMap {
		bm := breakToMap(ef)
		//fmt.Println(sf, "--cycleMap>", bm)
		cycleMap[sf] = func(c int) map[string]int {
			cm := map[string]int{}
			for k, v := range bm {
				cm[k] = v * c
			}
			return cm
		}
	}
	cycleMap["A"] = func(c int) map[string]int {
		return map[string]int{"A": c}
	}

	complexity := 0
	for _, in1 := range inputs {
		//if in1.s != "029A" {
		//	break
		//}
		results := numpadCalc(np, in1)
		minResult := 99999999999999999
		for _, r := range results {
			bm := breakToMap(r.s)
			//fmt.Println(r.s, "-->", bm)

			first := "A" + string(r.s[0])

			for i := 0; i < 25; i++ {
				first, bm = breed(first, bm, cycleMap, breedMap)
				//fmt.Println(in1, r, bm, acc)
			}
			acc := 0
			for _, v := range bm {
				acc += v
			}
			minResult = min(minResult, acc)
		}
		fmt.Println(in1.s, minResult, in1.NumPart()*minResult)
		complexity += in1.NumPart() * minResult
	}
	fmt.Println(complexity)
}

func breakToPairs(s string) (pairs []string) {
	for i := 0; i < len(s)-1; i++ {
		pairs = append(pairs, s[i:i+2])
	}
	return pairs
}

func breakToMap(s string) (m map[string]int) {
	if s == "A" {
		return map[string]int{"A": 1}
	}
	m = map[string]int{}
	for i := 0; i < len(s)-1; i++ {
		m[s[i:i+2]]++
	}
	return m
}

func breed(first string, pool map[string]int, cycleMap map[string]func(int) map[string]int, breedMap map[string]string) (nextFirst string, nextPool map[string]int) {
	nextCycle := map[string]int{}
	for k, v := range pool {
		t := cycleMap[k](v)
		for k, v := range t {
			nextCycle[k] += v
		}
		//fmt.Println(k, v, "==>", t)
	}

	nextFirst = breedMap[first]
	if len(nextFirst) > 1 {
		newThisCycleFromFirst := breakToMap(nextFirst)
		for k, v := range newThisCycleFromFirst {
			//fmt.Println(first, "==>", k, v)
			nextCycle[k] += v
		}
	}
	nextFirst = "A" + string(nextFirst[0])

	return nextFirst, nextCycle
}

func NewNumPad() Pad {
	g := NewGrid([]string{"7", "8", "9", "4", "5", "6", "1", "2", "3", "", "0", "A"}, 3)
	return Pad{Grid: g}
}

func NewArrowPad() Pad {
	g := NewGrid([]string{"", "^", "A", "<", "v", ">"}, 3)
	return Pad{Grid: g}
}

func fmtRoute(dirs []Vec2) string {
	var b strings.Builder
	for _, d := range dirs {
		b.WriteString(d.Dir())
	}
	return b.String()
}

func routeSliceToStringSlice(rs [][]Vec2) []string {
	ss := []string{}
	for _, r := range rs {
		ss = append(ss, fmtRoute(r)+"A")
	}
	return ss
}

type Pad struct {
	Grid[string]
}

func (p Pad) ShortestRoutes(sc, ec string) (paths [][]Vec2) {
	s := p.FindFirst(sc)
	e := p.FindFirst(ec)

	if s == e {
		return [][]Vec2{{}}
	}

	shortestDist := taxiDist(s, e)

	var find func(pos Vec2, d int, path []Vec2, visited []Vec2)
	find = func(pos Vec2, d int, path []Vec2, visited []Vec2) {
		//fmt.Println(sc, ec, s, e, pos, d, fmtRoute(path), visited)
		if d > shortestDist {
			return
		}
		if pos == e {
			paths = append(paths, append([]Vec2{}, path...))
			return
		}
		for _, dir := range dirs {
			n := pos.Add(dir)
			if p.InBounds(n) && p.At(n) != "" && !slices.Contains(visited, n) {
				find(n, d+1, append(path, dir), append(visited, pos))
			}
		}
	}
	find(s, 0, []Vec2{}, []Vec2{s})

	return
}

type Input struct {
	s       string
	depth   int
	history []string
}

func (i Input) Iter() func(func([2]string) bool) {
	items := append([]string{"A"}, strings.Split(i.s, "")...)
	return func(yield func([2]string) bool) {
		for i := 1; i < len(items); i++ {
			if !yield([2]string{items[i-1], items[i]}) {
				return
			}
		}
	}
}

func (i Input) Iter2() func(func([2]string) bool) {
	items := strings.Split(i.s, "")
	return func(yield func([2]string) bool) {
		for i := 1; i < len(items); i++ {
			if !yield([2]string{items[i-1], items[i]}) {
				return
			}
		}
	}
}

func (i Input) NumPart() int {
	num, err := strconv.Atoi(strings.TrimSuffix(i.s, "A"))
	if err != nil {
		panic(err)
	}
	return num
}

func (i Input) Len() int {
	return len(i.s)
}

func parseInput(s string) (inputs []Input) {
	lines := strings.Split(s, "\n")
	for _, line := range lines {
		inputs = append(inputs, Input{line, 0, nil})
	}
	return
}

type Vec2 struct {
	X, Y int
}

func (p Vec2) Add(p2 Vec2) Vec2 { return Vec2{p.X + p2.X, p.Y + p2.Y} }
func (p Vec2) Sub(p2 Vec2) Vec2 { return Vec2{p.X - p2.X, p.Y - p2.Y} }
func (p Vec2) Inv() Vec2        { return Vec2{-p.X, -p.Y} }
func (p Vec2) RotR() Vec2       { return Vec2{p.Y, -p.X} }
func (p Vec2) RotL() Vec2       { return Vec2{-p.Y, p.X} }
func (p Vec2) Dir() string {
	switch p {
	case N:
		return "^"
	case S:
		return "v"
	case W:
		return "<"
	case E:
		return ">"
	case DirNil:
		return "A"
	}
	return ""
}

var (
	N      = Vec2{0, -1}
	S      = Vec2{0, 1}
	W      = Vec2{-1, 0}
	E      = Vec2{1, 0}
	DirNil = Vec2{0, 0}
)
var dirs = []Vec2{N, E, S, W}

type Grid[T comparable] struct {
	Data   []T
	Stride int
	Height int
}

func (g Grid[T]) At(p Vec2) T           { return g.Data[p.Y*g.Stride+p.X] }
func (g Grid[T]) Set(p Vec2, v T)       { g.Data[p.Y*g.Stride+p.X] = v }
func (g Grid[T]) IdxToVec(idx int) Vec2 { return Vec2{idx % g.Stride, idx / g.Stride} }
func (g Grid[T]) VecToIdx(p Vec2) int   { return p.Y*g.Stride + p.X }
func (g Grid[T]) InBounds(p Vec2) bool {
	return p.X >= 0 && p.X < g.Stride && p.Y >= 0 && p.Y < g.Height
}
func (g Grid[T]) FindFirst(t T) Vec2 {
	for i, v := range g.Data {
		if v == t {
			return g.IdxToVec(i)
		}
	}
	panic("not found")
}

func (g Grid[T]) String() string {
	var b strings.Builder
	for i, v := range g.Data {
		if i%g.Stride == 0 {
			b.WriteString("\n")
		}
		b.WriteString(fmt.Sprintf("%s", v))
	}
	return b.String()
}

func NewGrid[T comparable](data []T, stride int) Grid[T] {
	return Grid[T]{Data: data, Stride: stride, Height: len(data) / stride}
}

func ptsWithinTaxiDist(p Vec2, distMin, distMax int) (points []Vec2) {
	for dx := -distMax; dx <= distMax; dx++ {
		for dy := -distMax; dy <= distMax; dy++ {
			dist := abs(dx) + abs(dy)
			if dist >= distMin && dist <= distMax {
				points = append(points, Vec2{p.X + dx, p.Y + dy})
			}
		}
	}
	return points
}

func taxiDist(p1, p2 Vec2) int {
	return abs(p1.X-p2.X) + abs(p1.Y-p2.Y)
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
