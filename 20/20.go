package main

import (
	_ "embed"
	"fmt"
	"slices"
	"strings"
)

//go:embed example.txt
var example string

//go:embed input.txt
var input string

func main() {
	star1(example)
	star1(input)
	star2(example, 50)
	star2(input, 100)
}

func star1(s string) {
	var (
		g     = parseInput(s)
		start = g.FindFirst(EntStart)
	)

	var (
		dist = map[Vec2]int{}
		path = []Vec2{}
	)

	var preProcess func(p Vec2, d int)

	preProcess = func(p Vec2, d int) {
		dist[p] = d
		path = append(path, p)
		for _, dir := range dirs {
			n := p.Add(dir)
			if g.At(n) == EntWall {
				continue
			} else if _, ok := dist[n]; ok {
				continue
			} else {
				preProcess(n, d+1)
			}
		}
	}

	preProcess(start, 0)

	type Cheat struct {
		P1, P2   Vec2
		TimeSave int
	}

	cheats := []Cheat{}

	for _, p := range path {
		for _, dir := range dirs {
			c1 := p.Add(dir)
			c2 := c1.Add(dir)

			if !g.InBounds(c1) || !g.InBounds(c2) {
				continue
			}
			if g.At(c1) != EntWall || (g.At(c2) != EntEnd && g.At(c2) != EntNil) {
				continue
			}
			if dist[c2] < dist[p] {
				continue
			}

			cheats = append(cheats, Cheat{c1, c2, dist[c2] - dist[p] - 2})
		}
	}

	slices.SortFunc(cheats, func(a, b Cheat) int {
		return a.TimeSave - b.TimeSave
	})

	acc := 0
	for _, c := range cheats {
		if c.TimeSave >= 100 {
			acc++
		}
	}
	fmt.Println(acc)
}

func star2(s string, atLeast int) {
	var (
		g     = parseInput(s)
		start = g.FindFirst(EntStart)
	)

	var (
		dist = map[Vec2]int{}
		path = []Vec2{}
	)

	var preProcess func(p Vec2, d int)

	preProcess = func(p Vec2, d int) {
		dist[p] = d
		path = append(path, p)
		for _, dir := range dirs {
			n := p.Add(dir)
			if g.At(n) == EntWall {
				continue
			} else if _, ok := dist[n]; ok {
				continue
			} else {
				preProcess(n, d+1)
			}
		}
	}
	preProcess(start, 0)

	type Cheat struct {
		P1, P2 Vec2
	}

	cheatsSeen := map[Cheat]bool{}
	cheatCount := map[int]int{}
	cheatAcc := 0

	for _, here := range path {
		for _, pt := range ptsWithinTaxiDist(here, 2, 20) {
			if !g.InBounds(pt) {
				continue
			}
			if _, ok := dist[pt]; !ok {
				continue
			}
			if dist[pt] < dist[here] {
				continue
			}

			saved := dist[pt] - dist[here] - taxiDist(here, pt)

			if saved < atLeast {
				continue
			}
			if _, ok := cheatsSeen[Cheat{here, pt}]; ok {
				continue
			}

			cheatsSeen[Cheat{here, pt}] = true
			cheatCount[saved]++
			cheatAcc++
		}
	}

	//fmt.Println(cheatCount)
	// 975317, 1105610
	fmt.Println(cheatAcc)
}

func parseInput(s string) Grid[Ent] {
	var (
		lines   = strings.Split(s, "\n")
		stride  = len(lines[0])
		dataStr = strings.Join(lines, "")
		data    = make([]Ent, len(dataStr))
	)
	for i, c := range dataStr {
		switch c {
		case '.':
			data[i] = EntNil
		case 'S':
			data[i] = EntStart
		case 'E':
			data[i] = EntEnd
		case '#':
			data[i] = EntWall
		default:
			panic(fmt.Sprintf("unknown char %c", c))
		}
	}
	return NewGrid(data, stride)
}

type Ent string

var (
	EntNil   Ent = "."
	EntStart Ent = "S"
	EntEnd   Ent = "E"
	EntWall  Ent = "#"
)

type Vec2 struct {
	X, Y int
}

var (
	N      = Vec2{0, -1}
	S      = Vec2{0, 1}
	W      = Vec2{-1, 0}
	E      = Vec2{1, 0}
	DirNil = Vec2{0, 0}
)
var dirs = []Vec2{N, E, S, W}

func (p Vec2) Add(p2 Vec2) Vec2 { return Vec2{p.X + p2.X, p.Y + p2.Y} }
func (p Vec2) Sub(p2 Vec2) Vec2 { return Vec2{p.X - p2.X, p.Y - p2.Y} }
func (p Vec2) Inv() Vec2        { return Vec2{-p.X, -p.Y} }
func (p Vec2) RotR() Vec2       { return Vec2{p.Y, -p.X} }
func (p Vec2) RotL() Vec2       { return Vec2{-p.Y, p.X} }

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
