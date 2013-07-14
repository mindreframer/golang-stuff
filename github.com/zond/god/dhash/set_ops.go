package dhash

import (
	"bytes"
	"github.com/zond/god/common"
	"github.com/zond/god/radix"
	"github.com/zond/setop"
)

const (
	setOpBufferSize = 128
)

type treeSkipper struct {
	key          []byte
	tree         *radix.Tree
	remote       common.Remote
	buffer       []setop.SetOpResult
	currentIndex int
}

func (self *treeSkipper) Skip(min []byte, inc bool) (result *setop.SetOpResult, err error) {
	lt := 1
	if inc {
		lt = 0
	}
	if len(self.buffer) == 0 {
		if err = self.refill(min, inc); err != nil {
			return
		}
		if len(self.buffer) == 0 {
			return
		}
	}
	result = &self.buffer[self.currentIndex]
	for min != nil && bytes.Compare(result.Key, min) < lt {
		if result, err = self.nextWithRefill(min, inc); result == nil || err != nil {
			return
		}
	}
	return
}

func (self *treeSkipper) nextWithRefill(refillMin []byte, inc bool) (result *setop.SetOpResult, err error) {
	result = self.nextFromBuf()
	if result == nil {
		if err = self.refill(refillMin, inc); err != nil {
			return
		}
		if len(self.buffer) == 0 {
			return
		}
		result = &self.buffer[self.currentIndex]
	}
	return
}

func (self *treeSkipper) nextFromBuf() (result *setop.SetOpResult) {
	self.currentIndex++
	if self.currentIndex < len(self.buffer) {
		result = &self.buffer[self.currentIndex]
		return
	}
	return
}

func (self *treeSkipper) refill(min []byte, inc bool) (err error) {
	self.buffer = make([]setop.SetOpResult, 0, setOpBufferSize)
	self.currentIndex = 0
	if self.tree == nil {
		if err = self.remoteRefill(min, inc); err != nil {
			return
		}
	} else {
		if err = self.treeRefill(min, inc); err != nil {
			return
		}
	}
	return
}

func (self *treeSkipper) remoteRefill(min []byte, inc bool) (err error) {
	r := common.Range{
		Key:    self.key,
		Min:    min,
		MinInc: inc,
		Len:    setOpBufferSize,
	}
	var items []common.Item
	if err = self.remote.Call("DHash.SliceLen", r, &items); err != nil {
		return
	}
	for _, item := range items {
		self.buffer = append(self.buffer, setop.SetOpResult{item.Key, [][]byte{item.Value}})
	}
	return
}

func (self *treeSkipper) treeRefill(min []byte, inc bool) error {
	filler := func(key, value []byte, timestamp int64) bool {
		self.buffer = append(self.buffer, setop.SetOpResult{key, [][]byte{value}})
		return len(self.buffer) < setOpBufferSize
	}
	self.tree.SubEachBetween(self.key, min, nil, inc, false, filler)
	return nil
}
