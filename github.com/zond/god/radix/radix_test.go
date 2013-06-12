package radix

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/zond/god/common"
	"github.com/zond/god/murmur"
	"math/big"
	"math/rand"
	"reflect"
	"runtime"
	"testing"
	"time"
)

var benchmarkTestTree *Tree
var benchmarkTestKeys [][]byte
var benchmarkTestValues [][]byte

func init() {
	rand.Seed(time.Now().UnixNano())
	benchmarkTestTree = NewTree()
	benchmarkTestTree.Log("benchmarklogs")
	benchmarkTestTree.logger.Clear()
}

func TestSyncVersions(t *testing.T) {
	tree1 := NewTree()
	tree3 := NewTree()
	n := 10
	var k []byte
	var v []byte
	for i := 0; i < n; i++ {
		k = []byte{byte(i)}
		v = []byte(fmt.Sprint(i))
		tree1.Put(k, v, 1)
		if i == 2 {
			tree3.Put(k, []byte("other timestamp"), 2)
		} else {
			tree3.Put(k, v, 1)
		}
	}
	tree2 := NewTree()
	tree2.Put([]byte{2}, []byte("other timestamp"), 2)
	s := NewSync(tree1, tree2)
	s.Run()
	if bytes.Compare(tree1.Hash(), tree2.Hash()) == 0 {
		t.Errorf("%v and %v have hashes\n%v\n%v\nand they should not be equal!", tree1.Describe(), tree2.Describe(), tree1.Hash(), tree2.Hash())
	}
	if tree1.deepEqual(tree2) {
		t.Errorf("%v and %v are equal", tree1, tree2)
	}
	if bytes.Compare(tree3.Hash(), tree2.Hash()) != 0 {
		t.Errorf("%v and %v have hashes\n%v\n%v\nand they should be equal!", tree3.Describe(), tree2.Describe(), tree3.Hash(), tree2.Hash())
	}
	if !tree3.deepEqual(tree2) {
		t.Errorf("%v and %v are unequal", tree3, tree2)
	}
	tree1.Put([]byte{2}, []byte("yet another timestamp"), 3)
	s.Run()
	if bytes.Compare(tree3.Hash(), tree2.Hash()) == 0 {
		t.Errorf("%v and %v have hashes\n%v\n%v\nand they should not be equal!", tree3.Describe(), tree2.Describe(), tree3.Hash(), tree2.Hash())
	}
	if tree3.deepEqual(tree2) {
		t.Errorf("%v and %v are equal", tree3, tree2)
	}
	if bytes.Compare(tree1.Hash(), tree2.Hash()) != 0 {
		t.Errorf("%v and %v have hashes\n%v\n%v\nand they should be equal!", tree1.Describe(), tree2.Describe(), tree1.Hash(), tree2.Hash())
	}
	if !tree1.deepEqual(tree2) {
		t.Errorf("%v and %v are unequal", tree1, tree2)
	}
}

func TestSyncLimits(t *testing.T) {
	tree1 := NewTree()
	tree3 := NewTree()
	n := 10
	from := 3
	to := 7
	fromKey := []byte{byte(from)}
	toKey := []byte{byte(to)}
	var k []byte
	var v []byte
	for i := 0; i < n; i++ {
		k = []byte{byte(i)}
		v = []byte(fmt.Sprint(i))
		tree1.Put(k, v, 1)
		if i >= from && i < to {
			tree3.Put(k, v, 1)
		}
	}
	tree2 := NewTree()
	s := NewSync(tree1, tree2).From(fromKey).To(toKey)
	s.Run()
	if bytes.Compare(tree3.Hash(), tree2.Hash()) != 0 {
		t.Errorf("%v and %v have hashes\n%v\n%v\nand they should be equal!", tree3.Describe(), tree2.Describe(), tree3.Hash(), tree2.Hash())
	}
	if !tree3.deepEqual(tree2) {
		t.Errorf("%v and %v are unequal", tree3, tree2)
	}
}

func TestripStitch(t *testing.T) {
	var b []byte
	for i := 0; i < 1000; i++ {
		b = make([]byte, rand.Int()%30)
		for j := 0; j < len(b); j++ {
			b[j] = byte(rand.Int())
		}
		if bytes.Compare(Stitch(Rip(b)), b) != 0 {
			t.Errorf("%v != %v", Stitch(Rip(b)), b)
		}
	}
}

func TestSyncSubTreeDestructive(t *testing.T) {
	tree1 := NewTree()
	tree3 := NewTree()
	n := 10
	var k, sk []byte
	var v []byte
	for i := 0; i < n; i++ {
		k = []byte(murmur.HashString(fmt.Sprint(i)))
		v = []byte(fmt.Sprint(i))
		if i%2 == 0 {
			tree1.Put(k, v, 1)
			tree3.Put(k, v, 1)
		} else {
			for j := 0; j < 10; j++ {
				sk = []byte(fmt.Sprint(j))
				tree1.SubPut(k, sk, v, 1)
				tree3.SubPut(k, sk, v, 1)
			}
		}
	}
	tree2 := NewTree()
	s := NewSync(tree3, tree2).Destroy()
	s.Run()
	if bytes.Compare(tree1.Hash(), tree2.Hash()) != 0 {
		t.Errorf("%v and %v have hashes\n%v\n%v\nand they should be equal!", tree1.Describe(), tree2.Describe(), tree1.Hash(), tree2.Hash())
	}
	if !tree1.deepEqual(tree2) {
		t.Errorf("\n%v and \n%v are unequal", tree1.Describe(), tree2.Describe())
	}
	if tree3.Size() != 0 {
		t.Errorf("%v should be empty", tree3.Describe())
	}
	if !tree3.deepEqual(NewTree()) {
		t.Errorf("%v and %v should be equal", tree3.Describe(), NewTree().Describe())
	}
}

