package main

import (
	"bytes"
	"fmt"
	"io"
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
			"1",
			args{strings.NewReader(`20
2 1 1
10 1 1 2 3 4 5 6 7 8 9
2 1 -1
2 1 0
4 3 4 1 1
4 4 3 1 1
4 1 1 3 4
4 1 1 4 3
4 -3 -4 -1 -1
4 -4 -3 -1 -1
4 -1 -1 -3 -4
4 -1 -1 -4 -3
4 2 1 3 3
4 3 3 1 2
4 -1 -2 -3 -3
4 -3 -3 -1 -2
5 100 100 100 100 100
5 -100 -100 -100 -100 -100
3 1 2 3
5 -1 1 0 -2 3
`)},
			`1 1
8 9
-1 1
0 1
3 4
3 4
3 4
3 4
-4 -3
-4 -3
-4 -3
-4 -3
3 3
3 3
-3 -3
-3 -3
100 100
-100 -100
2 3
1 3
`,
			true,
		},
		// {
		// 	"2",
		// 	args{strings.NewReader(``)},
		// 	``,
		// 	true,
		// },
		// {
		// 	"3",
		// 	args{strings.NewReader(``)},
		// 	``,
		// 	true,
		// },
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

func Test_solve(t *testing.T) {
	type args struct {
		aa []int
	}
	tests := []struct {
		args  args
		want  int
		want1 int
	}{
		{args{[]int{1, 1}}, 1, 1},                         // проверяем чтение (что 2 не попадает)
		{args{[]int{1, 1, 2, 3, 4, 5, 6, 7, 8, 9}}, 8, 9}, // проверяем чтение (что 10 не попадает а 9 попадает)

		{args{[]int{1, -1}}, -1, 1}, // положительный + отрицательный
		{args{[]int{1, 0}}, 0, 1},   // любой + ноль

		{args{[]int{3, 4, 1, 1}}, 3, 4}, // порядок первых двух
		{args{[]int{4, 3, 1, 1}}, 3, 4},
		{args{[]int{1, 1, 3, 4}}, 3, 4}, // порядок не первых двух
		{args{[]int{1, 1, 4, 3}}, 3, 4},

		{args{[]int{-3, -4, -1, -1}}, -4, -3}, // тоже для минимума
		{args{[]int{-4, -3, -1, -1}}, -4, -3},
		{args{[]int{-1, -1, -3, -4}}, -4, -3},
		{args{[]int{-1, -1, -4, -3}}, -4, -3},

		{args{[]int{3, 3, 1, 1}}, 3, 3}, // максимальные равны
		{args{[]int{1, 1, 3, 3}}, 3, 3},

		{args{[]int{-3, -3, -1, -1}}, -3, -3}, // минимальные равны
		{args{[]int{-1, -1, -3, -3}}, -3, -3},

		{args{[]int{1, 1, 3, 100}}, 3, 100},
		{args{[]int{-1, -1, -3, -100}}, -100, -3},

		// {args{[]int{2, 3}}, 2, 3},
		// {args{[]int{3, 2}}, 2, 3},
		// {args{[]int{2, 2}}, 2, 2},
		// {args{[]int{0, 1}}, 0, 1},
		// {args{[]int{0, 0}}, 0, 0},
		// {args{[]int{1, 100}}, 1, 100},
		// {args{[]int{100, 100}}, 100, 100},
		// {args{[]int{1, -100}}, 1, -100},
		// {args{[]int{-100, 1}}, 1, -100},

		// {args{[]int{1, 2, 3, 4}}, 3, 4},
		// {args{[]int{1, 2, 4, 3}}, 3, 4},
		// {args{[]int{1, 2, -3, -4}}, 3, 4},
		// {args{[]int{1, 2, -4, -3}}, 3, 4},
		// {args{[]int{1, 100}}, 1, 100},
		// {args{[]int{4, 2, 3, 1}}, 3, 4},
		// {args{[]int{4, 2, 1, 3}}, 3, 4},
		// {args{[]int{1, 3, 4, 2}}, 3, 4},
		// {args{[]int{-1, -2, -3, -4}}, -4, -3},
		// {args{[]int{1, 2, -3, -4}}, -4, -3},
		// {args{[]int{1, 2, 3, 3}}, 3, 3},
		// {args{[]int{2, 3, -3, -3}}, -3, -3},
		// {args{[]int{1, 2, 3, 100}}, 3, 100},
		// {args{[]int{-1, -2, -3, -100}}, -100, -3},
		// {args{[]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}}, 9, 10},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(fmt.Sprint(tt.args.aa), func(t *testing.T) {
			got, got1 := solve(tt.args.aa)
			if got != tt.want {
				t.Errorf("solve() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("solve() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
