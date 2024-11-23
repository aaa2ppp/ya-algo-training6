package main

import "math/rand"

// func genTreeLine(n int) (Edges, Graph) {
// 	nodes := make([]int, n+1)

// 	for i := range nodes {
// 		nodes[i] = int(i)
// 	}

// 	rand.Shuffle(n, func(i, j int) {
// 		i++
// 		j++
// 		nodes[i], nodes[j] = nodes[j], nodes[i]
// 	})

// 	graph := make(Graph, n+1)
// 	edges := make(EdgeSet, n)

// 	for i := n; i > 1; i-- {
// 		b := nodes[i]
// 		a := nodes[i-1]
// 		edges.Add(Edge{a, b})
// 		graph[a] = append(graph[a], b)
// 		graph[b] = append(graph[b], a)
// 	}

// 	return edges, graph
// }

func genTree(rand *rand.Rand, n int) (Edges, Graph) {

	nodes := make([]int, n+1)

	for i := range nodes {
		nodes[i] = int(i)
	}

	rand.Shuffle(n, func(i, j int) {
		i++
		j++
		nodes[i], nodes[j] = nodes[j], nodes[i]
	})

	graph := make(Graph, n+1)
	edges := make(Edges, 0, n)

	for i := n; i > 1; i-- {
		b := nodes[i]
		a := nodes[rand.Intn(i-1)+1]
		if rand.Intn(2) == 1 {
			a, b = b, a
		}
		edges = append(edges, Edge{a, b})
		graph[a].down = append(graph[a].down, b)
		graph[b].up = append(graph[b].up, a)
	}

	return edges, graph
}
