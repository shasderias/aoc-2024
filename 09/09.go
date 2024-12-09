package main

import (
	_ "embed"
	"fmt"
	"strconv"
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
	id := 0
	disk := []int{}

	for i, c := range s {
		l, err := strconv.Atoi(string(c))
		errPanic(err)
		for j := 0; j < l; j++ {
			if i%2 == 0 {
				disk = append(disk, id)
			} else {
				disk = append(disk, -1)
			}
		}
		if i%2 == 0 {
			id++
		}
	}

	for i := 0; i < len(disk); i++ {
		if disk[i] != -1 {
			continue
		}
		for disk[len(disk)-1] == -1 {
			disk = disk[:len(disk)-1]
		}
		disk[i] = disk[len(disk)-1]
		disk = disk[:len(disk)-1]
	}

	fmt.Println(checksum(disk))
}

type Interval struct {
	id   int
	a, b int
}

func (i Interval) len() int {
	return i.b - i.a
}

func star2(s string) {
	id, last := 0, 0
	disk := []Interval{}

	for i, c := range s {
		l, err := strconv.Atoi(string(c))
		errPanic(err)
		if i%2 == 0 {
			disk = append(disk, Interval{id, last, last + l})
			id++
		} else {
			disk = append(disk, Interval{-1, last, last + l})
		}
		last += l
	}

	id--

	for r := len(disk) - 1; r > 0; r-- {
		if disk[r].id != id {
			continue
		}
		for l := 0; l < r; l++ {
			// find free space from left
			if disk[l].id != -1 {
				continue
			}
			ll, rl := disk[l].len(), disk[r].len()
			switch {
			case ll > rl:
				diff := ll - rl
				disk[l].id = disk[r].id
				disk[l].b -= diff
				if disk[l+1].id == -1 {
					disk[l+1].a = disk[l].b
				} else {
					disk = append(disk[:l+1], append([]Interval{{-1, disk[l].b, disk[l].b + diff}}, disk[l+1:]...)...)
					r++
				}
				disk[r].id = -1
				goto next
			case disk[l].len() == disk[r].len():
				disk[l].id = disk[r].id
				disk[r].id = -1
				goto next
			case disk[r].len() < disk[l].len():
			}
		}
	next:
		id--
	}

	fmt.Println(checksum2(disk))
}

func checksum(disk []int) int {
	acc := 0
	for i, n := range disk {
		acc += i * n
	}
	return acc
}

func checksum2(disk []Interval) int {
	acc := 0
	for _, intv := range disk {
		if intv.len() == 0 || intv.id == -1 {
			continue
		}
		for i := intv.a; i < intv.b; i++ {
			acc += i * intv.id
		}

	}
	return acc
}

func printDisk(disk []Interval) {
	for _, intv := range disk {
		if intv.len() == 0 {
			continue
		}
		if intv.id == -1 {
			for i := intv.a; i < intv.b; i++ {
				fmt.Print(".")
			}
			continue
		}
		for i := intv.a; i < intv.b; i++ {
			fmt.Printf("%d", intv.id)
		}

	}
	fmt.Println("")
}

func errPanic(err error) {
	if err != nil {
		panic(err)
	}
}
