package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"strconv"
	"unsafe"
)

type (
	Idx     int32
	Edge    [2]Idx
	EdgeSet map[Edge]struct{}
	Graph   [][]Idx
)

func (es EdgeSet) Add(e Edge) {
	es[e] = struct{}{}
}

func solve(edges EdgeSet, graph Graph) int {
	n := len(graph)

	type item struct {
		node Idx
		val  int
	}

	maxHH := make([][3]item, n)
	maxLL := make([][2]item, n)

	var dfs1 func(node, prev Idx) (h, l int)

	dfs1 = func(node, prev Idx) (h, l int) {
		neigs := graph[node]

		var maxH0, maxH1, maxH2 item
		var maxL0, maxL1 item

		for _, neig := range neigs {
			if neig == prev {
				continue
			}

			h, l := dfs1(neig, node)
			h++

			switch {
			case h > maxH0.val:
				maxH2 = maxH1
				maxH1 = maxH0
				maxH0 = item{neig, h}
			case h > maxH1.val:
				maxH2 = maxH1
				maxH1 = item{neig, h}
			case h > maxH2.val:
				maxH2 = item{neig, h}
			}

			switch {
			case l > maxL0.val:
				maxL1 = maxL0
				maxL0 = item{neig, l}
			case l > maxL1.val:
				maxL1 = item{neig, l}
			}
		}

		maxHH[node] = [3]item{maxH0, maxH1, maxH2}
		maxLL[node] = [2]item{maxL0, maxL1}

		return maxH0.val, max(maxL0.val, maxH0.val+maxH1.val)
	}

	var dfs2 func(node, prev Idx, prevH, prevL int)

	dfs2 = func(node, prev Idx, prevH, prevL int) {
		neigs := graph[node]

		maxH0, maxH1, maxH2 := maxHH[node][0], maxHH[node][1], maxHH[node][2]
		maxL0, maxL1 := maxLL[node][0], maxLL[node][1]

		for _, neig := range neigs {
			if neig == prev {
				h, l := prevH+1, prevL

				switch {
				case h > maxH0.val:
					maxH2 = maxH1
					maxH1 = maxH0
					maxH0 = item{neig, h}
				case h > maxH1.val:
					maxH2 = maxH1
					maxH1 = item{neig, h}
				case h > maxH2.val:
					maxH2 = item{neig, h}
				}

				switch {
				case l > maxL0.val:
					maxL1 = maxL0
					maxL0 = item{neig, l}
				case l > maxL1.val:
					maxL1 = item{neig, l}
				}

				maxHH[node] = [3]item{maxH0, maxH1, maxH2}
				maxLL[node] = [2]item{maxL0, maxL1}
				break
			}
		}

		for _, neig := range neigs {
			if neig == prev {
				continue
			}

			var h, l int
			switch neig {
			case maxH0.node:
				h = maxH1.val
				l = maxH2.val + maxH1.val
			case maxH1.node:
				h = maxH0.val
				l = maxH0.val + maxH2.val
			default:
				h = maxH0.val
				l = maxH0.val + maxH1.val
			}

			switch neig {
			case maxL0.node:
				l = max(l, maxL1.val)
			default:
				l = max(l, maxL0.val)
			}

			dfs2(neig, node, h, l)
		}
	}

	dfs1(1, 0)
	if debugEnable {
		log.Printf("1: hh: %v", maxHH)
		log.Printf("1: ll: %v", maxLL)
	}

	dfs2(1, 0, 0, 0)
	if debugEnable {
		log.Printf("1: hh: %v", maxHH)
		log.Printf("1: ll: %v", maxLL)
	}

	maximum := 0
	for edge := range edges {
		a, b := edge[0], edge[1]

		aHH := maxHH[a]
		bHH := maxHH[b]

		aLL := maxLL[a]
		bLL := maxLL[b]

		var xa, xb int

		switch b {
		case aHH[0].node:
			xa = aHH[2].val + aHH[1].val
		case aHH[1].node:
			xa = aHH[0].val + aHH[2].val
		default:
			xa = aHH[0].val + aHH[1].val
		}

		switch b {
		case aLL[0].node:
			xa = max(xa, aLL[1].val)
		default:
			xa = max(xa, aLL[0].val)
		}

		switch a {
		case bHH[0].node:
			xb = bHH[2].val + bHH[1].val
		case bHH[1].node:
			xb = bHH[0].val + bHH[2].val
		default:
			xb = bHH[0].val + bHH[1].val
		}

		switch a {
		case bLL[0].node:
			xb = max(xb, bLL[1].val)
		default:
			xb = max(xb, bLL[0].val)
		}

		if debugEnable {
			log.Printf("3: %d-%d: %d %d", a, b, xa, xb)
		}

		maximum = max(maximum, xa*xb)
	}

	return maximum
}

func run(in io.Reader, out io.Writer, solve func(edges EdgeSet, graph Graph) int) {
	sc := bufio.NewScanner(in)
	sc.Split(bufio.ScanWords)
	bw := bufio.NewWriter(out)
	defer bw.Flush()

	n, err := scanInt(sc)
	if err != nil {
		panic(err)
	}

	graph := make(Graph, n+1)
	edges := make(EdgeSet, n)

	for i := 0; i < n-1; i++ {
		a, b, err := scanTwoIntX[Idx](sc)
		if err != nil {
			panic(err)
		}
		if a > b {
			a, b = b, a
		}
		edges.Add(Edge{a, b})
		graph[a] = append(graph[a], b)
		graph[b] = append(graph[b], a)
	}

	if debugEnable {
		log.Println("edges:", edges)
		log.Println("graph:", graph)
	}

	res := solve(edges, graph)

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
