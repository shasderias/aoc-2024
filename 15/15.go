package main

import (
	_ "embed"
	"fmt"
	"slices"
	"strings"
)

//go:embed example1.txt
var example1 string

//go:embed example2.txt
var example2 string

//go:embed example3.txt
var example3 string

//go:embed input.txt
var input string

func main() {
	star1(example1)
	star1(example2)
	star1(input)
	star2(toStar2Input(example3))
	star2(toStar2Input(example1))
	star2(toStar2Input(example2))
	star2(toStar2Input(input))
}

func star1(s string) {
	g, insts, robotPos := parseInput(s)
	for _, inst := range insts {
		nextPos := robotPos.Add(inst)
		switch g.At(nextPos) {
		case EntNil:
			g.Set(robotPos, EntNil)
			g.Set(nextPos, EntRobot)
			robotPos = nextPos
		case EntWall:
			continue
		case EntBox:
			for castPos := nextPos.Add(inst); ; castPos = castPos.Add(inst) {
				switch g.At(castPos) {
				case EntWall:
					goto endCast
				case EntNil:
					g.Set(robotPos, EntNil)
					g.Set(nextPos, EntRobot)
					g.Set(castPos, EntBox)
					robotPos = nextPos
					goto endCast
				case EntBox:
					continue
				default:
					panic(fmt.Sprintf("unexpected entity %s", g.At(castPos)))
				}
			}
		endCast:
			break
		default:
			panic(fmt.Sprintf("unexpected entity %s", g.At(nextPos)))
		}
	}

	acc := 0
	for i, ent := range g.Data {
		if ent != EntBox {
			continue
		}
		acc += g.IdxToPoint(i).Score()
	}

	fmt.Println(g, acc)
}

func star2(s string) {
	g, insts, robotPos := parseInput(s)

	for _, inst := range insts {
		nextPos := robotPos.Add(inst)
		switch g.At(nextPos) {
		case EntNil:
			g.Set(robotPos, EntNil)
			g.Set(nextPos, EntRobot)
			robotPos = nextPos
		case EntWall:
			continue
		case EntBoxL, EntBoxR:
			switch inst {
			case N, S:
				var (
					que    = make([]Vec2, 0, 1000)
					seen   = map[Vec2]struct{}{}
					n      = 0
					movQue = que
				)

				que4s := func(v ...Vec2) {
					for _, pos := range v {
						if _, ok := seen[pos]; ok {
							continue
						}
						seen[pos] = struct{}{}
						que = append(que, pos)
						n++
					}
				}

				if g.At(nextPos) == EntBoxL {
					que4s(nextPos, nextPos.Add(E))
				} else {
					que4s(nextPos, nextPos.Add(W))
				}

				for ; len(que) > 0; que = que[1:] {
					castPos := que[0].Add(inst)
					switch g.At(castPos) {
					case EntWall:
						goto nsBreak
					case EntBoxL:
						que4s(castPos, castPos.Add(E))
					case EntBoxR:
						que4s(castPos, castPos.Add(W))
					case EntNil:
					// do nothing
					default:
						panic(fmt.Sprintf("unexpected entity %s", g.At(castPos)))
					}
				}

				movQue = movQue[:n]
				slices.Reverse(movQue)

				for _, pos := range movQue {
					g.Set(pos.Add(inst), g.At(pos))
					g.Set(pos, EntNil)
				}
				g.Set(nextPos, EntRobot)
				g.Set(robotPos, EntNil)
				robotPos = nextPos
			nsBreak:
				break
			case W, E:
				for castPos := nextPos.Add(inst); ; castPos = castPos.Add(inst) {
					switch g.At(castPos) {
					case EntWall:
						goto ewBreak
					case EntNil:
						g.Shift(castPos, robotPos, inst)
						g.Set(robotPos, EntNil)
						robotPos = nextPos
						goto ewBreak
					case EntBoxL, EntBoxR:
						continue
					default:
						panic(fmt.Sprintf("unexpected entity %s", g.At(castPos)))
					}
				}
			ewBreak:
			}
		default:
			panic(fmt.Sprintf("unexpected entity %s", g.At(nextPos)))
		}
	}

	acc := 0
	for i, ent := range g.Data {
		if ent != EntBoxL {
			continue
		}
		acc += g.IdxToPoint(i).Score()
	}

	fmt.Println(g, acc)
}

