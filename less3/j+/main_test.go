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
			args{strings.NewReader(`4 7
1 4 1 2
1 4 2 3
`)},
			`2`,
			true,
		},
		{
			"2",
			args{strings.NewReader(`5 6
1 3 5 4 2
5 4 3 2 1
`)},
			`1`,
			true,
		},
		{
			"3???",
			args{strings.NewReader(`6 15
1 2 5 7 8 11
1 1 4 2 3 6
`)},
			`3`,
			true,
		},
		{
			"4",
			args{strings.NewReader(`1 10
10
10`)},
			`0`,
			true,
		},
		{
			"5",
			args{strings.NewReader(`3 10
1 2 10
1 1 10`)},
			`0`,
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
