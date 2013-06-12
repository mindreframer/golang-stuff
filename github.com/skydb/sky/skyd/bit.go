package skyd

// Condenses the even bits of a 64-bit integer into 32 bits.
func CondenseUint64Even(value uint64) uint32 {
	return condenseUint64(value, 0)
}

// Condenses the odd bits of a 64-bit integer into 32 bits.
func CondenseUint64Odd(value uint64) uint32 {
	return condenseUint64(value, 1)
}

// Condenses the bits of a 64-bit integer into 32 bits.
func condenseUint64(value uint64, offset uint) uint32 {
	var i, j uint
	var ret uint32 = 0
	for i = 0; i < 4; i++ {
		var x uint16 = uint16((value >> (i * 16)) & 0xFFFF)
		var y uint8 = 0
		for j = 0; j < 8; j++ {
			y |= uint8(((x >> ((j * 2) + offset)) & 1) << j)
		}
		ret |= (uint32(y) << (i * 8))
	}
	return ret
}
