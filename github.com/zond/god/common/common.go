package common

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"runtime"
	"sort"
	"strconv"
	"testing"
	"time"
)

const (
	PingInterval = time.Second
)

var (
	Redundancy int = 3
)

func SetRedundancy(r int) {
	Redundancy = r
}

func MustParseFloat64(s string) (result float64) {
	var err error
	if result, err = strconv.ParseFloat(s, 64); err != nil {
		panic(err)
	}
	return
}

func MustJSONEncode(i interface{}) []byte {
	result, err := json.Marshal(i)
	if err != nil {
		panic(err)
	}
	return result
}
func MustJSONDecode(b []byte, i interface{}) {
	err := json.Unmarshal(b, i)
	if err != nil {
		panic(err)
	}
}
func EncodeBigInt(i *big.Int) []byte {
	return i.Bytes()
}
func DecodeBigInt(b []byte) (result *big.Int) {
	result = new(big.Int).SetBytes(b)
	return
}
func EncodeInt64(i int64) []byte {
	result := new(bytes.Buffer)
	if err := binary.Write(result, binary.BigEndian, i); err != nil {
		panic(err)
	}
	return result.Bytes()
}
func MustDecodeInt64(b []byte) (result int64) {
	result, err := DecodeInt64(b)
	if err != nil {
		panic(err)
	}
	return
}
func DecodeInt64(b []byte) (result int64, err error) {
	err = binary.Read(bytes.NewBuffer(b), binary.BigEndian, &result)
	return
}
func EncodeFloat64(f float64) []byte {
	result := new(bytes.Buffer)
	if err := binary.Write(result, binary.BigEndian, f); err != nil {
		panic(err)
	}
	return result.Bytes()
}
func MustDecodeFloat64(b []byte) (result float64) {
	result, err := DecodeFloat64(b)
	if err != nil {
		panic(err)
	}
	return
}
func DecodeFloat64(b []byte) (result float64, err error) {
	err = binary.Read(bytes.NewBuffer(b), binary.BigEndian, &result)
	return
}

func Max64(i ...int64) (result int64) {
	for _, x := range i {
		if x > result {
			result = x
		}
	}
	return
}

func Min64(i ...int64) (result int64) {
	result = i[0]
	for _, x := range i {
		if x < result {
			result = x
		}
	}
	return
}

func Max(i ...int) (result int) {
	for _, x := range i {
		if x > result {
			result = x
		}
	}
	return
}

func Min(i ...int) (result int) {
	result = i[0]
	for _, x := range i {
		if x < result {
			result = x
		}
	}
	return
}

// AssertWithin asserts that the given func returns a true bool within d, or it fires an error through t.
func AssertWithin(t *testing.T, f func() (string, bool), d time.Duration) {
	deadline := time.Now().Add(d)
	var ok bool
	var msg string
	for time.Now().Before(deadline) {
		if msg, ok = f(); ok {
			return
		}
		time.Sleep(time.Second / 5)
	}
	var file string
	var line int
	_, file, line, _ = runtime.Caller(1)
	t.Errorf("%v:%v: Wanted %v to be true within %v, but it never happened: %v", file, line, f, d, msg)
}

// HexEncode will encode the given bytes to a string, and pad it with 0 until it is at least twice the length of b.
func HexEncode(b []byte) (result string) {
	encoded := hex.EncodeToString(b)
	buffer := new(bytes.Buffer)
	for i := len(encoded); i < len(b)*2; i++ {
		fmt.Fprint(buffer, "00")
	}
	fmt.Fprint(buffer, encoded)
	return string(buffer.Bytes())
}

// BetweenII returns whether needle is between fromInc and toInc, inclusive.
func BetweenII(needle, fromInc, toInc []byte) (result bool) {
	switch bytes.Compare(fromInc, toInc) {
	case 0:
		result = true
	case -1:
		result = bytes.Compare(fromInc, needle) < 1 && bytes.Compare(needle, toInc) < 1
	case 1:
		result = bytes.Compare(fromInc, needle) < 1 || bytes.Compare(needle, toInc) < 1
	default:
		panic("Shouldn't happen")
	}
	return
}

// BetweenIE returns whether needle is between fromInc, inclusive and toExc, exclusive.
func BetweenIE(needle, fromInc, toExc []byte) (result bool) {
	switch bytes.Compare(fromInc, toExc) {
	case 0:
		result = true
	case -1:
		result = bytes.Compare(fromInc, needle) < 1 && bytes.Compare(needle, toExc) < 0
	case 1:
		result = bytes.Compare(fromInc, needle) < 1 || bytes.Compare(needle, toExc) < 0
	default:
		panic("Shouldn't happen")
	}
	return
}

// MergeItems will merge the given slices of Items into a slice with the newest version of each Item, with regard to their keys.
func MergeItems(arys []*[]Item, up bool) (result []Item) {
	result = *arys[0]
	var items []Item
	for j := 1; j < len(arys); j++ {
		items = *arys[j]
		for _, item := range items {
			i := sort.Search(len(result), func(i int) bool {
				cmp := bytes.Compare(item.Key, result[i].Key)
				if up {
					return cmp < 1
				}
				return cmp > -1
			})
			if i == len(result) {
				result = append(result, item)
			} else {
				if bytes.Compare(result[i].Key, item.Key) == 0 {
					if result[i].Timestamp < item.Timestamp {
						result[i] = item
					}
				} else {
					result = append(result[:i], append([]Item{item}, result[i:]...)...)
				}
			}
		}
	}
	return
}

// DHashDescription contains a description of a dhash node.
type DHashDescription struct {
	Addr         string
	Pos          []byte
	LastReroute  time.Time
	LastSync     time.Time
	LastMigrate  time.Time
	Timer        time.Time
	OwnedEntries int
	HeldEntries  int
	Load         float64
	Nodes        Remotes
}

// Describe will return a humanly readable string description of the dhash node.
func (self DHashDescription) Describe() string {
	return fmt.Sprintf("%+v", struct {
		Addr         string
		Pos          string
		LastReroute  time.Time
		LastSync     time.Time
		LastMigrate  time.Time
		Timer        time.Time
		OwnedEntries int
		HeldEntries  int
		Load         float64
		Nodes        string
	}{
		Addr:         self.Addr,
		Pos:          HexEncode(self.Pos),
		LastReroute:  self.LastReroute,
		LastSync:     self.LastSync,
		LastMigrate:  self.LastMigrate,
		Timer:        self.Timer,
		OwnedEntries: self.OwnedEntries,
		HeldEntries:  self.HeldEntries,
		Load:         self.Load,
		Nodes:        fmt.Sprintf("\n%v", self.Nodes.Describe()),
	})
}
