package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"unicode"
)

func scanTokens(r io.ByteScanner) ([]string, bool) {
	var tokens []string

	first := true
	for {
		c, err := r.ReadByte()
		if debugEnable {
			log.Printf("c: '%c'", c)
		}
		if err == io.EOF {
			break
		}

		if unicode.IsSpace(rune(c)) {
			continue
		}

		if c == '(' {
			token := "("
			if debugEnable {
				log.Printf("token: '%v'", token)
			}
			tokens = append(tokens, token)
			first = true
			continue
		}

		if first {
			first = false
			if err := r.UnreadByte(); err != nil {
				panic(err)
			}
			token, ok := scanIntStr(r)
			if !ok {
				return nil, false
			}
			if debugEnable {
				log.Printf("token: '%v'", token)
			}
			tokens = append(tokens, token)
			continue
		}

		switch c {
		case '+', '-', '*', ')':
			token := string(c)
			if debugEnable {
				log.Printf("token: '%v'", token)
			}
			tokens = append(tokens, token)
		default:
			if err := r.UnreadByte(); err != nil {
				panic(err)
			}
			token, ok := scanIntStr(r)
			if !ok {
				return nil, false
			}
			if debugEnable {
				log.Printf("token: '%v'", token)
			}
			tokens = append(tokens, token)
		}
	}

	return tokens, true
}

func scanIntStr(r io.ByteScanner) (string, bool) {
	var sb strings.Builder

	// first character can be a sign
	c, err := r.ReadByte()
	if err == io.EOF {
		if debugEnable {
			log.Println("scanIntStr: unexpected EOF")
		}
		return "", false
	}
	if err != nil {
		panic(err)
	}
	switch c {
	case '-', '+':
		sb.WriteByte(c)
	default:
		if err := r.UnreadByte(); err != nil {
			panic(err)
		}
	}

	// must be at least one digit
	c, err = r.ReadByte()
	if err == io.EOF {
		if debugEnable {
			log.Println("scanIntStr: unexpected EOF")
		}
		return "", false
	}
	if err != nil {
		panic(err)
	}
	if !unicode.IsDigit(rune(c)) {
		if debugEnable {
			log.Printf("scanIntStr: expected digit, got '%c'", c)
		}
		return "", false
	}
	sb.WriteByte(c)

	// next characters must be digits
	for {
		c, err = r.ReadByte()
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}
		if !unicode.IsDigit(rune(c)) {
			if err := r.UnreadByte(); err != nil {
				panic(err)
			}
			break
		}
		sb.WriteByte(c)
	}

	return sb.String(), true
}

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
	c := old[n-1]
	*s = old[:n-1]
	return c
}

// чем больше значение, тем выше приоритет
func getPriority(op string) int {
	switch op {
	case "*":
		return 200
	case "+", "-":
		return 100
	default:
		panic(fmt.Errorf("unknown operation: %v", op))
	}
}

func toPolish(tokens []string) ([]string, bool) {
	type item struct {
		token    string
		priority int
	}

	var stack stack[item]
	res := make([]string, 0, len(tokens))

	for _, token := range tokens {
		switch token {
		case "(":
			// всегда кладем на стек, с наименьшим приоритетом
			stack.push(item{token: "(", priority: 0})

		case ")":
			// пререкладываем из стека в результат, пока не встертим скобку
			for {
				if stack.empty() {
					return nil, false
				}

				token := stack.pop().token

				if token == "(" {
					break
				}

				res = append(res, token)
			}

		case "+", "-", "*":
			// вытаскивает со стека все, чей приоритет выше или равен
			prioryty := getPriority(token)
			for !stack.empty() && stack.top().priority >= prioryty {
				res = append(res, stack.pop().token)
			}

			// и ложится на стек
			stack.push(item{token: token, priority: prioryty})

		default:
			// всегда ложится в результат
			res = append(res, token)
		}
	}

	for !stack.empty() {
		res = append(res, stack.pop().token)
	}

	return res, true
}

func calc(tokens []string) (int, bool) {
	var stack stack[int]

	for _, token := range tokens {
		switch token {
		case "+":
			if len(stack) < 2 {
				return 0, false
			}
			b, a := stack.pop(), stack.pop()
			stack.push(a + b)

		case "-":
			if len(stack) < 2 {
				return 0, false
			}
			b, a := stack.pop(), stack.pop()
			stack.push(a - b)

		case "*":
			if len(stack) < 2 {
				return 0, false
			}
			b, a := stack.pop(), stack.pop()
			stack.push(a * b)

		default:
			v, err := strconv.Atoi(token)
			if err != nil {
				return 0, false
			}
			stack.push(v)
		}
	}

	if len(stack) != 1 {
		return 0, false
	}

	return stack.pop(), true
}

func run(in io.Reader, out io.Writer) {
	br := bufio.NewReader(in)
	bw := bufio.NewWriter(out)
	defer bw.Flush()

	tokens, ok := scanTokens(br)
	if !ok {
		log.Println("can't sanToken")
		bw.WriteString("WRONG\n")
		return
	}

	tokens, ok = toPolish(tokens)
	if !ok {
		log.Println("can't toPolish")
		bw.WriteString("WRONG\n")
		return
	}

	res, ok := calc(tokens)
	if !ok {
		log.Println("can't calc")
		bw.WriteString("WRONG\n")
		return
	}

	bw.WriteString(strconv.Itoa(res))
	bw.WriteByte('\n')
}

// ----------------------------------------------------------------------------

var _, debugEnable = os.LookupEnv("DEBUG")

func main() {
	_ = debugEnable
	run(os.Stdin, os.Stdout)
}