func TestTreeSizeBetween(t *testing.T) {
	tree := NewTree()
	for i := 11; i < 20; i++ {
		tree.Put([]byte(fmt.Sprint(i)), []byte(fmt.Sprint(i)), 1)
	}
	for i := 10; i < 21; i++ {
		for j := i; j < 21; j++ {
			expected := common.Max(0, common.Min(j+1, 20)-common.Max(11, i))
			val := tree.SizeBetween([]byte(fmt.Sprint(i)), []byte(fmt.Sprint(j)), true, true)
			if val != expected {
				t.Errorf("%v.SizeBetween(%v, %v, true, true) should be %v but was %v", tree.Describe(), common.HexEncode([]byte(fmt.Sprint(i))), common.HexEncode([]byte(fmt.Sprint(j))), expected, val)
			}
			expected = common.Max(0, common.Min(j+1, 20)-common.Max(11, i+1))
			val = tree.SizeBetween([]byte(fmt.Sprint(i)), []byte(fmt.Sprint(j)), false, true)
			if val != expected {
				t.Errorf("%v.SizeBetween(%v, %v, false, true) should be %v but was %v", tree.Describe(), common.HexEncode([]byte(fmt.Sprint(i))), common.HexEncode([]byte(fmt.Sprint(j))), expected, val)
			}
			expected = common.Max(0, common.Min(j, 20)-common.Max(11, i))
			val = tree.SizeBetween([]byte(fmt.Sprint(i)), []byte(fmt.Sprint(j)), true, false)
			if val != expected {
				t.Errorf("%v.SizeBetween(%v, %v, true, false) should be %v but was %v", tree.Describe(), common.HexEncode([]byte(fmt.Sprint(i))), common.HexEncode([]byte(fmt.Sprint(j))), expected, val)
			}
			expected = common.Max(0, common.Min(j, 20)-common.Max(11, i+1))
			val = tree.SizeBetween([]byte(fmt.Sprint(i)), []byte(fmt.Sprint(j)), false, false)
			if val != expected {
				t.Errorf("%v.SizeBetween(%v, %v, false, false) should be %v but was %v", tree.Describe(), common.HexEncode([]byte(fmt.Sprint(i))), common.HexEncode([]byte(fmt.Sprint(j))), expected, val)
			}
		}
	}
	for i := 0; i < 10; i++ {
		tree.SubPut([]byte{50}, []byte{byte(i)}, []byte{byte(i)}, 1)
	}
	ary := []byte{50}
	if s := tree.SizeBetween(ary, ary, true, true); s != 10 {
		t.Errorf("wrong size calculated for %v\nbetween %v and %v\nwanted %v but got %v", tree.Describe(), common.HexEncode(ary), common.HexEncode(ary), 10, s)
	}
}

func TestSubTree(t *testing.T) {
	tree := NewTree()
	assertSize(t, tree, 0)
	tree.Put([]byte("a"), []byte("a"), 1)
	assertSize(t, tree, 1)
	tree.SubPut([]byte("b"), []byte("c"), []byte("d"), 1)
	assertSize(t, tree, 2)
	tree.SubPut([]byte("b"), []byte("d"), []byte("e"), 2)
	assertSize(t, tree, 3)
	if v, ver, e := tree.SubGet([]byte("b"), []byte("d")); bytes.Compare(v, []byte("e")) != 0 || ver != 2 || !e {
		t.Errorf("wrong result, wanted %v, %v, %v got %v, %v, %v", []byte("e"), 2, true, v, ver, e)
	}
}

func TestSyncSubTreeVersions(t *testing.T) {
	tree1 := NewTree()
	tree3 := NewTree()
	n := 10
	var k, sk []byte
	var v []byte
	for i := 0; i < n; i++ {
		k = []byte(murmur.HashString(fmt.Sprint(i)))
		v = []byte(fmt.Sprint(i))
		if i%2 == 0 {
			tree3.Put(k, v, 1)
			tree1.Put(k, v, 1)
		} else {
			for j := 0; j < 10; j++ {
				sk = []byte(fmt.Sprint(j))
				tree1.SubPut(k, sk, v, 1)
				if i == 1 && j == 3 {
					tree3.SubPut(k, sk, []byte("another value"), 2)
				} else {
					tree3.SubPut(k, sk, v, 1)
				}
			}
		}
	}
	tree2 := NewTree()
	tree2.SubPut([]byte(murmur.HashString(fmt.Sprint(1))), []byte(fmt.Sprint(3)), []byte("another value"), 2)
	s := NewSync(tree1, tree2)
	s.Run()
	if bytes.Compare(tree1.Hash(), tree2.Hash()) == 0 {
		t.Errorf("should not be equal")
	}
	if tree1.deepEqual(tree2) {
		t.Errorf("should not be equal")
	}
	if bytes.Compare(tree3.Hash(), tree2.Hash()) != 0 {
		t.Errorf("%v and %v have hashes\n%v\n%v\nand they should be equal!", tree3.Describe(), tree2.Describe(), tree3.Hash(), tree2.Hash())
	}
	if !tree3.deepEqual(tree2) {
		t.Errorf("\n%v and \n%v are unequal", tree3.Describe(), tree2.Describe())
	}
	tree1.SubPut([]byte(murmur.HashString(fmt.Sprint(1))), []byte(fmt.Sprint(3)), []byte("another value again"), 3)
	s.Run()
	if bytes.Compare(tree3.Hash(), tree2.Hash()) == 0 {
		t.Errorf("should not be equal")
	}
	if tree3.deepEqual(tree2) {
		t.Errorf("should not be equal")
	}
	if bytes.Compare(tree1.Hash(), tree2.Hash()) != 0 {
		t.Errorf("%v and %v have hashes\n%v\n%v\nand they should be equal!", tree1.Describe(), tree2.Describe(), tree1.Hash(), tree2.Hash())
	}
	if !tree1.deepEqual(tree2) {
		t.Errorf("\n%v and \n%v are unequal", tree1.Describe(), tree2.Describe())
	}
}
func TestSyncSubTree(t *testing.T) {
	tree1 := NewTree()
	n := 10
	var k, sk []byte
	var v []byte
	for i := 0; i < n; i++ {
		k = []byte(murmur.HashString(fmt.Sprint(i)))
		v = []byte(fmt.Sprint(i))
		if i%2 == 0 {
			tree1.Put(k, v, 1)
		} else {
			for j := 0; j < 10; j++ {
				sk = []byte(fmt.Sprint(j))
				tree1.SubPut(k, sk, v, 1)
			}
		}
	}
	tree2 := NewTree()
	s := NewSync(tree1, tree2)
	s.Run()
	if bytes.Compare(tree1.Hash(), tree2.Hash()) != 0 {
		t.Errorf("%v and %v have hashes\n%v\n%v\nand they should be equal!", tree1.Describe(), tree2.Describe(), tree1.Hash(), tree2.Hash())
	}
	if !tree1.deepEqual(tree2) {
		t.Errorf("\n%v and \n%v are unequal", tree1.Describe(), tree2.Describe())
	}
}

