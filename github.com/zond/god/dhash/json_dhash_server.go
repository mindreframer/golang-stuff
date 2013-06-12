package dhash

import (
	"github.com/zond/god/common"
	"github.com/zond/god/setop"
)

type Nothing struct{}
type SubValueRes struct {
	Key    []byte
	SubKey []byte
	Value  []byte
	Exists bool
}
type SubValueIndexRes struct {
	Key    []byte
	SubKey []byte
	Value  []byte
	Index  int
	Exists bool
}
type SubValueOp struct {
	Key    []byte
	SubKey []byte
	Value  []byte
	Sync   bool
}
type SubKeyOp struct {
	Key    []byte
	SubKey []byte
	Sync   bool
}
type SubKeyReq struct {
	Key    []byte
	SubKey []byte
}
type SubIndex struct {
	Key   []byte
	Index int
}
type ValueOp struct {
	Key   []byte
	Value []byte
	Sync  bool
}
type ValueRes struct {
	Key    []byte
	Value  []byte
	Exists bool
}
type KeyOp struct {
	Key  []byte
	Sync bool
}
type KeyReq struct {
	Key []byte
}
type KeyRange struct {
	Key    []byte
	Min    []byte
	Max    []byte
	MinInc bool
	MaxInc bool
}
type IndexRange struct {
	Key      []byte
	MinIndex *int
	MaxIndex *int
}
type PageRange struct {
	Key     []byte
	From    []byte
	FromInc bool
	Len     int
}
type SubConf struct {
	TreeKey []byte
	Key     string
	Value   string
}
type Conf struct {
	Key   string
	Value string
}

type JSONApi Node

func (self *JSONApi) convert(items []common.Item, result *[]ValueRes) {
	for _, item := range items {
		*result = append(*result, ValueRes{
			Key:    item.Key,
			Value:  item.Value,
			Exists: true,
		})
	}
}
func (self *JSONApi) forwardUnlessMe(cmd string, key []byte, in, out interface{}) (forwarded bool, err error) {
	succ := (*Node)(self).node.GetSuccessorFor(key)
	if succ.Addr != (*Node)(self).node.GetBroadcastAddr() {
		forwarded, err = true, succ.Call(cmd, in, out)
	}
	return
}

