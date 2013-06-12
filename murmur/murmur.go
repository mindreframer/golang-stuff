package murmur

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/spaolacci/murmur3"
)

const (
	BlockSize = 1
	Size      = 16
	seed      = 42
)

// HashString will return the hash for the provided string.
func HashString(s string) []byte {
	return HashBytes([]byte(s))
}

// HashInt will return the hash for the provided int.
func HashInt(i int) []byte {
	b := new(bytes.Buffer)
	if err := binary.Write(b, binary.BigEndian, i); err != nil {
		panic(err)
	}
	return HashBytes(b.Bytes())
}

// HashInt64 will return the hash for the provided int64.
func HashInt64(i int64) []byte {
	b := new(bytes.Buffer)
	if err := binary.Write(b, binary.BigEndian, i); err != nil {
		panic(err)
	}
	return HashBytes(b.Bytes())
}

// HashBytes will return the hash for the provided byte slice.
func HashBytes(b []byte) []byte {
	h1, h2 := murmur3.Sum128(b)
	return []byte{
		byte(h1 >> 56), byte(h1 >> 48), byte(h1 >> 40), byte(h1 >> 32),
		byte(h1 >> 24), byte(h1 >> 16), byte(h1 >> 8), byte(h1),
		byte(h2 >> 56), byte(h2 >> 48), byte(h2 >> 40), byte(h2 >> 32),
		byte(h2 >> 24), byte(h2 >> 16), byte(h2 >> 8), byte(h2),
	}
}

type Hash bytes.Buffer

func New() *Hash {
	return new(Hash)
}

func (self *Hash) MustWrite(b []byte) {
	n, err := (*bytes.Buffer)(self).Write(b)
	if n != len(b) || err != nil {
		panic(fmt.Errorf("Wanted to write %v bytes, but wrote %v and got %v", len(b), n, err))
	}
}

func (self *Hash) Get() []byte {
	return HashBytes((*bytes.Buffer)(self).Bytes())
}

func NewBytes(b []byte) *Hash {
	return (*Hash)(bytes.NewBuffer(b))
}

func NewString(s string) *Hash {
	return (*Hash)(bytes.NewBufferString(s))
}

func (self *Hash) Write(b []byte) (int, error) {
	self.MustWrite(b)
	return len(b), nil
}

func (self *Hash) Extrude(result []byte) {
	h1, h2 := murmur3.Sum128((*bytes.Buffer)(self).Bytes())
	result[0] = byte(h1 >> 56)
	result[1] = byte(h1 >> 48)
	result[2] = byte(h1 >> 40)
	result[3] = byte(h1 >> 32)
	result[4] = byte(h1 >> 24)
	result[5] = byte(h1 >> 16)
	result[6] = byte(h1 >> 8)
	result[7] = byte(h1)
	result[8] = byte(h2 >> 56)
	result[9] = byte(h2 >> 48)
	result[10] = byte(h2 >> 40)
	result[11] = byte(h2 >> 32)
	result[12] = byte(h2 >> 24)
	result[13] = byte(h2 >> 16)
	result[14] = byte(h2 >> 8)
	result[15] = byte(h2)
}