func TestSyncDestructive(t *testing.T) {
	tree1 := NewTree()
	tree3 := NewTree()
	n := 1000
	var k []byte
	var v []byte
	for i := 0; i < n; i++ {
		k = []byte(fmt.Sprint(i))
		v = []byte(fmt.Sprint(i))
		tree1.Put(k, v, 1)
		tree3.Put(k, v, 1)
	}
	tree2 := NewTree()
	s := NewSync(tree3, tree2).Destroy()
	s.Run()
	if bytes.Compare(tree1.Hash(), tree2.Hash()) != 0 {
		t.Errorf("%v and %v have hashes\n%v\n%v\nand they should be equal!", tree1.Describe(), tree2.Describe(), tree1.Hash(), tree2.Hash())
	}
	if !tree1.deepEqual(tree2) {
		t.Errorf("%v and %v are unequal", tree1, tree2)
	}
	if tree3.Size() != 0 {
		t.Errorf("%v should be size 0, is size %v", tree3, tree3.Size())
	}
	if !tree3.deepEqual(NewTree()) {
		t.Errorf("should be equal")
	}
}

func TestSyncDestructiveMatching(t *testing.T) {
	tree1 := NewTree()
	tree2 := NewTree()
	tree3 := NewTree()
	n := 1000
	var k []byte
	var v []byte
	for i := 0; i < n; i++ {
		k = murmur.HashString(fmt.Sprint(i))
		v = []byte(fmt.Sprint(i))
		tree1.Put(k, v, 1)
		tree2.Put(k, v, 1)
		tree3.Put(k, v, 1)
	}
	NewSync(tree1, tree2).Destroy().Run()
	if !tree2.deepEqual(tree3) {
		t.Errorf("should be equal")
	}
	if tree1.Size() != 0 {
		t.Errorf("should be empty!")
	}
	tree4 := NewTree()
	if !tree1.deepEqual(tree4) {
		t.Errorf("should be equal!")
	}
}

func TestSyncComplete(t *testing.T) {
	tree1 := NewTree()
	n := 1000
	var k []byte
	var v []byte
	for i := 0; i < n; i++ {
		k = []byte(fmt.Sprint(i))
		v = []byte(fmt.Sprint(i))
		tree1.Put(k, v, 1)
	}
	tree2 := NewTree()
	s := NewSync(tree1, tree2)
	s.Run()
	if bytes.Compare(tree1.Hash(), tree2.Hash()) != 0 {
		t.Errorf("%v and %v have hashes\n%v\n%v\nand they should be equal!", tree1.Describe(), tree2.Describe(), tree1.Hash(), tree2.Hash())
	}
	if !tree1.deepEqual(tree2) {
		t.Errorf("%v and %v are unequal", tree1, tree2)
	}
}

func TestSyncRandomLimits(t *testing.T) {
	tree1 := NewTree()
	n := 10
	var k []byte
	var v []byte
	for i := 0; i < n; i++ {
		k = murmur.HashString(fmt.Sprint(i))
		v = []byte(fmt.Sprint(i))
		tree1.Put(k, v, 1)
	}
	var keys [][]byte
	tree1.Each(func(key []byte, byteValue []byte, timestamp int64) bool {
		keys = append(keys, key)
		return true
	})
	var fromKey []byte
	var toKey []byte
	var tree2 *Tree
	var tree3 *Tree
	var s *Sync
	for fromIndex, _ := range keys {
		for toIndex, _ := range keys {
			if fromIndex != toIndex {
				fromKey = keys[fromIndex]
				toKey = keys[toIndex]
				if bytes.Compare(fromKey, toKey) < 0 {
					tree2 = NewTree()
					tree1.Each(func(key []byte, byteValue []byte, timestamp int64) bool {
						if common.BetweenIE(key, fromKey, toKey) {
							tree2.Put(key, byteValue, 1)
						}
						return true
					})
					tree3 = NewTree()
					s = NewSync(tree1, tree3).From(fromKey).To(toKey)
					s.Run()
					if !tree3.deepEqual(tree2) {
						t.Errorf("when syncing from %v to %v, %v and %v have hashes\n%v\n%v\nand they should be equal!", common.HexEncode(fromKey), common.HexEncode(toKey), tree3.Describe(), tree2.Describe(), tree3.Hash(), tree2.Hash())
					}
				}
			}
		}
	}
}

func TestSyncPartial(t *testing.T) {
	tree1 := NewTree()
	tree2 := NewTree()
	mod := 2
	n := 100
	var k []byte
	var v []byte
	for i := 0; i < n; i++ {
		k = []byte(fmt.Sprint(i))
		v = []byte(fmt.Sprint(i))
		tree1.Put(k, v, 1)
		if i%mod != 0 {
			tree2.Put(k, v, 1)
		}
	}
	s := NewSync(tree1, tree2)
	s.Run()
	if bytes.Compare(tree1.Hash(), tree2.Hash()) != 0 {
		t.Errorf("%v and %v have hashes\n%v\n%v\nand they should be equal!", tree1.Describe(), tree2.Describe(), tree1.Hash(), tree2.Hash())
	}
	if !tree1.deepEqual(tree2) {
		t.Errorf("%v and %v are unequal", tree1, tree2)
	}
}

func TestSyncSubConf(t *testing.T) {
	tree1 := NewTree()
	tree2 := NewTree()
	tree1.SubAddConfiguration([]byte("a"), 1, "blapp", "blepp")
	s := NewSync(tree1, tree2)
	s.Run()
	c1, _ := tree1.SubConfiguration([]byte("a"))
	c2, _ := tree2.SubConfiguration([]byte("a"))
	if !reflect.DeepEqual(c1, c2) {
		t.Errorf("%v and %v should be equal", c1, c2)
	}
}

func TestSyncConf(t *testing.T) {
	tree1 := NewTree()
	tree2 := NewTree()
	tree1.AddConfiguration(1, "blapp", "blepp")
	s := NewSync(tree1, tree2)
	s.Run()
	c1, _ := tree1.Configuration()
	c2, _ := tree2.Configuration()
	if !reflect.DeepEqual(c1, c2) {
		t.Errorf("%v and %v should be equal", c1, c2)
	}
}

