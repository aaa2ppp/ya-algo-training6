package main

import (
	"bufio"
	"container/heap"
	"io"
	"log"
	"os"
	"strconv"
	"unsafe"
)

type Item struct {
	id         int
	a, b       int
	idxA, idxB int
}

type QueueA []*Item

func (pq QueueA) Len() int { return len(pq) }

func (pq QueueA) Less(i, j int) bool {
	return pq[i].a > pq[j].a ||
		pq[i].a == pq[j].a && (pq[i].b > pq[j].b ||
			pq[i].b == pq[j].b && pq[i].id < pq[j].id)
}

func (pq QueueA) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].idxA = i
	pq[j].idxA = j
}

func (pq *QueueA) Push(x any) {
	n := len(*pq)
	item := x.(*Item)
	item.idxA = n
	*pq = append(*pq, item)
}

func (pq *QueueA) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil // don't stop the GC from reclaiming the item eventually
	item.idxA = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

type QueueB []*Item

func (pq QueueB) Len() int { return len(pq) }

func (pq QueueB) Less(i, j int) bool {
	return pq[i].b > pq[j].b ||
		pq[i].b == pq[j].b && (pq[i].a > pq[j].a ||
			pq[i].a == pq[j].a && pq[i].id < pq[j].id)
}

func (pq QueueB) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].idxB = i
	pq[j].idxB = j
}

func (pq *QueueB) Push(x any) {
	n := len(*pq)
	item := x.(*Item)
	item.idxB = n
	*pq = append(*pq, item)
}

func (pq *QueueB) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil // don't stop the GC from reclaiming the item eventually
	item.idxB = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

func run(in io.Reader, out io.Writer) {
	sc := bufio.NewScanner(in)
	sc.Split(bufio.ScanWords)
	bw := bufio.NewWriter(out)
	defer bw.Flush()

	n, err := scanInt(sc)
	if err != nil {
		panic(err)
	}

	items := make([]Item, n)
	for i := 0; i < n; i++ {
		items[i] = Item{
			id:   i + 1,
			idxA: i,
			idxB: i,
		}
	}

	for i := 0; i < n; i++ {
		v, err := scanInt(sc)
		if err != nil {
			panic(err)
		}
		items[i].a = v
	}

	for i := 0; i < n; i++ {
		v, err := scanInt(sc)
		if err != nil {
			panic(err)
		}
		items[i].b = v
	}

	pp := make([]int, n)
	if err := scanInts(sc, pp); err != nil {
		panic(err)
	}

	qa := make(QueueA, n)
	qb := make(QueueB, n)

	for i := 0; i < n; i++ {
		qa[i] = &items[i]
		qb[i] = &items[i]
	}

	heap.Init(&qa)
	heap.Init(&qb)

	res := make([]int, 0, n)
	for i := 0; i < n; i++ {
		if debugEnable {
			log.Println(items)
		}
		switch pp[i] {
		case 0:
			it := heap.Pop(&qa).(*Item)
			res = append(res, it.id)
			heap.Remove(&qb, it.idxB)
		case 1:
			it := heap.Pop(&qb).(*Item)
			res = append(res, it.id)
			heap.Remove(&qa, it.idxA)
		}
	}

	if debugEnable {
		log.Println(items)
	}

	wo := defaultWriteOpts()
	writeInts(bw, res, wo)
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
