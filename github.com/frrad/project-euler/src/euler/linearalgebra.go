package euler

//returns the square of a matrix
func SqrIntMatrix(A [][]int) [][]int {
	n := len(A)
	square := make([][]int, n)
	for i := 0; i < n; i++ {
		square[i] = make([]int, n)
	}

	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			answer := 0
			for k := 0; k < n; k++ {
				answer += A[j][k] * A[k][i]
			}
			square[j][i] = answer
		}
	}

	return square
}
