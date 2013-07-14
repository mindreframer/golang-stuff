package stats

type HeapType interface {
	SetIndex(i, j int)
}

type Heap struct {
	HeapType
	h []interface{}
}

func (x *Heap) Len() int {
	return len(x.h)
}

func (x *Heap) Swap(i, j int) {
	x.h[i], x.h[j] = x.h[j], x.h[i]
	x.SetIndex(i, i)
	x.SetIndex(j, j)
}

func (x *Heap) Push(a interface{}) {
	x.h = append(x.h, a)
	n := len(x.h)
	x.SetIndex(n-1, n-1)
}

func (x *Heap) Pop() interface{} {
	n := len(x.h)
	x.SetIndex(n-1, -1)
	y := x.h[n-1]
	x.h = x.h[0 : n-1]
	return y
}

func (x *Heap) Copy() Heap {
	y := *x
	y.h = make([]interface{}, len(x.h))
	copy(y.h, x.h)
	return y
}