func TestTreeHash(t *testing.T) {
	tree1 := NewTree()
	var keys [][]byte
	var vals [][]byte
	n := 10
	var k []byte
	var v []byte
	for i := 0; i < n; i++ {
		k = []byte(fmt.Sprint(rand.Int63()))
		v = []byte(fmt.Sprint(rand.Int63()))
		keys = append(keys, k)
		vals = append(vals, v)
		tree1.Put(k, v, 1)
	}
	keybackup := keys
	tree2 := NewTree()
	for i := 0; i < n; i++ {
		index := rand.Int() % len(keys)
		k = keys[index]
		v = vals[index]
		tree2.Put(k, v, 1)
		keys = append(keys[:index], keys[index+1:]...)
		vals = append(vals[:index], vals[index+1:]...)
	}
	if bytes.Compare(tree1.Hash(), tree2.Hash()) != 0 {
		t.Errorf("%v and %v have hashes\n%v\n%v\nand they should be equal!", tree1.Describe(), tree2.Describe(), tree1.Hash(), tree2.Hash())
	}
	if !reflect.DeepEqual(tree1.Finger(nil), tree2.Finger(nil)) {
		t.Errorf("%v and %v have prints\n%v\n%v\nand they should be equal!", tree1.Describe(), tree2.Describe(), tree1.Finger(nil), tree2.Finger(nil))
	}
	tree1.Each(func(key []byte, byteValue []byte, timestamp int64) bool {
		f1 := tree1.Finger(Rip(key))
		f2 := tree2.Finger(Rip(key))
		if f1 == nil || f2 == nil {
			t.Errorf("should not be nil!")
		}
		if !reflect.DeepEqual(f1, f2) {
			t.Errorf("should be equal!")
		}
		return true
	})
	var deletes []int
	for i := 0; i < n/10; i++ {
		index := rand.Int() % len(keybackup)
		deletes = append(deletes, index)
	}
	var successes []bool
	for i := 0; i < n/10; i++ {
		_, ok := tree1.Del(keybackup[deletes[i]])
		successes = append(successes, ok)
	}
	for i := 0; i < n/10; i++ {
		if _, ok := tree2.Del(keybackup[deletes[i]]); ok != successes[i] {
			t.Errorf("delete success should be %v", successes[i])
		}
	}
	if bytes.Compare(tree1.Hash(), tree2.Hash()) != 0 {
		t.Errorf("%v and %v have hashes\n%v\n%v\nand they should be equal!", tree1.Describe(), tree2.Describe(), tree1.Hash(), tree2.Hash())
	}
	if !reflect.DeepEqual(tree1.Finger(nil), tree2.Finger(nil)) {
		t.Errorf("%v and %v have prints\n%v\n%v\nand they should be equal!", tree1.Describe(), tree2.Describe(), tree1.Finger(nil), tree2.Finger(nil))
	}
	tree1.Each(func(key []byte, byteValue []byte, timestamp int64) bool {
		f1 := tree1.Finger(Rip(key))
		f2 := tree2.Finger(Rip(key))
		if f1 == nil || f2 == nil {
			t.Errorf("should not be nil!")
		}
		if !reflect.DeepEqual(f1, f2) {
			t.Errorf("should be equal!")
		}
		return true
	})
}

func createKVArraysDown(from, to int) (keys [][]byte, values [][]byte) {
	for i := to - 1; i >= from; i-- {
		if i >= from {
			keys = append(keys, []byte(fmt.Sprint(i)))
			values = append(values, []byte(fmt.Sprint(i)))
		}
	}
	return
}
func createKVArraysUp(from, to int) (keys [][]byte, values [][]byte) {
	for i := from; i < to; i++ {
		if i < to {
			keys = append(keys, []byte(fmt.Sprint(i)))
			values = append(values, []byte(fmt.Sprint(i)))
		}
	}
	return
}

func TestTreeReverseEach(t *testing.T) {
	tree := NewTree()
	for i := 100; i < 200; i++ {
		tree.Put([]byte(fmt.Sprint(i)), []byte(fmt.Sprint(i)), 1)
	}
	var foundKeys [][]byte
	var foundValues [][]byte
	tree.ReverseEach(func(key []byte, byteValue []byte, timestamp int64) bool {
		foundKeys = append(foundKeys, key)
		foundValues = append(foundValues, byteValue)
		return true
	})
	cmpKeys, cmpValues := createKVArraysDown(100, 200)
	if !reflect.DeepEqual(cmpKeys, foundKeys) {
		t.Errorf("%v should be %v", foundKeys, cmpKeys)
	}
	if !reflect.DeepEqual(cmpValues, foundValues) {
		t.Errorf("%v should be %v", foundValues, cmpValues)
	}
	foundKeys = nil
	foundValues = nil
	count := 10
	tree.ReverseEach(func(key []byte, byteValue []byte, timestamp int64) bool {
		foundKeys = append(foundKeys, key)
		foundValues = append(foundValues, byteValue)
		count--
		return count > 0
	})
	cmpKeys, cmpValues = createKVArraysDown(190, 200)
	if !reflect.DeepEqual(cmpKeys, foundKeys) {
		t.Errorf("%v should be %v", foundKeys, cmpKeys)
	}
	if !reflect.DeepEqual(cmpValues, foundValues) {
		t.Errorf("%v should be %v", foundValues, cmpValues)
	}
}

func TestTreeEach(t *testing.T) {
	tree := NewTree()
	for i := 100; i < 200; i++ {
		tree.Put([]byte(fmt.Sprint(i)), []byte(fmt.Sprint(i)), 1)
	}
	var foundKeys [][]byte
	var foundValues [][]byte
	tree.Each(func(key []byte, byteValue []byte, timestamp int64) bool {
		foundKeys = append(foundKeys, key)
		foundValues = append(foundValues, byteValue)
		return true
	})
	cmpKeys, cmpValues := createKVArraysUp(100, 200)
	if !reflect.DeepEqual(cmpKeys, foundKeys) {
		t.Errorf("%v should be %v", foundKeys, cmpKeys)
	}
	if !reflect.DeepEqual(cmpValues, foundValues) {
		t.Errorf("%v should be %v", foundValues, cmpValues)
	}
	foundKeys = nil
	foundValues = nil
	count := 10
	tree.Each(func(key []byte, byteValue []byte, timestamp int64) bool {
		foundKeys = append(foundKeys, key)
		foundValues = append(foundValues, byteValue)
		count--
		return count > 0
	})
	cmpKeys, cmpValues = createKVArraysUp(100, 110)
	if !reflect.DeepEqual(cmpKeys, foundKeys) {
		t.Errorf("%v should be %v", foundKeys, cmpKeys)
	}
	if !reflect.DeepEqual(cmpValues, foundValues) {
		t.Errorf("%v should be %v", foundValues, cmpValues)
	}
}

