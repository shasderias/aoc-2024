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

func star1(s string) {
	avail, towels := parseInput(s)
	acc := 0
	for _, t := range towels {
		if possible(avail, t) {
			acc++
		}
	}
	fmt.Println(acc)
}

func star2(s string) {
	avail, towels := parseInput(s)
	acc := 0
	for _, t := range towels {
		c := combinations(avail, t)
		acc += c
	}
	fmt.Println(acc)
}

func possible(avail []string, towel string) bool {
	que := []string{towel}
	memo := make(map[string]bool)

	for ; len(que) > 0; que = que[1:] {
		cur := que[0]
		memo[cur] = true

		for _, a := range avail {
			if cur == a {
				return true
			}
			if strings.HasPrefix(cur, a) && !memo[strings.TrimPrefix(cur, a)] {
				que = append(que, strings.TrimPrefix(cur, a))
			}
		}
	}
	return false
}

func combinations(avail []string, towel string) int {
	memo := make(map[string]int)

	return combi(memo, avail, towel)
}

func combi(memo map[string]int, avail []string, towel string) int {
	if v, ok := memo[towel]; ok {
		return v
	}
	acc := 0
	for _, a := range avail {
		if towel == a {
			acc += 1
		} else if strings.HasPrefix(towel, a) {
			acc += combi(memo, avail, strings.TrimPrefix(towel, a))
		}
	}
	memo[towel] = acc
	return acc
}

func parseInput(s string) (avail, towels []string) {
	parts := strings.Split(s, "\n\n")

	avail = strings.Split(parts[0], ", ")
	towels = strings.Split(parts[1], "\n")

	return avail, towels
}
