package main

import "math/rand"

func genTreeLine(n int) (EdgeSet, Graph) {
	nodes := make([]Idx, n+1)

	for i := range nodes {
		nodes[i] = Idx(i)
	}

	rand.Shuffle(n, func(i, j int) {
		i++
		j++
		nodes[i], nodes[j] = nodes[j], nodes[i]
	})

	graph := make(Graph, n+1)
	edges := make(EdgeSet, n)

	for i := n; i > 1; i-- {
		b := nodes[i]
		a := nodes[i-1]
		edges.Add(Edge{a, b})
		graph[a] = append(graph[a], b)
		graph[b] = append(graph[b], a)
	}

	return edges, graph
}

func genTree(rand *rand.Rand, n int) (EdgeSet, Graph) {

	nodes := make([]Idx, n+1)

	for i := range nodes {
		nodes[i] = Idx(i)
	}

	rand.Shuffle(n, func(i, j int) {
		i++
		j++
		nodes[i], nodes[j] = nodes[j], nodes[i]
	})

	graph := make(Graph, n+1)
	edges := make(EdgeSet, n)

	for i := n; i > 1; i-- {
		b := nodes[i]
		a := nodes[rand.Intn(i-1)+1]
		edges.Add(Edge{a, b})
		graph[a] = append(graph[a], b)
		graph[b] = append(graph[b], a)
	}

	return edges, graph
}