func TestTreeReverseEachBetween(t *testing.T) {
	tree := NewTree()
	for i := 11; i < 20; i++ {
		tree.Put([]byte(fmt.Sprint(i)), []byte(fmt.Sprint(i)), 1)
	}
	var foundKeys, cmpKeys [][]byte
	var foundValues, cmpValues [][]byte
	for i := 10; i < 21; i++ {
		for j := i; j < 21; j++ {
			foundKeys, cmpKeys, foundValues, cmpValues = nil, nil, nil, nil
			tree.ReverseEachBetween([]byte(fmt.Sprint(i)), []byte(fmt.Sprint(j)), true, true, func(key []byte, byteValue []byte, timestamp int64) bool {
				foundKeys = append(foundKeys, key)
				foundValues = append(foundValues, byteValue)
				return true
			})
			cmpKeys, cmpValues = createKVArraysDown(common.Max(11, i), common.Min(j+1, 20))
			if !reflect.DeepEqual(cmpKeys, foundKeys) || !reflect.DeepEqual(cmpValues, foundValues) {
				t.Errorf("%v.ReverseEachBetween(%v, %v, true, true) => %v should be %v", tree, i, j, foundValues, cmpValues)
			}

			foundKeys, cmpKeys, foundValues, cmpValues = nil, nil, nil, nil
			tree.ReverseEachBetween([]byte(fmt.Sprint(i)), []byte(fmt.Sprint(j)), false, true, func(key []byte, byteValue []byte, timestamp int64) bool {
				foundKeys = append(foundKeys, key)
				foundValues = append(foundValues, byteValue)
				return true
			})
			cmpKeys, cmpValues = createKVArraysDown(common.Max(11, i+1), common.Min(j+1, 20))
			if !reflect.DeepEqual(cmpKeys, foundKeys) || !reflect.DeepEqual(cmpValues, foundValues) {
				t.Errorf("%v.ReverseEachBetween(%v, %v, false, true) => %v should be %v", tree, i, j, foundValues, cmpValues)
			}

			foundKeys, cmpKeys, foundValues, cmpValues = nil, nil, nil, nil
			tree.ReverseEachBetween([]byte(fmt.Sprint(i)), []byte(fmt.Sprint(j)), true, false, func(key []byte, byteValue []byte, timestamp int64) bool {
				foundKeys = append(foundKeys, key)
				foundValues = append(foundValues, byteValue)
				return true
			})
			cmpKeys, cmpValues = createKVArraysDown(common.Max(11, i), common.Min(j, 20))
			if !reflect.DeepEqual(cmpKeys, foundKeys) || !reflect.DeepEqual(cmpValues, foundValues) {
				t.Errorf("%v.ReverseEachBetween(%v, %v, true, false) => %v should be %v", tree, i, j, foundValues, cmpValues)
			}

			foundKeys, cmpKeys, foundValues, cmpValues = nil, nil, nil, nil
			tree.ReverseEachBetween([]byte(fmt.Sprint(i)), []byte(fmt.Sprint(j)), false, false, func(key []byte, byteValue []byte, timestamp int64) bool {
				foundKeys = append(foundKeys, key)
				foundValues = append(foundValues, byteValue)
				return true
			})
			cmpKeys, cmpValues = createKVArraysDown(common.Max(11, i+1), common.Min(j, 20))
			if !reflect.DeepEqual(cmpKeys, foundKeys) || !reflect.DeepEqual(cmpValues, foundValues) {
				t.Errorf("%v.ReverseEachBetween(%v, %v, false, false) => %v should be %v", tree, i, j, foundValues, cmpValues)
			}
		}
	}
}

func TestTreeEachBetween(t *testing.T) {
	tree := NewTree()
	for i := 11; i < 20; i++ {
		tree.Put([]byte(fmt.Sprint(i)), []byte(fmt.Sprint(i)), 1)
	}
	var foundKeys, cmpKeys [][]byte
	var foundValues, cmpValues [][]byte
	for i := 10; i < 21; i++ {
		for j := i; j < 21; j++ {
			foundKeys, cmpKeys, foundValues, cmpValues = nil, nil, nil, nil
			tree.EachBetween([]byte(fmt.Sprint(i)), []byte(fmt.Sprint(j)), true, true, func(key []byte, byteValue []byte, timestamp int64) bool {
				foundKeys = append(foundKeys, key)
				foundValues = append(foundValues, byteValue)
				return true
			})
			cmpKeys, cmpValues = createKVArraysUp(common.Max(11, i), common.Min(j+1, 20))
			if !reflect.DeepEqual(cmpKeys, foundKeys) || !reflect.DeepEqual(cmpValues, foundValues) {
				t.Errorf("%v.EachBetween(%v, %v, true, true) => %v should be %v", tree, i, j, foundValues, cmpValues)
			}

			foundKeys, cmpKeys, foundValues, cmpValues = nil, nil, nil, nil
			tree.EachBetween([]byte(fmt.Sprint(i)), []byte(fmt.Sprint(j)), false, true, func(key []byte, byteValue []byte, timestamp int64) bool {
				foundKeys = append(foundKeys, key)
				foundValues = append(foundValues, byteValue)
				return true
			})
			cmpKeys, cmpValues = createKVArraysUp(common.Max(11, i+1), common.Min(j+1, 20))
			if !reflect.DeepEqual(cmpKeys, foundKeys) || !reflect.DeepEqual(cmpValues, foundValues) {
				t.Errorf("%v.EachBetween(%v, %v, false, true) => %v should be %v", tree, i, j, foundValues, cmpValues)
			}

			foundKeys, cmpKeys, foundValues, cmpValues = nil, nil, nil, nil
			tree.EachBetween([]byte(fmt.Sprint(i)), []byte(fmt.Sprint(j)), true, false, func(key []byte, byteValue []byte, timestamp int64) bool {
				foundKeys = append(foundKeys, key)
				foundValues = append(foundValues, byteValue)
				return true
			})
			cmpKeys, cmpValues = createKVArraysUp(common.Max(11, i), common.Min(j, 20))
			if !reflect.DeepEqual(cmpKeys, foundKeys) || !reflect.DeepEqual(cmpValues, foundValues) {
				t.Errorf("%v.EachBetween(%v, %v, true, false) => %v should be %v", tree, i, j, foundValues, cmpValues)
			}

			foundKeys, cmpKeys, foundValues, cmpValues = nil, nil, nil, nil
			tree.EachBetween([]byte(fmt.Sprint(i)), []byte(fmt.Sprint(j)), false, false, func(key []byte, byteValue []byte, timestamp int64) bool {
				foundKeys = append(foundKeys, key)
				foundValues = append(foundValues, byteValue)
				return true
			})
			cmpKeys, cmpValues = createKVArraysUp(common.Max(11, i+1), common.Min(j, 20))
			if !reflect.DeepEqual(cmpKeys, foundKeys) || !reflect.DeepEqual(cmpValues, foundValues) {
				t.Errorf("%v.EachBetween(%v, %v, false, false) => %v should be %v", tree, i, j, foundValues, cmpValues)
			}
		}
	}
}

