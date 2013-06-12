package timenet

import (
	"math"
)

type times []int64

func (self times) stats() (mean, deviation int64) {
	if len(self) > 1 {
		var sum int64
		for _, latency := range self {
			sum += latency
		}
		mean = sum / int64(len(self))
		var squareSum int64
		var diff int64
		for _, latency := range self {
			diff = (latency - mean) >> 10
			squareSum += diff * diff
		}
		deviation = int64(math.Sqrt(float64(squareSum/(int64(len(self))-1)))) << 10
	} else {
		mean, deviation = -1, -1
	}
	return
}
