package main

import (
	_ "embed"
	"fmt"
	"strconv"
	"strings"
)

//go:embed example_01.txt
var example01 string

//go:embed example_02.txt
var example02 string

//go:embed input.txt
var input string

func main() {
	star1(example01, 1)
	star1(example02, 6)
	star1(input, 25)
	star2(example01, 1)
	star2(example02, 6)
	star2(input, 25)
	star2(input, 75)
}

func parseInput(s string) []int {
	stonesStrs := strings.Split(s, " ")
	stones := make([]int, 0, len(stonesStrs))
	for _, ss := range stonesStrs {
		si, err := strconv.Atoi(ss)
		panicErr(err)
		stones = append(stones, si)
	}
	return stones
}

func star1(s string, blinks int) {
	stones := parseInput(s)

	for i := 0; i < blinks; i++ {
		stones = blink(stones)
	}

	fmt.Println(len(stones))
}

func star2(s string, blinks int) {
	stoneInts := parseInput(s)
	stones := make([]Stone, 0, len(stoneInts))
	for _, si := range stoneInts {
		stones = append(stones, getStone(si, 0))
	}
	for i := 0; i < blinks+1; i++ {
		acc := 0
		for _, s := range stones {
			acc += s.Len(i)
		}
		if i == blinks {
			fmt.Println(acc)
		}
	}
}

var stoneMap = map[int]*RealStone{}

func getStone(stone, gen int) Stone {
	if s, ok := stoneMap[stone]; ok {
		return VirtualStone{
			RealStone: s,
			at:        gen,
		}
	}
	s := &RealStone{stone, [][]Stone{}, gen, map[int]int{}}
	stoneMap[stone] = s
	s.children = append(s.children, []Stone{s})
	return s
}

type Stone interface {
	Len(gen int) int
}

type RealStone struct {
	stone    int
	children [][]Stone
	at       int
	memo     map[int]int
}

func (s *RealStone) Len(gen int) int {
	gen -= s.at
	if gen == 0 {
		return 1
	}
	if l, ok := s.memo[gen]; ok {
		return l
	}
	children := s.Children(gen)
	acc := 0
	for _, c := range children {
		acc += c.Len(gen)
	}
	s.memo[gen] = acc
	return acc
}

func (s *RealStone) Children(gen int) []Stone {
	if len(s.children) > gen {
		return s.children[gen]
	}

	lastGen := s.Children(gen - 1)
	thisGen := make([]Stone, 0, len(lastGen))
	for _, lg := range lastGen {
		switch v := lg.(type) {
		case *RealStone:
			if v.stone == 0 {
				thisGen = append(thisGen, getStone(1, gen))
				continue
			}
			lastGenStr := strconv.Itoa(v.stone)
			if len(lastGenStr)%2 == 0 {
				s2a, err := strconv.Atoi(lastGenStr[:len(lastGenStr)/2])
				panicErr(err)
				s2b, err := strconv.Atoi(lastGenStr[len(lastGenStr)/2:])
				panicErr(err)
				thisGen = append(thisGen, getStone(s2a, gen), getStone(s2b, gen))
				continue
			}
			thisGen = append(thisGen, getStone(v.stone*2024, gen))
		case VirtualStone:
			thisGen = append(thisGen, v)
		}
	}
	s.children = append(s.children, thisGen)
	return thisGen
}

func (s *RealStone) String() string {
	return fmt.Sprintf("{%d}", s.stone)
}

type VirtualStone struct {
	*RealStone
	at int
}

func (s VirtualStone) Len(gen int) int {
	return s.RealStone.Len(gen - s.at + s.RealStone.at)
}

func blink(stones []int) []int {
	s2 := make([]int, 0, len(stones))
	for _, s1 := range stones {
		if s1 == 0 {
			s2 = append(s2, 1)
			continue
		}

		s1Str := strconv.Itoa(s1)
		if len(s1Str)%2 == 0 {
			s2a, err := strconv.Atoi(s1Str[:len(s1Str)/2])
			panicErr(err)
			s2b, err := strconv.Atoi(s1Str[len(s1Str)/2:])
			panicErr(err)
			s2 = append(s2, s2a, s2b)
			continue
		}

		s2 = append(s2, s1*2024)
	}
	return s2
}

func (s VirtualStone) String() string {
	return fmt.Sprintf("{%d, %d}", s.RealStone.stone, s.at)
}

func panicErr(err error) {
	if err != nil {
		panic(err)
	}
}
