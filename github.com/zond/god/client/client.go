package client

import (
	"bytes"
	"fmt"
	"github.com/zond/god/common"
	"github.com/zond/setop"
	"net/rpc"
	"sync"
	"sync/atomic"
	"time"
)

const (
	created = iota
	started
	stopped
)

func findKeys(op *setop.SetOp) (result map[string]bool) {
	result = make(map[string]bool)
	for _, source := range op.Sources {
		if source.Key != nil {
			result[string(source.Key)] = true
		} else {
			for key, _ := range findKeys(source.SetOp) {
				result[key] = true
			}
		}
	}
	return
}

// Conn is the client connection.
//
// A god database has two data types: byte values and sub trees. All values are sorted in ascending order.
//
// Byte values are simply byte slices indexed by a key byte slice.
//
// Sub trees are trees containing byte values.
//
// Sub trees can be 'mirrored', which means that they contain a tree mirroring its values as keys and its keys as values.
// To mirror a sub tree, call SubAddConfiguration for the sub tree and set 'mirrored' to 'yes'.
//
// Naming conventions:
//
// If there are two methods with similar names except that one has a capital S prefixed, that means that the method with the capital S will not return until all nodes responsible for the written data has received the data, while the one without the capital S will return as soon as the owner of the data has received it.
//
// Methods prefixed Sub will work on sub trees.
//
// Methods prefixed Reverse will work in reverse order. Return slices in reverse order and indices from the end instead of the start etc.
//
// Methods prefixed Mirror will work on the mirror trees of the sub trees in question.
//
// Parameters named mininc and maxinc paired with parameters min and max of []byte type defined whether the min and max parameters are inclusive as opposed to exclusive.
//
// To install: go get github.com/zond/god/client
//
// Usage: https://github.com/zond/god/blob/master/client/client_test.go
type Conn struct {
	ring  *common.Ring
	state int32
}

// NewConnRing creates a new Conn from a given set of known nodes. For internal usage.
func NewConnRing(ring *common.Ring) *Conn {
	return &Conn{ring: ring}
}

// NewConn creates a new Conn to a cluster defined by the address of one of its members.
func NewConn(addr string) (result *Conn, err error) {
	result = &Conn{ring: common.NewRing()}
	var newNodes common.Remotes
	err = common.Switch.Call(addr, "Discord.Nodes", 0, &newNodes)
	result.ring.SetNodes(newNodes)
	return
}

// MustConn creates a working Conn or panics.
func MustConn(addr string) (result *Conn) {
	var err error
	if result, err = NewConn(addr); err != nil {
		panic(err)
	}
	return
}
func (self *Conn) hasState(s int32) bool {
	return atomic.LoadInt32(&self.state) == s
}
func (self *Conn) changeState(old, neu int32) bool {
	return atomic.CompareAndSwapInt32(&self.state, old, neu)
}
func (self *Conn) removeNode(node common.Remote) {
	self.ring.Remove(node)
	self.Reconnect()
}

// Nodes returns the set of known nodes for this Conn.
func (self *Conn) Nodes() common.Remotes {
	return self.ring.Nodes()
}
func (self *Conn) update() {
	myRingHash := self.ring.Hash()
	var otherRingHash []byte
	node := self.ring.Random()
	if err := node.Call("DHash.RingHash", 0, &otherRingHash); err != nil {
		self.removeNode(node)
		return
	}
	if bytes.Compare(myRingHash, otherRingHash) != 0 {
		var newNodes common.Remotes
		if err := node.Call("Discord.Nodes", 0, &newNodes); err != nil {
			self.removeNode(node)
			return
		}
		self.ring.SetNodes(newNodes)
	}
}
func (self *Conn) updateRegularly() {
	for self.hasState(started) {
		self.update()
		time.Sleep(common.PingInterval)
	}
}

// Start will begin to regularly update the set of known nodes for this Conn.
func (self *Conn) Start() {
	if self.changeState(created, started) {
		go self.updateRegularly()
	}
}

// Reconnect will try to refetch the set of known nodes from a randomly chosen currently known node.
func (self *Conn) Reconnect() {
	node := self.ring.Random()
	var err error
	for {
		var newNodes common.Remotes
		if err = node.Call("Discord.Nodes", 0, &newNodes); err == nil {
			self.ring.SetNodes(newNodes)
			return
		}
		self.ring.Remove(node)
		if self.ring.Size() == 0 {
			panic(fmt.Errorf("%v doesn't know of any live nodes!", self))
		}
		node = self.ring.Random()
	}
}

