package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strconv"
	"unsafe"
)

type stack[T any] []T

func (s stack[T]) empty() bool {
	return len(s) == 0
}

func (s *stack[T]) push(v T) {
	*s = append(*s, v)
}

func (s stack[T]) top() T {
	n := len(s)
	return s[n-1]
}

func (s *stack[T]) pop() T {
	old := *s
	n := len(old)
	v := old[n-1]
	*s = old[:n-1]
	return v
}

type queue[T any] struct {
	input  stack[T]
	output stack[T]
}

func (q *queue[T]) empty() bool {
	return q.input.empty() && q.output.empty()
}

func (q *queue[T]) len() int {
	return len(q.input) + len(q.output)
}

func (q *queue[T]) push(v T) {
	q.input.push(v)
}

func (q *queue[T]) front() T {
	q.pour()
	return q.output.top()
}

func (q *queue[T]) pop() T {
	q.pour()
	return q.output.pop()
}

func (q *queue[T]) pour() {
	if len(q.output) == 0 {
		for !q.input.empty() {
			q.output.push(q.input.pop())
		}
	}
}

func (q queue[T]) String() string {
	a := make([]T, 0, q.len())
	for i := len(q.output) - 1; i >= 0; i-- {
		a = append(a, q.output[i])
	}
	for i := 0; i < len(q.input); i++ {
		a = append(a, q.input[i])
	}
	return fmt.Sprint(a)
}

type rover struct {
	idx  int
	dir  int
	time int
}

// (!) function expected 0-indexed direcitions (0 <= a, b, rovers[i].dir < 4)
func solve(a, b int, rovers []rover) []int {
	n := len(rovers)

	// в условии не сказано, что список упорядочен по времени
	sort.Slice(rovers, func(i, j int) bool {
		return rovers[i].time < rovers[j].time
	})

	times := make([]int, n)

	prio := make([]byte, 4) // приоритет направления
	prio[a] = 1
	prio[b] = 1

	dqs := make([]queue[int], 4) // очереди на проезд в каждом направлении
	dqsN := 0                    // обще кол-во ожидающих проезда
	nice := make([]bool, 4)      // кто уступает на текущем тике
	t := 0                       // текущее время (тик) (роверы подъезжают начиная с t=1)

	for dqsN > 0 || len(rovers) > 0 {
		t++

		// скипаем тики, когда на перекрестке пусто
		if dqsN == 0 {
			t = rovers[0].time
		}

		// подъезжают роверы и вcтают в свои очереди
		for len(rovers) > 0 && rovers[0].time <= t {
			r := rovers[0]
			rovers = rovers[1:]
			q := &dqs[r.dir]
			q.push(r.idx)
			dqsN++
		}

		if debugEnable {
			log.Printf("%d: dqs : %v", t, dqs)
		}

		// принимаем решение, кто уступает
		for dir := 0; dir < 4; dir++ {
			nice[dir] = false

			q := &dqs[dir]
			if q.empty() {
				nice[dir] = true
				continue
			}

			// проверяем помеху справа
			opDir := (dir + 3) % 4
			opQ := &dqs[opDir]
			if !opQ.empty() && prio[opDir] >= prio[dir] {
				nice[dir] = true
				continue
			}

			// проверяем главную слева
			opDir = (dir + 1) % 4
			opQ = &dqs[opDir]
			if !opQ.empty() && prio[opDir] > prio[dir] {
				nice[dir] = true
				continue
			}

			// (!) проверяем главную прямо
			opDir = (dir + 2) % 4
			opQ = &dqs[opDir]
			if !opQ.empty() && prio[opDir] > prio[dir] {
				nice[dir] = true
				continue
			}
		}

		if debugEnable {
			log.Printf("%d: nice: %v", t, nice)
		}

		// проезжают те, кому не надо уступать
		count := 0
		for dir := 0; dir < 4; dir++ {
			if nice[dir] {
				continue
			}
			q := &dqs[dir] // это безопасно, если очередь пустая она уступает
			times[q.pop()] = t
			dqsN--
			count++
		}

		if count == 0 {
			// это никогда не должно произойти
			panic("no one passed the intersection")
		}

		if debugEnable {
			log.Printf("%d: dqs : %v", t, dqs)
		}
	}

	return times
}

func run(in io.Reader, out io.Writer) {
	sc := bufio.NewScanner(in)
	sc.Split(bufio.ScanWords)
	bw := bufio.NewWriter(out)
	defer bw.Flush()

	n, a, b, err := scanThreeInt(sc)
	if err != nil {
		panic(err)
	}

	// to 0-indexing
	a--
	b--

	rovers := make([]rover, n)
	for i := 0; i < n; i++ {
		d, t, err := scanTwoInt(sc)
		if err != nil {
			panic(err)
		}
		rovers[i] = rover{
			idx:  i,
			dir:  d - 1, // to 0-indexing
			time: t,
		}
	}

	times := solve(a, b, rovers)

	wo := writeOpts{sep: '\n', end: '\n'}
	writeInts(bw, times, wo)
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
