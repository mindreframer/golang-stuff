package euler

//reverses a slice of int64s (in place)
func ReverseLInts(list []int64) []int64 {
	for i := 0; i < len(list)/2; i++ {
		list[i], list[len(list)-1-i] = list[len(list)-1-i], list[i]
	}
	return list
}

func MaxFloat(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}
