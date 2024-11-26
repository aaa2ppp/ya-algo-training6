package main

import "math/rand"

func genEdges(rand *rand.Rand, n int) []Edge {
	nodes := make([]Idx, n+1)

	for i := range nodes {
		nodes[i] = Idx(i)
	}

	rand.Shuffle(n, func(i, j int) {
		i++
		j++
		nodes[i], nodes[j] = nodes[j], nodes[i]
	})

	edges := make([]Edge, 0, n)

	for i := n; i > 1; i-- {
		b := nodes[i]
		a := nodes[rand.Intn(i-1)+1]
		if rand.Intn(2) == 1 {
			a, b = b, a
		}
		edges = append(edges, Edge{a, b})
	}

	return edges
}

func genTree(rand *rand.Rand, n int) (_Graph, EdgeDirs) {

	edges := genEdges(rand, n)

	graph := make(_Graph, n+1)
	dirs := make(EdgeDirs, n)

	for _, e := range edges {
		a := e[0]
		b := e[1]
		dirs[Edge{a, b}] = true
		dirs[Edge{b, a}] = false
		graph[a] = append(graph[a], b)
		graph[b] = append(graph[b], a)
	}

	return graph, dirs
}