func (self *Conn) subClear(key []byte, sync bool) {
	data := common.Item{
		Key:  key,
		Sync: sync,
	}
	_, _, successor := self.ring.Remotes(key)
	var x int
	if err := successor.Call("DHash.SubClear", data, &x); err != nil {
		self.removeNode(*successor)
		self.subClear(key, sync)
	}
}
func (self *Conn) subDel(key, subKey []byte, sync bool) {
	data := common.Item{
		Key:    key,
		SubKey: subKey,
		Sync:   sync,
	}
	_, _, successor := self.ring.Remotes(key)
	var x int
	if err := successor.Call("DHash.SubDel", data, &x); err != nil {
		self.removeNode(*successor)
		self.subDel(key, subKey, sync)
	}
}
func (self *Conn) subPutVia(succ *common.Remote, key, subKey, value []byte, sync bool) {
	data := common.Item{
		Key:    key,
		SubKey: subKey,
		Value:  value,
		Sync:   sync,
	}
	var x int
	if err := succ.Call("DHash.SubPut", data, &x); err != nil {
		self.removeNode(*succ)
		_, _, newSuccessor := self.ring.Remotes(key)
		*succ = *newSuccessor
		self.subPutVia(succ, key, subKey, value, sync)
	}
}
func (self *Conn) subPut(key, subKey, value []byte, sync bool) {
	_, _, successor := self.ring.Remotes(key)
	self.subPutVia(successor, key, subKey, value, sync)
}
func (self *Conn) del(key []byte, sync bool) {
	data := common.Item{
		Key:  key,
		Sync: sync,
	}
	_, _, successor := self.ring.Remotes(key)
	var x int
	if err := successor.Call("DHash.Del", data, &x); err != nil {
		self.removeNode(*successor)
		self.del(key, sync)
	}
}
func (self *Conn) putVia(succ *common.Remote, key, value []byte, sync bool) {
	data := common.Item{
		Key:   key,
		Value: value,
		Sync:  sync,
	}
	var x int
	if err := succ.Call("DHash.Put", data, &x); err != nil {
		self.removeNode(*succ)
		_, _, newSuccessor := self.ring.Remotes(key)
		*succ = *newSuccessor
		self.putVia(succ, key, value, sync)
	}
}
func (self *Conn) put(key, value []byte, sync bool) {
	_, _, successor := self.ring.Remotes(key)
	self.putVia(successor, key, value, sync)
}
func (self *Conn) mergeRecent(operation string, r common.Range, up bool) (result []common.Item) {
	currentRedundancy := self.ring.Redundancy()
	futures := make([]*rpc.Call, currentRedundancy)
	results := make([]*[]common.Item, currentRedundancy)
	nodes := make(common.Remotes, currentRedundancy)
	nextKey := r.Key
	var nextSuccessor *common.Remote
	for i := 0; i < currentRedundancy; i++ {
		_, _, nextSuccessor = self.ring.Remotes(nextKey)
		var thisResult []common.Item
		nodes[i] = *nextSuccessor
		results[i] = &thisResult
		futures[i] = nextSuccessor.Go(operation, r, &thisResult)
		nextKey = nextSuccessor.Pos
	}
	for index, future := range futures {
		<-future.Done
		if future.Error != nil {
			self.removeNode(nodes[index])
			return self.mergeRecent(operation, r, up)
		}
	}
	result = common.MergeItems(results, up)
	return
}
func (self *Conn) findRecent(operation string, data common.Item) (result *common.Item) {
	currentRedundancy := self.ring.Redundancy()
	futures := make([]*rpc.Call, currentRedundancy)
	results := make([]*common.Item, currentRedundancy)
	nodes := make(common.Remotes, currentRedundancy)
	nextKey := data.Key
	var nextSuccessor *common.Remote
	for i := 0; i < currentRedundancy; i++ {
		_, _, nextSuccessor = self.ring.Remotes(nextKey)
		thisResult := &common.Item{}
		nodes[i] = *nextSuccessor
		results[i] = thisResult
		futures[i] = nextSuccessor.Go(operation, data, thisResult)
		nextKey = nextSuccessor.Pos
	}
	for index, future := range futures {
		<-future.Done
		if future.Error != nil {
			self.removeNode(nodes[index])
			return self.findRecent(operation, data)
		}
		if result == nil || result.Timestamp < results[index].Timestamp {
			result = results[index]
		}
	}
	return
}
func (self *Conn) consume(c chan [2][]byte, wait *sync.WaitGroup, successor *common.Remote) {
	for pair := range c {
		self.putVia(successor, pair[0], pair[1], false)
	}
	wait.Done()
}
func (self *Conn) dump(c chan [2][]byte, wait *sync.WaitGroup) {
	var succ *common.Remote
	dumps := make(map[string]chan [2][]byte)
	for pair := range c {
		_, _, succ = self.ring.Remotes(pair[0])
		if dump, ok := dumps[succ.Addr]; ok {
			dump <- pair
		} else {
			newDump := make(chan [2][]byte, 16)
			wait.Add(1)
			go self.consume(newDump, wait, succ)
			newDump <- pair
			dumps[succ.Addr] = newDump
		}
	}
	for _, dump := range dumps {
		close(dump)
	}
	wait.Done()
}
func (self *Conn) subDump(key []byte, c chan [2][]byte, wait *sync.WaitGroup) {
	_, _, succ := self.ring.Remotes(key)
	for pair := range c {
		self.subPutVia(succ, key, pair[0], pair[1], false)
	}
	wait.Done()
}

