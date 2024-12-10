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
	star1(example)
	star1(input)
	star2(example)
	star2(input)
}

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
func (g Grid[T]) InBounds(p Point) bool {
	return p.X >= 0 && p.X < g.Stride && p.Y >= 0 && p.Y < g.Height
}

func NewGrid[T any](data []T, stride int) Grid[T] {
	return Grid[T]{Data: data, Stride: stride, Height: len(data) / stride}
}

func star1(s string) {
	lines := strings.Split(s, "\n")
	stride := len(lines[0])
	g := NewGrid([]rune(strings.Join(lines, "")), stride)

	trailheads := []Point{}

	for i, c := range g.Data {
		if c == '0' {
			trailheads = append(trailheads, g.IdxToPoint(i))
		}
	}

	acc := 0
	for _, th := range trailheads {
		peaks := map[Point]struct{}{}
		score1(g, th, peaks)
		acc += len(peaks)
	}
	fmt.Println(acc)
}

func score1(g Grid[rune], pt Point, peaks map[Point]struct{}) {
	if g.At(pt) == '9' {
		peaks[pt] = struct{}{}
		return
	}

	for _, dir := range dirs {
		npt := pt.Add(dir)
		if !g.InBounds(npt) {
			continue
		}
		if g.At(pt)+1 == g.At(npt) {
			score1(g, npt, peaks)
		}
	}
}

func star2(s string) {
	lines := strings.Split(s, "\n")
	stride := len(lines[0])
	g := NewGrid([]rune(strings.Join(lines, "")), stride)

	trailheads := []Point{}

	for i, c := range g.Data {
		if c == '0' {
			trailheads = append(trailheads, g.IdxToPoint(i))
		}
	}

	acc := 0
	for _, th := range trailheads {
		acc += score2(g, th)
	}
	fmt.Println(acc)
}

func score2(g Grid[rune], pt Point) int {
	if g.At(pt) == '9' {
		return 1
	}

	acc := 0
	for _, dir := range dirs {
		npt := pt.Add(dir)
		if !g.InBounds(npt) {
			continue
		}
		if g.At(pt)+1 == g.At(npt) {
			acc += score2(g, npt)
		}
	}
	return acc
}
