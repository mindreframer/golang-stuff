package dhash

import (
	"github.com/zond/god/common"
	"github.com/zond/god/setop"
)

type dhashServer Node

func (self *dhashServer) Clear(x int, y *int) error {
	(*Node)(self).Clear()
	return nil
}
func (self *dhashServer) SlaveSubPut(data common.Item, x *int) error {
	return (*Node)(self).subPut(data)
}
func (self *dhashServer) SlaveSubClear(data common.Item, x *int) error {
	return (*Node)(self).subClear(data)
}
func (self *dhashServer) SlaveSubDel(data common.Item, x *int) error {
	return (*Node)(self).subDel(data)
}
func (self *dhashServer) SlaveDel(data common.Item, x *int) error {
	return (*Node)(self).del(data)
}
func (self *dhashServer) SlavePut(data common.Item, x *int) error {
	return (*Node)(self).put(data)
}
func (self *dhashServer) SubDel(data common.Item, x *int) error {
	return (*Node)(self).SubDel(data)
}
func (self *dhashServer) SubClear(data common.Item, x *int) error {
	return (*Node)(self).SubClear(data)
}
func (self *dhashServer) SubPut(data common.Item, x *int) error {
	return (*Node)(self).SubPut(data)
}
func (self *dhashServer) Del(data common.Item, x *int) error {
	return (*Node)(self).Del(data)
}
func (self *dhashServer) Put(data common.Item, x *int) error {
	return (*Node)(self).Put(data)
}
func (self *dhashServer) RingHash(x int, result *[]byte) error {
	return (*Node)(self).RingHash(x, result)
}
func (self *dhashServer) MirrorCount(r common.Range, result *int) error {
	return (*Node)(self).MirrorCount(r, result)
}
func (self *dhashServer) Count(r common.Range, result *int) error {
	return (*Node)(self).Count(r, result)
}
func (self *dhashServer) Next(data common.Item, result *common.Item) error {
	return (*Node)(self).Next(data, result)
}
func (self *dhashServer) Prev(data common.Item, result *common.Item) error {
	return (*Node)(self).Prev(data, result)
}
func (self *dhashServer) SubGet(data common.Item, result *common.Item) error {
	return (*Node)(self).SubGet(data, result)
}
func (self *dhashServer) Get(data common.Item, result *common.Item) error {
	return (*Node)(self).Get(data, result)
}
func (self *dhashServer) Size(x int, result *int) error {
	*result = (*Node)(self).Size()
	return nil
}
func (self *dhashServer) SubSize(key []byte, result *int) error {
	return (*Node)(self).SubSize(key, result)
}
func (self *dhashServer) Owned(x int, result *int) error {
	*result = (*Node)(self).Owned()
	return nil
}
func (self *dhashServer) Describe(x int, result *common.DHashDescription) error {
	*result = (*Node)(self).Description()
	return nil
}
func (self *dhashServer) DescribeTree(x int, result *string) error {
	*result = (*Node)(self).DescribeTree()
	return nil
}
func (self *dhashServer) MirrorPrevIndex(data common.Item, result *common.Item) error {
	return (*Node)(self).MirrorPrevIndex(data, result)
}
func (self *dhashServer) MirrorNextIndex(data common.Item, result *common.Item) error {
	return (*Node)(self).MirrorNextIndex(data, result)
}
func (self *dhashServer) PrevIndex(data common.Item, result *common.Item) error {
	return (*Node)(self).PrevIndex(data, result)
}
func (self *dhashServer) NextIndex(data common.Item, result *common.Item) error {
	return (*Node)(self).NextIndex(data, result)
}
func (self *dhashServer) MirrorReverseIndexOf(data common.Item, result *common.Index) error {
	return (*Node)(self).MirrorReverseIndexOf(data, result)
}
func (self *dhashServer) MirrorIndexOf(data common.Item, result *common.Index) error {
	return (*Node)(self).MirrorIndexOf(data, result)
}
func (self *dhashServer) ReverseIndexOf(data common.Item, result *common.Index) error {
	return (*Node)(self).ReverseIndexOf(data, result)
}
func (self *dhashServer) IndexOf(data common.Item, result *common.Index) error {
	return (*Node)(self).IndexOf(data, result)
}
func (self *dhashServer) SubMirrorPrev(data common.Item, result *common.Item) error {
	return (*Node)(self).SubMirrorPrev(data, result)
}
func (self *dhashServer) SubMirrorNext(data common.Item, result *common.Item) error {
	return (*Node)(self).SubMirrorNext(data, result)
}
func (self *dhashServer) SubPrev(data common.Item, result *common.Item) error {
	return (*Node)(self).SubPrev(data, result)
}
func (self *dhashServer) SubNext(data common.Item, result *common.Item) error {
	return (*Node)(self).SubNext(data, result)
}
func (self *dhashServer) MirrorFirst(data common.Item, result *common.Item) error {
	return (*Node)(self).MirrorFirst(data, result)
}
func (self *dhashServer) MirrorLast(data common.Item, result *common.Item) error {
	return (*Node)(self).MirrorLast(data, result)
}
func (self *dhashServer) First(data common.Item, result *common.Item) error {
	return (*Node)(self).First(data, result)
}
func (self *dhashServer) Last(data common.Item, result *common.Item) error {
	return (*Node)(self).Last(data, result)
}
func (self *dhashServer) MirrorReverseSlice(r common.Range, result *[]common.Item) error {
	return (*Node)(self).MirrorReverseSlice(r, result)
}
func (self *dhashServer) MirrorSlice(r common.Range, result *[]common.Item) error {
	return (*Node)(self).MirrorSlice(r, result)
}
func (self *dhashServer) MirrorSliceIndex(r common.Range, result *[]common.Item) error {
	return (*Node)(self).MirrorSliceIndex(r, result)
}
func (self *dhashServer) MirrorReverseSliceIndex(r common.Range, result *[]common.Item) error {
	return (*Node)(self).MirrorReverseSliceIndex(r, result)
}
func (self *dhashServer) MirrorSliceLen(r common.Range, result *[]common.Item) error {
	return (*Node)(self).MirrorSliceLen(r, result)
}
func (self *dhashServer) MirrorReverseSliceLen(r common.Range, result *[]common.Item) error {
	return (*Node)(self).MirrorReverseSliceLen(r, result)
}
func (self *dhashServer) ReverseSlice(r common.Range, result *[]common.Item) error {
	return (*Node)(self).ReverseSlice(r, result)
}
func (self *dhashServer) Slice(r common.Range, result *[]common.Item) error {
	return (*Node)(self).Slice(r, result)
}
func (self *dhashServer) SliceIndex(r common.Range, result *[]common.Item) error {
	return (*Node)(self).SliceIndex(r, result)
}
func (self *dhashServer) ReverseSliceIndex(r common.Range, result *[]common.Item) error {
	return (*Node)(self).ReverseSliceIndex(r, result)
}
func (self *dhashServer) SliceLen(r common.Range, result *[]common.Item) error {
	return (*Node)(self).SliceLen(r, result)
}
func (self *dhashServer) ReverseSliceLen(r common.Range, result *[]common.Item) error {
	return (*Node)(self).ReverseSliceLen(r, result)
}
func (self *dhashServer) SetExpression(expr setop.SetExpression, items *[]setop.SetOpResult) error {
	return (*Node)(self).SetExpression(expr, items)
}

func (self *dhashServer) AddConfiguration(c common.ConfItem, x *int) error {
	(*Node)(self).AddConfiguration(c)
	return nil
}
func (self *dhashServer) SlaveSubAddConfiguration(c common.ConfItem, x *int) error {
	(*Node)(self).subAddConfiguration(c)
	return nil
}
func (self *dhashServer) SubAddConfiguration(c common.ConfItem, x *int) error {
	(*Node)(self).SubAddConfiguration(c)
	return nil
}
func (self *dhashServer) Configuration(x int, result *common.Conf) error {
	*result = common.Conf{}
	(*result).Data, (*result).Timestamp = (*Node)(self).tree.Configuration()
	return nil
}
func (self *dhashServer) SubConfiguration(key []byte, result *common.Conf) error {
	*result = common.Conf{TreeKey: key}
	(*result).Data, (*result).Timestamp = (*Node)(self).tree.SubConfiguration(key)
	return nil
}
