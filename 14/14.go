package main

import (
	_ "embed"
	"fmt"
	"os"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

//go:embed example.txt
var example string

//go:embed input.txt
var input string

type Vec2 struct {
	X, Y int
}

func main() {
	//star1(example, 11, 7)
	//star1(input, 101, 103)
	//star2(example, 11, 7)
	star2(input, 101, 103)
}

func star1(s string, w, h int) {
	robots := parseInput(s)
	field := Field{Stride: w, Height: h}
	mid := Vec2{field.Stride / 2, field.Height / 2}
	quad := make([]int, 4)

	for _, r := range robots {
		posAt100 := field.NormalizePos(r.AtTick(100))
		switch {
		case posAt100.X < mid.X && posAt100.Y < mid.Y:
			quad[0]++
		case posAt100.X > mid.X && posAt100.Y < mid.Y:
			quad[1]++
		case posAt100.X < mid.X && posAt100.Y > mid.Y:
			quad[2]++
		case posAt100.X > mid.X && posAt100.Y > mid.Y:
			quad[3]++
		default:
		}
	}
	fmt.Println(quad[0], quad[1], quad[2], quad[3], quad[0]*quad[1]*quad[2]*quad[3])
}

func star2(s string, w, h int) {
	robots := parseInput(s)
	field := Field{Stride: w, Height: h}

	for t := 0; t < 1000000; t++ {
		mid := Vec2{field.Stride / 2, field.Height / 2}
		quad := make([]int, 4)
		for _, r := range robots {
			pos := field.NormalizePos(r.AtTick(t))
			switch {
			case pos.X < mid.X && pos.Y < mid.Y:
				quad[0]++
			case pos.X > mid.X && pos.Y < mid.Y:
				quad[1]++
			case pos.X < mid.X && pos.Y > mid.Y:
				quad[2]++
			case pos.X > mid.X && pos.Y > mid.Y:
				quad[3]++
			default:
			}
		}
		if quad[0]*quad[1]*quad[2]*quad[3] < 100000000 {
			fmt.Println(render(field, robots, t))
		}
		time.Sleep(1 * time.Millisecond)
	}
	return

	prog := tea.NewProgram(Model{
		n:      0,
		field:  field,
		robots: robots,
	})

	if _, err := prog.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

type Model struct {
	n      int
	field  Field
	robots []Robot
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "right", "k":
			m.n++
		case "left", "j":
			m.n--
		case "l":
			m.n += 100
		case "h":
			m.n -= 100
		case "ctrl+c", "q", "c":
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m Model) View() string {
	occ := make(map[Vec2]struct{})
	for _, r := range m.robots {
		occ[m.field.NormalizePos(r.AtTick(m.n))] = struct{}{}
	}

	buf := strings.Builder{}

	for y := 0; y < m.field.Height; y++ {
		for x := 0; x < m.field.Stride; x++ {
			if _, ok := occ[Vec2{x, y}]; ok {
				buf.WriteString("*")
			} else {
				buf.WriteString(".")
			}
		}
		buf.WriteString("\n")
	}
	buf.WriteString("\n")
	buf.WriteString(fmt.Sprintf("t=%d", m.n))

	return buf.String()
}

func render(field Field, robots []Robot, t int) string {
	occ := make(map[Vec2]struct{})
	for _, r := range robots {
		occ[field.NormalizePos(r.AtTick(t))] = struct{}{}
	}

	buf := strings.Builder{}

	for y := 0; y < field.Height; y++ {
		for x := 0; x < field.Stride; x++ {
			if _, ok := occ[Vec2{x, y}]; ok {
				buf.WriteString("*")
			} else {
				buf.WriteString(" ")
			}
		}
		buf.WriteString("\n")
	}
	buf.WriteString("\n")
	buf.WriteString(fmt.Sprintf("t=%d", t))

	return buf.String()
}

func NewModel() Model {
	return Model{}
}

type Robot struct {
	Pos Vec2
	Vel Vec2
}

type Field struct {
	Stride int
	Height int
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

func (f Field) NormalizePos(p Vec2) Vec2 {
	return Vec2{
		X: (p.X + (((abs(p.X) / f.Stride) + 1) * f.Stride)) % f.Stride,
		Y: (p.Y + (((abs(p.Y) / f.Height) + 1) * f.Height)) % f.Height,
	}
}

func (r Robot) AtTick(t int) Vec2 {
	return Vec2{
		X: r.Pos.X + r.Vel.X*t,
		Y: r.Pos.Y + r.Vel.Y*t,
	}
}

func parseInput(s string) []Robot {
	lines := strings.Split(s, "\n")
	robots := make([]Robot, 0, len(lines))
	for _, line := range lines {
		robot := Robot{}
		_, err := fmt.Sscanf(line, "p=%d,%d v=%d,%d",
			&robot.Pos.X, &robot.Pos.Y, &robot.Vel.X, &robot.Vel.Y)
		panicErr(err)
		robots = append(robots, robot)

	}
	return robots
}

func panicErr(err error) {
	if err != nil {
		panic(err)
	}
}
