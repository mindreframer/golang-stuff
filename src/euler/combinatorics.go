package euler

func Choose(N, K int64) int64 {
	factors := make(map[int64]int64)

	if K == 0 || N == K || N <= 1 {
		return 1
	}

	if N < K {
		return 0
	}

	for n := N; n > N-K; n-- {
		nfactors := Factors(n)
		for i := 0; i < len(nfactors); i++ {
			factors[nfactors[i][0]] += nfactors[i][1]
		}
	}

	for k := K; k >= 2; k-- {

		kfactors := Factors(k)
		for i := 0; i < len(kfactors); i++ {
			factors[kfactors[i][0]] -= kfactors[i][1]
		}

	}

	answer := int64(1)

	for prime, multiplicity := range factors {

		for i := int64(0); i < multiplicity; i++ {
			answer *= prime
		}

	}
	return answer
}

//returns the nth permutation of the given slice
func Permutation(n int, list []int) []int {
	if len(list) == 1 {
		return list
	}

	k := n % len(list)

	first := []int{list[k]}
	next := make([]int, len(list)-1)

	copy(next, append(list[:k], list[k+1:]...))

	return append(first, Permutation(n/len(list), next)...)

}

func PermuteFloats(n int, list []float64) []float64 {
	if len(list) == 1 {
		return list
	}

	k := n % len(list)

	first := []float64{list[k]}
	next := make([]float64, len(list)-1)

	copy(next, append(list[:k], list[k+1:]...))

	return append(first, PermuteFloats(n/len(list), next)...)

}

func PermuteString(n int, word string) string {
	if len(word) == 1 {
		return word
	}

	k := n % len(word)

	return word[k:k+1] + PermuteString(n/len(word), word[:k]+word[k+1:])
}

func SplitInts(list []int, K, N int) (a, b []int) {
	a, b = make([]int, 0), make([]int, 0)

	indices := make(map[int]bool)

	for k := K; k >= 1; k-- {

		n := k - 1

		if Choose(int64(n), int64(k)) <= int64(N) {
			for ; Choose(int64(n), int64(k)) <= int64(N); n++ {

			}
			n--
		}

		indices[n] = true

		N = N - int(Choose(int64(n), int64(k)))
	}

	a, b = make([]int, 0), make([]int, 0)

	for i := 0; i < len(list); i++ {
		if indices[i] {
			a = append(a, list[i])
		} else {
			b = append(b, list[i])
		}
	}

	return a, b
}

func SplitSeq(K, N int) (a []int) {

	indices := make([]int, 0)

	for k := K; k >= 1; k-- {

		n := k - 1

		if Choose(int64(n), int64(k)) <= int64(N) {
			for ; Choose(int64(n), int64(k)) <= int64(N); n++ {

			}
			n--
		}

		indices = append(indices, n)

		N = N - int(Choose(int64(n), int64(k)))
	}

	return indices
}

//returns which permutation takes a->b, or -1
//NOTE: THIS IS A TERRIBLE ALGORITHM -- Fix later
func UnPermuteStrings(a, b string) int {
	for i := 0; i < int(Factorial(int64(len(a)))); i++ {
		if PermuteString(i, a) == b {
			return i
		}
	}
	return -1

}
