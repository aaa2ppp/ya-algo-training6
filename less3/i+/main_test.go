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
			args{strings.NewReader(`4
1 3
1 1
3 1
2 1
2 2
`)},
			`1
1
2
3
`,
			true,
		},
		{
			"2",
			args{strings.NewReader(`4
1 2
1 1
2 1
3 1
4 1
`)},
			`1
2
3
4
`,
			true,
		},
		{
			"3",
			args{strings.NewReader(`1
1 4
1 1
`)},
			`1
`,
			true,
		},
		{
			"4",
			args{strings.NewReader(`3
1 3
1 1
1 3
1 10`)},
			`1
3
10`,
			true,
		},
		{
			"5?",
			args{strings.NewReader(`25
1 2
1 1
4 2
1 2
4 3
4 4
3 4
1 5
4 5
3 5
1 6
4 6
2 6
3 6
2 7
1 8
2 8
1 9
2 9
3 9
1 10
2 10
3 10
4 10
2 11
3 11
`)},
			`1
3
2
21
22
4
5
23
16
6
24
7
17
11
8
12
9
13
18
10
14
19
25
15
20
`,
			true,
		},
		{
			"4",
			args{strings.NewReader(`10
1 3
1 3
4 2
1 1
1 2
2 4
3 3
1 5
2 5
3 5
4 5
`)},
			`3
4
1
2
4
3
5
6
5
6`,
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
