package main

import (
	"math"
	"sort"
)

// O(n) в среднем на рандомных данных, O(n^2) в худшем случае. см. тест 999

func fastSolve(vasyaH int, hh, ww []int) int {
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
	for i := 0; i < n; {
		width := hw[i].w
		maximum := 0
		next := i + 1
		for j := i + 1; j < n && width < vasyaH; j++ {
			width += hw[j].w
			if d := hw[j].h-hw[j-1].h; d >= maximum {
				maximum = d
				next = j
			}
		}
		if width >= vasyaH {
			minmax = min(minmax, maximum)
		}
		i = next
	}

	return minmax
}
