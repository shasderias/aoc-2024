package main

import (
	_ "embed"
	"fmt"
	"strconv"
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
	star1(input)
	star2(example2)
	star2(input)
}

func star1(s string) {
	input := parseInput(s)

	acc := 0
	for _, in := range input {
		n := in
		for i := 0; i < 2000; i++ {
			n = evolve(n)
		}
		acc += n
	}
	fmt.Println(acc)
}

func star2(s string) {
	input := parseInput(s)

	allMonkeyPrices := []map[string]int{}
	for _, in := range input {
		var (
			n            = in
			prices       = []int{}
			changes      = []int{}
			monkeyPrices = map[string]int{}
		)

		for i := 0; i < 2000; i++ {
			price := n % 10
			prices = append(prices, price)
			if len(prices) > 1 {
				changes = append(changes, prices[len(prices)-1]-prices[len(prices)-2])
			}
			if len(changes) > 3 {
				key := fmt.Sprintf("%v", changes[len(changes)-4:])
				if _, ok := monkeyPrices[key]; !ok {
					monkeyPrices[key] = price
				}
			}
			n = evolve(n)
		}
		allMonkeyPrices = append(allMonkeyPrices, monkeyPrices)
	}

	allPossibleSeqs := map[string]bool{}

	for _, mp := range allMonkeyPrices {
		for k := range mp {
			allPossibleSeqs[k] = true
		}
	}

	maxPrice := -1
	maxSeq := ""
	for seq := range allPossibleSeqs {

		acc := 0
		for _, mp := range allMonkeyPrices {
			acc += mp[seq]
		}

		if acc > maxPrice {
			maxPrice = acc
			maxSeq = seq
		}
	}

	fmt.Println(maxPrice, maxSeq)

}

type seqPrice struct {
}

func parseInput(s string) (o []int) {
	for _, l := range strings.Split(s, "\n") {
		n, err := strconv.Atoi(l)
		if err != nil {
			panic(err)
		}
		o = append(o, n)
	}
	return o
}

func evolve(n int) int {
	s1 := ((n * 64) ^ n) % 16777216
	s2 := ((s1 / 32) ^ s1) % 16777216
	s3 := ((s2 * 2048) ^ s2) % 16777216
	return s3
}
