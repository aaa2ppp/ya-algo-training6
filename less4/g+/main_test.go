package main

import (
	"bytes"
	"io"
	"math/rand"
	"strconv"
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
			args{strings.NewReader(`3 2 10
1 2
1 3
`)},
			`4`,
			true,
		},
		{
			"2",
			args{strings.NewReader(`4 4 17
1 2
1 3
4 2
3 4
`)},
			`0`,
			true,
		},
		{
			"3",
			args{strings.NewReader(`4 2 13
1 2
2 3
`)},
			`7`,
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

func Benchmark_fact(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_fact = _fact[:2]
		fact(1000000, int(1e9)+7)
	}
}

// slow-slow
func Benchmark_paskal(b *testing.B) {
	for i := 0; i < b.N; i++ {
		paskal(1e4, 5e3, 1e9+7)
	}
}

func Benchmark_paskal2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_fact = _fact[:2]
		paskal2(1e6, 5e5, 1e9+7)
	}
}

func Test_paskal2(t *testing.T) {
	rand := rand.New(rand.NewSource(1))
	for i := 0; i < 100; i++ {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			m := rand.Intn(1000) + 1
			n := rand.Intn(1000) + 1
			if n > m {
				m, n = n, m
			}
			if got, want := paskal2(m, n, 1e9+7), paskal(m, n, 1e9+7); got != want {
				t.Errorf("paskal2(%d, %d) = %d, want %d", m, n, got, want)
			}
		})
	}
}

func Test_paskal3(t *testing.T) {
	rand := rand.New(rand.NewSource(1))
	for i := 0; i < 100; i++ {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			m := rand.Intn(1000) + 1
			n := rand.Intn(1000) + 1
			if n > m {
				m, n = n, m
			}
			if got, want := paskal3(m, n, 1e9+7), paskal(m, n, 1e9+7); got != want {
				t.Errorf("paskal3(%d, %d) = %d, want %d", m, n, got, want)
			}
		})
	}
}

func Test_paskal4(t *testing.T) {
	rand := rand.New(rand.NewSource(1))
	for i := 0; i < 100; i++ {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			m := rand.Intn(1000) + 1
			n := rand.Intn(1000) + 1
			if n > m {
				m, n = n, m
			}
			if got, want := paskal4(m, n, 1e9+7), paskal(m, n, 1e9+7); got != want {
				t.Errorf("paskal4(%d, %d) = %d, want %d", m, n, got, want)
			}
		})
	}
}

func Test_paskal5(t *testing.T) {
	defer func(v bool) { debugEnable = v }(debugEnable)
	debugEnable = true

	rand := rand.New(rand.NewSource(1))
	for i := 0; i < 3; i++ {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			m := rand.Intn(10) + 1
			n := rand.Intn(10) + 1
			if n > m {
				m, n = n, m
			}
			if got, want := paskal5(m, n, 1e9+7), paskal(m, n, 1e9+7); got != want {
				t.Fatalf("paskal5(%d, %d) = %d, want %d", m, n, got, want)
			}
		})
	}
}
