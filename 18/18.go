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

func main() {
	star1(example, 7, 7, 11)
	star1(input, 71, 71, 1024)
	star2(example, 7, 7)
	star2(input, 71, 71)

}

type State struct {
	Pos  Vec2
	Cost int
}

func star1(s string, w, h, first int) {
	g := parseInput(s, w, h, first)
	e := Vec2{w - 1, h - 1}
	memo := make(map[Vec2]bool)
	que := []State{{Pos: Vec2{0, 0}, Cost: 0}}
	memo[Vec2{0, 0}] = true

	for ; len(que) > 0; que = que[1:] {
		cur := que[0]
		if cur.Pos == e {
			fmt.Println(cur.Cost)
			break
		}
		for _, d := range dirs {
			next := cur.Pos.Add(d)
			if !g.InBounds(next) || g.At(next) == EntWall || memo[next] {
				continue
			}
			memo[next] = true
			que = append(que, State{Pos: next, Cost: cur.Cost + 1})
		}
	}
}

func star2(s string, w, h int) {
	e := Vec2{w - 1, h - 1}

	coords := parseCoords(s)

	makeGrid := func(t int) Grid[Ent] {
		data := make([]Ent, w*h)
		for i := range data {
			data[i] = EntNil
		}
		g := NewGrid(data, w)
		for i := 0; i < t; i++ {
			g.Set(coords[i], EntWall)
		}
		return g
	}

	solve := func(t int) bool {
		g := makeGrid(t)
		fmt.Println(g)
		memo := make(map[Vec2]bool)
		que := []State{{Pos: Vec2{0, 0}, Cost: 0}}
		memo[Vec2{0, 0}] = true

		for ; len(que) > 0; que = que[1:] {
			cur := que[0]
			if cur.Pos == e {
				return true
			}
			for _, d := range dirs {
				next := cur.Pos.Add(d)
				if !g.InBounds(next) || g.At(next) == EntWall || memo[next] {
					continue
				}
				memo[next] = true
				que = append(que, State{Pos: next, Cost: cur.Cost + 1})
			}
		}
		return false
	}

	for i := range coords {
		if !solve(i) {
			fmt.Println(i, coords[i-1])
			break
		}
	}
}

func parseInput(s string, w, h, first int) Grid[Ent] {
	data := make([]Ent, w*h)
	for i := range data {
		data[i] = EntNil
	}
	g := NewGrid(data, w)

	for i, coord := range strings.Split(s, "\n") {
		var x, y int
		fmt.Sscanf(coord, "%d,%d", &x, &y)
		g.Set(Vec2{x, y}, EntWall)
		if i == first {
			break
		}
	}
	return g
}

func parseCoords(s string) (coords []Vec2) {
	for _, coord := range strings.Split(s, "\n") {
		var x, y int
		fmt.Sscanf(coord, "%d,%d", &x, &y)
		coords = append(coords, Vec2{x, y})
	}
	return coords
}

type Ent string

var (
	EntNil  Ent = "."
	EntWall Ent = "#"
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
