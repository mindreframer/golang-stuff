package common

import (
	"sync"
	"syscall"
	"time"
)

const RefreshInterval time.Duration = time.Second * 1

type ProcessStatus struct {
	sync.Mutex
	rusage      *syscall.Rusage
	lastCpuTime int64
	stopSignal  chan bool
	stopped     bool

	CpuUsage float64
	MemRss   int64
}

func NewProcessStatus() *ProcessStatus {
	p := new(ProcessStatus)
	p.rusage = new(syscall.Rusage)

	go func() {
		timer := time.Tick(RefreshInterval)
		for {
			select {
			case <-timer:
				p.Update()
			case <-p.stopSignal:
				return
			}
		}
	}()

	return p
}

func (p *ProcessStatus) Update() {
	e := syscall.Getrusage(syscall.RUSAGE_SELF, p.rusage)
	if e != nil {
		log.Fatal(e.Error())
	}

	p.MemRss = int64(p.rusage.Maxrss)

	t := p.rusage.Utime.Nano() + p.rusage.Stime.Nano()
	p.CpuUsage = float64(t-p.lastCpuTime) / float64(RefreshInterval.Nanoseconds())
	p.lastCpuTime = t
}

func (p *ProcessStatus) StopUpdate() {
	p.Lock()
	defer p.Unlock()
	if !p.stopped {
		p.stopped = true
		p.stopSignal <- true
		p.stopSignal = nil
	}
}