func TestTreeReverseIndexOf(t *testing.T) {
	tree := NewTree()
	for i := 100; i < 200; i += 2 {
		tree.Put([]byte(fmt.Sprint(i)), []byte(fmt.Sprint(i)), 1)
	}
	for i := 100; i < 200; i++ {
		shouldExist := i%2 == 0
		wantedIndex := 49 - (i-100)/2
		if ind, e := tree.ReverseIndexOf([]byte(fmt.Sprint(i))); ind != wantedIndex || e != shouldExist {
			t.Errorf("%v.ReverseIndexOf(%v) => %v, %v should be %v, %v", tree.Describe(), i, ind, e, wantedIndex, shouldExist)
		}
	}
	if ind, e := tree.ReverseIndexOf([]byte("1991")); ind != 0 || e {
		t.Errorf("%v.IndexOf(%v) => %v, %v should be %v, %v", tree.Describe(), "1991", ind, e, 0, false)
	}
	if ind, e := tree.ReverseIndexOf([]byte("099")); ind != 50 || e {
		t.Errorf("%v.IndexOf(%v) => %v, %v should be %v, %v", tree.Describe(), "099", ind, e, 50, false)
	}
}

func TestTreeIndexOf(t *testing.T) {
	tree := NewTree()
	for i := 100; i < 200; i += 2 {
		tree.Put([]byte(fmt.Sprint(i)), []byte(fmt.Sprint(i)), 1)
	}
	for i := 100; i < 200; i++ {
		shouldExist := i%2 == 0
		wantedIndex := (i - 99) / 2
		if ind, e := tree.IndexOf([]byte(fmt.Sprint(i))); ind != wantedIndex || e != shouldExist {
			t.Errorf("%v.IndexOf(%v) => %v, %v should be %v, %v", tree.Describe(), common.HexEncode([]byte(fmt.Sprint(i))), ind, e, wantedIndex, shouldExist)
		}
	}
	if ind, e := tree.IndexOf([]byte("1991")); ind != 50 || e {
		t.Errorf("%v.IndexOf(%v) => %v, %v should be %v, %v", tree.Describe(), "1991", ind, e, 50, false)
	}
	if ind, e := tree.IndexOf([]byte("099")); ind != 0 || e {
		t.Errorf("%v.IndexOf(%v) => %v, %v should be %v, %v", tree.Describe(), "099", ind, e, 0, false)
	}
}

func TestTreeReverseEachBetweenIndex(t *testing.T) {
	tree := NewTree()
	for i := 11; i < 20; i++ {
		tree.Put([]byte(fmt.Sprint(i)), []byte(fmt.Sprint(i)), 1)
	}
	var foundKeys, cmpKeys [][]byte
	var foundValues, cmpValues [][]byte
	for i := -1; i < 10; i++ {
		for j := i; j < 10; j++ {
			foundKeys, cmpKeys, foundValues, cmpValues = nil, nil, nil, nil
			tree.ReverseEachBetweenIndex(&i, &j, func(key []byte, byteValue []byte, timestamp int64, index int) bool {
				foundKeys = append(foundKeys, key)
				foundValues = append(foundValues, byteValue)
				return true
			})

			cmpKeys, cmpValues = createKVArraysDown(common.Max(11, 19-j), common.Min(20, 20-i))
			if !reflect.DeepEqual(cmpKeys, foundKeys) || !reflect.DeepEqual(cmpValues, foundValues) {
				t.Errorf("%v.EachBetweenIndex(%v, %v) => %v should be %v", tree, i, j, foundValues, cmpValues)
			}
		}
	}
}

func TestTreeEachBetweenIndex(t *testing.T) {
	tree := NewTree()
	for i := 11; i < 20; i++ {
		tree.Put([]byte(fmt.Sprint(i)), []byte(fmt.Sprint(i)), 1)
	}
	var foundKeys, cmpKeys [][]byte
	var foundValues, cmpValues [][]byte
	for i := -1; i < 10; i++ {
		for j := i; j < 10; j++ {
			foundKeys, cmpKeys, foundValues, cmpValues = nil, nil, nil, nil
			tree.EachBetweenIndex(&i, &j, func(key []byte, byteValue []byte, timestamp int64, index int) bool {
				foundKeys = append(foundKeys, key)
				foundValues = append(foundValues, byteValue)
				return true
			})

			cmpKeys, cmpValues = createKVArraysUp(common.Max(11, i+11), common.Min(j+12, 20))
			if !reflect.DeepEqual(cmpKeys, foundKeys) || !reflect.DeepEqual(cmpValues, foundValues) {
				t.Errorf("%v.EachBetweenIndex(%v, %v) => %v should be %v", tree, i, j, foundValues, cmpValues)
			}
		}
	}
}

func TestTreeNilKey(t *testing.T) {
	tree := NewTree()
	h := tree.Hash()
	if value, _, existed := tree.Get(nil); value != nil || existed {
		t.Errorf("should not exist")
	}
	if value, existed := tree.Put(nil, nil, 1); value != nil || existed {
		t.Errorf("should not exist")
	}
	if value, _, existed := tree.Get(nil); value != nil || !existed {
		t.Errorf("should exist")
	}
	if value, existed := tree.Del(nil); value != nil || !existed {
		t.Errorf("should exist")
	}
	if value, _, existed := tree.Get(nil); value != nil || existed {
		t.Errorf("nil should not exist in %v", tree.Describe())
	}
	if bytes.Compare(h, tree.Hash()) != 0 {
		t.Errorf("should be equal")
	}
}

func TestTreeFirstLast(t *testing.T) {
	tree := NewTree()
	for i := 10; i < 20; i++ {
		tree.Put([]byte(fmt.Sprint(i)), []byte(fmt.Sprint(i)), 1)
	}
	if key, value, timestamp, existed := tree.First(); !existed || bytes.Compare(key, []byte(fmt.Sprint(10))) != 0 || timestamp != 1 || bytes.Compare(value, []byte(fmt.Sprint(10))) != 0 {
		t.Errorf("%v.First() should be %v, %v, %v, %v but was %v, %v, %v, %v", tree.Describe(), []byte(fmt.Sprint(10)), []byte(fmt.Sprint(10)), 1, true, key, value, timestamp, existed)
	}
	if key, value, timestamp, existed := tree.Last(); !existed || bytes.Compare(key, []byte(fmt.Sprint(19))) != 0 || timestamp != 1 || bytes.Compare(value, []byte(fmt.Sprint(19))) != 0 {
		t.Errorf("%v.Last() should be %v, %v, %v, %v but was %v, %v, %v, %v", tree.Describe(), string([]byte(fmt.Sprint(19))), []byte(fmt.Sprint(19)), 1, true, string(key), string(value), timestamp, existed)
	}
}

