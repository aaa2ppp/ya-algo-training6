package main

import (
	"log"
	"slices"
)

func slowSolve(edges EdgeSet, graph Graph) int {
	maximum := 0

	for edge := range edges {
		graph := cloneGraph(graph)

		a := edge[0]
		b := edge[1]

		cutEdge(graph, a, b)
		if debugEnable {
			log.Printf("%d-%d: graph: %v", a, b, graph)
		}

		xa := calcDiameter(graph, a)
		xb := calcDiameter(graph, b)
		if debugEnable {
			log.Printf("%d-%d: %d %d", a, b, xa, xb)
		}
		maximum = max(maximum, xa*xb)

		// graph[a] = append(graph[a], b)
		// graph[b] = append(graph[b], a)
	}

	return maximum
}

func cloneGraph(graph Graph) Graph {
	n := len(graph)
	g := make(Graph, n)
	for i := 1; i < n; i++ {
		g[i] = slices.Clone(graph[i])
	}
	return g
}

func cutEdge(graph Graph, a, b Idx) {

	for i := len(graph[a]) - 1; i >= 0; i-- {
		if graph[a][i] == b {
			n := len(graph[a])
			graph[a][i] = graph[a][n-1]
			graph[a] = graph[a][:n-1]
		}
	}

	for i := len(graph[b]) - 1; i >= 0; i-- {
		if graph[b][i] == a {
			n := len(graph[b])
			graph[b][i] = graph[b][n-1]
			graph[b] = graph[b][:n-1]
		}
	}
}

func calcDiameter(graph Graph, node Idx) int {
	var dfs func(node, prev Idx) (height int, deep Idx)

	dfs = func(node, prev Idx) (height int, deep Idx) {
		deep = node

		for _, neig := range graph[node] {
			if neig == prev {
				continue
			}

			h, d := dfs(neig, node)

			if h >= height {
				height = h + 1
				deep = d
			}
		}

		return height, deep
	}

	_, node = dfs(node, 0)
	h, _ := dfs(node, 0)

	return h
}
