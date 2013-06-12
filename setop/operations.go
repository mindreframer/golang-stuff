package setop

import (
	"bytes"
	"fmt"
)

type Skipper interface {
	// skip returns a value matching the min and inclusive criteria.
	// If the last yielded value matches the criteria the same value will be returned again.
	Skip(min []byte, inc bool) (result *SetOpResult, err error)
}

func createSkippersAndWeights(r RawSourceCreator, sources []SetOpSource) (skippers []Skipper, weights []float64) {
	skippers = make([]Skipper, len(sources))
	weights = make([]float64, len(sources))
	for index, source := range sources {
		if source.Key != nil {
			skippers[index] = r(source.Key)
		} else {
			skippers[index] = createSkipper(r, source.SetOp)
		}
		if source.Weight != nil {
			weights[index] = *source.Weight
		} else {
			weights[index] = 1
		}
	}
	return
}

func createSkipper(r RawSourceCreator, op *SetOp) (result Skipper) {
	skippers, weights := createSkippersAndWeights(r, op.Sources)
	switch op.Type {
	case Union:
		result = &unionOp{
			skippers: skippers,
			weights:  weights,
			merger:   getMerger(op.Merge),
		}
	case Intersection:
		result = &interOp{
			skippers: skippers,
			weights:  weights,
			merger:   getMerger(op.Merge),
		}
	case Difference:
		result = &diffOp{
			skippers: skippers,
			weights:  weights,
			merger:   getMerger(op.Merge),
		}
	case Xor:
		result = &xorOp{
			skippers: skippers,
			weights:  weights,
			merger:   getMerger(op.Merge),
		}
	default:
		panic(fmt.Errorf("Unknown SetOp Type %v", op.Type))
	}
	return
}

type xorOp struct {
	skippers []Skipper
	weights  []float64
	curr     *SetOpResult
	merger   mergeFunc
}

func (self *xorOp) Skip(min []byte, inc bool) (result *SetOpResult, err error) {
	gt := 0
	if inc {
		gt = -1
	}

	if self.curr != nil && bytes.Compare(self.curr.Key, min) > gt {
		result = self.curr
		return
	}

	newSkippers := make([]Skipper, 0, len(self.skippers))

	var res *SetOpResult
	var cmp int
	var multi bool

	for result == nil {
		for index, thisSkipper := range self.skippers {
			if res, err = thisSkipper.Skip(min, inc); err != nil {
				result = nil
				self.curr = nil
				return
			}
			if res != nil {
				newSkippers = append(newSkippers, thisSkipper)
				if result == nil {
					result = res.ShallowCopy()
					result.Values = self.merger(nil, result.Values, self.weights[index])
					multi = false
				} else {
					cmp = bytes.Compare(res.Key, result.Key)
					if cmp < 0 {
						multi = false
						result = res.ShallowCopy()
						result.Values = self.merger(nil, result.Values, self.weights[index])
					} else if cmp == 0 {
						multi = true
					}
				}
			}
		}

		if len(newSkippers) == 0 {
			result = nil
			self.curr = nil
			return
		}

		if result != nil && multi {
			min = result.Key
			inc = false
			result = nil
		}

		self.skippers = newSkippers
		newSkippers = newSkippers[:0]

	}

	self.curr = result

	return
}

type unionOp struct {
	skippers []Skipper
	weights  []float64
	curr     *SetOpResult
	merger   mergeFunc
}

func (self *unionOp) Skip(min []byte, inc bool) (result *SetOpResult, err error) {
	gt := 0
	if inc {
		gt = -1
	}

	if self.curr != nil && bytes.Compare(self.curr.Key, min) > gt {
		result = self.curr
		return
	}

	newSkippers := make([]Skipper, 0, len(self.skippers))

	var cmp int
	var res *SetOpResult

	for index, thisSkipper := range self.skippers {
		if res, err = thisSkipper.Skip(min, inc); err != nil {
			result = nil
			self.curr = nil
			return
		}
		if res != nil {
			newSkippers = append(newSkippers, thisSkipper)
			if result == nil {
				result = res.ShallowCopy()
				result.Values = self.merger(nil, result.Values, self.weights[index])
			} else {
				cmp = bytes.Compare(res.Key, result.Key)
				if cmp < 0 {
					result = res.ShallowCopy()
					result.Values = self.merger(nil, result.Values, self.weights[index])
				} else if cmp == 0 {
					result.Values = self.merger(result.Values, res.Values, self.weights[index])
				}
			}
		}
	}

	self.skippers = newSkippers

	self.curr = result

	return
}

type interOp struct {
	skippers []Skipper
	weights  []float64
	curr     *SetOpResult
	merger   mergeFunc
}

func (self *interOp) Skip(min []byte, inc bool) (result *SetOpResult, err error) {
	gt := 0
	if inc {
		gt = -1
	}

	if self.curr != nil && bytes.Compare(self.curr.Key, min) > gt {
		result = self.curr
		return
	}

	var maxKey []byte
	var res *SetOpResult
	var cmp int

	for result == nil {
		maxKey = nil
		for index, thisSkipper := range self.skippers {
			if res, err = thisSkipper.Skip(min, inc); res == nil || err != nil {
				result = nil
				self.curr = nil
				return
			}
			if maxKey == nil {
				maxKey = res.Key
				result = res.ShallowCopy()
				result.Values = self.merger(nil, result.Values, self.weights[index])
			} else {
				cmp = bytes.Compare(res.Key, maxKey)
				if cmp != 0 {
					if cmp > 0 {
						maxKey = res.Key
					}
					result = nil
				} else {
					result.Values = self.merger(result.Values, res.Values, self.weights[index])
				}
			}
		}

		min = maxKey
		inc = true
	}

	self.curr = result

	return
}

type diffOp struct {
	skippers []Skipper
	weights  []float64
	curr     *SetOpResult
	merger   mergeFunc
}

func (self *diffOp) Skip(min []byte, inc bool) (result *SetOpResult, err error) {
	gt := 0
	if inc {
		gt = -1
	}

	if self.curr != nil && bytes.Compare(self.curr.Key, min) > gt {
		result = self.curr
		return
	}

	var newSkippers = make([]Skipper, 0, len(self.skippers))
	var res *SetOpResult

	for result == nil {
		for index, thisSkipper := range self.skippers {
			if res, err = thisSkipper.Skip(min, inc); err != nil {
				result = nil
				self.curr = nil
				return
			}
			if index == 0 {
				if res == nil {
					result = nil
					self.curr = nil
					return
				}
				result = res.ShallowCopy()
				result.Values = self.merger(nil, result.Values, self.weights[0])
				newSkippers = append(newSkippers, thisSkipper)
				min = res.Key
				inc = true
			} else {
				if res != nil {
					newSkippers = append(newSkippers, thisSkipper)
					if bytes.Compare(min, res.Key) == 0 {
						result = nil
						break
					}
				}
			}
		}

		self.skippers = newSkippers
		newSkippers = newSkippers[:0]
		inc = false

	}

	self.curr = result

	return
}
