package main

import (
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"reflect"
	"strings"
	"testing"
)

func Test_run(t *testing.T) {
	type args struct {
		in io.Reader
	}
	tests := []struct {
		name    string
		args    args
		wantOut string
		debug   bool
	}{
		{
			"1.0",
			args{strings.NewReader(`5
1 2 3 4 5
5 4 3 2 1
1 0 1 0 0
`)},
			`1 5 2 4 3`,
			true,
		},
		{
			"1.1",
			args{strings.NewReader(`5
5 4 3 2 1
1 2 3 4 5
0 1 0 1 1
`)},
			`1 5 2 4 3`,
			true,
		},
		{
			"2.0",
			args{strings.NewReader(`6
3 10 6 2 10 1
3 5 10 7 5 9
0 0 1 1 0 1
`)},
			`2 5 3 6 1 4`,
			true,
		},
		{
			"2.1",
			args{strings.NewReader(`6
3 5 10 7 5 9
3 10 6 2 10 1
1 1 0 0 1 0
`)},
			`2 5 3 6 1 4`,
			true,
		},
		{
			"3.0",
			args{strings.NewReader(`4
1 1 1 1
1 2 3 4
0 0 0 0`)},
			`4 3 2 1`,
			true,
		},
		{
			"3.0",
			args{strings.NewReader(`4
1 2 3 4
1 1 1 1
1 1 1 1`)},
			`4 3 2 1`,
			true,
		},
		{
			"4.0",
			args{strings.NewReader(`4
1 1 1 1
2 2 2 2
0 0 0 0`)},
			`1 2 3 4`,
			true,
		},
		{
			"4.1",
			args{strings.NewReader(`4
1 1 1 1
2 2 2 2
1 1 1 1`)},
			`1 2 3 4`,
			true,
		},
		{
			"5.0",
			args{strings.NewReader(`1
1
2
0`)},
			`1`,
			true,
		},
		{
			"5.1",
			args{strings.NewReader(`1
1
2
1`)},
			`1`,
			true,
		},
		{
			"6.0",
			args{strings.NewReader(`3
1 2 3
4 4 4
0 0 0`)},
			`3 2 1`,
			true,
		},
		{
			"6.1",
			args{strings.NewReader(`3
1 2 3
4 4 4
1 1 1`)},
			`3 2 1`,
			true,
		},
		{
			"6.2",
			args{strings.NewReader(`3
1 1 1
4 4 4
1 1 1`)},
			`1 2 3`,
			true,
		},
		{
			"37",
			args{strings.NewReader(`5
10 6 1 1 4
1 6 2 8 2
0 0 0 0 0`)},
			`1 2 5 4 3`,
			true,
		},
		{
			"37.1",
			args{strings.NewReader(`5
1 6 2 8 2
10 6 1 1 4
1 1 1 1 1`)},
			`1 2 5 4 3`,
			true,
		},

		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func(v bool) { debugEnable = v }(debugEnable)
			debugEnable = tt.debug
			out := &bytes.Buffer{}
			run(tt.args.in, out)
			if gotOut := out.String(); trimLines(gotOut) != trimLines(tt.wantOut) {
				t.Errorf("run() = %v, want %v", gotOut, tt.wantOut)
			}
		})
	}
}

func trimLines(text string) string {
	lines := strings.Split(text, "\n")
	for i, line := range lines {
		lines[i] = strings.TrimRight(line, " \t\r\n")
	}
	for n := len(lines); n > 0 && lines[n-1] == ""; n-- {
		lines = lines[:n-1]
	}
	return strings.Join(lines, "\n")
}

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