func TestTreeIndex(t *testing.T) {
	tree := NewTree()
	for i := 100; i < 200; i++ {
		tree.Put([]byte(fmt.Sprint(i)), []byte(fmt.Sprint(i)), 1)
	}
	for i := 0; i < 100; i++ {
		if key, value, timestamp, existed := tree.Index(i); !existed || bytes.Compare(key, []byte(fmt.Sprint(i+100))) != 0 || timestamp != 1 || bytes.Compare(value, []byte(fmt.Sprint(i+100))) != 0 {
			t.Errorf("%v.Index(%v) should be %v, %v, %v, %v but was %v, %v, %v, %v", tree.Describe(), i, []byte(fmt.Sprint(i+100)), []byte(fmt.Sprint(i+100)), 1, true, key, value, timestamp, existed)
		}
	}
}

func TestTreePrev(t *testing.T) {
	tree := NewTree()
	for i := 100; i < 200; i++ {
		tree.Put([]byte(fmt.Sprint(i)), []byte(fmt.Sprint(i)), 1)
	}
	for i := 101; i < 200; i++ {
		if key, _, _, existed := tree.Prev([]byte(fmt.Sprint(i))); string(key) != fmt.Sprint(i-1) || !existed {
			t.Errorf("%v, %v should be %v, %v", string(key), existed, fmt.Sprint(i-1), true)
		}
	}
	if key, _, _, existed := tree.Prev([]byte("100")); existed {
		t.Errorf("%v, %v should not exist!", key, existed)
	}
}

func TestTreeNext(t *testing.T) {
	tree := NewTree()
	for i := 100; i < 200; i++ {
		tree.Put([]byte(fmt.Sprint(i)), []byte(fmt.Sprint(i)), 1)
	}
	for i := 100; i < 199; i++ {
		if key, _, _, existed := tree.Next([]byte(fmt.Sprint(i))); string(key) != fmt.Sprint(i+1) || !existed {
			t.Errorf("%v, %v should be %v, %v", string(key), existed, fmt.Sprint(i+1), true)
		}
	}
	if key, _, _, existed := tree.Next([]byte("199")); existed {
		t.Errorf("%v, %v should not exist!", key, existed)
	}
}

func TestTreeBasicOps(t *testing.T) {
	tree := NewTree()
	assertSize(t, tree, 0)
	assertNewPut(t, tree, "apple", "stonefruit")
	assertSize(t, tree, 1)
	assertOldPut(t, tree, "apple", "fruit", "stonefruit")
	assertSize(t, tree, 1)
	assertNewPut(t, tree, "crab", "critter")
	assertSize(t, tree, 2)
	assertOldPut(t, tree, "crab", "animal", "critter")
	assertSize(t, tree, 2)
	assertNewPut(t, tree, "crabapple", "poop")
	assertSize(t, tree, 3)
	assertOldPut(t, tree, "crabapple", "fruit", "poop")
	assertSize(t, tree, 3)
	assertNewPut(t, tree, "banana", "yellow")
	assertSize(t, tree, 4)
	assertOldPut(t, tree, "banana", "fruit", "yellow")
	assertSize(t, tree, 4)
	assertNewPut(t, tree, "guava", "fart")
	assertSize(t, tree, 5)
	assertOldPut(t, tree, "guava", "fruit", "fart")
	assertSize(t, tree, 5)
	assertNewPut(t, tree, "guanabana", "man")
	assertSize(t, tree, 6)
	assertOldPut(t, tree, "guanabana", "city", "man")
	assertSize(t, tree, 6)
	m := make(map[string][]byte)
	tree.Each(func(key []byte, byteValue []byte, timestamp int64) bool {
		m[hex.EncodeToString(key)] = byteValue
		return true
	})
	comp := map[string][]byte{
		hex.EncodeToString([]byte("apple")):     []byte("fruit"),
		hex.EncodeToString([]byte("crab")):      []byte("animal"),
		hex.EncodeToString([]byte("crabapple")): []byte("fruit"),
		hex.EncodeToString([]byte("banana")):    []byte("fruit"),
		hex.EncodeToString([]byte("guava")):     []byte("fruit"),
		hex.EncodeToString([]byte("guanabana")): []byte("city"),
	}
	if !reflect.DeepEqual(m, comp) {
		t.Errorf("%+v and %+v should be equal!", m, comp)
	}
	if !reflect.DeepEqual(tree.ToMap(), comp) {
		t.Errorf("%v and %v should be equal!", tree.ToMap(), comp)
	}
	if old, existed := tree.Put(nil, []byte("nil"), 1); old != nil || existed {
		t.Error("should not exist yet")
	}
	if old, existed := tree.Put([]byte("nil"), nil, 1); old != nil || existed {
		t.Error("should not exist yet")
	}
	if value, _, existed := tree.Get(nil); !existed || bytes.Compare(value, []byte("nil")) != 0 {
		t.Errorf("%v should contain %v => %v, got %v, %v", tree.Describe(), nil, "nil", value, existed)
	}
	if value, _, existed := tree.Get([]byte("nil")); !existed || value != nil {
		t.Errorf("%v should contain %v => %v, got %v, %v", tree, "nil", nil, value, existed)
	}
	assertDelFailure(t, tree, "gua")
	assertSize(t, tree, 8)
	assertDelSuccess(t, tree, "apple", "fruit")
	assertSize(t, tree, 7)
	assertDelFailure(t, tree, "apple")
	assertSize(t, tree, 7)
	assertDelSuccess(t, tree, "crab", "animal")
	assertSize(t, tree, 6)
	assertDelFailure(t, tree, "crab")
	assertSize(t, tree, 6)
	assertDelSuccess(t, tree, "crabapple", "fruit")
	assertSize(t, tree, 5)
	assertDelFailure(t, tree, "crabapple")
	assertSize(t, tree, 5)
	assertDelSuccess(t, tree, "banana", "fruit")
	assertSize(t, tree, 4)
	assertDelFailure(t, tree, "banana")
	assertSize(t, tree, 4)
	assertDelSuccess(t, tree, "guava", "fruit")
	assertSize(t, tree, 3)
	assertDelFailure(t, tree, "guava")
	assertSize(t, tree, 3)
	assertDelSuccess(t, tree, "guanabana", "city")
	assertSize(t, tree, 2)
	assertDelFailure(t, tree, "guanabana")
	assertSize(t, tree, 2)
}