func (self *JSONApi) Clear(x Nothing, y *Nothing) (err error) {
	(*Node)(self).Clear()
	return nil
}
func (self *JSONApi) Nodes(x Nothing, result *common.Remotes) (err error) {
	*result = (*Node)(self).node.GetNodes()
	return nil
}
func (self *JSONApi) SubDel(d SubKeyOp, n *Nothing) (err error) {
	data := common.Item{
		Key:    d.Key,
		SubKey: d.SubKey,
		Sync:   d.Sync,
	}
	var x int
	var f bool
	if f, err = self.forwardUnlessMe("DHash.SubDel", data.Key, data, &x); !f {
		err = (*Node)(self).SubDel(data)
	}
	return
}
func (self *JSONApi) SubClear(d SubKeyOp, n *Nothing) (err error) {
	data := common.Item{
		Key:    d.Key,
		SubKey: d.SubKey,
		Sync:   d.Sync,
	}
	var x int
	var f bool
	if f, err = self.forwardUnlessMe("DHash.SubClear", data.Key, data, &x); !f {
		err = (*Node)(self).SubClear(data)
	}
	return
}
func (self *JSONApi) SubPut(d SubValueOp, n *Nothing) (err error) {
	data := common.Item{
		Key:    d.Key,
		SubKey: d.SubKey,
		Value:  d.Value,
		Sync:   d.Sync,
	}
	var x int
	var f bool
	if f, err = self.forwardUnlessMe("DHash.SubPut", data.Key, data, &x); !f {
		err = (*Node)(self).SubPut(data)
	}
	return
}
func (self *JSONApi) Del(d KeyOp, n *Nothing) (err error) {
	data := common.Item{
		Key:  d.Key,
		Sync: d.Sync,
	}
	var x int
	var f bool
	if f, err = self.forwardUnlessMe("DHash.Del", data.Key, data, &x); !f {
		err = (*Node)(self).Del(data)
	}
	return
}
func (self *JSONApi) Put(d ValueOp, n *Nothing) (err error) {
	data := common.Item{
		Key:   d.Key,
		Value: d.Value,
		Sync:  d.Sync,
	}
	var x int
	var f bool
	if f, err = self.forwardUnlessMe("DHash.Put", data.Key, data, &x); !f {
		err = (*Node)(self).Put(data)
	}
	return
}
func (self *JSONApi) MirrorCount(kr KeyRange, result *int) (err error) {
	r := common.Range{
		Key:    kr.Key,
		Min:    kr.Min,
		Max:    kr.Max,
		MinInc: kr.MinInc,
		MaxInc: kr.MaxInc,
	}
	var f bool
	if f, err = self.forwardUnlessMe("DHash.MirrorCount", r.Key, r, result); !f {
		err = (*Node)(self).MirrorCount(r, result)
	}
	return
}
func (self *JSONApi) Count(kr KeyRange, result *int) (err error) {
	r := common.Range{
		Key:    kr.Key,
		Min:    kr.Min,
		Max:    kr.Max,
		MinInc: kr.MinInc,
		MaxInc: kr.MaxInc,
	}
	var f bool
	if f, err = self.forwardUnlessMe("DHash.Count", r.Key, r, result); !f {
		err = (*Node)(self).Count(r, result)
	}
	return
}
func (self *JSONApi) Next(kr KeyReq, result *ValueRes) (err error) {
	k, v, e := (*Node)(self).client().Next(kr.Key)
	*result = ValueRes{
		Key:    k,
		Value:  v,
		Exists: e,
	}
	return nil
}
func (self *JSONApi) Prev(kr KeyReq, result *ValueRes) (err error) {
	k, v, e := (*Node)(self).client().Prev(kr.Key)
	*result = ValueRes{
		Key:    k,
		Value:  v,
		Exists: e,
	}
	return nil
}
func (self *JSONApi) SubGet(k SubKeyReq, result *SubValueRes) (err error) {
	data := common.Item{
		Key:    k.Key,
		SubKey: k.SubKey,
	}
	var item common.Item
	var f bool
	if f, err = self.forwardUnlessMe("DHash.SubGet", data.Key, data, &item); !f {
		err = (*Node)(self).SubGet(data, &item)
	}
	*result = SubValueRes{
		Key:    item.Key,
		SubKey: item.SubKey,
		Value:  item.Value,
		Exists: item.Exists,
	}
	return
}
func (self *JSONApi) Get(k KeyReq, result *ValueRes) (err error) {
	data := common.Item{
		Key: k.Key,
	}
	var item common.Item
	var f bool
	if f, err = self.forwardUnlessMe("DHash.Get", data.Key, data, &item); !f {
		err = (*Node)(self).Get(data, &item)
	}
	*result = ValueRes{
		Key:    item.Key,
		Value:  item.Value,
		Exists: item.Exists,
	}
	return
}
func (self *JSONApi) Size(x Nothing, result *int) (err error) {
	*result = (*Node)(self).Size()
	return nil
}
func (self *JSONApi) SubSize(k KeyReq, result *int) (err error) {
	key := k.Key
	var f bool
	if f, err = self.forwardUnlessMe("DHash.SubSize", key, key, result); !f {
		err = (*Node)(self).SubSize(key, result)
	}
	return
}
func (self *JSONApi) Owned(x Nothing, result *int) (err error) {
	*result = (*Node)(self).Owned()
	return nil
}
func (self *JSONApi) Describe(x Nothing, result *common.DHashDescription) (err error) {
	*result = (*Node)(self).Description()
	return nil
}
func (self *JSONApi) DescribeTree(x Nothing, result *string) (err error) {
	*result = (*Node)(self).DescribeTree()
	return nil
}
func (self *JSONApi) PrevIndex(i SubIndex, result *SubValueIndexRes) (err error) {
	data := common.Item{
		Key:   i.Key,
		Index: i.Index,
	}
	var item common.Item
	var f bool
	if f, err = self.forwardUnlessMe("DHash.PrevIndex", data.Key, data, &item); !f {
		err = (*Node)(self).PrevIndex(data, &item)
	}
	*result = SubValueIndexRes{
		Key:    item.Key,
		SubKey: item.SubKey,
		Value:  item.Value,
		Index:  item.Index,
		Exists: item.Exists,
	}
	return
}
func (self *JSONApi) MirrorPrevIndex(i SubIndex, result *SubValueIndexRes) (err error) {
	data := common.Item{
		Key:   i.Key,
		Index: i.Index,
	}
	var item common.Item
	var f bool
	if f, err = self.forwardUnlessMe("DHash.MirrorPrevIndex", data.Key, data, &item); !f {
		err = (*Node)(self).MirrorPrevIndex(data, &item)
	}
	*result = SubValueIndexRes{
		Key:    item.Key,
		SubKey: item.SubKey,
		Value:  item.Value,
		Index:  item.Index,
		Exists: item.Exists,
	}
	return
}
func (self *JSONApi) MirrorNextIndex(i SubIndex, result *SubValueIndexRes) (err error) {
	data := common.Item{
		Key:   i.Key,
		Index: i.Index,
	}
	var item common.Item
	var f bool
	if f, err = self.forwardUnlessMe("DHash.MirrorNextIndex", data.Key, data, &item); !f {
		err = (*Node)(self).MirrorNextIndex(data, &item)
	}
	*result = SubValueIndexRes{
		Key:    item.Key,
		SubKey: item.SubKey,
		Value:  item.Value,
		Index:  item.Index,
		Exists: item.Exists,
	}
	return
}
func (self *JSONApi) NextIndex(i SubIndex, result *SubValueIndexRes) (err error) {
	data := common.Item{
		Key:   i.Key,
		Index: i.Index,
	}
	var item common.Item
	var f bool
	if f, err = self.forwardUnlessMe("DHash.NextIndex", data.Key, data, &item); !f {
		err = (*Node)(self).NextIndex(data, &item)
	}
	*result = SubValueIndexRes{
		Key:    item.Key,
		SubKey: item.SubKey,
		Value:  item.Value,
		Index:  item.Index,
		Exists: item.Exists,
	}
	return
}
func (self *JSONApi) MirrorReverseIndexOf(i SubKeyReq, result *common.Index) (err error) {
	data := common.Item{
		Key:    i.Key,
		SubKey: i.SubKey,
	}
	var f bool
	if f, err = self.forwardUnlessMe("DHash.MirrorReverseIndexOf", data.Key, data, result); !f {
		err = (*Node)(self).MirrorReverseIndexOf(data, result)
	}
	return
}
func (self *JSONApi) MirrorIndexOf(i SubKeyReq, result *common.Index) (err error) {
	data := common.Item{
		Key:    i.Key,
		SubKey: i.SubKey,
	}
	var f bool
	if f, err = self.forwardUnlessMe("DHash.MirrorIndexOf", data.Key, data, result); !f {
		err = (*Node)(self).MirrorIndexOf(data, result)
	}
	return
}
func (self *JSONApi) ReverseIndexOf(i SubKeyReq, result *common.Index) (err error) {
	data := common.Item{
		Key:    i.Key,
		SubKey: i.SubKey,
	}
	var f bool
	if f, err = self.forwardUnlessMe("DHash.ReverseIndexOf", data.Key, data, result); !f {
		err = (*Node)(self).ReverseIndexOf(data, result)
	}
	return
}
func (self *JSONApi) IndexOf(i SubKeyReq, result *common.Index) (err error) {
	data := common.Item{
		Key:    i.Key,
		SubKey: i.SubKey,
	}
	var f bool
	if f, err = self.forwardUnlessMe("DHash.IndexOf", data.Key, data, result); !f {
		err = (*Node)(self).IndexOf(data, result)
	}
	return
}
func (self *JSONApi) SubMirrorPrev(k SubKeyReq, result *SubValueRes) (err error) {
	data := common.Item{
		Key:    k.Key,
		SubKey: k.SubKey,
	}
	var item common.Item
	var f bool
	if f, err = self.forwardUnlessMe("DHash.SubMirrorPrev", data.Key, data, &item); !f {
		err = (*Node)(self).SubMirrorPrev(data, &item)
	}
	*result = SubValueRes{
		Key:    item.Key,
		SubKey: item.SubKey,
		Value:  item.Value,
		Exists: item.Exists,
	}
	return
}
func (self *JSONApi) SubMirrorNext(k SubKeyReq, result *SubValueRes) (err error) {
	data := common.Item{
		Key:    k.Key,
		SubKey: k.SubKey,
	}
	var item common.Item
	var f bool
	if f, err = self.forwardUnlessMe("DHash.SubMirrorNext", data.Key, data, &item); !f {
		err = (*Node)(self).SubMirrorNext(data, &item)
	}
	*result = SubValueRes{
		Key:    item.Key,
		SubKey: item.SubKey,
		Value:  item.Value,
		Exists: item.Exists,
	}
	return
}
func (self *JSONApi) SubPrev(k SubKeyReq, result *SubValueRes) (err error) {
	data := common.Item{
		Key:    k.Key,
		SubKey: k.SubKey,
	}
	var item common.Item
	var f bool
	if f, err = self.forwardUnlessMe("DHash.SubPrev", data.Key, data, &item); !f {
		err = (*Node)(self).SubPrev(data, &item)
	}
	*result = SubValueRes{
		Key:    item.Key,
		SubKey: item.SubKey,
		Value:  item.Value,
		Exists: item.Exists,
	}
	return
}
func (self *JSONApi) SubNext(k SubKeyReq, result *SubValueRes) (err error) {
	data := common.Item{
		Key:    k.Key,
		SubKey: k.SubKey,
	}
	var item common.Item
	var f bool
	if f, err = self.forwardUnlessMe("DHash.SubNext", data.Key, data, &item); !f {
		err = (*Node)(self).SubNext(data, &item)
	}
	*result = SubValueRes{
		Key:    item.Key,
		SubKey: item.SubKey,
		Value:  item.Value,
		Exists: item.Exists,
	}
	return
}
func (self *JSONApi) MirrorFirst(k KeyReq, result *SubValueRes) (err error) {
	data := common.Item{
		Key: k.Key,
	}
	var item common.Item
	var f bool
	if f, err = self.forwardUnlessMe("DHash.MirrorFirst", data.Key, data, &item); !f {
		err = (*Node)(self).MirrorFirst(data, &item)
	}
	*result = SubValueRes{
		Key:    item.Key,
		SubKey: item.SubKey,
		Value:  item.Value,
		Exists: item.Exists,
	}
	return
}
func (self *JSONApi) MirrorLast(k KeyReq, result *SubValueRes) (err error) {
	data := common.Item{
		Key: k.Key,
	}
	var item common.Item
	var f bool
	if f, err = self.forwardUnlessMe("DHash.MirrorLast", data.Key, data, &item); !f {
		err = (*Node)(self).MirrorLast(data, &item)
	}
	*result = SubValueRes{
		Key:    item.Key,
		SubKey: item.SubKey,
		Value:  item.Value,
		Exists: item.Exists,
	}
	return
}
func (self *JSONApi) First(k KeyReq, result *SubValueRes) (err error) {
	data := common.Item{
		Key: k.Key,
	}
	var item common.Item
	var f bool
	if f, err = self.forwardUnlessMe("DHash.First", data.Key, data, &item); !f {
		err = (*Node)(self).First(data, &item)
	}
	*result = SubValueRes{
		Key:    item.Key,
		SubKey: item.SubKey,
		Value:  item.Value,
		Exists: item.Exists,
	}
	return
}
func (self *JSONApi) Last(k KeyReq, result *SubValueRes) (err error) {
	data := common.Item{
		Key: k.Key,
	}
	var item common.Item
	var f bool
	if f, err = self.forwardUnlessMe("DHash.Last", data.Key, data, &item); !f {
		err = (*Node)(self).Last(data, &item)
	}
	*result = SubValueRes{
		Key:    item.Key,
		SubKey: item.SubKey,
		Value:  item.Value,
		Exists: item.Exists,
	}
	return
}
func (self *JSONApi) MirrorReverseSlice(kr KeyRange, result *[]ValueRes) (err error) {
	r := common.Range{
		Key:    kr.Key,
		Min:    kr.Min,
		Max:    kr.Max,
		MinInc: kr.MinInc,
		MaxInc: kr.MaxInc,
	}
	var items []common.Item
	var f bool
	if f, err = self.forwardUnlessMe("DHash.MirrorReverseSlice", r.Key, r, &items); !f {
		err = (*Node)(self).MirrorReverseSlice(r, &items)
	}
	self.convert(items, result)
	return
}
func (self *JSONApi) MirrorSlice(kr KeyRange, result *[]ValueRes) (err error) {
	r := common.Range{
		Key:    kr.Key,
		Min:    kr.Min,
		Max:    kr.Max,
		MinInc: kr.MinInc,
		MaxInc: kr.MaxInc,
	}
	var items []common.Item
	var f bool
	if f, err = self.forwardUnlessMe("DHash.MirrorSlice", r.Key, r, &items); !f {
		err = (*Node)(self).MirrorSlice(r, &items)
	}
	self.convert(items, result)
	return
}
func (self *JSONApi) MirrorSliceIndex(ir IndexRange, result *[]ValueRes) (err error) {
	var mi int
	var ma int
	if ir.MinIndex != nil {
		mi = *ir.MinIndex
	}
	if ir.MaxIndex != nil {
		ma = *ir.MaxIndex
	}
	r := common.Range{
		Key:      ir.Key,
		MinIndex: mi,
		MaxIndex: ma,
		MinInc:   ir.MinIndex != nil,
		MaxInc:   ir.MaxIndex != nil,
	}
	var items []common.Item
	var f bool
	if f, err = self.forwardUnlessMe("DHash.MirrorSliceIndex", r.Key, r, &items); !f {
		err = (*Node)(self).MirrorSliceIndex(r, &items)
	}
	self.convert(items, result)
	return
}
func (self *JSONApi) MirrorReverseSliceIndex(ir IndexRange, result *[]ValueRes) (err error) {
	var mi int
	var ma int
	if ir.MinIndex != nil {
		mi = *ir.MinIndex
	}
	if ir.MaxIndex != nil {
		ma = *ir.MaxIndex
	}
	r := common.Range{
		Key:      ir.Key,
		MinIndex: mi,
		MaxIndex: ma,
		MinInc:   ir.MinIndex != nil,
		MaxInc:   ir.MaxIndex != nil,
	}
	var items []common.Item
	var f bool
	if f, err = self.forwardUnlessMe("DHash.MirrorReverseSliceIndex", r.Key, r, &items); !f {
		err = (*Node)(self).MirrorReverseSliceIndex(r, &items)
	}
	self.convert(items, result)
	return
}
func (self *JSONApi) MirrorSliceLen(pr PageRange, result *[]ValueRes) (err error) {
	r := common.Range{
		Key:    pr.Key,
		Min:    pr.From,
		MinInc: pr.FromInc,
		Len:    pr.Len,
	}
	var items []common.Item
	var f bool
	if f, err = self.forwardUnlessMe("DHash.MirrorSliceLen", r.Key, r, &items); !f {
		err = (*Node)(self).MirrorSliceLen(r, &items)
	}
	self.convert(items, result)
	return
}
func (self *JSONApi) MirrorReverseSliceLen(pr PageRange, result *[]ValueRes) (err error) {
	r := common.Range{
		Key:    pr.Key,
		Max:    pr.From,
		MaxInc: pr.FromInc,
		Len:    pr.Len,
	}
	var items []common.Item
	var f bool
	if f, err = self.forwardUnlessMe("DHash.MirrorReverseSliceLen", r.Key, r, &items); !f {
		err = (*Node)(self).MirrorReverseSliceLen(r, &items)
	}
	self.convert(items, result)
	return
}
func (self *JSONApi) ReverseSlice(kr KeyRange, result *[]ValueRes) (err error) {
	r := common.Range{
		Key:    kr.Key,
		Min:    kr.Min,
		Max:    kr.Max,
		MinInc: kr.MinInc,
		MaxInc: kr.MaxInc,
	}
	var items []common.Item
	var f bool
	if f, err = self.forwardUnlessMe("DHash.ReverseSlice", r.Key, r, &items); !f {
		err = (*Node)(self).ReverseSlice(r, &items)
	}
	self.convert(items, result)
	return
}
func (self *JSONApi) Slice(kr KeyRange, result *[]ValueRes) (err error) {
	r := common.Range{
		Key:    kr.Key,
		Min:    kr.Min,
		Max:    kr.Max,
		MinInc: kr.MinInc,
		MaxInc: kr.MaxInc,
	}
	var items []common.Item
	var f bool
	if f, err = self.forwardUnlessMe("DHash.Slice", r.Key, r, &items); !f {
		err = (*Node)(self).Slice(r, &items)
	}
	self.convert(items, result)
	return
}
func (self *JSONApi) SliceIndex(ir IndexRange, result *[]ValueRes) (err error) {
	var mi int
	var ma int
	if ir.MinIndex != nil {
		mi = *ir.MinIndex
	}
	if ir.MaxIndex != nil {
		ma = *ir.MaxIndex
	}
	r := common.Range{
		Key:      ir.Key,
		MinIndex: mi,
		MaxIndex: ma,
		MinInc:   ir.MinIndex != nil,
		MaxInc:   ir.MaxIndex != nil,
	}
	var items []common.Item
	var f bool
	if f, err = self.forwardUnlessMe("DHash.SliceIndex", r.Key, r, &items); !f {
		err = (*Node)(self).SliceIndex(r, &items)
	}
	self.convert(items, result)
	return
}
func (self *JSONApi) ReverseSliceIndex(ir IndexRange, result *[]ValueRes) (err error) {
	var mi int
	var ma int
	if ir.MinIndex != nil {
		mi = *ir.MinIndex
	}
	if ir.MaxIndex != nil {
		ma = *ir.MaxIndex
	}
	r := common.Range{
		Key:      ir.Key,
		MinIndex: mi,
		MaxIndex: ma,
		MinInc:   ir.MinIndex != nil,
		MaxInc:   ir.MaxIndex != nil,
	}
	var items []common.Item
	var f bool
	if f, err = self.forwardUnlessMe("DHash.ReverseSliceIndex", r.Key, r, &items); !f {
		err = (*Node)(self).ReverseSliceIndex(r, &items)
	}
	self.convert(items, result)
	return
}
func (self *JSONApi) SliceLen(pr PageRange, result *[]ValueRes) (err error) {
	r := common.Range{
		Key:    pr.Key,
		Min:    pr.From,
		MinInc: pr.FromInc,
		Len:    pr.Len,
	}
	var items []common.Item
	var f bool
	if f, err = self.forwardUnlessMe("DHash.SliceLen", r.Key, r, &items); !f {
		err = (*Node)(self).SliceLen(r, &items)
	}
	self.convert(items, result)
	return
}
func (self *JSONApi) ReverseSliceLen(pr PageRange, result *[]ValueRes) (err error) {
	r := common.Range{
		Key:    pr.Key,
		Max:    pr.From,
		MaxInc: pr.FromInc,
		Len:    pr.Len,
	}
	var items []common.Item
	var f bool
	if f, err = self.forwardUnlessMe("DHash.ReverseSliceLen", r.Key, r, &items); !f {
		err = (*Node)(self).ReverseSliceLen(r, &items)
	}
	self.convert(items, result)
	return
}
func (self *JSONApi) SetExpression(expr setop.SetExpression, items *[]setop.SetOpResult) (err error) {
	if expr.Op == nil {
		var err error
		if expr.Op, err = setop.NewSetOpParser(expr.Code).Parse(); err != nil {
			return err
		}
	}
	return (*Node)(self).SetExpression(expr, items)
}

func (self *JSONApi) AddConfiguration(co Conf, x *Nothing) (err error) {
	c := common.ConfItem{
		Key:   co.Key,
		Value: co.Value,
	}
	(*Node)(self).AddConfiguration(c)
	return nil
}
func (self *JSONApi) SubAddConfiguration(co SubConf, x *Nothing) (err error) {
	c := common.ConfItem{
		TreeKey: co.TreeKey,
		Key:     co.Key,
		Value:   co.Value,
	}
	(*Node)(self).SubAddConfiguration(c)
	return nil
}
func (self *JSONApi) Configuration(x Nothing, result *common.Conf) (err error) {
	*result = common.Conf{}
	(*result).Data, (*result).Timestamp = (*Node)(self).tree.Configuration()
	return nil
}
func (self *JSONApi) SubConfiguration(k KeyReq, result *common.Conf) (err error) {
	key := k.Key
	*result = common.Conf{TreeKey: key}
	(*result).Data, (*result).Timestamp = (*Node)(self).tree.SubConfiguration(key)
	return nil
}
