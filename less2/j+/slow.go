package main

import "fmt"

func slowSolve(aa, xx []int, k int) []int {
	// n := len(aa)
	m := len(xx)

	if k < 0 {
		panic(fmt.Errorf("k must be >= 0, got %v", k))
	}

	res := make([]int, 0, m)

	for _, i := range xx {

		k := k
		i := i - 1 // to 0-indexing

		for ; i > 0; i-- {

			if aa[i] < aa[i-1] {
				break
			}

			if aa[i] == aa[i-1] {
				if k == 0 {
					break
				}
				k--
			}
		}

		res = append(res, i+1) // to 1-indexing
	}

	return res
}
