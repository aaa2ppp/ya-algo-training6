package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"
)

const test_data = "./test_data/"

func Test_run_solenSolve(t *testing.T) {
	run_solve(t, stolenSolve)
}

func Test_run_solve(t *testing.T) {
	run_solve(t, solve)
}

func Test_solve100(t *testing.T) {
	for i := 0; i < 100; i++ {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			rand := rand.New(rand.NewSource(int64(i)))
			n := rand.Intn(100) + 1
			edges := genEdges(rand, n)
			want := stolenSolve(edges)
			got := solve(edges)
			if got != want {
				t.Errorf("solve() = %d, want %d", got, want)
			}
		})
	}
}

func Test_solve1000(t *testing.T) {
	for i := 0; i < 100; i++ {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			rand := rand.New(rand.NewSource(int64(i)))
			n := rand.Intn(1000) + 1
			edges := genEdges(rand, n)
			want := stolenSolve(edges)
			got := solve(edges)
			if got != want {
				t.Errorf("solve() = %d, want %d", got, want)
			}
		})
	}
}

func Test_solve3000(t *testing.T) {
	for i := 0; i < 30; i++ {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			rand := rand.New(rand.NewSource(int64(i)))
			n := rand.Intn(3000) + 1
			edges := genEdges(rand, n)
			want := stolenSolve(edges)
			got := solve(edges)
			if got != want {
				t.Errorf("solve() = %d, want %d", got, want)
			}
		})
	}
}

var (
	bench_res      int
	bench_edges    = genEdges(rand.New(rand.NewSource(1)), 3000)
	bench_edges_69 = func() []Edge {
		buf := readTestFile("69")
		sc := bufio.NewScanner(bytes.NewReader(buf))
		sc.Split(bufio.ScanWords)

		n, err := scanInt(sc)
		if err != nil {
			panic(err)
		}
		edges := make([]Edge, n-1)
		for i := range edges {
			a, b, err := scanTwoIntX[Idx](sc)
			if err != nil {
				panic(err)
			}
			edges[i] = Edge{a, b}
		}
		return edges
	}()
)

func Benchmark_solve(b *testing.B) {
	b.Run("stolenSolve", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			bench_res = stolenSolve(bench_edges)
		}
	})
	b.Run("solve", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			bench_res = solve(bench_edges)
		}
	})
}

func Benchmark_solve_69(b *testing.B) {
	b.Run("solve_69", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			bench_res = solve(bench_edges_69)
		}
	})
}

func readTestFile(filename string) []byte {
	buf, err := os.ReadFile(test_data + filename)
	if err != nil {
		panic(err)
	}
	return buf
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
		{
			"6 ?",
			args{strings.NewReader(`6
1 3
2 3
2 4
3 5
3 6
`)},
			`18`,
			true,
		},
		{
			"7 ?",
			args{strings.NewReader(`8
1 3
2 3
2 4
3 5
3 6
6 7
6 8
`)},
			`104`,
			true,
		},
		{
			"8 ?",
			args{strings.NewReader(`8
4 5
6 7
1 4
1 3
2 8
8 1
1 6`)},
			`30`,
			false,
		},
		{
			"9 ?",
			args{strings.NewReader(`8
3 7
2 8
2 6
5 3
4 1
6 4
4 5`)},
			`28`,
			false,
		},
		{
			"10 ?",
			args{strings.NewReader(`20
5 7
12 13
17 16
11 3
4 10
18 6
19 7
14 13
1 14
18 10
2 11
9 7
8 18
15 17
11 20
8 1
20 8
9 17
17 20`)},
			`464151755`,
			false,
		},
		{
			"42 ?",
			args{strings.NewReader(`8
8 7
3 5
5 2
6 7
7 3
1 2
4 7`)},
			`42`,
			false,
		},
		// {
		// 	// Can't stolenSolve TL
		// 	"69",
		// 	args{bytes.NewReader(readTestFile("69"))},
		// 	unsafeString(readTestFile("69.a")),
		// 	false,
		// },
		{
			"999",
			args{bytes.NewReader(readTestFile("999"))},
			unsafeString(readTestFile("999.a")),
			false,
		},
		{
			"998",
			args{bytes.NewReader(readTestFile("998"))},
			unsafeString(readTestFile("998.a")),
			false,
		},
		{
			"997",
			args{bytes.NewReader(readTestFile("997"))},
			unsafeString(readTestFile("997.a")),
			false,
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

func Test_gentest(t *testing.T) {
	edges := genEdges(rand.New(rand.NewSource(time.Now().UnixNano())), 8)
	fmt.Println(len(edges) + 1)
	for _, e := range edges {
		fmt.Println(e[0], e[1])
	}
	res := stolenSolve(edges)
	fmt.Println("a:", res)
}

func Test_gentest42(t *testing.T) {
	for {
		edges := genEdges(rand.New(rand.NewSource(time.Now().UnixNano())), 8)
		res := stolenSolve(edges)
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

// 0     1
// 1    1 1
// 2   1 2 1
// 3  1 3 3 1
// 4 1 4 6 4 1
func Test_paskal(t *testing.T) {
	type args struct {
		i int
		j int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			"0 0",
			args{0, 0},
			1,
		},
		{
			"1 0",
			args{1, 0},
			1,
		},
		{
			"1 1",
			args{1, 1},
			1,
		},
		{
			"2 0",
			args{2, 0},
			1,
		},
		{
			"2 1",
			args{2, 1},
			2,
		},
		{
			"2 2",
			args{2, 2},
			1,
		},
		{
			"3 0",
			args{3, 0},
			1,
		},
		{
			"3 1",
			args{3, 1},
			3,
		},
		{
			"3 2",
			args{3, 2},
			3,
		},
		{
			"3 3",
			args{3, 3},
			1,
		},
		{
			"4 0",
			args{4, 0},
			1,
		},
		{
			"4 1",
			args{4, 1},
			4,
		},
		{
			"4 2",
			args{4, 2},
			6,
		},
		{
			"4 3",
			args{4, 3},
			4,
		},
		{
			"4 4",
			args{4, 4},
			1,
		},

		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := paskal(tt.args.i, tt.args.j); got != tt.want {
				t.Errorf("paskal() = %v, want %v", got, tt.want)
				for i := range _paskal {
					log.Printf("%2d %v", i, _paskal[i])
				}
			}
		})
	}
}

func Test_paskal2(t *testing.T) {
	type args struct {
		i int
		j int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			"0 0",
			args{0, 0},
			1,
		},
		{
			"1 0",
			args{1, 0},
			1,
		},
		{
			"1 1",
			args{1, 1},
			1,
		},
		{
			"2 0",
			args{2, 0},
			1,
		},
		{
			"2 1",
			args{2, 1},
			2,
		},
		{
			"2 2",
			args{2, 2},
			1,
		},
		{
			"3 0",
			args{3, 0},
			1,
		},
		{
			"3 1",
			args{3, 1},
			3,
		},
		{
			"3 2",
			args{3, 2},
			3,
		},
		{
			"3 3",
			args{3, 3},
			1,
		},
		{
			"4 0",
			args{4, 0},
			1,
		},
		{
			"4 1",
			args{4, 1},
			4,
		},
		{
			"4 2",
			args{4, 2},
			6,
		},
		{
			"4 3",
			args{4, 3},
			4,
		},
		{
			"4 4",
			args{4, 4},
			1,
		},

		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := paskal2(tt.args.i, tt.args.j); got != tt.want {
				t.Errorf("paskal2() = %v, want %v", got, tt.want)
			}
		})
	}
}
