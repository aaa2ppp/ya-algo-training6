package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"unsafe"
)

func solve(matrix [][]byte) byte {
	area, count := seachArea(matrix, '#')
	if debugEnable {
		log.Printf("area(%d):\n", count)
		for _, line := range area {
			log.Printf("%s\n", line)
		}
	}
	if area == nil {
		return 'X'
	}

	n, m := len(area), len(area[0])
	if n*m == count {
		return 'I'
	}

	var rects [][4]int // список *непересекающихся* примоугольников

	idx := -1
	for y, row := range area {
		for x, c := range row {
			if c == '.' {
				if len(rects) == 2 {
					return 'X'
				}

				idx++
				let := byte(idx) + '1'
				x0, y0, x1, y1, ok := fillAndCheckReact(area, x, y, let)
				if !ok {
					if debugEnable {
						log.Println("is not rect:", x0, y0, x1, y1)
					}
					return 'X'
				}

				rects = append(rects, [4]int{x0, y0, x1, y1})
			}
		}
	}

	if debugEnable {
		log.Println("rects:", rects)
		for _, line := range area {
			log.Printf("%s\n", line)
		}
	}

	switch len(rects) {
	case 1:
		x0, y0, x1, y1 := rects[0][0], rects[0][1], rects[0][2], rects[0][3]
		if x0 != 0 && x1 != m-1 && y0 != 0 && y1 != n-1 {
			return 'O'
		}
		if x0 != 0 && x1 == m-1 && y0 != 0 && y1 != n-1 {
			return 'C'
		}
		if x0 != 0 && x1 == m-1 && y0 == 0 && y1 != n-1 {
			return 'L'
		}
	case 2:
		x00, y00, x01, y01 := rects[0][0], rects[0][1], rects[0][2], rects[0][3]
		x10, y10, x11, y11 := rects[1][0], rects[1][1], rects[1][2], rects[1][3]
		if x00 != 0 && x01 != m-1 && y00 == 0 && y01 != n-1 &&
			x10 != 0 && x11 != m-1 && y10 != 0 && y11 == n-1 &&
			x00 == x10 && x01 == x11 {
			return 'H'
		}
		if x00 != 0 && x01 != m-1 && y00 != 0 && y01 != n-1 &&
			x10 != 0 && x11 == m-1 && y10 != 0 && y11 == n-1 &&
			x00 == x10 {
			return 'P'
		}
	}

	return 'X'
}

func seachArea(matrix [][]byte, let byte) ([][]byte, int) {
	n, m := len(matrix), len(matrix[0])

	x0, y0 := n, m
	x1, y1 := -1, -1
	count := 0

	for y, row := range matrix {
		for x, c := range row {
			if c == let {
				count++
				x0, y0 = min(x0, x), min(y0, y)
				x1, y1 = max(x1, x), max(y1, y)
			}
		}
	}

	if x1 == -1 {
		// no any
		return nil, 0
	}

	res := make([][]byte, y1-y0+1)
	for i := range res {
		res[i] = matrix[y0+i][x0 : x1+1]
	}

	return res, count
}

func fillAndCheckReact(area [][]byte, x, y int, let byte) (x0, y0, x1, y1 int, ok bool) {
	n, m := len(area), len(area[0])

	if !(0 <= x && x < m && 0 <= y && y < n) {
		return -1, -1, -1, -1, false
	}

	x0, y0 = x, y
	x1, y1 = x, y
	count := 1

	var frontier [][2]int
	c := area[y][x]
	area[y][x] = let
	frontier = append(frontier, [2]int{x, y})

	for len(frontier) > 0 {
		p := frontier[0]
		frontier = frontier[1:]
		for _, of := range [][2]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}} {
			x, y := p[0]+of[0], p[1]+of[1]
			if 0 <= x && x < m && 0 <= y && y < n && area[y][x] == c {
				count++
				x0, y0 = min(x0, x), min(y0, y)
				x1, y1 = max(x1, x), max(y1, y)
				area[y][x] = let
				frontier = append(frontier, [2]int{x, y})
			}
		}
	}

	isRect := (x1-x0+1)*(y1-y0+1) == count

	return x0, y0, x1, y1, isRect
}

func run(in io.Reader, out io.Writer) {
	br := bufio.NewReader(in)
	bw := bufio.NewWriter(out)
	defer bw.Flush()

	var n int

	if _, err := fmt.Fscanln(br, &n); err != nil {
		panic(err)
	}
	_ = n

	buf, err := io.ReadAll(br)
	if err != nil {
		panic(err)
	}

	matrix := bytes.Split(buf, []byte{'\n'})[:n] // не збываем усеч (может быть финальный \n или мусор)
	for i, line := range matrix {
		matrix[i] = bytes.TrimSpace(line)
	}

	if debugEnable {
		log.Println("matrix:")
		for _, line := range matrix {
			log.Printf("%s\n", line)
		}
	}

	let := solve(matrix)
	bw.WriteByte(let)
}

// ----------------------------------------------------------------------------

var _, debugEnable = os.LookupEnv("DEBUG")

func main() {
	_ = debugEnable
	run(os.Stdin, os.Stdout)
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
