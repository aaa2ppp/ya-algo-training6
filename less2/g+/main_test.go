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
			"1:aab,1",
			args{strings.NewReader(`3 1
aab
`)},
			`2`,
			true,
		},
		{
			"2:aabcbb,2",
			args{strings.NewReader(`6 2
aabcbb
`)},
			`4`,
			true,
		},
		{
			"3,a,0",
			args{strings.NewReader(`1 0
a`)},
			`1`,
			true,
		},
		{
			"4,ab,0",
			args{strings.NewReader(`2 0
ab`)},
			`1`,
			true,
		},
		{
			"5,aa,0",
			args{strings.NewReader(`2 0
aa`)},
			`2`,
			true,
		},
		{
			"16",
			args{strings.NewReader(`50 159
aaabaababaabaaaabaabbxaazbaaaababaababbbaaaabaabbb
`)},
			`39`, // 35 - не верно
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
