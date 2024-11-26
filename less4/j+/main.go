package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"unsafe"
)

const (
	modulo = 1e9 + 7
	maxN   = 3000
)

type (
	Idx  int32
	Size int32

	Edge  [2]Idx
	Graph []Node

	Node struct {
		up   []Idx
		down []Idx
	}

	SolveFunc func([]Edge) int
)

func solve(edges []Edge) int {
	n := len(edges) + 1
	graph := make(Graph, n+1)

	for _, e := range edges {
		a := e[0]
		b := e[1]
		graph[a].down = append(graph[a].down, b)
		graph[b].up = append(graph[b].up, a)
	}

	return _solve(graph)
}

func _solve(graph Graph) int {
	n := len(graph)
	dp := makeMatrix[Size](n, n) // (!) both dir 1-indexing
	sizes := make([]Size, n)

	nCr64 := func(n, r int) int64 {
		if r > n {
			return 0
		}
		return int64(paskal(n, r))
	}

	var dfs func(node, prev Idx, dir bool) (int, int, int)

	dfs = func(node, prev Idx, dir bool) (int, int, int) {
		var (
			size      = 1
			upSize    = 0
			downSize  = 0
			upCount   = 0
			downCount = 0
		)

		for _, child := range graph[node].up {
			if child == prev {
				continue
			}

			childSize, childUpCount, _ := dfs(child, node, true)
			size += childSize
			upSize += childSize
			if dir {
				upCount += childUpCount
			}

			for i := 1; i <= childSize; i++ {
				dp[child][i] += dp[child][i-1]
				dp[child][i] %= modulo
			}
		}

		for _, child := range graph[node].down {
			if child == prev {
				continue
			}

			childSize, _, childDownCount := dfs(child, node, false)
			size += childSize
			downSize += childSize
			if !dir {
				downCount += childDownCount
			}

			for i := childSize; i >= 1; i-- {
				dp[child][i] += dp[child][i+1]
				dp[child][i] %= modulo
			}
		}

		if size-1 != upSize+downSize {
			panic("incorrectly summed sizes")
		}

		sizes[node] = Size(size)

		if size == 1 {
			dp[node][1] = 1
			return 1, 0, 0
		}

		var (
			buf    = make([]Size, max(upSize, downSize)+1)
			upDP   = make([]Size, len(buf))
			downDP = make([]Size, len(buf))
			count  int
		)

		count, upDP[0] = upCount, 1
		for _, child := range graph[node].up {
			if child == prev {
				continue
			}

			childSize := int(sizes[child])
			count += childSize

			for i := upCount; i <= count; i++ {
				if upDP[i] == 0 {
					continue
				}

				for j := 1; j <= childSize; j++ {
					p := i + j // positions after merge with child subtree

					v := int64(upDP[i])
					v *= int64(dp[child][j])
					v %= modulo
					v *= nCr64(p, j)
					v %= modulo
					v *= nCr64(count-p, childSize-j)
					v %= modulo

					buf[p] += Size(v)
					buf[p] %= modulo
				}
			}

			upDP, buf = buf, upDP
			clear(buf[:count])
		}

		clear(buf)
		count, downDP[0] = downCount, 1
		for _, child := range graph[node].down {
			if child == prev {
				continue
			}

			childSize := int(sizes[child])
			count += childSize

			for i := downCount; i <= count; i++ {
				if downDP[i] == 0 {
					continue
				}

				for j := 1; j <= childSize; j++ {
					p := i + j // positions after merge with child subtree

					v := int64(downDP[i])
					v *= int64(dp[child][childSize-(j-1)])
					v %= modulo
					v *= nCr64(p, j)
					v %= modulo
					v *= nCr64(count-p, childSize-j)
					v %= modulo

					buf[p] += Size(v)
					buf[p] %= modulo
				}
			}

			downDP, buf = buf, downDP
			clear(buf[:count])
		}

		for i := 1; i <= size; i++ {
			dp[node][i] = 0

			for j, n := upCount, min(i-1, upSize); j <= n; j++ {
				a := size - i
				b := a - (upSize - j)

				if !(0 <= b && b <= downSize) {
					continue
				}

				v := int64(upDP[j])
				v *= nCr64(i-1, j)
				v %= modulo
				v *= int64(downDP[b])
				v %= modulo
				v *= nCr64(a, b)
				v %= modulo

				dp[node][i] += Size(v)
				dp[node][i] %= modulo
			}
		}

		return size, upCount, downCount
	}

	dfs(1, 0, false)

	var count int
	for i := 1; i < n; i++ {
		count += int(dp[1][i])
		count %= modulo
	}

	return count
}

