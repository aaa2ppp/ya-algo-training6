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

type Idx int32

type Graph [][]Idx

func solve(graph Graph, modulo int) int {
	const op = "solve"

	res := 1

	comps, ok := searchComponents(graph)

	if !ok {
		return 0
	}

	// посчитаем и выкенем одинокие поинты
	lonerCount := 0
	totalCount := 0
	for i := len(comps) - 1; i >= 0; i-- {
		n := comps[i].size
		totalCount += n
		if n == 1 {
			lonerCount++
			n := len(comps)
			comps[i] = comps[n-1]
			comps = comps[:n-1]
		}
	}

	if debugEnable {
		log.Printf("%s: found components: %v", op, comps)
		log.Printf("%s: loner points %d", op, lonerCount)
		log.Printf("%s: total points %d", op, totalCount)
	}

	// перестановки одиночек
	res *= fact(lonerCount, modulo)
	res %= modulo

	// ситаем внутреннюю жизнь компонент
	for _, comp := range comps {
		n := calcComponent(graph, comp, modulo)
		if n == 0 {
			return 0
		}
		res *= n
		res %= modulo
	}

	// перестановки компонентов
	res *= fact(len(comps), modulo)
	res %= modulo

	// одиночки могут вклиниватся в любое место гупп. +1 так как два дерева
	// за счет этого добавляется одно место (перегородка)
	res *= paskal5(totalCount+1, lonerCount, modulo)
	res %= modulo

	return res
}

type Component struct {
	node Idx
	size int
}

func calcComponent(graph Graph, comp Component, k int) int {
	const op = "calcComponent"

	// check point или p2p - симетрия по горизонтали
	if comp.size == 1 {
		if debugEnable {
			log.Printf("%s: %v is point", op, comp)
		}
		return 2
	}
	if comp.size == 2 {
		if debugEnable {
			log.Printf("%s: %v is p2p", op, comp)
		}
		return 2
	}

	type hub struct {
		hubN   int
		pointN int
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
				hubN:   hubN,
				pointN: pointN,
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

		res *= fact(hub.pointN, k)
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
			n, ok := dfs(Idx(node), -1)
			if !ok {
				return nil, false
			}
			comps = append(comps, Component{
				node: Idx(node),
				size: n,
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
	f := _fact[i-1]

	for ; i <= n; i++ {
		f *= i
		f %= modulo
		_fact = append(_fact, f)
	}

	return f
}

func paskal5(m, n int, modulo int) int {
	// if n > m/2 {
	// 	n = m - n
	// }

	bb := make([]int, 0, n)
	for i := n; i > 1; i-- {
		bb = append(bb, i)
	}
	if debugEnable {
		log.Printf("packal5: bb: %v", bb)
	}

	a := 1
mainLoop:
	for p := m - n + 1; p <= m; p++ {
		p := p
		for i := len(bb) - 1; i >= 0; i-- {
			if p == 1 {
				continue mainLoop
			}
			d := gcd(p, bb[i])
			if d != 1 {
				if debugEnable {
					log.Printf("packal4: a:%d/%d=%d b:%d/%d=%d", p, d, p/d, bb[i], d, bb[i]/d)
				}
				p /= d
				bb[i] /= d
				if bb[i] == 1 {
					n := len(bb)
					bb[i] = bb[n-1]
					bb = bb[:n-1]
					if debugEnable {
						log.Printf("packal5: bb: %v", bb)
					}
				}
			}
		}
		a *= p
		a %= modulo
	}

	if debugEnable {
		log.Printf("paskal5: bb: %v", bb)
	}

	return a
}

func paskal4(m, n int, modulo int) int {
	// Пройтись по ряду, найти наибольшие наибольшие делители,
	// а потом сократить на них числитель. По идее должно быть NlogN?
	if n > m/2 {
		n = m - n
	}

	b := 1
	var ds []int
	for i := 2; i <= n; i++ {
		i := i
		d := gcd(i, modulo)
		if d != 1 {
			ds = append(ds, d)
			i /= d
		}
		b *= i
		b %= modulo
	}

	a := 1
	for p := m - n + 1; p <= m; p++ {
		for i := len(ds) - 1; i >= 0; i-- {
			d := gcd(p, ds[i])
			if d != 1 {
				p /= d
				ds[i] /= d
				if ds[i] == 1 {
					n := len(ds)
					ds[i] = ds[n-1]
					ds = ds[:n-1]
				}
			}
		}
		a *= p
		a %= modulo
	}

	a *= invmod(b, modulo)
	a %= modulo

	return a
}

func paskal3(m, n int, modulo int) int {

	a := 1
	for p := m - n + 1; p <= m; p++ {
		a *= p
		a %= modulo
	}

	b := fact(n, modulo)

	a *= invmod(b, modulo)
	a %= modulo

	return a
}

func paskal2(m, n int, modulo int) int {
	if n > m/2 {
		n = m - n
	}
	a := fact(m, modulo)
	b := fact(n, modulo)
	b *= fact(m-n, modulo)
	b %= modulo
	a *= invmod(b, modulo)
	a %= modulo
	return a
}

func paskal(m, n int, modulo int) int {
	if debugEnable {
		log.Printf("paskal: %d %d", m, n)
	}

	// std::cin >> m >> n;
	// if (n > m / 2)
	//    n = m - n;
	if n > m/2 {
		n = m - n
	}

	// std::vector<unsigned> v(m + 1, 1);
	v := make([]int, m+1)
	for i := range v {
		v[i] = 1
	}

	// for (unsigned cm = 0; cm <= m; cm++)
	//    for (unsigned i = m + 1 - cm; i < m; i++)
	// 	  v[i] = (v[i] + v[i + 1]) % modulus;
	for cm := 0; cm <= m; cm++ {
		for i := m + 1 - cm; i < m; i++ {
			v[i] = (v[i] + v[i+1]) % modulo
		}
	}

	// std::cout << v[n] << std::endl;
	return v[n]
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

	graph := make(Graph, n+1)

	for i := 0; i < m; i++ {
		a, b, err := scanTwoInt(sc)
		if err != nil {
			panic(err)
		}
		graph[a] = append(graph[a], Idx(b))
		graph[b] = append(graph[b], Idx(a))
	}

	if debugEnable {
		log.Println("graph:", graph)
	}

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
