package main

import (
	"fmt"
	"log"
)

type NodeSizes struct {
	up   int
	down int
}

func slowSolve(edges Edges, graph Graph) int {
	if debugEnable {
		log.Println("=> slowSolve")
	}

	n := len(graph) - 1 // 1-indexing

	sizes := make([]NodeSizes, n+1)

	countUpSizes(graph, sizes)
	countDownSizes(graph, sizes)

	if debugEnable {
		log.Printf("sizes: %v", sizes)
	}

	// count permutations with bit matrix

	matrix := make(map[int32]int, 0)
	curMatrix := make(map[int32]int, 0)
	matrix[0] = 1

	for node := 1; node <= n; node++ {
		clear(curMatrix)

		begin := sizes[node].up
		end := n - sizes[node].down

		for k, v := range matrix {
			for j := begin; j < end; j++ {
				if k&(1<<j) == 0 {
					x := curMatrix[k|(1<<j)]
					x += v
					x %= modulo
					curMatrix[k|(1<<j)] = x
				}
			}
		}

		matrix, curMatrix = curMatrix, matrix

		if debugEnable {
			logBitMatrix(fmt.Sprintf("%d [%d..%d]", node, begin, end-1), matrix, n)
		}
	}

	if len(matrix) != 1 {
		panic(fmt.Errorf("len of prev must be 1, got %d", len(matrix)))
	}

	var count int
	for _, v := range matrix {
		count = v
		break
	}

	return count
}

func logBitMatrix(label string, matrix map[int32]int, n int) {
	log.Printf("%s:\n", label)
	i := 0
	for k, v := range matrix {
		buf := make([]byte, n)
		for i := 0; i < n; i++ {
			if k&(1<<i) != 0 {
				buf[i] = '1'
			} else {
				buf[i] = '0'
			}
		}
		log.Printf("  %d: %s %d\n", i, buf, v)
		i++
	}
}

func countUpSizes(graph Graph, sizes []NodeSizes) {
	type FrontierItem struct {
		node   int
		heigth int
	}

	frontier := make([]FrontierItem, 0) // TODO optimize prealloc
	count := make([]int, len(graph))

	for node := 1; node < len(graph); node++ {
		count[node] = len(graph[node].up)
		if count[node] == 0 {
			count[node] = 1
			frontier = append(frontier, FrontierItem{node: node, heigth: 0})
		}
	}

	if debugEnable {
		log.Println("frontier:", frontier)
	}

	for len(frontier) > 0 {
		it := frontier[0]
		frontier = frontier[1:]

		sizes[it.node].up += it.heigth
		count[it.node]--

		if count[it.node] == 0 {
			for _, neig := range graph[it.node].down {
				frontier = append(frontier, FrontierItem{
					node:   neig,
					heigth: sizes[it.node].up + 1,
				})
			}
		}
	}
}

func countDownSizes(graph Graph, sizes []NodeSizes) {
	type FrontierItem struct {
		node   int
		heigth int
	}

	frontier := make([]FrontierItem, 0) // TODO optimize prealloc
	count := make([]int, len(graph))

	for node := 1; node < len(graph); node++ {
		count[node] = len(graph[node].down)
		if count[node] == 0 {
			count[node] = 1
			frontier = append(frontier, FrontierItem{node: node, heigth: 0})
		}
	}

	for len(frontier) > 0 {
		it := frontier[0]
		frontier = frontier[1:]

		sizes[it.node].down += it.heigth
		count[it.node]--

		if count[it.node] == 0 {
			for _, neig := range graph[it.node].up {
				frontier = append(frontier, FrontierItem{
					node:   neig,
					heigth: sizes[it.node].down + 1,
				})
			}
		}

	}
}
