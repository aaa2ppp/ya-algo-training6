package main

import (
	"bytes"
	"io"
	"math/rand"
	"strconv"
	"strings"
	"testing"
)

type runArgs struct {
	in io.Reader
}

type runTest struct {
	name    string
	args    runArgs
	wantOut string
	debug   bool
}

func runTests() []runTest {

	return []runTest{
		{
			"1",
			runArgs{strings.NewReader(`4 7
1 4 1 2
1 4 2 3
`)},
			`2`,
			true,
		},
		{
			"2",
			runArgs{strings.NewReader(`5 6
1 3 5 4 2
5 4 3 2 1
`)},
			`1`,
			true,
		},
		{
			"3???",
			runArgs{strings.NewReader(`6 15
1 2 5 7 8 11
1 1 4 2 3 6
`)},
			`3`,
			true,
		},
		{
			"4",
			runArgs{strings.NewReader(`1 10
10
10`)},
			`0`,
			true,
		},
		{
			"5",
			runArgs{strings.NewReader(`3 10
1 2 10
1 1 10`)},
			`0`,
			true,
		},
		// TODO: Add test cases.
	}
}

func Test_run(t *testing.T) {
	for _, tt := range runTests() {
		t.Run(tt.name, func(t *testing.T) {
			defer func(v bool) { debugEnable = v }(debugEnable)
			debugEnable = tt.debug
			out := &bytes.Buffer{}
			run(tt.args.in, out, solve)
			if gotOut := out.String(); trimLines(gotOut) != trimLines(tt.wantOut) {
				t.Errorf("run() = %v, want %v", gotOut, tt.wantOut)
			}
		})
	}
}

func Test_slowRun(t *testing.T) {
	for _, tt := range runTests() {
		t.Run(tt.name, func(t *testing.T) {
			defer func(v bool) { debugEnable = v }(debugEnable)
			debugEnable = tt.debug
			out := &bytes.Buffer{}
			run(tt.args.in, out, slowSolve)
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

func genData(a uint, maxN, maxH, maxW int) (h int, hh, ww []int) {
	rand := rand.New(rand.NewSource(int64(a)))

	n := rand.Intn(maxN) + 1

	hh = make([]int, n)
	ww = make([]int, n)

	maxVasyaH := 0
	for i := 0; i < n; i++ {
		hh[i] = rand.Intn(maxH) + 1
		ww[i] = rand.Intn(maxW) + 1
		maxVasyaH += ww[i]
	}

	h = rand.Intn(maxVasyaH)

	return h, hh, ww
}

func joinIntSlice(aa []int, del string) string {
	n := len(aa)
	if n == 0 {
		return ""
	}

	var sb strings.Builder
	sb.WriteString(strconv.Itoa(aa[0]))
	for i := 1; i < n; i++ {
		sb.WriteString(del)
		sb.WriteString(strconv.Itoa(aa[i]))
	}

	return sb.String()
}

func testSolveN(f *testing.F, n int) {
	for i := uint(1); i < 10; i++ {
		f.Add(i)
	}
	f.Fuzz(func(t *testing.T, a uint) {
		h, hh, ww := genData(a, n, n, n)

		res1 := slowSolve(h, hh, ww)
		res2 := solve(h, hh, ww)

		if res1 != res2 {
			t.Log("h:", h)
			t.Log("hh:", joinIntSlice(hh, ", "))
			t.Log("ww:", joinIntSlice(ww, ", "))
			t.Errorf("Oops!.. slowSolve()=%d solve()=%d", res1, res2)
		}
	})
}

func Fuzz_solve10(f *testing.F) {
	testSolveN(f, 10)
}

func Fuzz_solve100(f *testing.F) {
	testSolveN(f, 10000)
}

func Fuzz_solve1000(f *testing.F) {
	testSolveN(f, 1000)
}

func Fuzz_solve10000(f *testing.F) {
	testSolveN(f, 10000)
}