// Clear will remove all data from all currently known database nodes.
func (self *Conn) Clear() {
	var x int
	for _, node := range self.ring.Nodes() {
		if err := node.Call("DHash.Clear", 0, &x); err != nil {
			self.removeNode(node)
		}
	}
}

// SSubPut will put value under subKey in the sub tree defined by key.
func (self *Conn) SSubPut(key, subKey, value []byte) {
	self.subPut(key, subKey, value, true)
}

// SubPut will put value under subKey in the sub tree defined by key.
func (self *Conn) SubPut(key, subKey, value []byte) {
	self.subPut(key, subKey, value, false)
}

// SPut will put value under key.
func (self *Conn) SPut(key, value []byte) {
	self.put(key, value, true)
}

// Put will put value under key.
func (self *Conn) Put(key, value []byte) {
	self.put(key, value, false)
}

// Dump will return a channel to send multiple key/value pairs through. When finished, close the channel and #Wait for the *sync.WaitGroup.
func (self *Conn) Dump() (c chan [2][]byte, wait *sync.WaitGroup) {
	wait = new(sync.WaitGroup)
	c = make(chan [2][]byte, 16)
	wait.Add(1)
	go self.dump(c, wait)
	return
}

// SubDump will return a channel to send multiple key/value pairs to a given sub tree through. When finished, close the channel and #Wait for the *sync.WaitGroup.
func (self *Conn) SubDump(key []byte) (c chan [2][]byte, wait *sync.WaitGroup) {
	wait = new(sync.WaitGroup)
	c = make(chan [2][]byte)
	wait.Add(1)
	go self.subDump(key, c, wait)
	return
}

// SubClear will remove all byte values from the sub tree defined by key. It will retain delete markers for all deleted values.
func (self *Conn) SubClear(key []byte) {
	self.subClear(key, false)
}

// SubClear will remove all byte values from the sub tree defined by key. It will retain delete markers for all deleted values.
func (self *Conn) SSubClear(key []byte) {
	self.subClear(key, true)
}

// SubDel will remove the value under subKey from the sub tree defined by key.
func (self *Conn) SubDel(key, subKey []byte) {
	self.subDel(key, subKey, false)
}

// SSubDel will remove the value under subKey from the sub tree defined by key.
func (self *Conn) SSubDel(key, subKey []byte) {
	self.subDel(key, subKey, true)
}

// SDel will remove the byte value under key.
func (self *Conn) SDel(key []byte) {
	self.del(key, true)
}

// Del will remove the byte value under key.
func (self *Conn) Del(key []byte) {
	self.del(key, false)
}

// MirrorReverseIndexOf will return the the distance from the end for subKey, looking at the mirror tree of the sub tree defined by key.
func (self *Conn) MirrorReverseIndexOf(key, subKey []byte) (index int, existed bool) {
	data := common.Item{
		Key:    key,
		SubKey: subKey,
	}
	_, _, successor := self.ring.Remotes(key)
	var result common.Index
	if err := successor.Call("DHash.MirrorReverseIndexOf", data, &result); err != nil {
		self.removeNode(*successor)
		return self.MirrorReverseIndexOf(key, subKey)
	}
	index, existed = result.N, result.Existed
	return
}

