package main

import "math/rand"

func genTree(rand *rand.Rand, n int) (Graph, EdgeDirs) {

	nodes := make([]Node, n+1)

	for i := range nodes {
		nodes[i] = Node(i)
	}

	rand.Shuffle(n, func(i, j int) {
		i++
		j++
		nodes[i], nodes[j] = nodes[j], nodes[i]
	})

	adj := make(Graph, n+1)
	dir := make(EdgeDirs, n)

	for i := n; i > 1; i-- {
		b := nodes[i]
		a := nodes[rand.Intn(i-1)+1]
		if rand.Intn(2) == 1 {
			a, b = b, a
		}
		dir[Edge{a, b}] = true
		dir[Edge{b, a}] = false
		adj[a] = append(adj[a], b)
		adj[b] = append(adj[b], a)
	}

	return adj, dir
}
