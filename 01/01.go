package main

import (
	"bufio"
	"bytes"
	_ "embed"
	"fmt"
	"io"
	"slices"
)

const example = `3   4
4   3
2   5
1   3
3   9
3   3
`

//go:embed input.txt
var input string

func main() {
	task1(bytes.NewBufferString(example))
	task1(bytes.NewBufferString(input))
	task2(bytes.NewBufferString(example))
	task2(bytes.NewBufferString(input))
}

func task1(input io.Reader) {
	var l1, l2 []int

	s := bufio.NewScanner(input)

	for s.Scan() {
		var n1, n2 int

		fmt.Sscanf(s.Text(), "%d   %d", &n1, &n2)

		l1 = append(l1, n1)
		l2 = append(l2, n2)
	}

	slices.Sort(l1)
	slices.Sort(l2)

	totalDist := 0

	for i := 0; i < len(l1); i++ {
		totalDist += abs(l1[i] - l2[i])
	}

	fmt.Println(totalDist)
}

func task2(input io.Reader) {
	var l1, l2 []int

	s := bufio.NewScanner(input)

	for s.Scan() {
		var n1, n2 int

		fmt.Sscanf(s.Text(), "%d   %d", &n1, &n2)

		l1 = append(l1, n1)
		l2 = append(l2, n2)
	}

	freqMap := make(map[int]int)

	for _, n := range l2 {
		freqMap[n]++
	}

	totalSim := 0

	for _, n := range l1 {
		totalSim += n * freqMap[n]
	}

	fmt.Println(totalSim)
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}