// MirrorIndexOf will return the the distance from the start for subKey, looking at the mirror tree of the sub tree defined by key.
func (self *Conn) MirrorIndexOf(key, subKey []byte) (index int, existed bool) {
	data := common.Item{
		Key:    key,
		SubKey: subKey,
	}
	_, _, successor := self.ring.Remotes(key)
	var result common.Index
	if err := successor.Call("DHash.MirrorIndexOf", data, &result); err != nil {
		self.removeNode(*successor)
		return self.MirrorIndexOf(key, subKey)
	}
	index, existed = result.N, result.Existed
	return
}

// ReverseIndexOf will return the the distance from the end for subKey, looking at the sub tree defined by key.
func (self *Conn) ReverseIndexOf(key, subKey []byte) (index int, existed bool) {
	data := common.Item{
		Key:    key,
		SubKey: subKey,
	}
	_, _, successor := self.ring.Remotes(key)
	var result common.Index
	if err := successor.Call("DHash.ReverseIndexOf", data, &result); err != nil {
		self.removeNode(*successor)
		return self.ReverseIndexOf(key, subKey)
	}
	index, existed = result.N, result.Existed
	return
}

// IndexOf will return the the distance from the start for subKey, looking at the sub tree defined by key.
func (self *Conn) IndexOf(key, subKey []byte) (index int, existed bool) {
	data := common.Item{
		Key:    key,
		SubKey: subKey,
	}
	_, _, successor := self.ring.Remotes(key)
	var result common.Index
	if err := successor.Call("DHash.IndexOf", data, &result); err != nil {
		self.removeNode(*successor)
		return self.IndexOf(key, subKey)
	}
	index, existed = result.N, result.Existed
	return
}

// Next will return the next key and value after key.
func (self *Conn) Next(key []byte) (nextKey, nextValue []byte, existed bool) {
	data := common.Item{
		Key: key,
	}
	result := &common.Item{}
	_, _, successor := self.ring.Remotes(key)
	firstAddr := successor.Addr
	for {
		if err := successor.Call("DHash.Next", data, result); err != nil {
			self.removeNode(*successor)
			return self.Next(key)
		}
		if result.Exists {
			break
		}
		_, _, successor = self.ring.Remotes(successor.Pos)
		if successor.Addr == firstAddr {
			break
		}
	}
	nextKey, nextValue, existed = result.Key, result.Value, result.Exists
	return
}

// Prev will return the previous key and value before key.
func (self *Conn) Prev(key []byte) (prevKey, prevValue []byte, existed bool) {
	data := common.Item{
		Key: key,
	}
	result := &common.Item{}
	_, _, successor := self.ring.Remotes(key)
	firstAddr := successor.Addr
	for {
		if err := successor.Call("DHash.Prev", data, result); err != nil {
			self.removeNode(*successor)
			return self.Prev(key)
		}
		if result.Exists {
			break
		}
		successor, _, _ = self.ring.Remotes(successor.Pos)
		if successor.Addr == firstAddr {
			break
		}
	}
	prevKey, prevValue, existed = result.Key, result.Value, result.Exists
	return
}

// MirrorCount will count the number of keys between min and max in the mirror tree of the sub tree defined by key.
func (self *Conn) MirrorCount(key, min, max []byte, mininc, maxinc bool) (result int) {
	r := common.Range{
		Key:    key,
		Min:    min,
		Max:    max,
		MinInc: mininc,
		MaxInc: maxinc,
	}
	_, _, successor := self.ring.Remotes(key)
	if err := successor.Call("DHash.MirrorCount", r, &result); err != nil {
		self.removeNode(*successor)
		return self.MirrorCount(key, min, max, mininc, maxinc)
	}
	return
}

// Count will count the number of keys between min and max in the sub tree defined by key.
func (self *Conn) Count(key, min, max []byte, mininc, maxinc bool) (result int) {
	r := common.Range{
		Key:    key,
		Min:    min,
		Max:    max,
		MinInc: mininc,
		MaxInc: maxinc,
	}
	_, _, successor := self.ring.Remotes(key)
	if err := successor.Call("DHash.Count", r, &result); err != nil {
		self.removeNode(*successor)
		return self.Count(key, min, max, mininc, maxinc)
	}
	return
}

// MirrorNextIndex will return the key, value and index of the first key after index in the mirror tree of the sub tree defined by key.
func (self *Conn) MirrorNextIndex(key []byte, index int) (foundKey, foundValue []byte, foundIndex int, existed bool) {
	data := common.Item{
		Key:   key,
		Index: index,
	}
	result := &common.Item{}
	_, _, successor := self.ring.Remotes(key)
	if err := successor.Call("DHash.MirrorNextIndex", data, result); err != nil {
		self.removeNode(*successor)
		return self.MirrorNextIndex(key, index)
	}
	foundKey, foundValue, foundIndex, existed = result.Key, result.Value, result.Index, result.Exists
	return
}

