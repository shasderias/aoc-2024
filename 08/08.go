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

type Point struct {
	X, Y int
}

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

func main() {
	star1(example)
	star1(input)
	star2(example)
	star2(input)
}

func star1(s string) {
	lines := strings.Split(s, "\n")
	stride := len(lines[0])
	g := NewGrid([]rune(strings.Join(lines, "")), stride)

	antennas := map[rune][]Point{}

	for i, c := range g.Data {
		switch c {
		case '.':
			continue
		default:
			antennas[c] = append(antennas[c], g.IdxToPoint(i))
		}
	}

	antinodes := map[Point]struct{}{}

	for _, ant := range antennas {
		for j, pt1 := range ant {
			for k, pt2 := range ant {
				if j == k {
					continue
				}
				diff := pt2.Sub(pt1)
				antinodes[pt1.Sub(diff)] = struct{}{}
				antinodes[pt2.Add(diff)] = struct{}{}
			}
		}
	}

	acc := 0
	for an := range antinodes {
		if g.InBounds(an) {
			acc++
		}
	}
	fmt.Println(acc)
}

func star2(s string) {
	lines := strings.Split(s, "\n")
	stride := len(lines[0])
	g := NewGrid([]rune(strings.Join(lines, "")), stride)

	antennas := map[rune][]Point{}

	for i, c := range g.Data {
		switch c {
		case '.':
			continue
		default:
			antennas[c] = append(antennas[c], g.IdxToPoint(i))
		}
	}

	antinodes := map[Point]struct{}{}

	for _, ant := range antennas {
		for j, pt1 := range ant {
			for k, pt2 := range ant {
				if j == k {
					continue
				}
				diff := pt2.Sub(pt1)
				for pt := pt1; g.InBounds(pt); pt = pt.Sub(diff) {
					antinodes[pt] = struct{}{}
				}
				for pt := pt2; g.InBounds(pt); pt = pt.Add(diff) {
					antinodes[pt] = struct{}{}
				}
			}
		}
	}

	acc := 0
	for an := range antinodes {
		if g.InBounds(an) {
			acc++
		}
	}
	fmt.Println(acc)
}
