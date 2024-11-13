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
			args{strings.NewReader(`1+(2*2 - 3)
`)},
			`2`,
			true,
		},
		{
			"2",
			args{strings.NewReader(`1+a+1
`)},
			`WRONG`,
			true,
		},
		{
			"3",
			args{strings.NewReader(`1 1 + 2
`)},
			`WRONG`,
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

func Test_scanIntStr(t *testing.T) {
	tests := []struct {
		name  string
		want  string
		want1 bool
	}{
		{
			"1",
			"1",
			true,
		},
		{
			"1*",
			"1",
			true,
		},
		{
			"12",
			"12",
			true,
		},
		{
			"12*",
			"12",
			true,
		},
		{
			"-",
			"",
			false,
		},
		{
			"-*",
			"",
			false,
		},
		{
			"-1",
			"-1",
			true,
		},
		{
			"+1",
			"+1",
			true,
		},
		{
			"*1",
			"",
			false,
		},
		{
			"1+(",
			"1",
			true,
		},

		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := scanIntStr(strings.NewReader(tt.name))
			if got != tt.want {
				t.Errorf("scanIntStr() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("scanIntStr() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
