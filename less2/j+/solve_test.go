package main

import (
	"fmt"
	"math/rand"
	"reflect"
	"slices"
	"testing"
)

func generate(seed int64) (aa, xx []int, k int) {
	rand := rand.New(rand.NewSource(seed))

	n := rand.Intn(10) + 1
	aa = make([]int, 0, n)

	for i := 0; i < n; i++ {
		aa = append(aa, rand.Intn(10)+1)
	}

	m := rand.Intn(10) + 1
	xx = make([]int, 0, m)

	for i := 0; i < m; i++ {
		xx = append(xx, rand.Intn(n)+1)
	}

	return
}

func Test_solve(t *testing.T) {
	for i := 1; i < 5000; i++ {
		aa, xx, k := generate(int64(i))
		t.Run(fmt.Sprintf("%d: %v %v %v", i, aa, xx, k), func(t *testing.T) {
			res1 := slowSolve(slices.Clone(aa), slices.Clone(xx), k)
			res2 := solve(aa, xx, k)
			if !reflect.DeepEqual(res1, res2) {
				t.Fatalf("slove:%v not equal clow solve:%v", res2, res1)
			}
		})
	}
}
