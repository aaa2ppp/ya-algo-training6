package main

import (
	"bytes"
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
			"4 1",
			args{4, 1},
			4,
		},
		{
			"4 2",
			args{4, 2},
			6,
		},

		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := paskal(tt.args.i, tt.args.j); got != tt.want {
				t.Errorf("paskal() = %v, want %v", got, tt.want)
			}
		})
	}
}