// ----------------------------------------------------------------------------

var (
	_fact    = [maxN + 1]int{1, 1}
	_factN   = 2
	_invfact = [maxN + 1]int{1, 1}
)

func fact(n int) int {
	if n < _factN {
		return _fact[n]
	}

	v := int64(_fact[_factN-1])
	for ; _factN <= n; _factN++ {
		v *= int64(_factN)
		v %= modulo
		_fact[_factN] = int(v)
	}

	return int(v)
}

func invfact(n int) int {
	v := _invfact[n]

	if v == 0 {
		v = invmod(fact(n), modulo)
		_invfact[n] = v
	}

	return v
}

func paskal2(m, n int) int {
	if n == 0 || n == m {
		return 1
	}
	v := int64(fact(m))
	v *= int64(invfact(n))
	v %= modulo
	v *= int64(invfact(m - n))
	v %= modulo
	return int(v)
}

var _paskal = [][]int{{1}, {1}}

func paskal(i, j int) int {
	if j > i/2 {
		j = i - j
	}

	for i >= len(_paskal) {
		n := len(_paskal)
		row := make([]int, n/2+1)
		row[0] = 1
		prev := _paskal[n-1]

		for j := 1; j < len(prev); j++ {
			row[j] = (prev[j-1] + prev[j]) % modulo
		}

		if n%2 == 0 {
			row[len(row)-1] = prev[len(prev)-1] * 2
		}

		_paskal = append(_paskal, row)
	}

	return _paskal[i][j]
}

// ----------------------------------------------------------------------------

func run(in io.Reader, out io.Writer, solve SolveFunc) {
	sc := bufio.NewScanner(in)
	sc.Split(bufio.ScanWords)
	bw := bufio.NewWriter(out)
	defer bw.Flush()

	n, err := scanInt(sc)
	if err != nil {
		panic(err)
	}

	edges := make([]Edge, 0, n-1)

	for i := 0; i < n-1; i++ {
		a, b, err := scanTwoIntX[Idx](sc)
		if err != nil {
			panic(err)
		}
		edges = append(edges, Edge{a, b})
	}

	if debugEnable {
		log.Println("edges:", edges)
	}

	res := solve(edges)

	writeInt(bw, res, defaultWriteOpts())
}

// ----------------------------------------------------------------------------

var _, debugEnable = os.LookupEnv("DEBUG")

func main() {
	_ = debugEnable
	run(os.Stdin, os.Stdout, solve)
}

// ----------------------------------------------------------------------------

func unsafeString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func scanWord(sc *bufio.Scanner) (string, error) {
	if sc.Scan() {
		return sc.Text(), nil
	}
	return "", io.EOF
}

func scanInt(sc *bufio.Scanner) (int, error)                  { return scanIntX[int](sc) }
func scanTwoInt(sc *bufio.Scanner) (_, _ int, _ error)        { return scanTwoIntX[int](sc) }
func scanThreeInt(sc *bufio.Scanner) (_, _, _ int, _ error)   { return scanThreeIntX[int](sc) }
func scanFourInt(sc *bufio.Scanner) (_, _, _, _ int, _ error) { return scanFourIntX[int](sc) }

func scanIntX[T Int](sc *bufio.Scanner) (res T, err error) {
	sc.Scan()
	v, err := strconv.ParseInt(unsafeString(sc.Bytes()), 0, int(unsafe.Sizeof(res))<<3)
	return T(v), err
}

func scanTwoIntX[T Int](sc *bufio.Scanner) (v1, v2 T, err error) {
	v1, err = scanIntX[T](sc)
	if err == nil {
		v2, err = scanIntX[T](sc)
	}
	return v1, v2, err
}

