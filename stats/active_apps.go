package stats

import (
	"container/heap"
	"sync"
	"time"
)

const (
	ActiveAppsTrimInterval  = 1 * time.Minute
	ActiveAppsEntryLifetime = 30 * time.Minute
)

type activeAppsEntry struct {
	t  int64 // Last update
	ti int   // Index in time min-heap
	tj int   // Index in time max-heap

	ApplicationId string
}

func (x *activeAppsEntry) Mark(t int64) {
	if x.t < t {
		x.t = t
	}
}

type byTimeMinHeap struct{ Heap }

func (x *byTimeMinHeap) Init() {
	x.Heap.HeapType = x
}

func (x *byTimeMinHeap) Less(i, j int) bool {
	yi := x.Heap.h[i].(*activeAppsEntry)
	yj := x.Heap.h[j].(*activeAppsEntry)
	return yi.t < yj.t
}

func (x *byTimeMinHeap) SetIndex(i, j int) {
	y := x.Heap.h[i].(*activeAppsEntry)
	y.ti = j
}

type byTimeMinHeapSnapshot struct{ byTimeMinHeap }

func (x *byTimeMinHeapSnapshot) Init() {
	x.Heap.HeapType = x
}

func (x *byTimeMinHeapSnapshot) SetIndex(i, j int) {
	// No-op
}

type byTimeMaxHeap struct{ Heap }

func (x *byTimeMaxHeap) Init() {
	x.Heap.HeapType = x
}

func (x *byTimeMaxHeap) Less(i, j int) bool {
	yi := x.Heap.h[i].(*activeAppsEntry)
	yj := x.Heap.h[j].(*activeAppsEntry)
	return yi.t > yj.t
}

func (x *byTimeMaxHeap) SetIndex(i, j int) {
	y := x.Heap.h[i].(*activeAppsEntry)
	y.tj = j
}

type byTimeMaxHeapSnapshot struct{ byTimeMaxHeap }

func (x *byTimeMaxHeapSnapshot) Init() {
	x.Heap.HeapType = x
}

func (x *byTimeMaxHeapSnapshot) SetIndex(i, j int) {
	// No-op
}

type ActiveApps struct {
	sync.Mutex

	t *time.Ticker

	m map[string]*activeAppsEntry
	i byTimeMinHeap
	j byTimeMaxHeap
}

func NewActiveApps() *ActiveApps {
	x := &ActiveApps{}

	x.t = time.NewTicker(1 * time.Minute)

	x.m = make(map[string]*activeAppsEntry)
	x.i.Init()
	x.j.Init()

	go func() {
		for {
			select {
			case <-x.t.C:
				x.Trim(time.Now().Add(-ActiveAppsEntryLifetime))
			}
		}
	}()

	return x
}

func (x *ActiveApps) Mark(ApplicationId string, z time.Time) {
	t := z.Unix()

	x.Lock()
	defer x.Unlock()

	y := x.m[ApplicationId]
	if y != nil {
		heap.Remove(&x.i, y.ti)
		heap.Remove(&x.j, y.tj)
	} else {
		// New entry
		y = &activeAppsEntry{ApplicationId: ApplicationId}
		x.m[ApplicationId] = y
	}

	y.Mark(t)

	heap.Push(&x.i, y)
	heap.Push(&x.j, y)
}

func (x *ActiveApps) Trim(y time.Time) {
	t := y.Unix()

	x.Lock()
	defer x.Unlock()

	for x.i.Len() > 0 {
		// Pop from the min-heap
		z := heap.Pop(&x.i).(*activeAppsEntry)
		if z.t > t {
			// Push back to the min-heap
			heap.Push(&x.i, z)
			break
		}

		// Remove from max-heap
		heap.Remove(&x.j, z.tj)

		// Remove from map
		delete(x.m, z.ApplicationId)
	}
}

func (x *ActiveApps) ActiveSince(y time.Time) []string {
	t := y.Unix()

	x.Lock()

	a := byTimeMaxHeapSnapshot{}
	a.Heap = x.j.Copy()
	a.Init()

	x.Unlock()

	// Collect active applications
	b := make([]string, 0)
	for a.Len() > 0 {
		z := heap.Pop(&a).(*activeAppsEntry)
		if z.t < t {
			break
		}

		// Add active application
		b = append(b, z.ApplicationId)
	}

	return b
}
