package main

import "testing"

func TestEquationSolve(t *testing.T) {

	testCases := []struct {
		eq   Equation
		ops  []Operator
		want bool
	}{
		{
			Equation{
				Val:  190,
				Nums: []int{10, 19},
			},
			[]Operator{OpMul},
			true,
		},
		{
			Equation{
				Val:  156,
				Nums: []int{15, 6},
			},
			[]Operator{OpCat},
			true,
		},
	}

	for _, tt := range testCases {
		if got := tt.eq.Solve(tt.ops); got != tt.want {
			t.Errorf("Equation{%v}.Solve(%v) = %v, want %v", tt.eq, tt.ops, got, tt.want)
		}
	}
}