func benchTreeSync(b *testing.B, size, delta int) {
	b.StopTimer()
	tree1 := NewTree()
	tree2 := NewTree()
	var k []byte
	var v []byte
	for i := 0; i < size; i++ {
		k = murmur.HashString(fmt.Sprint(i))
		v = []byte(fmt.Sprint(i))
		tree1.Put(k, v, 1)
		tree2.Put(k, v, 1)
	}
	var s *Sync
	for i := 0; i < b.N/delta; i++ {
		for j := 0; j < delta; j++ {
			tree2.Del(murmur.HashString(fmt.Sprint(j)))
		}
		b.StartTimer()
		s = NewSync(tree1, tree2)
		s.Run()
		b.StopTimer()
		if bytes.Compare(tree1.Hash(), tree2.Hash()) != 0 {
			b.Fatalf("%v != %v", tree1.Hash(), tree2.Hash())
		}
	}
}

func BenchmarkTreeSync10000_1(b *testing.B) {
	benchTreeSync(b, 10000, 1)
}
func BenchmarkTreeSync10000_10(b *testing.B) {
	benchTreeSync(b, 10000, 10)
}
func BenchmarkTreeSync10000_100(b *testing.B) {
	benchTreeSync(b, 10000, 100)
}
func BenchmarkTreeSync10000_1000(b *testing.B) {
	benchTreeSync(b, 10000, 1000)
}

func BenchmarkTreeSync100000_1(b *testing.B) {
	benchTreeSync(b, 100000, 1)
}
func BenchmarkTreeSync100000_10(b *testing.B) {
	benchTreeSync(b, 100000, 10)
}
func BenchmarkTreeSync100000_100(b *testing.B) {
	benchTreeSync(b, 100000, 100)
}
func BenchmarkTreeSync100000_1000(b *testing.B) {
	benchTreeSync(b, 100000, 1000)
}

func fillBenchTree(b *testing.B, n int) {
	b.StopTimer()
	for len(benchmarkTestKeys) < n {
		benchmarkTestKeys = append(benchmarkTestKeys, murmur.HashString(fmt.Sprint(len(benchmarkTestKeys))))
		benchmarkTestValues = append(benchmarkTestValues, []byte(fmt.Sprint(len(benchmarkTestValues))))
	}
	for benchmarkTestTree.Size() < n {
		benchmarkTestTree.Put(benchmarkTestKeys[benchmarkTestTree.Size()], benchmarkTestValues[benchmarkTestTree.Size()], 1)
	}
	b.StartTimer()
}

func benchTree(b *testing.B, n int, put, get bool) {
	fillBenchTree(b, n)
	b.StopTimer()
	oldprocs := runtime.GOMAXPROCS(runtime.NumCPU())
	defer runtime.GOMAXPROCS(oldprocs)
	var keys [][]byte
	var vals [][]byte
	for i := 0; i < b.N; i++ {
		keys = append(keys, murmur.HashString(fmt.Sprint(rand.Int63())))
		vals = append(vals, []byte(fmt.Sprint(rand.Int63())))
	}
	var k []byte
	var v []byte
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		k = benchmarkTestKeys[i%len(benchmarkTestKeys)]
		v = benchmarkTestValues[i%len(benchmarkTestValues)]
		if put {
			benchmarkTestTree.Put(k, v, 1)
		}
		if get {
			j, _, existed := benchmarkTestTree.Get(k)
			if bytes.Compare(j, v) != 0 {
				b.Fatalf("%v should contain %v, but got %v, %v", benchmarkTestTree.Describe(), v, j, existed)
			}
		}
	}
}

func BenchmarkTreePut10(b *testing.B) {
	benchTree(b, 10, true, false)
}

func BenchmarkTreeGet10(b *testing.B) {
	benchTree(b, 10, false, true)
}

func BenchmarkTreePut100(b *testing.B) {
	benchTree(b, 100, true, false)
}

func BenchmarkTreeGet100(b *testing.B) {
	benchTree(b, 100, false, true)
}

func BenchmarkTreePut1000(b *testing.B) {
	benchTree(b, 1000, true, false)
}

func BenchmarkTreeGet1000(b *testing.B) {
	benchTree(b, 1000, false, true)
}

func BenchmarkTreePut10000(b *testing.B) {
	benchTree(b, 10000, true, false)
}

func BenchmarkTreeGet10000(b *testing.B) {
	benchTree(b, 10000, false, true)
}

func BenchmarkTreePut100000(b *testing.B) {
	benchTree(b, 100000, true, false)
}

func BenchmarkTreeGet100000(b *testing.B) {
	benchTree(b, 100000, false, true)
}

func BenchmarkTreePut1000000(b *testing.B) {
	benchTree(b, 1000000, true, false)
}

func BenchmarkTreeGet1000000(b *testing.B) {
	benchTree(b, 1000000, false, true)
}

func BenchmarkTreeRealSizeBetween0_8_100000(b *testing.B) {
	fillBenchTree(b, 100000)
	max := new(big.Int).Div(new(big.Int).Lsh(big.NewInt(1), murmur.Size*8), big.NewInt(2))
	for i := 0; i < b.N; i++ {
		benchmarkTestTree.RealSizeBetween(nil, max.Bytes(), true, false)
	}
}

func BenchmarkTreeRealSizeBetween8_0_100000(b *testing.B) {
	fillBenchTree(b, 100000)
	max := new(big.Int).Div(new(big.Int).Lsh(big.NewInt(1), murmur.Size*8), big.NewInt(2))
	for i := 0; i < b.N; i++ {
		benchmarkTestTree.RealSizeBetween(max.Bytes(), nil, true, false)
	}
}

func BenchmarkTreeNextMarkerIndex1000000(b *testing.B) {
	fillBenchTree(b, 100000)
	s := benchmarkTestTree.RealSize()
	for i := 0; i < b.N; i++ {
		benchmarkTestTree.NextMarkerIndex(int(rand.Int31n(int32(s))))
	}
}

func BenchmarkTreeMirrorPut10(b *testing.B) {
	benchmarkTestTree = NewTree()
	benchmarkTestTree.AddConfiguration(1, mirrored, yes)
	benchTree(b, 10, true, false)
}

func BenchmarkTreeMirrorPut100(b *testing.B) {
	benchTree(b, 100, true, false)
}

func BenchmarkTreeMirrorPut1000(b *testing.B) {
	benchTree(b, 1000, true, false)
}

func BenchmarkTreeMirrorPut10000(b *testing.B) {
	benchTree(b, 10000, true, false)
}

func BenchmarkTreeMirrorPut100000(b *testing.B) {
	benchTree(b, 100000, true, false)
}

func BenchmarkTreeMirrorPut1000000(b *testing.B) {
	benchTree(b, 1000000, true, false)
}
