package main

import (
	"bytes"
	"fmt"
	"io"
	"maps"
	"math/rand"
	"strconv"
	"strings"
	"testing"
)

func Test_run_solve(t *testing.T) {
	run_test(t, solve)
}

func Test_run_slowSolve(t *testing.T) {
	run_test(t, slowSolve)
}

func run_test(t *testing.T, solve func(edges EdgeSet, graph Graph) int) {
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
			args{strings.NewReader(`4
1 2
2 3
3 4
`)},
			`1`,
			true,
		},
		{
			"2",
			args{strings.NewReader(`7
1 2
1 3
1 4
1 5
1 6
1 7
`)},
			`0`,
			true,
		},
		{
			"3",
			args{strings.NewReader(`6
1 2
2 3
2 4
5 4
6 4
`)},
			`4`,
			true,
		},
		// {
		// 	"4",
		// 	args{strings.NewReader(``)},
		// 	``,
		// 	true,
		// },
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

func Test_solve(t *testing.T) {
	for i := 0; i < 100; i++ {
		seed := int64(i)
		rand := rand.New(rand.NewSource(seed))
		edges, graph := genTree(rand, 12)

		t.Run(strconv.Itoa(i), func(t *testing.T) {
			want := slowSolve(maps.Clone(edges), cloneGraph(graph))
			got := solve(maps.Clone(edges), cloneGraph(graph))
			if got != want {
				t.Errorf("solve()=%d, want %d", got, want)
				for edge := range edges {
					fmt.Println(edge[0], edge[1])
				}
			} else {
				t.Log(got)
			}
		})
	}
}

func Fuzz_solve(f *testing.F) {
	for i := 1; i <= 10; i++ {
		f.Add(int64(i))
	}
	f.Fuzz(func(t *testing.T, seed int64) {
		rand := rand.New(rand.NewSource(seed))
		n := rand.Intn(1000) + 2
		edges, graph := genTree(rand, n)
		want := slowSolve(maps.Clone(edges), cloneGraph(graph))
		got := solve(maps.Clone(edges), cloneGraph(graph))
		if got != want {
			t.Errorf("solve()=%d, want %d", got, want)
			for edge := range edges {
				fmt.Println(edge[0], edge[1])
			}
		}
	})
}

func Benchmark_slowSolve1K(b *testing.B) {
	rand := rand.New(rand.NewSource(1))
	edges, graph := genTree(rand, 1000)
	for i := 0; i < b.N; i++ {
		slowSolve(maps.Clone(edges), cloneGraph(graph))
	}
}

func Benchmark_solve100K(b *testing.B) {
	rand := rand.New(rand.NewSource(1))
	edges, graph := genTree(rand, 100000)
	for i := 0; i < b.N; i++ {
		solve(maps.Clone(edges), cloneGraph(graph))
	}
}

func Benchmark_solve100KLine(b *testing.B) {
	edges, graph := genTreeLine(100000)
	for i := 0; i < b.N; i++ {
		solve(maps.Clone(edges), cloneGraph(graph))
	}
}
