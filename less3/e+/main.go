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
	const op = "scanTokens"

	var tokens []string
	var prevIsExpr bool // число или выражение в скобках

	for {
		c, err := r.ReadByte()
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}

		if debugEnable {
			log.Printf("%s: c: '%c'", op, c)
		}

		if unicode.IsSpace(rune(c)) {
			continue
		}

		if prevIsExpr {
			switch c {
			case ')':
				token := string(c)
				if debugEnable {
					log.Printf("%s: token: '%v'", op, token)
				}
				tokens = append(tokens, token)
			case '+', '-', '*':
				token := string(c)
				if debugEnable {
					log.Printf("%s: token: '%v'", op, token)
				}
				tokens = append(tokens, token)
				prevIsExpr = false
			default:
				return nil, false
			}
			continue
		}

		if c == '(' {
			token := "("
			if debugEnable {
				log.Printf("%s: token: '%v'", op, token)
			}
			tokens = append(tokens, token)
			continue
		}

		if err := r.UnreadByte(); err != nil {
			panic(err)
		}
		token, ok := scanIntStr(r)
		if !ok {
			return nil, false
		}
		if debugEnable {
			log.Printf("%s: token: '%v'", op, token)
		}
		tokens = append(tokens, token)
		prevIsExpr = true
	}

	return tokens, true
}

func scanIntStr(r io.ByteScanner) (string, bool) {
	const op = "scanIntStr"
	var sb strings.Builder

	// может быть унарный +/-
	sign := 1
	for {
		c, err := r.ReadByte()
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}

		if c == '-' {
			sign = -sign
		} else if c != '+' {
			if err := r.UnreadByte(); err != nil {
				panic(err)
			}
			break
		}
	}

	if sign == -1 {
		sb.WriteByte('-')
	}

	count := 0
	for {
		c, err := r.ReadByte()
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

		count++
		sb.WriteByte(c)
	}

	if count == 0 {
		if debugEnable {
			log.Printf("%s: digit was expected", op)
		}
		return "", false
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

// чем больше значение, тем выше приоритет, 0 зарезервирован для '('
func getPriority(op string) int {
	switch op {
	case "(":
		return 0
	case "+", "-":
		return 100
	case "*":
		return 200
	default:
		panic(fmt.Errorf("unknown operation: %v", op))
	}
}

func toPolish(tokens []string) ([]string, bool) {
	const op = "toPolish"

	type item struct {
		token    string
		priority int
	}

	var (
		stack      stack[item]
		prevIsExpr bool // значение или выражение в скобках
	)
	res := make([]string, 0, len(tokens))

	for _, token := range tokens {
		if debugEnable {
			log.Printf("%s: token: %v", op, token)
		}
		switch token {
		case "(":
			if prevIsExpr {
				if debugEnable {
					log.Printf("%s: want prev not expr", op)
				}
				return nil, false
			}

			// всегда кладем на стек, с наименьшим приоритетом
			stack.push(item{token: "(", priority: getPriority("(")})

		case ")":
			if !prevIsExpr {
				if debugEnable {
					log.Printf("%s: want prev expr", op)
				}
				return nil, false
			}
			prevIsExpr = true

			// пререкладываем из стека в результат, пока не встертим скобку
			for {
				if stack.empty() {
					if debugEnable {
						log.Printf("%s: not balanced brackets", op)
					}
					return nil, false
				}

				token := stack.pop().token

				if token == "(" {
					break
				}

				res = append(res, token)
			}

		case "+", "-", "*":
			if !prevIsExpr {
				if debugEnable {
					log.Printf("%s: want prev expr", op)
				}
				return nil, false
			}
			prevIsExpr = false

			// вытаскивает со стека все, чей приоритет выше или равен
			prioryty := getPriority(token)
			for !stack.empty() && stack.top().priority >= prioryty {
				res = append(res, stack.pop().token)
			}

			// и ложится на стек
			stack.push(item{token: token, priority: prioryty})

		default: // иначе это значение
			if prevIsExpr {
				if debugEnable {
					log.Printf("%s: want prev not expr", op)
				}
				return nil, false
			}
			prevIsExpr = true

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