func toStar2Input(s string) string {
	repl := strings.NewReplacer(
		"#", "##",
		"O", "[]",
		".", "..",
		"@", "@.",
	)
	return repl.Replace(s)
}

func parseInput(s string) (Grid[Ent], []Vec2, Vec2) {
	var (
		parts       = strings.Split(s, "\n\n")
		gridRows    = strings.Split(parts[0], "\n")
		stride      = len(gridRows[0])
		gridDataStr = strings.Join(gridRows, "")
		gridData    = make([]Ent, 0, len(gridRows)*stride)
		robotVec    Vec2
	)
	for i, r := range gridDataStr {
		switch r {
		case '#':
			gridData = append(gridData, EntWall)
		case 'O':
			gridData = append(gridData, EntBox)
		case '[':
			gridData = append(gridData, EntBoxL)
		case ']':
			gridData = append(gridData, EntBoxR)
		case '@':
			gridData = append(gridData, EntRobot)
			robotVec = Vec2{i % stride, i / stride}
		case '.':
			gridData = append(gridData, EntNil)
		default:
			panic(fmt.Sprintf("unknown entity %q", string(r)))
		}
	}
	grid := NewGrid(gridData, stride)

	insts := make([]Vec2, 0, len(parts[1]))
	for _, is := range parts[1] {
		switch is {
		case '^':
			insts = append(insts, N)
		case 'v':
			insts = append(insts, S)
		case '<':
			insts = append(insts, W)
		case '>':
			insts = append(insts, E)
		case '\r', '\n':
			continue
		default:
			panic(fmt.Sprintf("unknown instruction %q", is))

		}
	}

	return grid, insts, robotVec
}

type Ent string

var (
	EntNil   Ent = "."
	EntWall  Ent = "#"
	EntBox   Ent = "O"
	EntBoxL  Ent = "["
	EntBoxR  Ent = "]"
	EntRobot Ent = "@"
)

type Vec2 struct {
	X, Y int
}

var (
	N = Vec2{0, -1}
	S = Vec2{0, 1}
	W = Vec2{-1, 0}
	E = Vec2{1, 0}
)
var dirs = []Vec2{N, E, S, W}

func (p Vec2) Add(p2 Vec2) Vec2 { return Vec2{p.X + p2.X, p.Y + p2.Y} }
func (p Vec2) Sub(p2 Vec2) Vec2 { return Vec2{p.X - p2.X, p.Y - p2.Y} }
func (p Vec2) Inv() Vec2        { return Vec2{-p.X, -p.Y} }
func (p Vec2) Score() int       { return p.X + p.Y*100 }

type Grid[T any] struct {
	Data   []T
	Stride int
	Height int
}

func (g Grid[T]) At(p Vec2) T             { return g.Data[p.Y*g.Stride+p.X] }
func (g Grid[T]) Set(p Vec2, v T)         { g.Data[p.Y*g.Stride+p.X] = v }
func (g Grid[T]) IdxToPoint(idx int) Vec2 { return Vec2{idx % g.Stride, idx / g.Stride} }
func (g Grid[T]) PointToIdx(p Vec2) int   { return p.Y*g.Stride + p.X }
func (g Grid[T]) InBounds(p Vec2) bool {
	return p.X >= 0 && p.X < g.Stride && p.Y >= 0 && p.Y < g.Height
}
func (g Grid[T]) Shift(s, e, dir Vec2) {
	dir = dir.Inv()
	for v := s; v != e; v = v.Add(dir) {
		if !g.InBounds(v) {
			panic(fmt.Sprintf("out of bounds %v", v))
		}
		g.Set(v, g.At(v.Add(dir)))
	}
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

func NewGrid[T any](data []T, stride int) Grid[T] {
	return Grid[T]{Data: data, Stride: stride, Height: len(data) / stride}
}
