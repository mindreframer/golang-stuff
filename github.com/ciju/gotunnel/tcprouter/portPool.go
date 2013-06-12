package tcprouter

import (
	"fmt"
)

type PortPool struct {
	Start         int
	End           int
	lastAllocated int
	Pool          map[int]struct{}
}

type yes struct{}

func NewPortPool(start, end int) *PortPool {
	pool := &PortPool{Start: start, End: end, lastAllocated: start, Pool: make(map[int]struct{})}
	for i := start; i <= end; i++ {
		pool.Pool[i] = yes{}
	}
	return pool
}

func (p *PortPool) inRange(port int) bool {
	return p.Start <= port && port <= p.End
}
func (p *PortPool) IsAvailable(port int) bool {
	if !p.inRange(port) {
		return false
	}
	_, ok := p.Pool[port]
	return ok
}
func (p *PortPool) SetUnavailable(port int) {
	if !p.inRange(port) {
		return
	}
	delete(p.Pool, port)
	p.lastAllocated = port
}
func (p *PortPool) SetAvailable(port int) {
	if !p.inRange(port) {
		return
	}
	p.Pool[port] = yes{}
}
func (p *PortPool) GetAvailable() (int, bool) {
	for i := p.lastAllocated; i <= p.End; i++ {
		if p.IsAvailable(i) {
			p.SetUnavailable(i)
			return i, true
		}
	}

	for i := p.Start; i <= p.lastAllocated; i++ {
		if p.IsAvailable(i) {
			p.SetUnavailable(i)
			return i, true
		}
	}

	return 0, false
}

func (p *PortPool) String() string {
	return fmt.Sprintf("pool: %v", p.Pool)
}
