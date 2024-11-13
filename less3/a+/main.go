package main

import (
	"bufio"
	"bytes"
	"io"
	"os"
)

type stack []byte

func (s stack) empty() bool {
	return len(s) == 0
}

func (s *stack) push(c byte) {
	*s = append(*s, c)
}

func (s *stack) pop() byte {
	old := *s
	n := len(old)
	c := old[n-1]
	*s = old[:n-1]
	return c
}

func solve(line []byte) bool {
	stack := make(stack, 0, len(line))

	for _, c := range line {
		switch c {
		case '(', '[', '{':
			stack.push(c)
		case ')':
			if stack.empty() || stack.pop() != '(' {
				return false
			}
		case ']':
			if stack.empty() || stack.pop() != '[' {
				return false
			}
		case '}':
			if stack.empty() || stack.pop() != '{' {
				return false
			}
		}
	}

	return stack.empty()
}

func run(in io.Reader, out io.Writer) {
	br := bufio.NewReader(in)
	bw := bufio.NewWriter(out)
	defer bw.Flush()

	line, err := io.ReadAll(br)
	if err != nil {
		panic(err)
	}

	line = bytes.TrimSpace(line)

	if solve(line) {
		bw.WriteString("yes\n")
	} else {
		bw.WriteString("no\n")
	}
}

// ----------------------------------------------------------------------------

var _, debugEnable = os.LookupEnv("DEBUG")

func main() {
	_ = debugEnable
	run(os.Stdin, os.Stdout)
}