// MirrorPrevIndex will return the key, value and index of the first key before index in the mirror tree of the sub tree defined by key.
func (self *Conn) MirrorPrevIndex(key []byte, index int) (foundKey, foundValue []byte, foundIndex int, existed bool) {
	data := common.Item{
		Key:   key,
		Index: index,
	}
	result := &common.Item{}
	_, _, successor := self.ring.Remotes(key)
	if err := successor.Call("DHash.MirrorPrevIndex", data, result); err != nil {
		self.removeNode(*successor)
		return self.MirrorNextIndex(key, index)
	}
	foundKey, foundValue, foundIndex, existed = result.Key, result.Value, result.Index, result.Exists
	return
}

// NextIndex will return the key, value and index of the first key after index in the sub tree defined by key.
func (self *Conn) NextIndex(key []byte, index int) (foundKey, foundValue []byte, foundIndex int, existed bool) {
	data := common.Item{
		Key:   key,
		Index: index,
	}
	result := &common.Item{}
	_, _, successor := self.ring.Remotes(key)
	if err := successor.Call("DHash.NextIndex", data, result); err != nil {
		self.removeNode(*successor)
		return self.NextIndex(key, index)
	}
	foundKey, foundValue, foundIndex, existed = result.Key, result.Value, result.Index, result.Exists
	return
}

// PrevIndex will return the key, value and index of the first key before index in the sub tree defined by key.
func (self *Conn) PrevIndex(key []byte, index int) (foundKey, foundValue []byte, foundIndex int, existed bool) {
	data := common.Item{
		Key:   key,
		Index: index,
	}
	result := &common.Item{}
	_, _, successor := self.ring.Remotes(key)
	if err := successor.Call("DHash.PrevIndex", data, result); err != nil {
		self.removeNode(*successor)
		return self.NextIndex(key, index)
	}
	foundKey, foundValue, foundIndex, existed = result.Key, result.Value, result.Index, result.Exists
	return
}

// MirrorReverseSliceIndex will return the reverse slice between index min and max in the mirror tree of the sub tree defined by key.
// A min of nil will return from the end. A max of nil will return to the start.
func (self *Conn) MirrorReverseSliceIndex(key []byte, min, max *int) (result []common.Item) {
	var mi int
	var ma int
	if min != nil {
		mi = *min
	}
	if max != nil {
		ma = *max
	}
	r := common.Range{
		Key:      key,
		MinIndex: mi,
		MaxIndex: ma,
		MinInc:   min != nil,
		MaxInc:   max != nil,
	}
	result = self.mergeRecent("DHash.MirrorReverseSliceIndex", r, false)
	return
}

// MirrorSliceIndex will return the slice between index min and max in the mirror tree of the sub tree defined by key.
// A min of nil will return from the start. A max of nil will return to the end.
func (self *Conn) MirrorSliceIndex(key []byte, min, max *int) (result []common.Item) {
	var mi int
	var ma int
	if min != nil {
		mi = *min
	}
	if max != nil {
		ma = *max
	}
	r := common.Range{
		Key:      key,
		MinIndex: mi,
		MaxIndex: ma,
		MinInc:   min != nil,
		MaxInc:   max != nil,
	}
	result = self.mergeRecent("DHash.MirrorSliceIndex", r, true)
	return
}

// MirrorReverseSlice will return the reverse slice between min and max in the mirror tree of the sub tree defined by key.
// A min of nil will return from the end. A max of nil will return to the start.
func (self *Conn) MirrorReverseSlice(key, min, max []byte, mininc, maxinc bool) (result []common.Item) {
	r := common.Range{
		Key:    key,
		Min:    min,
		Max:    max,
		MinInc: mininc,
		MaxInc: maxinc,
	}
	result = self.mergeRecent("DHash.MirrorReverseSlice", r, false)
	return
}

// MirrorSlice will return the slice between min and max in the mirror tree of the sub tree defined by key.
// A min of nil will return from the start. A max of nil will return to the end.
func (self *Conn) MirrorSlice(key, min, max []byte, mininc, maxinc bool) (result []common.Item) {
	r := common.Range{
		Key:    key,
		Min:    min,
		Max:    max,
		MinInc: mininc,
		MaxInc: maxinc,
	}
	result = self.mergeRecent("DHash.MirrorSlice", r, true)
	return
}

