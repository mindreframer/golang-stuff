package common

type Range struct {
	Key      []byte
	Min      []byte
	Max      []byte
	MinInc   bool
	MaxInc   bool
	MinIndex int
	MaxIndex int
	Len      int
}