func scanThreeIntX[T Int](sc *bufio.Scanner) (v1, v2, v3 T, err error) {
	v1, err = scanIntX[T](sc)
	if err == nil {
		v2, err = scanIntX[T](sc)
	}
	if err == nil {
		v3, err = scanIntX[T](sc)
	}
	return v1, v2, v3, err
}

func scanFourIntX[T Int](sc *bufio.Scanner) (v1, v2, v3, v4 T, err error) {
	v1, err = scanIntX[T](sc)
	if err == nil {
		v2, err = scanIntX[T](sc)
	}
	if err == nil {
		v3, err = scanIntX[T](sc)
	}
	if err == nil {
		v4, err = scanIntX[T](sc)
	}
	return v1, v2, v3, v4, err
}

func scanInts[T Int](sc *bufio.Scanner, a []T) error {
	for i := range a {
		v, err := scanIntX[T](sc)
		if err != nil {
			return err
		}
		a[i] = v
	}
	return nil
}

type Int interface {
	~int | ~int64 | ~int32 | ~int16 | ~int8
}

type Number interface {
	Int | ~float32 | ~float64
}

type writeOpts struct {
	sep   byte
	begin byte
	end   byte
}

func defaultWriteOpts() writeOpts {
	return writeOpts{sep: ' ', end: '\n'}
}

func writeInt[I Int](bw *bufio.Writer, v I, opts writeOpts) error {
	var buf [32]byte

	var err error
	if opts.begin != 0 {
		err = bw.WriteByte(opts.begin)
	}

	if err == nil {
		_, err = bw.Write(strconv.AppendInt(buf[:0], int64(v), 10))
	}

	if err == nil && opts.end != 0 {
		err = bw.WriteByte(opts.end)
	}

	return err
}

func writeInts[I Int](bw *bufio.Writer, a []I, opts writeOpts) error {
	var err error
	if opts.begin != 0 {
		err = bw.WriteByte(opts.begin)
	}

	if len(a) != 0 {
		var buf [32]byte

		if opts.sep == 0 {
			opts.sep = ' '
		}

		_, err = bw.Write(strconv.AppendInt(buf[:0], int64(a[0]), 10))

		for i := 1; err == nil && i < len(a); i++ {
			err = bw.WriteByte(opts.sep)
			if err == nil {
				_, err = bw.Write(strconv.AppendInt(buf[:0], int64(a[i]), 10))
			}
		}
	}

	if err == nil && opts.end != 0 {
		err = bw.WriteByte(opts.end)
	}

	return err
}

// ----------------------------------------------------------------------------

func gcd[I Int](a, b I) I {
	if a > b {
		a, b = b, a
	}
	for a > 0 {
		a, b = b%a, a
	}
	return b
}

func gcdx(a, b int, x, y *int) int {
	if a == 0 {
		*x = 0
		*y = 1
		return b
	}
	var x1, y1 int
	d := gcdx(b%a, a, &x1, &y1)
	*x = y1 - (b/a)*x1
	*y = x1
	return d
}

func invmod(a, m int) int {
	var x, y int
	g := gcdx(a, m, &x, &y)
	if g != 1 {
		panic(fmt.Errorf("invmod %d %d: g=%d no solution", a, m, g))
	}
	x = (x%m + m) % m
	return x
}

func abs[N Number](a N) N {
	if a < 0 {
		return -a
	}
	return a
}

func sign[N Number](a N) N {
	if a < 0 {
		return -1
	} else if a > 0 {
		return 1
	}
	return 0
}

type Ordered interface {
	Number | ~string
}

func max[T Ordered](a, b T) T {
	if a < b {
		return b
	}
	return a
}

func min[T Ordered](a, b T) T {
	if a > b {
		return b
	}
	return a
}

// ----------------------------------------------------------------------------

func makeMatrix[T any](n, m int) [][]T {
	buf := make([]T, n*m)
	matrix := make([][]T, n)
	for i, j := 0, 0; i < n; i, j = i+1, j+m {
		matrix[i] = buf[j : j+m]
	}
	return matrix
}
