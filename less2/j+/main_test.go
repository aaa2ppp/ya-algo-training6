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
			args{strings.NewReader(`6
3 3 3 4 4 5
4 2
3 4 5 6
`)},
			`1 1 2 2 
`,
			true,
		},
		{
			"2",
			args{strings.NewReader(`7
1 5 7 2 10 10 6
7 0
1 2 3 4 5 6 7
`)},
			`1 1 1 4 4 6 7 
`,
			true,
		},
		{
			"3",
			args{strings.NewReader(`1
10
7 100
1 1 1 1 1 1 1`)},
			`1 1 1 1 1 1 1`,
			true,
		},
		{
			"4",
			args{strings.NewReader(`7
10 10 10 10 10 10 10
7 100
1 1 1 1 1 1 1`)},
			`1 1 1 1 1 1 1`,
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
