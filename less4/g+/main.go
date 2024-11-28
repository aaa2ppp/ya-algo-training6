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

type (
	Idx   int32
	Size  int32
	Graph [][]Idx
)

func solve(graph Graph, modulo int) int {
	const op = "solve"

	comps, ok := searchComponents(graph)
	if !ok {
		return 0
	}

	// посчитаем и выкенем одинокие поинты
	lonerCount := 0
	totalCount := 0
	for i := len(comps) - 1; i >= 0; i-- {
		compSize := int(comps[i].size)
		totalCount += compSize
		if compSize == 1 {
			lonerCount++
			comps = remove(comps, i)
		}
	}

	if debugEnable {
		log.Printf("%s: found components: %v", op, comps)
		log.Printf("%s: loner points %d", op, lonerCount)
		log.Printf("%s: total points %d", op, totalCount)
	}

	res := 1

	// считаем "внутреннюю жизнь" компонент
	for _, comp := range comps {
		n := calcComponent(graph, comp, modulo)
		if n == 0 {
			return 0
		}
		res *= n
		res %= modulo
	}

	// перестановки компонент
	res *= fact(len(comps), modulo)
	res %= modulo

	// одиночки могут вклиниватся в любое место. +1 так как два дерева
	// за счет этого добавляется одно место (перегородка)
	res *= nCk5(totalCount+1, lonerCount, modulo)
	res %= modulo

	// перестановки одиночек
	res *= fact(lonerCount, modulo)
	res %= modulo

	return res
}

type Component struct {
	node Idx
	size Size
}

func calcComponent(graph Graph, comp Component, k int) int {
	const op = "calcComponent"

	// check point - симетрия только по горизонтали
	if comp.size == 1 {
		if debugEnable {
			log.Printf("%s: %v is point", op, comp)
		}
		return 2
	}

	// check p2p - симетрия только по горизонтали
	if comp.size == 2 {
		if debugEnable {
			log.Printf("%s: %v is p2p", op, comp)
		}
		return 2
	}

	type hub struct {
		hubN   Size
		pointN Size
	}

	var hubs []hub

	var dfs func(node, prev Idx, prevIsHub bool) (nodeIsHub bool)

	dfs = func(node, prev Idx, prevIsHub bool) (nodeIsHub bool) {
		nodeIsHub = len(graph[node]) > 1

		var hubN, pointN int

		if prev != -1 {
			if prevIsHub {
				hubN++
			} else {
				pointN++
			}
		}

		for _, neig := range graph[node] {
			if neig == prev {
				continue
			}
			if dfs(neig, node, nodeIsHub) {
				hubN++
			} else {
				pointN++
			}
		}

		if nodeIsHub {
			hubs = append(hubs, hub{
				hubN:   Size(hubN),
				pointN: Size(pointN),
			})
		}

		return nodeIsHub
	}

	dfs(comp.node, -1, false)

	if debugEnable {
		log.Printf("%s: found hubs: %v", op, hubs)
	}

	var res int
	if len(hubs) == 1 {
		// только один хаб - симетрия только по горизонтали
		res = 2
	} else {
		// цепочка хабов - симетрия по горизонтали и вертикали
		res = 4
	}

	for _, hub := range hubs {
		if hub.hubN > 2 {
			// Все хабы должны быть выстроены в линию.
			// Т.е. два конечных, которые соеденине только с одним хабом,
			// остальные соеденены с двумя хабами.
			// Вырожденый случай с одним хабом, который соеденен только с поинтами.
			return 0
		}

		res *= fact(int(hub.pointN), k)
		res %= k
	}

	return res
}

func searchComponents(graph Graph) ([]Component, bool) {
	const op = "searchComponents"

	n := len(graph)

	var comps []Component
	var dfs func(node, prev Idx) (int, bool)

	visited := make([]bool, n)

	dfs = func(node, prev Idx) (int, bool) {
		if visited[node] {
			if debugEnable {
				log.Printf("%s: %d: found loop", op, node)
			}
			return 0, false
		}
		visited[node] = true
		size := 1

		for _, neig := range graph[node] {
			if neig != prev {
				n, ok := dfs(neig, node)
				if !ok {
					return 0, false
				}
				size += n
			}
		}

		return size, true
	}

	for node := 1; node < n; node++ {
		if !visited[node] {
			size, ok := dfs(Idx(node), -1)
			if !ok {
				return nil, false
			}
			comps = append(comps, Component{
				node: Idx(node),
				size: Size(size),
			})
		}
	}

	return comps, true
}

var _fact = []int{1, 1}

func fact(n int, modulo int) int {
	if n < len(_fact) {
		return _fact[n]
	}

	i := len(_fact)
	v := _fact[i-1]

	for ; i <= n; i++ {
		v *= i
		v %= modulo
		_fact = append(_fact, v)
	}

	return v
}

func remove[T any](aa []T, i int) []T {
	n := len(aa)
	if i < n-1 {
		aa[i] = aa[n-1]
	}
	return aa[:n-1]
}

func nCk5(n, k int, modulo int) int {
	if k > n/2 {
		k = n - k
	}

	aa := make([]int, k)
	for i, v := 0, n-k+1; v <= n; i, v = i+1, v+1 {
		aa[i] = v
	}

	for b := 2; b <= k; b++ {
		b := b
		for i := len(aa) - 1; i >= 0 && b != 1; i-- {
			d := gcd(aa[i], b)
			aa[i] /= d
			b /= d
			if aa[i] == 1 {
				aa = remove(aa, i)
			}
		}
	}

	v := 1
	for _, a := range aa {
		v *= a
		v %= modulo
	}

	return v
}

func nCk2(n, k int, modulo int) int {
	v := fact(k, modulo)
	v *= fact(n-k, modulo)
	v %= modulo
	v = invmod(v, modulo) * fact(n, modulo)
	v %= modulo
	return v
}

func run(in io.Reader, out io.Writer) {
	sc := bufio.NewScanner(in)
	sc.Split(bufio.ScanWords)
	bw := bufio.NewWriter(out)
	defer bw.Flush()

	n, m, k, err := scanThreeInt(sc)
	if err != nil {
		panic(err)
	}

	if m > n-1 {
		bw.WriteString("0\n")
		return
	}

	graph := make(Graph, n+1)

	for i := 0; i < m; i++ {
		a, b, err := scanTwoInt(sc)
		if err != nil {
			panic(err)
		}
		graph[a] = append(graph[a], Idx(b))
		graph[b] = append(graph[b], Idx(a))
	}

	// if debugEnable {
	// 	log.Println("graph:", graph)
	// }

	res := solve(graph, k)

	writeInt(bw, res, defaultWriteOpts())
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
