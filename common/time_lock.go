package common

import (
	"sync"
	"sync/atomic"
	"time"
)

const (
	logsize = 16
)

type TimeLock struct {
	lockdurations [logsize]int64
	locktimes     [logsize]int64
	lock          *sync.RWMutex
	index         int
}

func NewTimeLock() *TimeLock {
	return &TimeLock{
		lock: new(sync.RWMutex),
	}
}

func (self *TimeLock) Lock() {
	self.lock.Lock()
	atomic.StoreInt64(&self.locktimes[self.index], time.Now().UnixNano())
	atomic.StoreInt64(&self.lockdurations[self.index], -self.locktimes[self.index])
}

func (self *TimeLock) Unlock() {
	atomic.AddInt64(&self.lockdurations[self.index], time.Now().UnixNano())
	self.index = (self.index + 1) % logsize
	self.lock.Unlock()
}

func (self *TimeLock) RLock() {
	self.lock.RLock()
}

func (self *TimeLock) RUnlock() {
	self.lock.RUnlock()
}

func (self *TimeLock) Load() float64 {
	var sum int64
	var first int64
	var tmp int64
	for i := 0; i < logsize; i++ {
		tmp = atomic.LoadInt64(&self.lockdurations[i])
		if tmp > 0 {
			sum += tmp
		}
		tmp = atomic.LoadInt64(&self.locktimes[i])
		if first == 0 || tmp < first {
			first = tmp
		}
	}
	return float64(sum) / float64(time.Now().UnixNano()-first)
}
