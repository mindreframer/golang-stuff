package stats

import (
	"container/heap"
	"sync"
	"time"
)

const (
	TopAppsTrimInterval  = 10 * time.Second
	TopAppsEntryLifetime = 60 * time.Second
)

type topAppsEntryTimeSlot struct {
	t int64 // Unix time
	n int64 // Number of requests
}

type topAppsEntry struct {
	ti int // Index in time heap
	ni int // Index in top heap

	t []topAppsEntryTimeSlot
	n int64

	ApplicationId string
}

func (x *topAppsEntry) Mark(t int64) {
	// Add time slot if necessary
	n := len(x.t)
	if n == 0 || x.t[n-1].t < t {
		x.t = append(x.t, topAppsEntryTimeSlot{t: t, n: 0})
		n = len(x.t)
	}

	x.t[n-1].n++
	x.n++
}

// Trim time slots up to and including time t.
func (x *topAppsEntry) Trim(t int64) {
	var i int
	var n int64

	// Collect slots that can be removed
	l := len(x.t)
	for i = 0; i < l; i++ {
		if x.t[i].t > t {
			break
		}

		n += x.t[i].n
	}

	copy(x.t, x.t[i:])
	x.t = x.t[0:(l - i)]
	x.n -= n
}

type byFirstMarkTimeHeap struct{ Heap }

func (x *byFirstMarkTimeHeap) Init() {
	x.Heap.HeapType = x
}

func (x *byFirstMarkTimeHeap) Less(i, j int) bool {
	yi := x.Heap.h[i].(*topAppsEntry)
	yj := x.Heap.h[j].(*topAppsEntry)

	// This asserts the slice of time slots is non-empty
	return yi.t[0].t < yj.t[0].t
}

func (x *byFirstMarkTimeHeap) SetIndex(i, j int) {
	y := x.Heap.h[i].(*topAppsEntry)
	y.ti = j
}

type byFirstMarkTimeHeapSnapshot struct{ byFirstMarkTimeHeap }

func (x *byFirstMarkTimeHeapSnapshot) Init() {
	x.Heap.HeapType = x
}

func (x *byFirstMarkTimeHeapSnapshot) SetIndex(i, j int) {
	// No-op
}

type byRequestsHeap struct{ Heap }

func (x *byRequestsHeap) Init() {
	x.Heap.HeapType = x
}

func (x *byRequestsHeap) Less(i, j int) bool {
	yi := x.Heap.h[i].(*topAppsEntry)
	yj := x.Heap.h[j].(*topAppsEntry)
	return yi.n > yj.n
}

func (x *byRequestsHeap) SetIndex(i, j int) {
	y := x.Heap.h[i].(*topAppsEntry)
	y.ni = j
}

type byRequestsHeapSnapshot struct{ byRequestsHeap }

func (x *byRequestsHeapSnapshot) Init() {
	x.Heap.HeapType = x
}

func (x *byRequestsHeapSnapshot) SetIndex(i, j int) {
	// No-op
}

type TopApps struct {
	sync.Mutex

	*time.Ticker

	m map[string]*topAppsEntry
	t byFirstMarkTimeHeap
	n byRequestsHeap
}

func NewTopApps() *TopApps {
	x := &TopApps{}

	x.Ticker = time.NewTicker(TopAppsTrimInterval)

	x.m = make(map[string]*topAppsEntry)
	x.t.Init()
	x.n.Init()

	go func() {
		for {
			select {
			case <-x.C:
				x.Trim(time.Now().Add(-TopAppsEntryLifetime))
			}
		}
	}()

	return x
}

func (x *TopApps) Mark(ApplicationId string, z time.Time) {
	t := z.Unix()

	x.Lock()
	defer x.Unlock()

	y := x.m[ApplicationId]
	if y != nil {
		z1 := heap.Remove(&x.t, y.ti).(*topAppsEntry)
		if z1 != y {
			panic("z1 != y")
		}
		z2 := heap.Remove(&x.n, y.ni).(*topAppsEntry)
		if z2 != y {
			panic("z2 != y")
		}
	} else {
		// New entry
		y = &topAppsEntry{ApplicationId: ApplicationId}
		x.m[ApplicationId] = y
	}

	y.Mark(t)

	heap.Push(&x.t, y)
	heap.Push(&x.n, y)
}

func (x *TopApps) trim(y time.Time) {
	t := y.Unix()

	for x.t.Len() > 0 {
		// Pop from the time heap
		u := heap.Pop(&x.t).(*topAppsEntry)

		// Remove from the requests heap
		v := heap.Remove(&x.n, u.ni).(*topAppsEntry)
		if v != u {
			panic("v != u")
		}

		if u.t[0].t > t {
			heap.Push(&x.t, u)
			heap.Push(&x.n, u)
			break
		}

		u.Trim(t)

		if len(u.t) == 0 {
			delete(x.m, u.ApplicationId)
			continue
		}

		// The first time slot should have t' > t
		if u.t[0].t <= t {
			panic("expected u.t[0].t <= t")
		}

		// Push back
		heap.Push(&x.t, u)
		heap.Push(&x.n, u)
	}
}

func (x *TopApps) Trim(y time.Time) {
	x.Lock()
	defer x.Unlock()

	x.trim(y)
}

type topAppsTopEntry struct {
	ApplicationId string
	Requests      int64
}

func (x *TopApps) TopSince(y time.Time, n int) []topAppsTopEntry {
	x.Lock()

	x.trim(y.Add(-1 * time.Second))

	a := byRequestsHeapSnapshot{}
	a.Heap = x.n.Copy()
	a.Init()

	x.Unlock()

	// Collect the top N applications
	s := make([]topAppsTopEntry, 0, n)
	for a.Len() > 0 && len(s) < n {
		z := heap.Pop(&a).(*topAppsEntry)

		s = append(s, topAppsTopEntry{
			ApplicationId: z.ApplicationId,
			Requests:      z.n,
		})
	}

	return s
}
