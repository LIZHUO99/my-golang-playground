package main

func numWays(n int) int {
	if n == 0 {
		return 1
	} else if n == 1 {
		return 1
	}
	a := 1
	b := 1
	for i := 1; i < n; i++ {
		a = a + b
		if a > (1e9 + 7) {
			a = a % (1e9 + 7)
			b = b % (1e9 + 7)
		}
		a, b = b, a
	}
	return b % (1e9 + 7)
}
