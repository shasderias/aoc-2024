package main

import (
	_ "embed"
	"fmt"
	"slices"
	"strings"
)

//go:embed example_01.txt
var example01 string

//go:embed example_02.txt
var example02 string

//go:embed example_03.txt
var example03 string

//go:embed example_04.txt
var example04 string

//go:embed input.txt
var input string

type Point struct {
	X, Y int
}

var (
	N = Point{0, -1}
	S = Point{0, 1}
	W = Point{-1, 0}
	E = Point{1, 0}
)
var dirs = []Point{N, E, S, W}

func (p Point) Add(p2 Point) Point { return Point{p.X + p2.X, p.Y + p2.Y} }
func (p Point) Sub(p2 Point) Point { return Point{p.X - p2.X, p.Y - p2.Y} }

type Grid[T any] struct {
	Data   []T
	Stride int
	Height int
}

func (g Grid[T]) At(p Point) T             { return g.Data[p.Y*g.Stride+p.X] }
func (g Grid[T]) Set(p Point, v T)         { g.Data[p.X*g.Stride+p.Y] = v }
func (g Grid[T]) IdxToPoint(idx int) Point { return Point{idx % g.Stride, idx / g.Stride} }
func (g Grid[T]) PointToIdx(p Point) int   { return p.Y*g.Stride + p.X }
func (g Grid[T]) InBounds(p Point) bool {
	return p.X >= 0 && p.X < g.Stride && p.Y >= 0 && p.Y < g.Height
}

func NewGrid[T any](data []T, stride int) Grid[T] {
	return Grid[T]{Data: data, Stride: stride, Height: len(data) / stride}
}

func main() {
	star1(example01)
	star1(example02)
	star1(example03)
	star1(input)
	star2(example01)
	star2(example04)
	star2(example03)
	star2(input)
}

func star1(s string) {
	g := parseInput(s)

	visited := make([]bool, len(g.Data))

	var fill func(c rune, idx int) (int, int)
	fill = func(c rune, idx int) (p, a int) {
		a++
		visited[idx] = true
		for _, d := range dirs {
			pt := g.IdxToPoint(idx).Add(d)
			switch {
			case !g.InBounds(pt):
				p++
			case g.At(pt) != c:
				p++
			case g.At(pt) == c && !visited[g.PointToIdx(pt)]:
				p2, a2 := fill(c, g.PointToIdx(pt))
				p += p2
				a += a2
			}
		}
		return p, a
	}

	acc := 0
	for i, c := range g.Data {
		if visited[i] {
			continue
		}
		p, a := fill(c, i)
		acc += p * a
	}

	fmt.Println(acc)
}

type Side struct {
	facing Point
	at     Point
}

func nsSort(a, b Side) int {
	if a.at.Y == b.at.Y {
		return a.at.X - b.at.X
	}
	return a.at.Y - b.at.Y
}
func ewSort(a, b Side) int {
	if a.at.X == b.at.X {
		return a.at.Y - b.at.Y
	}
	return a.at.X - b.at.X
}

func star2(s string) {
	g := parseInput(s)

	visited := make([]bool, len(g.Data))

	var fill func(c rune, idx int) (int, int, []Side)
	fill = func(c rune, idx int) (p, a int, sides []Side) {
		a++
		visited[idx] = true
		for _, d := range dirs {
			pt := g.IdxToPoint(idx).Add(d)
			switch {
			case !g.InBounds(pt):
				p++
				sides = append(sides, Side{d, g.IdxToPoint(idx)})
			case g.At(pt) != c:
				p++
				sides = append(sides, Side{d, g.IdxToPoint(idx)})
			case g.At(pt) == c && !visited[g.PointToIdx(pt)]:
				p2, a2, s2 := fill(c, g.PointToIdx(pt))
				p += p2
				a += a2
				sides = append(sides, s2...)
			}
		}
		return p, a, sides
	}

	acc := 0
	for i, c := range g.Data {
		if visited[i] {
			continue
		}

		_, a, sides := fill(c, i)

		var nSides, sSides, eSides, wSides []Side
		for _, s := range sides {
			switch s.facing {
			case N:
				nSides = append(nSides, s)
			case S:
				sSides = append(sSides, s)
			case E:
				eSides = append(eSides, s)
			case W:
				wSides = append(wSides, s)
			}
		}
		slices.SortFunc(nSides, nsSort)
		slices.SortFunc(sSides, nsSort)
		slices.SortFunc(eSides, ewSort)
		slices.SortFunc(wSides, ewSort)

		sideCount := countNSSides(nSides) + countNSSides(sSides) + countEWSides(eSides) + countEWSides(wSides)
		acc += sideCount * a
	}

	fmt.Println(acc)
}

func parseInput(s string) Grid[rune] {
	lines := strings.Split(s, "\n")
	stride := len(lines[0])
	g := NewGrid([]rune(strings.Join(lines, "")), stride)
	return g
}

func countNSSides(sides []Side) (count int) {
	prev := Point{-1, -1}

	for _, side := range sides {
		if side.at.Y != prev.Y || side.at.X-1 != prev.X {
			count++
		}
		prev = side.at
	}

	return count
}

func countEWSides(sides []Side) (count int) {
	prev := Point{-1, -1}

	for _, side := range sides {
		if side.at.X != prev.X || side.at.Y-1 != prev.Y {
			count++
		}
		prev = side.at
	}

	return count
}