// MirrorSliceLen will return at most maxRes elements after min in the mirror tree of the sub tree defined by key.
// A min of nil will return from the start.
func (self *Conn) MirrorSliceLen(key, min []byte, mininc bool, maxRes int) (result []common.Item) {
	r := common.Range{
		Key:    key,
		Min:    min,
		MinInc: mininc,
		Len:    maxRes,
	}
	result = self.mergeRecent("DHash.MirrorSliceLen", r, true)
	return
}

// MirrorReverseSliceLen will return at most maxRes elements before max in the mirror tree of the sub tree defined by key.
// A min of nil will return from the end.
func (self *Conn) MirrorReverseSliceLen(key, max []byte, maxinc bool, maxRes int) (result []common.Item) {
	r := common.Range{
		Key:    key,
		Max:    max,
		MaxInc: maxinc,
		Len:    maxRes,
	}
	result = self.mergeRecent("DHash.MirrorReverseSliceLen", r, false)
	return
}

// ReverseSliceIndex will the reverse slice between index min and max in the sub tree defined by key.
// A min of nil will return from the end. A max of nil will return to the start.
func (self *Conn) ReverseSliceIndex(key []byte, min, max *int) (result []common.Item) {
	var mi int
	var ma int
	if min != nil {
		mi = *min
	}
	if max != nil {
		ma = *max
	}
	r := common.Range{
		Key:      key,
		MinIndex: mi,
		MaxIndex: ma,
		MinInc:   min != nil,
		MaxInc:   max != nil,
	}
	result = self.mergeRecent("DHash.ReverseSliceIndex", r, false)
	return
}

// SliceIndex will return the slice between index min and max in the sub tree defined by key.
// A min of nil will return from the start. A max of nil will return to the end.
func (self *Conn) SliceIndex(key []byte, min, max *int) (result []common.Item) {
	var mi int
	var ma int
	if min != nil {
		mi = *min
	}
	if max != nil {
		ma = *max
	}
	r := common.Range{
		Key:      key,
		MinIndex: mi,
		MaxIndex: ma,
		MinInc:   min != nil,
		MaxInc:   max != nil,
	}
	result = self.mergeRecent("DHash.SliceIndex", r, true)
	return
}

// ReverseSlice will return the reverse slice between min and max in the sub tree defined by key.
// A min of nil will return from the end. A max of nil will return to the start.
func (self *Conn) ReverseSlice(key, min, max []byte, mininc, maxinc bool) (result []common.Item) {
	r := common.Range{
		Key:    key,
		Min:    min,
		Max:    max,
		MinInc: mininc,
		MaxInc: maxinc,
	}
	result = self.mergeRecent("DHash.ReverseSlice", r, false)
	return
}

// ReverseSlice will return the slice between min and max in the sub tree defined by key.
// A min of nil will return from the start. A max of nil will return to the end.
func (self *Conn) Slice(key, min, max []byte, mininc, maxinc bool) (result []common.Item) {
	r := common.Range{
		Key:    key,
		Min:    min,
		Max:    max,
		MinInc: mininc,
		MaxInc: maxinc,
	}
	result = self.mergeRecent("DHash.Slice", r, true)
	return
}

// SliceLen will return at most maxRes elements after min in the sub tree defined by key.
// A min of nil will return from the start.
func (self *Conn) SliceLen(key, min []byte, mininc bool, maxRes int) (result []common.Item) {
	r := common.Range{
		Key:    key,
		Min:    min,
		MinInc: mininc,
		Len:    maxRes,
	}
	result = self.mergeRecent("DHash.SliceLen", r, true)
	return
}

// ReverseSliceLen will return at most maxRes elements before max in the sub tree defined by key.
// A min of nil will return from the end. A max of nil will return to the start.
func (self *Conn) ReverseSliceLen(key, max []byte, maxinc bool, maxRes int) (result []common.Item) {
	r := common.Range{
		Key:    key,
		Max:    max,
		MaxInc: maxinc,
		Len:    maxRes,
	}
	result = self.mergeRecent("DHash.ReverseSliceLen", r, false)
	return
}

// SubMirrorPrev will return the previous key and value before subKey in the sub tree defined by key.
func (self *Conn) SubMirrorPrev(key, subKey []byte) (prevKey, prevValue []byte, existed bool) {
	data := common.Item{
		Key:    key,
		SubKey: subKey,
	}
	result := self.findRecent("DHash.SubMirrorPrev", data)
	prevKey, prevValue, existed = result.Key, result.Value, result.Exists
	return
}

