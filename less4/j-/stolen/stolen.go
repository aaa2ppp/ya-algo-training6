package main

// stole solve from
// https://www.hackerrank.com/contests/hourrank-29/challenges/birthday-assignment/problem

func stolenSolve(graph Graph, dir EdgeDirs) int {
	n := len(graph)
	sz := make([]int, n)
	dp := makeMatrix[int](n, n)

	var dfs func(node, prev Node)

	dfs = func(node, prev Node) {
		sz[node] = 1
		var total_down, total_up int

		for _, neig := range graph[node] {
			if neig == prev {
				continue
			}

			dfs(neig, node)

			sz[node] += sz[neig]

			if dir[Edge{node, neig}] {
				total_down += sz[neig]
				for i := 1; i <= sz[neig]; i++ {
					dp[neig][i] += dp[neig][i-1]
					dp[neig][i] %= modulo
				}
			} else {
				total_up += sz[neig]
				for i := sz[neig]; i >= 1; i-- {
					dp[neig][i] += dp[neig][i+1]
					dp[neig][i] %= modulo
				}
			}
		}

		if sz[node] == 1 {
			dp[node][1] = 1
		} else {
			down_p := make([]int, total_down+1)
			down_n := make([]int, total_down+1)
			up_p := make([]int, total_up+1)
			up_n := make([]int, total_up+1)

			down_p[0] = 1
			up_p[0] = 1

			down_count := total_down
			up_count := total_up

			total_down := 0
			total_up := 0

			for _, neig := range graph[node] {
				if neig == prev {
					continue
				}

				if dir[Edge{node, neig}] {
					total_down += sz[neig]

					for i := 0; i <= down_count; i++ {
						if down_p[i] == 0 {
							continue
						}

						for j := 1; j <= sz[neig]; j++ {
							v := nCr(i+j, j) * dp[neig][j]
							v %= modulo
							v *= down_p[i]
							v %= modulo
							v *= nCr(total_down-(i+j), sz[neig]-j)
							v %= modulo
							down_n[i+j] += v
							down_n[i+j] %= modulo
						}

						down_p[i] = 0
					}

					down_p, down_n = down_n, down_p

				} else {
					total_up += sz[neig]

					for i := 0; i <= up_count; i++ {
						if up_p[i] == 0 {
							continue
						}

						for j := 1; j <= sz[neig]; j++ {
							v := nCr(i+j, j) * dp[neig][sz[neig]-j+1]
							v %= modulo
							v *= up_p[i]
							v %= modulo
							v *= nCr(total_up-(i+j), sz[neig]-j)
							v %= modulo
							up_n[i+j] += v
							up_n[i+j] %= modulo
						}

						up_p[i] = 0
					}

					up_p, up_n = up_n, up_p
				}
			}

			total_down = down_count
			total_up = up_count

			for i := 1; i <= sz[node]; i++ {
				dp[node][i] = 0

				for j := 0; j <= min(i-1, total_down); j++ {
					x := total_down - j
					v := down_p[j] % modulo
					v *= nCr(i-1, j)
					v %= modulo
					v *= nCr(sz[node]-i, x)
					v %= modulo

					if sz[node]-i-x >= 0 && sz[node]-i-x <= total_up {
						v *= up_p[sz[node]-i-x]
						v %= modulo
					} else {
						v = 0
					}

					dp[node][i] += v
					dp[node][1] %= modulo
				}
			}
		}
	}

	dfs(1, -1)

	count := 0
	for i := 1; i < n; i++ {
		count += dp[1][i]
		count %= modulo
	}

	return count
}

func nCr(n, r int) int {
	if r > n {
		return 0
	}
	return paskal(n, r)
}
