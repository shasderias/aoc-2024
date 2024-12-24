package main

import (
	_ "embed"
	"fmt"
	"strings"
)

//go:embed example1.txt
var example1 string

//go:embed example2.txt
var example2 string

//go:embed input.txt
var input string

func main() {
	star1(example1)
	star1(example2)
	star1(input)
	star2(example1)
	star2(example2)
	star2(input)
}

type State struct {
	Pos, Dir Vec2
}

type SearchState struct {
	State
	Score int
}

type SearchState2 struct {
	State
	Score int
	Path  []Vec2
}

func (ss2 SearchState2) Append(nextPos, nextDir Vec2, cost int) SearchState2 {
	path := make([]Vec2, len(ss2.Path), len(ss2.Path)+1)
	copy(path, ss2.Path)
	if nextPos != ss2.Path[len(ss2.Path)-1] {
		path = append(path, nextPos)
	}
	return SearchState2{State{nextPos, nextDir}, ss2.Score + cost, path}
}

func star1(s string) int {
	g := parseInput(s)
	st := g.FindFirst(EntStart)
	ed := g.FindFirst(EntEnd)
	memo := map[State]int{}

	memoIfBetter := func(st State, score int) bool {
		if s, ok := memo[st]; !ok || score < s {
			memo[st] = score
			return true
		}
		return false
	}

	que := []SearchState{{State{st, E}, 0}}

	for ; len(que) > 0; que = que[1:] {
		cur := que[0]
		memo[cur.State] = cur.Score

		fwd := cur.Pos.Add(cur.Dir)
		switch g.At(fwd) {
		case EntNil:
			if memoIfBetter(State{Pos: fwd, Dir: cur.Dir}, cur.Score+1) {
				que = append(que, SearchState{State{Pos: fwd, Dir: cur.Dir}, cur.Score + 1})
			}
		case EntEnd:
			memoIfBetter(State{Pos: fwd, Dir: DirNil}, cur.Score+1)
		}

		lef := cur.Pos.Add(cur.Dir.RotL())
		switch g.At(lef) {
		case EntNil, EntEnd:
			if memoIfBetter(State{Pos: cur.Pos, Dir: cur.Dir.RotL()}, cur.Score+1000) {
				que = append(que, SearchState{State{Pos: cur.Pos, Dir: cur.Dir.RotL()}, cur.Score + 1000})
			}
		}

		rig := cur.Pos.Add(cur.Dir.RotR())
		switch g.At(rig) {
		case EntNil, EntEnd:
			if memoIfBetter(State{Pos: cur.Pos, Dir: cur.Dir.RotR()}, cur.Score+1000) {
				que = append(que, SearchState{State{Pos: cur.Pos, Dir: cur.Dir.RotR()}, cur.Score + 1000})
			}
		}
	}
	fmt.Println(memo[State{ed, DirNil}])
	return memo[State{ed, DirNil}]
}

func star2(s string) {
	bestScore := star1(s)
	g := parseInput(s)
	st := g.FindFirst(EntStart)
	ed := g.FindFirst(EntEnd)
	memo := map[State]int{}
	visited := map[Vec2]bool{
		st: true,
		ed: true,
	}

	memoIfBetter := func(st State, score int) bool {
		if s, ok := memo[st]; !ok || score <= s {
			memo[st] = score
			return true
		}
		return false
	}

	que := []SearchState2{{State{st, E}, 0, []Vec2{st}}}

	for ; len(que) > 0; que = que[1:] {
		cur := que[0]
		memo[cur.State] = cur.Score

		fwd := cur.Pos.Add(cur.Dir)
		switch g.At(fwd) {
		case EntNil:
			if memoIfBetter(State{Pos: fwd, Dir: cur.Dir}, cur.Score+1) {
				que = append(que, cur.Append(fwd, cur.Dir, 1))
			}
		case EntEnd:
			//memoIfBetter(State{Pos: fwd, Dir: DirNil}, cur.Score+1000)
			if cur.Score+1 == bestScore {
				for _, p := range cur.Path {
					visited[p] = true
				}
			}
		}

		lef := cur.Pos.Add(cur.Dir.RotL())
		switch g.At(lef) {
		case EntNil, EntEnd:
			if memoIfBetter(State{Pos: cur.Pos, Dir: cur.Dir.RotL()}, cur.Score+1000) {
				que = append(que, cur.Append(cur.Pos, cur.Dir.RotL(), 1000))
			}
		}

		rig := cur.Pos.Add(cur.Dir.RotR())
		switch g.At(rig) {
		case EntNil, EntEnd:
			if memoIfBetter(State{Pos: cur.Pos, Dir: cur.Dir.RotR()}, cur.Score+1000) {
				que = append(que, cur.Append(cur.Pos, cur.Dir.RotR(), 1000))
			}
		}
	}
	fmt.Println(len(visited))
}

func parseInput(s string) Grid[Ent] {
	var (
		rowStrs = strings.Split(s, "\n")
		stride  = len(rowStrs[0])
		dataStr = strings.Join(rowStrs, "")
		data    = make([]Ent, 0, len(dataStr))
	)

	for _, r := range dataStr {
		switch r {
		case '#':
			data = append(data, EntWall)
		case 'S':
			data = append(data, EntStart)
		case 'E':
			data = append(data, EntEnd)
		case '.':
			data = append(data, EntNil)
		}
	}

	return NewGrid(data, stride)
}

type Ent string

var (
	EntNil   Ent = "."
	EntWall  Ent = "#"
	EntStart Ent = "S"
	EntEnd   Ent = "E"
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
