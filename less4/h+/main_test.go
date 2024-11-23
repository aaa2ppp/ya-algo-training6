package main

import (
	"bytes"
	"io"
	"math/rand"
	"strings"
	"testing"
)

func Test_run_solve(t *testing.T) {
	run_solve(t, solve)
}

// func Test_run_slowSolve(t *testing.T) {
// 	run_solve(t, slowSolve)
// }

func run_solve(t *testing.T, solve func([][]int32, []int32) (int, []int32)) {
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
			"1",
			args{strings.NewReader(`6
1 2
2 3
1 4
4 5
4 6
22 48 2 2 8 1
`)},
			`26 3
1 3 4 `,
			true,
		},
		{
			"2",
			args{strings.NewReader(`2
1 2
1 100`)},
			`1 1
1`,
			true,
		},
		{
			"3",
			args{strings.NewReader(`1

100`)},
			`100 1
1`,
			true,
		},
		{
			"4",
			args{strings.NewReader(`3
1 2
1 3
10 20 30`)},
			`10 1
1`,
			true,
		},
		{
			"4.2",
			args{strings.NewReader(`3
1 2
1 3
20 10 30`)},
			`20 1
1`,
			true,
		},
		{
			"5",
			args{strings.NewReader(`3
1 2
1 3
30 10 5`)},
			`15 2
2 3`,
			true,
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func(v bool) { debugEnable = v }(debugEnable)
			debugEnable = tt.debug
			out := &bytes.Buffer{}
			run(tt.args.in, out, solve)
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

func genTestN(seed int64, n int) ([][]int, []int) {
	rand := rand.New(rand.NewSource(seed))
	graph := make([][]int, n+1)
	aa := make([]int, n+1)
	for i := n; i > 1; i-- {
		a := i
		b := rand.Intn(i-1) + 1
		graph[a] = append(graph[a], b)
		graph[b] = append(graph[b], a)
		aa[i] = rand.Intn(100)
	}
	return graph, aa
}
