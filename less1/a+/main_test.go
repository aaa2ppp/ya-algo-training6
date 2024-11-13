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
			args{strings.NewReader(`-1
-2
5
3
-4
6
`)},
			`NW`,
			true,
		},
		{
			"2",
			args{strings.NewReader(`1 1 3 3 2 0`)},
			`S`,
			true,
		},
		{
			"3",
			args{strings.NewReader(`1 1 3 3 2 4`)},
			`N`,
			true,
		},
		{
			"4",
			args{strings.NewReader(`1 1 3 3 0 2`)},
			`W`,
			true,
		},
		{
			"5",
			args{strings.NewReader(`1 1 3 3 4 2`)},
			`E`,
			true,
		},
		{
			"6",
			args{strings.NewReader(`1 1 3 3 0 0`)},
			`SW`,
			true,
		},
		{
			"7",
			args{strings.NewReader(`1 1 3 3 0 4`)},
			`NW`,
			true,
		},
		{
			"8",
			args{strings.NewReader(`1 1 3 3 4 0`)},
			`SE`,
			true,
		},
		{
			"9",
			args{strings.NewReader(`1 1 3 3 4 4`)},
			`NE`,
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
