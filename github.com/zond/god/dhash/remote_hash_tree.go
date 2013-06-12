package dhash

import (
	"github.com/zond/god/common"
	"github.com/zond/god/radix"
)

type remoteHashTree struct {
	destination common.Remote
	source      common.Remote
	node        *Node
}

func (self remoteHashTree) Configuration() (conf map[string]string, timestamp int64) {
	var result common.Conf
	if err := self.destination.Call("DHash.Configuration", 0, &result); err != nil {
		conf = make(map[string]string)
	} else {
		conf, timestamp = result.Data, result.Timestamp
	}
	return
}
func (self remoteHashTree) SubConfiguration(key []byte) (conf map[string]string, timestamp int64) {
	var result common.Conf
	if err := self.destination.Call("DHash.SubConfiguration", key, &result); err != nil {
		conf = make(map[string]string)
	} else {
		conf, timestamp = result.Data, result.Timestamp
	}
	return
}
func (self remoteHashTree) Configure(conf map[string]string, timestamp int64) {
	var x int
	self.destination.Call("HashTree.Configure", common.Conf{
		Data:      conf,
		Timestamp: timestamp,
	}, &x)
}
func (self remoteHashTree) SubConfigure(key []byte, conf map[string]string, timestamp int64) {
	var x int
	self.destination.Call("HashTree.SubConfigure", common.Conf{
		TreeKey:   key,
		Data:      conf,
		Timestamp: timestamp,
	}, &x)
}
func (self remoteHashTree) Hash() (result []byte) {
	self.destination.Call("HashTree.Hash", 0, &result)
	return
}
func (self remoteHashTree) Finger(key []radix.Nibble) (result *radix.Print) {
	result = &radix.Print{}
	self.destination.Call("HashTree.Finger", key, result)
	return
}
func (self remoteHashTree) GetTimestamp(key []radix.Nibble) (value []byte, timestamp int64, present bool) {
	result := HashTreeItem{}
	self.destination.Call("HashTree.GetTimestamp", key, &result)
	value, timestamp, present = result.Value, result.Timestamp, result.Exists
	return
}
func (self remoteHashTree) PutTimestamp(key []radix.Nibble, value []byte, present bool, expected, timestamp int64) (changed bool) {
	data := HashTreeItem{
		Key:       key,
		Value:     value,
		Exists:    present,
		Expected:  expected,
		Timestamp: timestamp,
	}
	op := "HashTree.PutTimestamp"
	if self.node.hasCommListeners() {
		self.node.triggerCommListeners(Comm{
			Key:         radix.Stitch(data.Key),
			Source:      self.source,
			Destination: self.destination,
			Type:        op,
		})
	}
	self.destination.Call(op, data, &changed)
	return
}
func (self remoteHashTree) DelTimestamp(key []radix.Nibble, expected int64) (changed bool) {
	data := HashTreeItem{
		Key:      key,
		Expected: expected,
	}
	op := "HashTree.DelTimestamp"
	if self.node.hasCommListeners() {
		self.node.triggerCommListeners(Comm{
			Key:         radix.Stitch(data.Key),
			Source:      self.source,
			Destination: self.destination,
			Type:        op,
		})
	}
	self.destination.Call(op, data, &changed)
	return
}
func (self remoteHashTree) SubFinger(key, subKey []radix.Nibble) (result *radix.Print) {
	data := HashTreeItem{
		Key:    key,
		SubKey: subKey,
	}
	result = &radix.Print{}
	self.destination.Call("HashTree.SubFinger", data, result)
	return
}
func (self remoteHashTree) SubGetTimestamp(key, subKey []radix.Nibble) (value []byte, timestamp int64, present bool) {
	data := HashTreeItem{
		Key:    key,
		SubKey: subKey,
	}
	self.destination.Call("HashTree.SubGetTimestamp", data, &data)
	value, timestamp, present = data.Value, data.Timestamp, data.Exists
	return
}
func (self remoteHashTree) SubPutTimestamp(key, subKey []radix.Nibble, value []byte, present bool, subExpected, subTimestamp int64) (changed bool) {
	data := HashTreeItem{
		Key:       key,
		SubKey:    subKey,
		Value:     value,
		Exists:    present,
		Expected:  subExpected,
		Timestamp: subTimestamp,
	}
	op := "HashTree.SubPutTimestamp"
	if self.node.hasCommListeners() {
		self.node.triggerCommListeners(Comm{
			Key:         radix.Stitch(data.Key),
			SubKey:      radix.Stitch(data.SubKey),
			Source:      self.source,
			Destination: self.destination,
			Type:        op,
		})
	}
	self.destination.Call(op, data, &changed)
	return
}
func (self remoteHashTree) SubDelTimestamp(key, subKey []radix.Nibble, subExpected int64) (changed bool) {
	data := HashTreeItem{
		Key:      key,
		SubKey:   subKey,
		Expected: subExpected,
	}
	op := "HashTree.SubDelTimestamp"
	if self.node.hasCommListeners() {
		self.node.triggerCommListeners(Comm{
			Key:         radix.Stitch(data.Key),
			SubKey:      radix.Stitch(data.SubKey),
			Source:      self.source,
			Destination: self.destination,
			Type:        op,
		})
	}
	self.destination.Call(op, data, &changed)
	return
}
func (self remoteHashTree) SubClearTimestamp(key []radix.Nibble, expected, timestamp int64) (deleted int) {
	data := HashTreeItem{
		Key:       key,
		Expected:  expected,
		Timestamp: timestamp,
	}
	op := "HashTree.SubClearTimestamp"
	if self.node.hasCommListeners() {
		self.node.triggerCommListeners(Comm{
			Key:         radix.Stitch(data.Key),
			Source:      self.source,
			Destination: self.destination,
			Type:        op,
		})
	}
	self.destination.Call(op, data, &deleted)
	return
}
func (self remoteHashTree) SubKillTimestamp(key []radix.Nibble, expected int64) (deleted int) {
	data := HashTreeItem{
		Key:      key,
		Expected: expected,
	}
	op := "HashTree.SubKillTimestamp"
	if self.node.hasCommListeners() {
		self.node.triggerCommListeners(Comm{
			Key:         radix.Stitch(data.Key),
			Source:      self.source,
			Destination: self.destination,
			Type:        op,
		})
	}
	self.destination.Call(op, data, &deleted)
	return
}
