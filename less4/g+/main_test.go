package main

import (
	"bytes"
	"io"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"
)

const test_data = "./test_data/"

func Test_run(t *testing.T) {
	readFile := func(filename string) []byte {
		buf, err := os.ReadFile(test_data + filename)
		if err != nil {
			panic(err)
		}
		return buf
	}

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
		{
			"19",
			args{bytes.NewReader(readFile("19"))},
			unsafeString(readFile("19.a")),
			false,
		},
		{
			"24",
			args{bytes.NewReader(readFile("24"))},
			unsafeString(readFile("24.a")),
			false,
		},
		{
			"26",
			args{bytes.NewReader(readFile("26"))},
			unsafeString(readFile("26.a")),
			false,
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

func Benchmark_fact(b *testing.B) {
	defer func(v bool) { debugEnable = v }(debugEnable)
	debugEnable = false

	for i := 0; i < b.N; i++ {
		_fact = _fact[:2]
		fact(1000000, int(1e9)+7)
	}
}

func nCk_slow(n, k int, modulo int) int {
	defer func(v bool) { debugEnable = v }(debugEnable)
	debugEnable = false

	if k > n/2 {
		k = n - k
	}

	v := make([]int, n+1)
	for i := range v {
		v[i] = 1
	}

	for cm := 0; cm <= n; cm++ {
		for i := n + 1 - cm; i < n; i++ {
			v[i] = (v[i] + v[i+1]) % modulo
		}
	}

	return v[k]
}

func Test_nCk2(t *testing.T) {
	defer func(v bool) { debugEnable = v }(debugEnable)
	debugEnable = false

	const modulo = 1e9 + 7
	rand := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < 100; i++ {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			m := rand.Intn(1000) + 1
			n := rand.Intn(1000) + 1
			if n > m {
				m, n = n, m
			}
			if got, want := nCk2(m, n, modulo), nCk_slow(m, n, modulo); got != want {
				t.Errorf("nCk2(%d, %d) = %d, want %d", m, n, got, want)
			}
		})
	}
}

func Test_nCk5(t *testing.T) {
	defer func(v bool) { debugEnable = v }(debugEnable)
	debugEnable = false

	const modulo = 1e9 + 7
	rand := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < 100; i++ {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			m := rand.Intn(1000) + 1
			n := rand.Intn(1000) + 1
			if n > m {
				m, n = n, m
			}
			if got, want := nCk5(m, n, modulo), nCk_slow(m, n, modulo); got != want {
				t.Fatalf("nCk5(%d, %d) = %d, want %d", m, n, got, want)
			}
		})
	}
}

func Benchmark_nCk(b *testing.B) {
	defer func(v bool) { debugEnable = v }(debugEnable)
	debugEnable = false

	const modulo = 1e9 + 7
	rand := rand.New(rand.NewSource(time.Now().UnixNano()))

	n := rand.Intn(1000) + 9000
	k := rand.Intn(1000) + 4500

	b.Run("nCk_slow 10K", func(b *testing.B) {

		for i := 0; i < b.N; i++ {
			nCk_slow(n, k, modulo)
		}
	})

	b.Run("nCk2 1M", func(b *testing.B) {
		_fact = []int{1, 1}
		for i := 0; i < b.N; i++ {
			nCk2(n, k, modulo)
		}
	})

	b.Run("nCk5 1M", func(b *testing.B) {
		_fact = []int{1, 1}
		for i := 0; i < b.N; i++ {
			nCk5(n, k, modulo)
		}
	})

}
