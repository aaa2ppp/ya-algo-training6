def solve(aa, xx, k):
	n, m = len(aa), len(xx)

	bb = [0]*n
	l, r, count = 0, 1, 0

	while r < n:
		if aa[r] < aa[r-1]:
			l = r
			count = 0
		elif aa[r] == aa[r-1]:
			count += 1

		while count > k:
			l +=1
			if aa[l] == aa[l-1]:
				count -= 1

		bb[r] = l
		r += 1

	res = []
	for i in xx:
		res.append(bb[i-1]+1) # to 0-indexing, to 1-indexing

	return res


def main():
	n = int(input())
	aa = list(map(int, input().split()))
	m, k = map(int, input().split())
	xx = list(map(int, input().split()))
	res = solve(aa, xx, k)
	print(*res)


if __name__ == '__main__':
    main()