// SubMirrorPrev will return the next key and value after subKey in the sub tree defined by key.
func (self *Conn) SubMirrorNext(key, subKey []byte) (nextKey, nextValue []byte, existed bool) {
	data := common.Item{
		Key:    key,
		SubKey: subKey,
	}
	result := self.findRecent("DHash.SubMirrorNext", data)
	nextKey, nextValue, existed = result.Key, result.Value, result.Exists
	return
}

// SubPrev will return the previous key and value before subKey in the sub tree defined by key.
func (self *Conn) SubPrev(key, subKey []byte) (prevKey, prevValue []byte, existed bool) {
	data := common.Item{
		Key:    key,
		SubKey: subKey,
	}
	result := self.findRecent("DHash.SubPrev", data)
	prevKey, prevValue, existed = result.Key, result.Value, result.Exists
	return
}

// SubNext will return the next key and value after subKey in the sub tree defined by key.
func (self *Conn) SubNext(key, subKey []byte) (nextKey, nextValue []byte, existed bool) {
	data := common.Item{
		Key:    key,
		SubKey: subKey,
	}
	result := self.findRecent("DHash.SubNext", data)
	nextKey, nextValue, existed = result.Key, result.Value, result.Exists
	return
}

// MirrorLast will return the last key and valuei in the mirror tree of the sub tree defined by key.
func (self *Conn) MirrorLast(key []byte) (lastKey, lastValue []byte, existed bool) {
	data := common.Item{
		Key: key,
	}
	result := self.findRecent("DHash.MirrorLast", data)
	lastKey, lastValue, existed = result.Key, result.Value, result.Exists
	return
}

// MirrorFirst will return the first key and value in the mirror tree of the sub tree defined by key.
func (self *Conn) MirrorFirst(key []byte) (firstKey, firstValue []byte, existed bool) {
	data := common.Item{
		Key: key,
	}
	result := self.findRecent("DHash.MirrorFirst", data)
	firstKey, firstValue, existed = result.Key, result.Value, result.Exists
	return
}

// Last will return the last key and value in the sub tree defined by key.
func (self *Conn) Last(key []byte) (lastKey, lastValue []byte, existed bool) {
	data := common.Item{
		Key: key,
	}
	result := self.findRecent("DHash.Last", data)
	lastKey, lastValue, existed = result.Key, result.Value, result.Exists
	return
}

// First will return the first key and value in the sub tree defined by key.
func (self *Conn) First(key []byte) (firstKey, firstValue []byte, existed bool) {
	data := common.Item{
		Key: key,
	}
	result := self.findRecent("DHash.First", data)
	firstKey, firstValue, existed = result.Key, result.Value, result.Exists
	return
}

// SubGet will return the value under subKey in the sub tree defined by key.
func (self *Conn) SubGet(key, subKey []byte) (value []byte, existed bool) {
	data := common.Item{
		Key:    key,
		SubKey: subKey,
	}
	result := self.findRecent("DHash.SubGet", data)
	if result.Value != nil {
		value, existed = result.Value, result.Exists
	} else {
		value, existed = nil, false
	}
	return
}

// Get will return the value under key.
func (self *Conn) Get(key []byte) (value []byte, existed bool) {
	data := common.Item{
		Key: key,
	}
	result := self.findRecent("DHash.Get", data)
	if result.Value != nil {
		value, existed = result.Value, result.Exists
	} else {
		value, existed = nil, false
	}
	return
}

// DescribeTree will return a string representation of the complete tree in the node at pos.
// Used for debug purposes, don't do it on big databases!
func (self *Conn) DescribeTree(pos []byte) (result string, err error) {
	_, match, _ := self.ring.Remotes(pos)
	if match == nil {
		err = fmt.Errorf("No node with position %v found", common.HexEncode(pos))
		return
	}
	err = match.Call("DHash.DescribeTree", 0, &result)
	return
}

// DescribeTree will return a string representation of the complete trees of all known nodes.
// Used for debug purposes, don't do it on big databases!
func (self *Conn) DescribeAllTrees() string {
	buf := new(bytes.Buffer)
	for _, rem := range self.ring.Nodes() {
		if res, err := self.DescribeTree(rem.Pos); err == nil {
			fmt.Fprintln(buf, res)
		}
	}
	return string(buf.Bytes())
}

// DescribeAllNodes will return the description structures of all known nodes.
func (self *Conn) DescribeAllNodes() (result []common.DHashDescription) {
	for _, rem := range self.ring.Nodes() {
		if res, err := self.DescribeNode(rem.Pos); err == nil {
			result = append(result, res)
		}
	}
	return
}

