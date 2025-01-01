package main

import (
	_ "embed"
	"fmt"
	"slices"
	"strings"
)

//go:embed example.txt
var example string

//go:embed 	input.txt
var input string

func main() {
	star1(example)
	star1(input)
	star2(example)
	star2(input)
}

type Set map[string]struct{}

func NewSet(v ...string) Set {
	s := Set{}
	for _, vv := range v {
		s[vv] = struct{}{}
	}
	return s
}

func (s Set) String() string {
	return strings.Join(s.Slice(), "-")
}

func (s Set) StringAOC() string {
	return strings.Join(s.Slice(), ",")
}

func (s Set) Slice() []string {
	slc := make([]string, 0, len(s))
	for v := range s {
		slc = append(slc, v)
	}
	slices.Sort(slc)
	return slc
}

func (s1 Set) Intersect(s2 Set) Set {
	res := Set{}
	for v := range s1 {
		if _, ok := s2[v]; ok {
			res[v] = struct{}{}
		}
	}
	return res

}

func (s1 Set) Remove(v string) Set {
	s2 := Set{}
	for vv := range s1 {
		if vv != v {
			s2[vv] = struct{}{}
		}
	}
	return s2
}

func (s1 Set) Add(v string) Set {
	s2 := Set{}
	for vv := range s1 {
		s2[vv] = struct{}{}
	}
	s2[v] = struct{}{}
	return s2
}

func (s1 Set) Union(s2 Set) Set {
	res := Set{}
	for v := range s1 {
		res[v] = struct{}{}
	}
	for v := range s2 {
		res[v] = struct{}{}
	}
	return res
}

func star1(s string) {
	n := parseInput(s)
	found := map[string]Set{}

	for _, c := range n.computers {
		for _, l1 := range n.links[c] {
			for _, l2 := range n.links[l1] {
				if n.Connected(c, l2) && n.Connected(c, l1) && n.Connected(l1, l2) {
					set := NewSet(c, l1, l2)
					found[set.String()] = set
				}
			}
		}
	}

	acc := 0
	for _, f := range found {
		for c := range f {
			if strings.HasPrefix(c, "t") {
				acc++
				break
			}
		}
	}
	fmt.Println(acc)
}

func star2(s string) {
	n := parseInput(s)

	cliques := map[string]Set{}
	BronKerbosch(nil, NewSet(n.computers...), nil, n, cliques)
	maxLen := 0
	for _, c := range cliques {
		if len(c) > maxLen {
			fmt.Println(c.StringAOC())
			maxLen = len(c)
		}
	}
}

func BronKerbosch(r Set, p Set, x Set, n Network, cliques map[string]Set) {
	if len(p) == 0 && len(x) == 0 {
		cliques[r.String()] = r
		return
	}

	var (
		pivotCandidates = p.Union(x)
		pivot           string
		pivotMax        int
	)
	for c := range pivotCandidates {
		if len(n.links[c]) > pivotMax {
			pivotMax = len(n.links[c])
			pivot = c
		}
	}

	for u := range p {
		if !n.Connected(u, pivot) {
			BronKerbosch(r.Add(u), p.Intersect(NewSet(n.links[u]...)), x.Intersect(NewSet(n.links[u]...)), n, cliques)
			p = p.Remove(u)
			x = x.Add(u)
		}
	}
}

type Network struct {
	links     map[string][]string
	computers []string
}

func (n Network) Connected(a, b string) bool {
	return slices.Contains(n.links[a], b)
}

func parseInput(s string) Network {
	links := map[string][]string{}
	computers := []string{}

	for _, line := range strings.Split(s, "\n") {
		conn := strings.Split(line, "-")
		links[conn[0]] = append(links[conn[0]], conn[1])
		links[conn[1]] = append(links[conn[1]], conn[0])
		computers = append(computers, conn[0])
	}

	return Network{links, computers}
}

func parseInput2(s string) string {
	links := [][2]string{}
	computers := []string{}

	for _, line := range strings.Split(s, "\n") {
		conn := strings.Split(line, "-")
		links = append(links, [2]string{conn[0], conn[1]})
		computers = append(computers, conn[0])
	}

	buf := strings.Builder{}

	buf.WriteString("[")

	for _, c := range computers {
		buf.WriteString(fmt.Sprintf("{ data: { id: '%s'}},", c))
	}
	for _, l := range links {
		buf.WriteString(fmt.Sprintf("{ data: { id: '%s-%s', source: '%s', target: '%s'}},", l[0], l[1], l[0], l[1]))
	}

	buf.WriteString("]")

	return buf.String()
}
