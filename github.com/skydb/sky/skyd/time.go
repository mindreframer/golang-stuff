package skyd

import (
	"time"
)

// Shifts a Go time into Sky timestamp format.
func ShiftTime(value time.Time) int64 {
	timestamp := value.UnixNano() / 1000
	usec := timestamp % 1000000
	sec := timestamp / 1000000
	return (sec << 20) + usec
}

// Shifts a Sky timestamp format into a Go time.
func UnshiftTime(value int64) time.Time {
	usec := value & 0xFFFFF
	sec := value >> 20
	return time.Unix(sec, usec*1000)
}
