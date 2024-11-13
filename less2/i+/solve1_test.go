package main

import (
	"fmt"
	"math/rand"
	"reflect"
	"testing"
)

func generate(seed int64) (aa, bb, pp []int) {
	rand := rand.New(rand.NewSource(seed))
	n := rand.Intn(10) + 1
	aa = make([]int, n)
	bb = make([]int, n)
	pp = make([]int, n)
	for i := 0; i < n; i++ {
		aa[i] = rand.Intn(10) + 1
		bb[i] = rand.Intn(10) + 1
		pp[i] = rand.Intn(2)
	}
	return
}

func Test_solve(t *testing.T) {
	for i := 1; i < 1000; i++ {
		aa, bb, pp := generate(int64(i))
		t.Run(fmt.Sprintf("%d: %v %v %v", i, aa, bb, pp), func(t *testing.T) {
			res1 := solve1(aa, bb, pp)
			res2 := solve2(aa, bb, pp)
			if !reflect.DeepEqual(res1, res2) {
				t.Fatalf("res1:%v not equal res2:%v", res1, res2)
			}
		})
	}
}
