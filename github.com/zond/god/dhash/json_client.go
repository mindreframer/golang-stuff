package dhash

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/zond/god/common"
	"github.com/zond/setop"
	"net/http"
)

// JSONClient is used in the tests to ensure that the JSON API provides roughly the same functionality as the gob API.
// It is also a demonstration and example of how the JSON API can be used.
// It is NOT meant to be used as a real client, since if you are using Go anyway the client.Conn type is much more efficient.
type JSONClient string

func (self JSONClient) call(action string, params, result interface{}) {
	client := new(http.Client)
	buf := new(bytes.Buffer)
	if params != nil {
		err := json.NewEncoder(buf).Encode(params)
		if err != nil {
			panic(err)
		}
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("http://%v/rpc/DHash.%v", self, action), buf)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Accept", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	err = json.NewDecoder(resp.Body).Decode(result)
	if err != nil {
		panic(err)
	}
	resp.Body.Close()
}
func (self JSONClient) SSubPut(key, subKey, value []byte) {
	var x Nothing
	item := SubValueOp{
		Key:    key,
		SubKey: subKey,
		Value:  value,
		Sync:   true,
	}
	self.call("SubPut", item, &x)
}
func (self JSONClient) SubPut(key, subKey, value []byte) {
	var x Nothing
	item := SubValueOp{
		Key:    key,
		SubKey: subKey,
		Value:  value,
	}
	self.call("SubPut", item, &x)
}
func (self JSONClient) SPut(key, value []byte) {
	var x Nothing
	item := ValueOp{
		Key:   key,
		Value: value,
		Sync:  true,
	}
	self.call("Put", item, &x)
}
func (self JSONClient) Put(key, value []byte) {
	var x Nothing
	item := ValueOp{
		Key:   key,
		Value: value,
	}
	self.call("Put", item, &x)
}
func (self JSONClient) SubClear(key []byte) {
	var x Nothing
	item := KeyOp{
		Key: key,
	}
	self.call("SubClear", item, &x)
}
func (self JSONClient) SSubClear(key []byte) {
	var x Nothing
	item := KeyOp{
		Key:  key,
		Sync: true,
	}
	self.call("SubClear", item, &x)
}
func (self JSONClient) SubDel(key, subKey []byte) {
	var x Nothing
	item := SubKeyOp{
		Key:    key,
		SubKey: subKey,
	}
	self.call("SubDel", item, &x)
}
func (self JSONClient) SSubDel(key, subKey []byte) {
	var x Nothing
	item := SubKeyOp{
		Key:    key,
		SubKey: subKey,
		Sync:   true,
	}
	self.call("SubDel", item, &x)
}
func (self JSONClient) SDel(key []byte) {
	var x Nothing
	item := KeyOp{
		Key:  key,
		Sync: true,
	}
	self.call("Del", item, &x)
}
func (self JSONClient) Del(key []byte) {
	var x Nothing
	item := KeyOp{
		Key: key,
	}
	self.call("Del", item, &x)
}
func (self JSONClient) MirrorReverseIndexOf(key, subKey []byte) (index int, existed bool) {
	item := SubKeyReq{
		Key:    key,
		SubKey: subKey,
	}
	var result common.Index
	self.call("MirrorReverseIndexOf", item, &result)
	return result.N, result.Existed
}
func (self JSONClient) MirrorIndexOf(key, subKey []byte) (index int, existed bool) {
	item := SubKeyReq{
		Key:    key,
		SubKey: subKey,
	}
	var result common.Index
	self.call("MirrorIndexOf", item, &result)
	return result.N, result.Existed
}
func (self JSONClient) ReverseIndexOf(key, subKey []byte) (index int, existed bool) {
	item := SubKeyReq{
		Key:    key,
		SubKey: subKey,
	}
	var result common.Index
	self.call("ReverseIndexOf", item, &result)
	return result.N, result.Existed
}
func (self JSONClient) IndexOf(key, subKey []byte) (index int, existed bool) {
	item := SubKeyReq{
		Key:    key,
		SubKey: subKey,
	}
	var result common.Index
	self.call("IndexOf", item, &result)
	return result.N, result.Existed
}
func (self JSONClient) Next(key []byte) (nextKey, nextValue []byte, existed bool) {
	item := KeyReq{
		Key: key,
	}
	var result common.Item
	self.call("Next", item, &result)
	return result.Key, result.Value, result.Exists
}
func (self JSONClient) Prev(key []byte) (prevKey, prevValue []byte, existed bool) {
	item := KeyReq{
		Key: key,
	}
	var result common.Item
	self.call("Prev", item, &result)
	return result.Key, result.Value, result.Exists
}
func (self JSONClient) MirrorCount(key, min, max []byte, mininc, maxinc bool) (result int) {
	item := KeyRange{
		Key:    key,
		Min:    min,
		Max:    max,
		MinInc: mininc,
		MaxInc: maxinc,
	}
	self.call("MirrorCount", item, &result)
	return result
}
func (self JSONClient) Count(key, min, max []byte, mininc, maxinc bool) (result int) {
	item := KeyRange{
		Key:    key,
		Min:    min,
		Max:    max,
		MinInc: mininc,
		MaxInc: maxinc,
	}
	self.call("Count", item, &result)
	return result
}
func (self JSONClient) MirrorNextIndex(key []byte, index int) (foundKey, foundValue []byte, foundIndex int, existed bool) {
	item := SubIndex{
		Key:   key,
		Index: index,
	}
	var result common.Item
	self.call("MirrorNextIndex", item, &result)
	return result.Key, result.Value, result.Index, result.Exists
}
func (self JSONClient) MirrorPrevIndex(key []byte, index int) (foundKey, foundValue []byte, foundIndex int, existed bool) {
	item := SubIndex{
		Key:   key,
		Index: index,
	}
	var result common.Item
	self.call("MirrorPrevIndex", item, &result)
	return result.Key, result.Value, result.Index, result.Exists
}
func (self JSONClient) NextIndex(key []byte, index int) (foundKey, foundValue []byte, foundIndex int, existed bool) {
	item := SubIndex{
		Key:   key,
		Index: index,
	}
	var result common.Item
	self.call("NextIndex", item, &result)
	return result.Key, result.Value, result.Index, result.Exists
}
func (self JSONClient) PrevIndex(key []byte, index int) (foundKey, foundValue []byte, foundIndex int, existed bool) {
	item := SubIndex{
		Key:   key,
		Index: index,
	}
	var result common.Item
	self.call("PrevIndex", item, &result)
	return result.Key, result.Value, result.Index, result.Exists
}
func (self JSONClient) MirrorReverseSliceIndex(key []byte, min, max *int) (result []common.Item) {
	item := IndexRange{
		Key:      key,
		MinIndex: min,
		MaxIndex: max,
	}
	self.call("MirrorReverseSliceIndex", item, &result)
	return result
}
func (self JSONClient) MirrorSliceIndex(key []byte, min, max *int) (result []common.Item) {
	item := IndexRange{
		Key:      key,
		MinIndex: min,
		MaxIndex: max,
	}
	self.call("MirrorSliceIndex", item, &result)
	return result
}
func (self JSONClient) MirrorReverseSlice(key, min, max []byte, mininc, maxinc bool) (result []common.Item) {
	item := KeyRange{
		Key:    key,
		Min:    min,
		Max:    max,
		MinInc: mininc,
		MaxInc: maxinc,
	}
	self.call("MirrorReverseSlice", item, &result)
	return result
}
func (self JSONClient) MirrorSlice(key, min, max []byte, mininc, maxinc bool) (result []common.Item) {
	item := KeyRange{
		Key:    key,
		Min:    min,
		Max:    max,
		MinInc: mininc,
		MaxInc: maxinc,
	}
	self.call("MirrorSlice", item, &result)
	return result
}
func (self JSONClient) MirrorSliceLen(key, min []byte, mininc bool, maxRes int) (result []common.Item) {
	item := PageRange{
		Key:     key,
		From:    min,
		FromInc: mininc,
		Len:     maxRes,
	}
	self.call("MirrorSliceLen", item, &result)
	return result
}
func (self JSONClient) MirrorReverseSliceLen(key, max []byte, maxinc bool, maxRes int) (result []common.Item) {
	item := PageRange{
		Key:     key,
		From:    max,
		FromInc: maxinc,
		Len:     maxRes,
	}
	self.call("MirrorReverseSliceLen", item, &result)
	return result
}
func (self JSONClient) ReverseSliceIndex(key []byte, min, max *int) (result []common.Item) {
	item := IndexRange{
		Key:      key,
		MinIndex: min,
		MaxIndex: max,
	}
	self.call("ReverseSliceIndex", item, &result)
	return result
}
func (self JSONClient) SliceIndex(key []byte, min, max *int) (result []common.Item) {
	item := IndexRange{
		Key:      key,
		MinIndex: min,
		MaxIndex: max,
	}
	self.call("SliceIndex", item, &result)
	return result
}
func (self JSONClient) ReverseSlice(key, min, max []byte, mininc, maxinc bool) (result []common.Item) {
	item := KeyRange{
		Key:    key,
		Min:    min,
		Max:    max,
		MinInc: mininc,
		MaxInc: maxinc,
	}
	self.call("ReverseSlice", item, &result)
	return result
}
func (self JSONClient) Slice(key, min, max []byte, mininc, maxinc bool) (result []common.Item) {
	item := KeyRange{
		Key:    key,
		Min:    min,
		Max:    max,
		MinInc: mininc,
		MaxInc: maxinc,
	}
	self.call("Slice", item, &result)
	return result
}
func (self JSONClient) SliceLen(key, min []byte, mininc bool, maxRes int) (result []common.Item) {
	item := PageRange{
		Key:     key,
		From:    min,
		FromInc: mininc,
		Len:     maxRes,
	}
	self.call("SliceLen", item, &result)
	return result
}
func (self JSONClient) ReverseSliceLen(key, max []byte, maxinc bool, maxRes int) (result []common.Item) {
	item := PageRange{
		Key:     key,
		From:    max,
		FromInc: maxinc,
		Len:     maxRes,
	}
	self.call("ReverseSliceLen", item, &result)
	return result
}
func (self JSONClient) SubMirrorPrev(key, subKey []byte) (prevKey, prevValue []byte, existed bool) {
	item := SubKeyReq{
		Key:    key,
		SubKey: subKey,
	}
	var result common.Item
	self.call("SubMirrorPrev", item, &result)
	return result.Key, result.Value, result.Exists
}
func (self JSONClient) SubMirrorNext(key, subKey []byte) (nextKey, nextValue []byte, existed bool) {
	item := SubKeyReq{
		Key:    key,
		SubKey: subKey,
	}
	var result common.Item
	self.call("SubMirrorNext", item, &result)
	return result.Key, result.Value, result.Exists
}
func (self JSONClient) SubPrev(key, subKey []byte) (prevKey, prevValue []byte, existed bool) {
	item := SubKeyReq{
		Key:    key,
		SubKey: subKey,
	}
	var result common.Item
	self.call("SubPrev", item, &result)
	return result.Key, result.Value, result.Exists
}
func (self JSONClient) SubNext(key, subKey []byte) (nextKey, nextValue []byte, existed bool) {
	item := SubKeyReq{
		Key:    key,
		SubKey: subKey,
	}
	var result common.Item
	self.call("SubNext", item, &result)
	return result.Key, result.Value, result.Exists
}
func (self JSONClient) MirrorLast(key []byte) (lastKey, lastValue []byte, existed bool) {
	item := KeyReq{
		Key: key,
	}
	var result common.Item
	self.call("MirrorLast", item, &result)
	return result.Key, result.Value, result.Exists
}
func (self JSONClient) MirrorFirst(key []byte) (firstKey, firstValue []byte, existed bool) {
	item := KeyReq{
		Key: key,
	}
	var result common.Item
	self.call("MirrorFirst", item, &result)
	return result.Key, result.Value, result.Exists
}
func (self JSONClient) Last(key []byte) (lastKey, lastValue []byte, existed bool) {
	item := KeyReq{
		Key: key,
	}
	var result common.Item
	self.call("Last", item, &result)
	return result.Key, result.Value, result.Exists
}
func (self JSONClient) First(key []byte) (firstKey, firstValue []byte, existed bool) {
	item := KeyReq{
		Key: key,
	}
	var result common.Item
	self.call("First", item, &result)
	return result.Key, result.Value, result.Exists
}
func (self JSONClient) SubGet(key, subKey []byte) (value []byte, existed bool) {
	item := SubKeyReq{
		Key:    key,
		SubKey: subKey,
	}
	var result common.Item
	self.call("SubGet", item, &result)
	return result.Value, result.Exists
}
func (self JSONClient) Get(key []byte) (value []byte, existed bool) {
	item := KeyReq{
		Key: key,
	}
	var result common.Item
	self.call("Get", item, &result)
	return result.Value, result.Exists
}
func (self JSONClient) SubSize(key []byte) (result int) {
	self.call("SubSize", KeyReq{Key: key}, &result)
	return result
}
func (self JSONClient) Size() (result int) {
	self.call("Size", nil, &result)
	return result
}
func (self JSONClient) SetExpression(expr setop.SetExpression) (result []setop.SetOpResult) {
	self.call("SetExpression", expr, &result)
	return
}
func (self JSONClient) AddConfiguration(key, value string) {
	conf := SubConf{
		Key:   key,
		Value: value,
	}
	var x Nothing
	self.call("AddConfiguration", conf, &x)
}
func (self JSONClient) SubAddConfiguration(treeKey []byte, key, value string) {
	conf := SubConf{
		TreeKey: treeKey,
		Key:     key,
		Value:   value,
	}
	var x Nothing
	self.call("SubAddConfiguration", conf, &x)
}