// DescribeNode will return the description structure of the node at pos.
func (self *Conn) DescribeNode(pos []byte) (result common.DHashDescription, err error) {
	_, match, _ := self.ring.Remotes(pos)
	if match == nil {
		err = fmt.Errorf("No node with position %v found", common.HexEncode(pos))
		return
	}
	err = match.Call("DHash.Describe", 0, &result)
	return
}

// SubSize will return the size of the sub tree defined by key.
func (self *Conn) SubSize(key []byte) (result int) {
	_, _, successor := self.ring.Remotes(key)
	if err := successor.Call("DHash.SubSize", key, &result); err != nil {
		self.removeNode(*successor)
		return self.SubSize(key)
	}
	return
}

// Size will return the total size of all known nodes.
func (self *Conn) Size() (result int) {
	var tmp int
	for _, node := range self.ring.Nodes() {
		if err := node.Call("DHash.Size", 0, &tmp); err != nil {
			self.removeNode(node)
			return self.Size()
		}
		result += tmp
	}
	return
}

// Describe will return a string representation of the known cluster of nodes.
func (self *Conn) Describe() string {
	return self.ring.Describe()
}

// SetExpression will execute the given expr.
//
// If expr.Dest is set it will store the result under the sub tree defined by expr.Dest.
//
// If expr.Dest is nil it will return the result.
//
// Either expr.Op or expr.Code has to be set.
//
// If expr.Op is nil expr.Code will be parsed using SetOpParser to provide expr.Op.
func (self *Conn) SetExpression(expr setop.SetExpression) (result []setop.SetOpResult) {
	if expr.Op == nil {
		expr.Op = setop.MustParse(expr.Code)
	}
	var biggestKey []byte
	biggestSize := 0
	var thisSize int

	for key, _ := range findKeys(expr.Op) {
		thisSize = self.SubSize([]byte(key))
		if biggestKey == nil {
			biggestKey = []byte(key)
			biggestSize = thisSize
		} else if thisSize > biggestSize {
			biggestKey = []byte(key)
			biggestSize = thisSize
		}
	}
	_, _, successor := self.ring.Remotes(biggestKey)
	var results []setop.SetOpResult
	err := successor.Call("DHash.SetExpression", expr, &results)
	for err != nil {
		self.removeNode(*successor)
		_, _, successor = self.ring.Remotes(biggestKey)
		err = successor.Call("DHash.SetExpression", expr, &results)
	}
	return results
}

// Configuration will return the configuration for the entire cluster.
// Not internally used for anything right now.
func (self *Conn) Configuration() (conf map[string]string) {
	var result common.Conf
	_, _, successor := self.ring.Remotes(nil)
	if err := successor.Call("DHash.Configuration", 0, &result); err != nil {
		self.removeNode(*successor)
		return self.Configuration()
	}
	return result.Data
}

// SubConfiguratino will return the configuration for the sub tree defined by key.
//
// mirrored=yes means that the sub tree is currently mirrored.
func (self *Conn) SubConfiguration(key []byte) (conf map[string]string) {
	var result common.Conf
	_, _, successor := self.ring.Remotes(nil)
	if err := successor.Call("DHash.SubConfiguration", key, &result); err != nil {
		self.removeNode(*successor)
		return self.Configuration()
	}
	return result.Data
}

// AddConfiguration will set a key and value to the cluster configuration.
// Not internally used for anything right now.
func (self *Conn) AddConfiguration(key, value string) {
	conf := common.ConfItem{
		Key:   key,
		Value: value,
	}
	_, _, successor := self.ring.Remotes(nil)
	var x int
	if err := successor.Call("DHash.AddConfiguration", conf, &x); err != nil {
		self.removeNode(*successor)
		self.AddConfiguration(key, value)
	}
}

// SubAddConfiguration will set a key and value to the configuration of the sub tree defined by key.
//
// To mirror a sub tree, set mirrored=yes. To turn off mirroring of a sub tree, set mirrored!=yes.
func (self *Conn) SubAddConfiguration(treeKey []byte, key, value string) {
	conf := common.ConfItem{
		TreeKey: treeKey,
		Key:     key,
		Value:   value,
	}
	_, _, successor := self.ring.Remotes(nil)
	var x int
	if err := successor.Call("DHash.SubAddConfiguration", conf, &x); err != nil {
		self.removeNode(*successor)
		self.AddConfiguration(key, value)
	}
}
