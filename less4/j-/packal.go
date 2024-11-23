package main

import "log"

var _paskal = [][]int{{1}, {1, 1}}

func paskal(i, j int) int {
	if debugEnable {
		log.Printf("paskal: %d %d", i, j)
	}

	for i >= len(_paskal) {
		n := len(_paskal)
		row := make([]int, n+1)
		row[0] = 1
		row[n] = 1
		prev := _paskal[n-1]

		for j, v := range prev[:n-1] {
			row[j+1] = (v + prev[j+1]) % modulo
		}

		if debugEnable {
			log.Printf("packal: add row %d %v", n, row)
		}

		_paskal = append(_paskal, row)
	}

	return _paskal[i][j]
}
