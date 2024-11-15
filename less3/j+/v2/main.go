package main

import (
	"bufio"
	"io"
	"log"
	"math"
	"os"
	"sort"
	"strconv"
	"unsafe"
)

// usage: [SOLVE={slow|fast|queue}] [DEBUG=] ./main < ../test_data/<test_num>
// queue by default

// ----------------------------------------------------------------------------
// O(n^2)

func slowSolve(vasyaH int, hh, ww []int) int {
	if debugEnable {
		log.Print("slowSolve()")
	}

	n := len(hh)

	type item struct {
		h int
		w int
	}

	hw := make([]item, 0, n)
	for i := 0; i < n; i++ {
		hw = append(hw, item{h: hh[i], w: ww[i]})
	}

	sort.Slice(hw, func(i, j int) bool {
		return hw[i].h < hw[j].h
	})

	minmax := math.MaxInt
	for i := 0; i < n; i++ {
		width := hw[i].w
		maximum := 0
		for j := i + 1; j < n && width < vasyaH; j++ {
			width += hw[j].w
			maximum = max(maximum, hw[j].h-hw[j-1].h)
		}
		if width >= vasyaH {
			minmax = min(minmax, maximum)
		}
	}

	return minmax
}

// ----------------------------------------------------------------------------
// O(n) в среднем на рандомных данных, O(n^2) в худшем случае. см. тест 999

func fastSolve(vasyaH int, hh, ww []int) int {
	if debugEnable {
		log.Print("fastSolve()")
	}

	n := len(hh)

	type item struct {
		h int
		w int
	}

	hw := make([]item, 0, n)
	for i := 0; i < n; i++ {
		hw = append(hw, item{h: hh[i], w: ww[i]})
	}

	sort.Slice(hw, func(i, j int) bool {
		return hw[i].h < hw[j].h
	})

	minmax := math.MaxInt
	for i := 0; i < n; {
		width := hw[i].w
		maximum := 0
		next := i + 1
		for j := i + 1; j < n && width < vasyaH; j++ {
			width += hw[j].w
			if d := hw[j].h - hw[j-1].h; d >= maximum {
				maximum = d
				next = j
			}
		}
		if width >= vasyaH {
			minmax = min(minmax, maximum)
		}
		i = next
	}

	return minmax
}

// ----------------------------------------------------------------------------
// O(n) с помощью очереди на двух стеках с поддержкой максимума

type item struct {
	w      int
	dh     int // always >= 0
	max_dh int
}

type stack []item

func (s stack) empty() bool {
	return len(s) == 0
}

func (s stack) maxDh() int {
	if s.empty() {
		return 0
	}
	return s.top().max_dh
}

func (s *stack) push(dh, w int) {
	max_dh := abs(dh)
	if !s.empty() {
		max_dh = max(max_dh, s.top().max_dh)
	}
	*s = append(*s, item{
		w:      w,
		dh:     dh,
		max_dh: max_dh,
	})
}

func (s stack) top() *item {
	n := len(s)
	return &s[n-1]
}

func (s *stack) pop() (dh, w int) {
	old := *s
	n := len(old)
	it := &old[n-1]
	dh, w = it.dh, it.w
	*s = old[:n-1]
	return dh, w
}

type queue struct {
	input  stack
	output stack
	width  int
}

func (q *queue) empty() bool {
	return q.input.empty() && q.output.empty()
}

func (q *queue) maxDh() int {
	return max(q.input.maxDh(), q.output.maxDh())
}

func (q *queue) push(dh, w int) {
	q.input.push(dh, w)
	q.width += w
}

func (q *queue) pop() (dh, w int) {
	if q.output.empty() {
		q.pour()
	}
	q.width -= q.output.top().w
	return q.output.pop()
}

func (q *queue) pour() {
	for !q.input.empty() {
		dh, w := q.input.pop()
		q.output.push(dh, w)
	}
}

func solve(vasyaH int, hh, ww []int) int {
	if debugEnable {
		log.Print("solve()")
	}

	n := len(hh)

	oo := make([]int, n)
	for i := range oo {
		oo[i] = i
	}

	sort.Slice(oo, func(i, j int) bool {
		i, j = oo[i], oo[j]
		return hh[i] < hh[j]
	})

	var q queue
	minMaxDh := math.MaxInt

	i := 0
	j := oo[i]
	prevH, prevW := hh[j], ww[j]

	for {
		width := q.width + prevW
		if width >= vasyaH {
			if q.empty() {
				minMaxDh = 0
				break
			}
			minMaxDh = min(minMaxDh, q.maxDh())
			q.pop()
			continue
		}

		i++
		if i == n {
			break
		}
		j := oo[i]
		h, w := hh[j], ww[j]

		q.push(h-prevH, prevW)
		prevH, prevW = h, w
	}

	return minMaxDh
}

// ----------------------------------------------------------------------------

func run(in io.Reader, out io.Writer, solve func(h int, hh, ww []int) int) {
	sc := bufio.NewScanner(in)
	sc.Split(bufio.ScanWords)
	bw := bufio.NewWriter(out)
	defer bw.Flush()

	n, h, err := scanTwoInt(sc)
	if err != nil {
		panic(err)
	}

	hh := make([]int, n)
	ww := make([]int, n)

	if err := scanInts(sc, hh); err != nil {
		panic(err)
	}

	if err := scanInts(sc, ww); err != nil {
		panic(err)
	}

	res := solve(h, hh, ww)

	writeInt(bw, res, defaultWriteOpts())
}

// ----------------------------------------------------------------------------

var _, debugEnable = os.LookupEnv("DEBUG")

func main() {
	_ = debugEnable
	solveEnv := os.Getenv("SOLVE")
	solve := solve // queue by default
	switch solveEnv {
	case "slow":
		solve = slowSolve
	case "fast":
		solve = fastSolve
	}
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
