package dhash

import (
	"github.com/zond/god/common"
	"github.com/zond/god/radix"
	"sync/atomic"
	"time"
)

type HashTreeItem struct {
	Key       []radix.Nibble
	SubKey    []radix.Nibble
	Timestamp int64
	Expected  int64
	Value     []byte
	Exists    bool
}

type hashTreeServer Node

func (self *hashTreeServer) Configure(conf common.Conf, x *int) error {
	atomic.StoreInt64(&(*Node)(self).lastSync, time.Now().UnixNano())
	(*Node)(self).tree.Configure(conf.Data, conf.Timestamp)
	return nil
}
func (self *hashTreeServer) SubConfigure(conf common.Conf, x *int) error {
	atomic.StoreInt64(&(*Node)(self).lastSync, time.Now().UnixNano())
	(*Node)(self).tree.SubConfigure(conf.TreeKey, conf.Data, conf.Timestamp)
	return nil
}
func (self *hashTreeServer) Hash(x int, result *[]byte) error {
	*result = (*Node)(self).tree.Hash()
	return nil
}
func (self *hashTreeServer) Finger(key []radix.Nibble, result *radix.Print) error {
	*result = *((*Node)(self).tree.Finger(key))
	return nil
}
func (self *hashTreeServer) GetTimestamp(key []radix.Nibble, result *HashTreeItem) error {
	atomic.StoreInt64(&(*Node)(self).lastSync, time.Now().UnixNano())
	*result = HashTreeItem{Key: key}
	result.Value, result.Timestamp, result.Exists = (*Node)(self).tree.GetTimestamp(key)
	return nil
}
func (self *hashTreeServer) PutTimestamp(data HashTreeItem, changed *bool) error {
	atomic.StoreInt64(&(*Node)(self).lastSync, time.Now().UnixNano())
	*changed = (*Node)(self).tree.PutTimestamp(data.Key, data.Value, data.Exists, data.Expected, data.Timestamp)
	return nil
}
func (self *hashTreeServer) DelTimestamp(data HashTreeItem, changed *bool) error {
	atomic.StoreInt64(&(*Node)(self).lastSync, time.Now().UnixNano())
	*changed = (*Node)(self).tree.DelTimestamp(data.Key, data.Expected)
	return nil
}
func (self *hashTreeServer) SubFinger(data HashTreeItem, result *radix.Print) error {
	*result = *((*Node)(self).tree.SubFinger(data.Key, data.SubKey))
	return nil
}
func (self *hashTreeServer) SubGetTimestamp(data HashTreeItem, result *HashTreeItem) error {
	atomic.StoreInt64(&(*Node)(self).lastSync, time.Now().UnixNano())
	*result = data
	result.Value, result.Timestamp, result.Exists = (*Node)(self).tree.SubGetTimestamp(data.Key, data.SubKey)
	return nil
}
func (self *hashTreeServer) SubPutTimestamp(data HashTreeItem, changed *bool) error {
	atomic.StoreInt64(&(*Node)(self).lastSync, time.Now().UnixNano())
	*changed = (*Node)(self).tree.SubPutTimestamp(data.Key, data.SubKey, data.Value, data.Exists, data.Expected, data.Timestamp)
	return nil
}
func (self *hashTreeServer) SubDelTimestamp(data HashTreeItem, changed *bool) error {
	atomic.StoreInt64(&(*Node)(self).lastSync, time.Now().UnixNano())
	*changed = (*Node)(self).tree.SubDelTimestamp(data.Key, data.SubKey, data.Expected)
	return nil
}
func (self *hashTreeServer) SubClearTimestamp(data HashTreeItem, changed *int) error {
	atomic.StoreInt64(&(*Node)(self).lastSync, time.Now().UnixNano())
	*changed = (*Node)(self).tree.SubClearTimestamp(data.Key, data.Expected, data.Timestamp)
	return nil
}
func (self *hashTreeServer) SubKillTimestamp(data HashTreeItem, changed *int) error {
	atomic.StoreInt64(&(*Node)(self).lastSync, time.Now().UnixNano())
	*changed = (*Node)(self).tree.SubKillTimestamp(data.Key, data.Expected)
	return nil
}
