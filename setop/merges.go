package setop

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/zond/god/common"
	"math/big"
)

type mergeFunc func(oldValues [][]byte, newValues [][]byte, weightForNewValues float64) (result [][]byte)

func getMerger(m SetOpMerge) mergeFunc {
	switch m {
	case Append:
		return _append
	case ConCat:
		return conCat
	case IntegerSum:
		return integerSum
	case IntegerDiv:
		return integerDiv
	case IntegerMul:
		return integerMul
	case FloatSum:
		return floatSum
	case FloatDiv:
		return floatDiv
	case FloatMul:
		return floatMul
	case BigIntAnd:
		return bigIntAnd
	case BigIntAdd:
		return bigIntAdd
	case BigIntAndNot:
		return bigIntAndNot
	case BigIntDiv:
		return bigIntDiv
	case BigIntMod:
		return bigIntMod
	case BigIntMul:
		return bigIntMul
	case BigIntOr:
		return bigIntOr
	case BigIntRem:
		return bigIntRem
	case BigIntXor:
		return bigIntXor
	case First:
		return first
	case Last:
		return last
	}
	panic(fmt.Errorf("Unknown SetOpType %v", int(m)))
}

func last(oldValues [][]byte, newValues [][]byte, w float64) (result [][]byte) {
	return [][]byte{newValues[len(newValues)-1]}
}
func first(oldValues [][]byte, newValues [][]byte, w float64) (result [][]byte) {
	if oldValues == nil {
		return [][]byte{newValues[0]}
	}
	return [][]byte{oldValues[0]}
}
func _append(oldValues [][]byte, newValues [][]byte, w float64) (result [][]byte) {
	result = oldValues
	for i := 0; i < int(w); i++ {
		result = append(result, newValues...)
	}
	return
}
func conCat(oldValues [][]byte, newValues [][]byte, w float64) (result [][]byte) {
	var res []byte
	for _, b := range oldValues {
		res = append(res, b...)
	}
	for i := 0; i < int(w); i++ {
		for _, b := range newValues {
			res = append(res, b...)
		}
	}
	return [][]byte{res}
}
func integerSum(oldValues [][]byte, newValues [][]byte, w float64) (result [][]byte) {
	var sum int64
	var tmp int64
	var err error
	for _, b := range oldValues {
		if tmp, err = common.DecodeInt64(b); err == nil {
			sum += tmp
		}
	}
	for _, b := range newValues {
		if tmp, err = common.DecodeInt64(b); err == nil {
			sum += (tmp * int64(w))
		}
	}
	res := new(bytes.Buffer)
	binary.Write(res, binary.BigEndian, sum)
	return [][]byte{res.Bytes()}
}
func integerDiv(oldValues [][]byte, newValues [][]byte, w float64) (result [][]byte) {
	var sum int64
	var tmp int64
	var err error
	if oldValues != nil {
		if tmp, err = common.DecodeInt64(oldValues[0]); err == nil {
			sum = tmp
		}
		for _, b := range newValues {
			if tmp, err = common.DecodeInt64(b); err == nil {
				sum /= (tmp * int64(w))
			}
		}
	} else {
		if tmp, err = common.DecodeInt64(newValues[0]); err == nil {
			sum = (tmp * int64(w))
		}
		for _, b := range newValues[1:] {
			if tmp, err = common.DecodeInt64(b); err == nil {
				sum /= (tmp * int64(w))
			}
		}
	}
	res := new(bytes.Buffer)
	binary.Write(res, binary.BigEndian, sum)
	return [][]byte{res.Bytes()}
}
func integerMul(oldValues [][]byte, newValues [][]byte, w float64) (result [][]byte) {
	var sum int64 = 1
	var tmp int64
	var err error
	for _, b := range oldValues {
		if tmp, err = common.DecodeInt64(b); err == nil {
			sum *= tmp
		}
	}
	for _, b := range newValues {
		if tmp, err = common.DecodeInt64(b); err == nil {
			sum *= (tmp * int64(w))
		}
	}
	res := new(bytes.Buffer)
	binary.Write(res, binary.BigEndian, sum)
	return [][]byte{res.Bytes()}
}
func floatSum(oldValues [][]byte, newValues [][]byte, w float64) (result [][]byte) {
	var sum float64
	var tmp float64
	var err error
	for _, b := range oldValues {
		if tmp, err = common.DecodeFloat64(b); err == nil {
			sum += tmp
		}
	}
	for _, b := range newValues {
		if tmp, err = common.DecodeFloat64(b); err == nil {
			sum += (tmp * w)
		}
	}
	res := new(bytes.Buffer)
	binary.Write(res, binary.BigEndian, sum)
	return [][]byte{res.Bytes()}
}
func floatDiv(oldValues [][]byte, newValues [][]byte, w float64) (result [][]byte) {
	var sum float64
	var tmp float64
	var err error
	if oldValues != nil {
		if tmp, err = common.DecodeFloat64(oldValues[0]); err == nil {
			sum = tmp
		}
		for _, b := range newValues {
			if tmp, err = common.DecodeFloat64(b); err == nil {
				sum /= (tmp * w)
			}
		}
	} else {
		if tmp, err = common.DecodeFloat64(newValues[0]); err == nil {
			sum = (tmp * w)
		}
		for _, b := range newValues[1:] {
			if tmp, err = common.DecodeFloat64(b); err == nil {
				sum /= (tmp * w)
			}
		}
	}
	res := new(bytes.Buffer)
	binary.Write(res, binary.BigEndian, sum)
	return [][]byte{res.Bytes()}
}
func floatMul(oldValues [][]byte, newValues [][]byte, w float64) (result [][]byte) {
	var sum float64 = 1
	var tmp float64
	var err error
	for _, b := range oldValues {
		if tmp, err = common.DecodeFloat64(b); err == nil {
			sum *= tmp
		}
	}
	for _, b := range newValues {
		if tmp, err = common.DecodeFloat64(b); err == nil {
			sum *= (tmp * w)
		}
	}
	res := new(bytes.Buffer)
	binary.Write(res, binary.BigEndian, sum)
	return [][]byte{res.Bytes()}
}
func bigIntAnd(oldValues [][]byte, newValues [][]byte, w float64) (result [][]byte) {
	var sum *big.Int
	if oldValues != nil {
		sum = new(big.Int).SetBytes(oldValues[0])
		for _, b := range newValues {
			sum.And(sum, new(big.Int).Mul(common.DecodeBigInt(b), big.NewInt(int64(w))))
		}
	} else {
		sum = new(big.Int).Mul(new(big.Int).SetBytes(newValues[0]), big.NewInt(int64(w)))
		for _, b := range newValues[1:] {
			sum.And(sum, new(big.Int).Mul(common.DecodeBigInt(b), big.NewInt(int64(w))))
		}
	}
	return [][]byte{sum.Bytes()}
}
func bigIntAdd(oldValues [][]byte, newValues [][]byte, w float64) (result [][]byte) {
	sum := new(big.Int)
	for _, b := range oldValues {
		sum.Add(sum, common.DecodeBigInt(b))
	}
	for _, b := range newValues {
		sum.Add(sum, new(big.Int).Mul(common.DecodeBigInt(b), big.NewInt(int64(w))))
	}
	return [][]byte{sum.Bytes()}
}
func bigIntAndNot(oldValues [][]byte, newValues [][]byte, w float64) (result [][]byte) {
	var sum *big.Int
	if oldValues != nil {
		sum = new(big.Int).SetBytes(oldValues[0])
		for _, b := range newValues {
			sum.AndNot(sum, new(big.Int).Mul(common.DecodeBigInt(b), big.NewInt(int64(w))))
		}
	} else {
		sum = new(big.Int).Mul(new(big.Int).SetBytes(newValues[0]), big.NewInt(int64(w)))
		for _, b := range newValues[1:] {
			sum.AndNot(sum, new(big.Int).Mul(common.DecodeBigInt(b), big.NewInt(int64(w))))
		}
	}
	return [][]byte{sum.Bytes()}
}
func bigIntDiv(oldValues [][]byte, newValues [][]byte, w float64) (result [][]byte) {
	var sum *big.Int
	if oldValues != nil {
		sum = new(big.Int).SetBytes(oldValues[0])
		for _, b := range newValues {
			sum.Div(sum, new(big.Int).Mul(common.DecodeBigInt(b), big.NewInt(int64(w))))
		}
	} else {
		sum = new(big.Int).Mul(new(big.Int).SetBytes(newValues[0]), big.NewInt(int64(w)))
		for _, b := range newValues[1:] {
			sum.Div(sum, new(big.Int).Mul(common.DecodeBigInt(b), big.NewInt(int64(w))))
		}
	}
	return [][]byte{sum.Bytes()}
}
func bigIntMod(oldValues [][]byte, newValues [][]byte, w float64) (result [][]byte) {
	var sum *big.Int
	if oldValues != nil {
		sum = new(big.Int).SetBytes(oldValues[0])
		for _, b := range newValues {
			sum.Mod(sum, new(big.Int).Mul(common.DecodeBigInt(b), big.NewInt(int64(w))))
		}
	} else {
		sum = new(big.Int).Mul(new(big.Int).SetBytes(newValues[0]), big.NewInt(int64(w)))
		for _, b := range newValues[1:] {
			sum.Mod(sum, new(big.Int).Mul(common.DecodeBigInt(b), big.NewInt(int64(w))))
		}
	}
	return [][]byte{sum.Bytes()}
}
func bigIntMul(oldValues [][]byte, newValues [][]byte, w float64) (result [][]byte) {
	sum := big.NewInt(1)
	for _, b := range oldValues {
		sum.Mul(sum, common.DecodeBigInt(b))
	}
	for _, b := range newValues {
		sum.Mul(sum, new(big.Int).Mul(common.DecodeBigInt(b), big.NewInt(int64(w))))
	}
	return [][]byte{sum.Bytes()}
}
func bigIntOr(oldValues [][]byte, newValues [][]byte, w float64) (result [][]byte) {
	var sum *big.Int
	if oldValues != nil {
		sum = new(big.Int).SetBytes(oldValues[0])
		for _, b := range newValues {
			sum.Or(sum, new(big.Int).Mul(common.DecodeBigInt(b), big.NewInt(int64(w))))
		}
	} else {
		sum = new(big.Int).Mul(new(big.Int).SetBytes(newValues[0]), big.NewInt(int64(w)))
		for _, b := range newValues[1:] {
			sum.Or(sum, new(big.Int).Mul(common.DecodeBigInt(b), big.NewInt(int64(w))))
		}
	}
	return [][]byte{sum.Bytes()}
}
func bigIntRem(oldValues [][]byte, newValues [][]byte, w float64) (result [][]byte) {
	var sum *big.Int
	if oldValues != nil {
		sum = new(big.Int).SetBytes(oldValues[0])
		for _, b := range newValues {
			sum.Rem(sum, new(big.Int).Mul(common.DecodeBigInt(b), big.NewInt(int64(w))))
		}
	} else {
		sum = new(big.Int).Mul(new(big.Int).SetBytes(newValues[0]), big.NewInt(int64(w)))
		for _, b := range newValues[1:] {
			sum.Rem(sum, new(big.Int).Mul(common.DecodeBigInt(b), big.NewInt(int64(w))))
		}
	}
	return [][]byte{sum.Bytes()}
}
func bigIntXor(oldValues [][]byte, newValues [][]byte, w float64) (result [][]byte) {
	var sum *big.Int
	if oldValues != nil {
		sum = new(big.Int).SetBytes(oldValues[0])
		for _, b := range newValues {
			sum.Xor(sum, new(big.Int).Mul(common.DecodeBigInt(b), big.NewInt(int64(w))))
		}
	} else {
		sum = new(big.Int).Mul(new(big.Int).SetBytes(newValues[0]), big.NewInt(int64(w)))
		for _, b := range newValues[1:] {
			sum.Xor(sum, new(big.Int).Mul(common.DecodeBigInt(b), big.NewInt(int64(w))))
		}
	}
	return [][]byte{sum.Bytes()}
}
