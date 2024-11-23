package main

import (
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"strings"
	"testing"
	"time"
)

func Test_run_slowSolve(t *testing.T) {
	run_solve(t, slowSolve)
}

func Test_run_solve(t *testing.T) {
	run_solve(t, solve)
}

func run_solve(t *testing.T, solve SolveFunc) {
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
			args{strings.NewReader(`3
1 2
1 3
`)},
			`2`,
			true,
		},
		{
			"2",
			args{strings.NewReader(`4
2 3
3 1
2 4
`)},
			`3`,
			true,
		},
		{
			"3",
			args{strings.NewReader(`10
1 2
2 3
3 4
4 5
5 6
6 7
7 8
8 9
9 10`)},
			`1`,
			true,
		},
		{
			"4",
			args{strings.NewReader(`1`)},
			`1`,
			true,
		},
		{
			"5", // sam https://t.me/c/1322266617/201496/287700
			args{strings.NewReader(`4
1 2
1 3
4 3`)},
			`5`,
			true,
		},
		// 		{
		// 			"6 ?",
		// 			args{strings.NewReader(`8
		// 4 5
		// 6 7
		// 1 4
		// 1 3
		// 2 8
		// 8 1
		// 1 6`)},
		// 			`48`,
		// 			false,
		// 		},
		// 		{
		// 			"7 ?",
		// 			args{strings.NewReader(`8
		// 3 7
		// 2 8
		// 2 6
		// 5 3
		// 4 1
		// 6 4
		// 4 5`)},
		// 			`42`,
		// 			false,
		// 		},
		// 		{
		// 			"8 ?",
		// 			args{strings.NewReader(`20
		// 5 7
		// 12 13
		// 17 16
		// 11 3
		// 4 10
		// 18 6
		// 19 7
		// 14 13
		// 1 14
		// 18 10
		// 2 11
		// 9 7
		// 8 18
		// 15 17
		// 11 20
		// 8 1
		// 20 8
		// 9 17
		// 17 20`)},
		// 			`247130905`,
		// 			false,
		// 		},
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

func Test_gentest(t *testing.T) {
	edges, graph := genTree(rand.New(rand.NewSource(time.Now().UnixNano())), 20)
	fmt.Println(len(edges) + 1)
	for _, e := range edges {
		fmt.Println(e[0], e[1])
	}

	fmt.Println("a:", slowSolve(edges, graph))
}
func Test_gentest42(t *testing.T) {
	for {
		edges, graph := genTree(rand.New(rand.NewSource(time.Now().UnixNano())), 8)
		res := slowSolve(edges, graph)
		if res == 42 {
			fmt.Println(len(edges) + 1)
			for _, e := range edges {
				fmt.Println(e[0], e[1])
			}
			fmt.Println("a:", res)
			break
		}
	}
}

var (
	bench_res                int
	bench_edges, bench_graph = genTree(rand.New(rand.NewSource(1)), 20)
)

func Benchmark_slowSolve(b *testing.B) {
	for i := 0; i < b.N; i++ {
		bench_res = slowSolve(bench_edges, bench_graph)
	}
}
