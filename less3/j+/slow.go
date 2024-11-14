package main

import (
	"math"
	"sort"
)

func slowSolve(vasyaH int, hh, ww []int) int {
	n := len(hh)

	type item struct {
		h int
		w int
	}

	hw := make([]item, 0, n)
	for i := 0; i < n; i++ {
		hw = append(hw, item{h: hh[i], w: ww[i]})
	}

	sort.Slice(hw, func(i, j int) bool {
		return hw[i].h < hw[j].h
	})

	minmax := math.MaxInt
	for i := 0; i < n; i++ {
		width := hw[i].w
		maximum := 0
		for j := i + 1; j < n && width < vasyaH; j++ {
			width += hw[j].w
			maximum = max(maximum, hw[j].h-hw[j-1].h)
		}
		if width >= vasyaH {
			minmax = min(minmax, maximum)
		}
	}

	return minmax
}
