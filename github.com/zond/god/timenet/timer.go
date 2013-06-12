package timenet

import (
	"math"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

const (
	loglen = 10
)

const (
	created = iota
	started
	stopped
)

type Peer interface {
	ActualTime() (time time.Time)
}

type PeerProducer interface {
	Peers() map[string]Peer
}

// Timer is an abstract time synchronization structure that needs a PeerProducer
// to provide it with peers from a set of other Timers. It will randomly contact
// the Peers produced by the PeerProducer to synchronize its time with other peers
// in the network.
//
// It calculates deviation from _normal_ network latency, and only uses time measurements
// that respond within standard deviation from the normal response times when adjusting
// its timer.
type Timer struct {
	lock          *sync.RWMutex
	state         int32
	offset        int64
	dilations     *dilations
	peerProducer  PeerProducer
	peerErrors    map[string]int64
	peerLatencies map[string]times
}

func NewTimer(producer PeerProducer) *Timer {
	return &Timer{
		lock:          &sync.RWMutex{},
		peerProducer:  producer,
		state:         created,
		dilations:     &dilations{},
		peerErrors:    make(map[string]int64),
		peerLatencies: make(map[string]times),
	}
}
func (self *Timer) adjustments() int64 {
	return self.offset + self.dilations.delta()
}

// ActualTime will return what time this Timer thinks it is. It may change forwards or backwards
// non continuously depending on network interactions.
func (self *Timer) ActualTime() time.Time {
	self.lock.RLock()
	defer self.lock.RUnlock()
	return time.Unix(0, time.Now().UnixNano()+self.adjustments())
}

// ContinuousTime will return a continous nice version of the time this Timer thinks it is. It us guaranteed
// to never move backwards, and to only move forwards in a smooth fashion.
func (self *Timer) ContinuousTime() (result int64) {
	self.lock.RLock()
	temporaryEffect, permanentEffect := self.dilations.effect()
	result = time.Now().UnixNano() + self.offset + permanentEffect + temporaryEffect
	self.lock.RUnlock()
	if permanentEffect != 0 {
		self.lock.Lock()
		defer self.lock.Unlock()
		self.offset += permanentEffect
	}
	return
}

// Error returns the deviation of the error of this Timer.
func (self *Timer) Error() (err time.Duration) {
	self.lock.RLock()
	defer self.lock.RUnlock()
	if len(self.peerErrors) > 1 {
		var thisErr int64
		for _, e := range self.peerErrors {
			thisErr = e >> 10
			err += time.Duration(thisErr * thisErr)
		}
		err = time.Duration(math.Sqrt(float64(err/time.Duration(len(self.peerErrors))))) << 10
	} else {
		err = -1
	}
	return
}

// Stability returns the deviation of the latency between this Timer and its peers.
func (self *Timer) Stability() (result time.Duration) {
	self.lock.RLock()
	defer self.lock.RUnlock()
	if len(self.peerLatencies) > 1 {
		var deviation int64
		for _, latencies := range self.peerLatencies {
			_, deviation = latencies.stats()
			result += time.Duration(deviation * deviation)
		}
		result = time.Duration(math.Sqrt(float64(result / time.Duration(len(self.peerLatencies)))))
	} else {
		result = -1
	}
	return
}
func (self *Timer) adjust(id string, adjustment int64) {
	self.peerErrors[id] = adjustment
	self.dilations.add(adjustment)
}
func (self *Timer) randomPeer() (id string, peer Peer, peerLatencies times) {
	currentPeers := self.peerProducer.Peers()
	chosenIndex := rand.Int() % len(currentPeers)
	for thisId, theseLatencies := range self.peerLatencies {
		if currentPeer, ok := currentPeers[thisId]; ok {
			if chosenIndex == 0 {
				peer = currentPeer
				id = thisId
				peerLatencies = theseLatencies
			}
			chosenIndex--
			delete(currentPeers, thisId)
		} else {
			delete(self.peerLatencies, thisId)
			delete(self.peerErrors, thisId)
		}
	}
	for thisId, thisPeer := range currentPeers {
		if chosenIndex == 0 {
			peer = thisPeer
			id = thisId
		}
		chosenIndex--
	}
	return
}
func (self *Timer) timeAndLatency(peer Peer) (peerTime, latency, myTime int64) {
	latency = -time.Now().UnixNano()
	peerTime = peer.ActualTime().UnixNano()
	latency += time.Now().UnixNano()
	peerTime += latency / 2
	self.lock.RLock()
	defer self.lock.RUnlock()
	myTime = self.ActualTime().UnixNano()
	return
}

// Conform will adjust this Timer to be as exactly as possible that of the provided peer.
func (self *Timer) Conform(peer Peer) {
	peerTime, _, myTime := self.timeAndLatency(peer)
	self.lock.Lock()
	defer self.lock.Unlock()
	self.offset += (peerTime - myTime)
}

// Skew will change the actual time of this timer with delta, adjusting the continuous time in a smooth fashion.
func (self *Timer) Skew(delta time.Duration) {
	self.lock.Lock()
	defer self.lock.Unlock()
	self.offset += int64(delta)
}

// Sample will make this Timer sample a random Peer produced by the PeerProducer, and Skew according to the delta with the time of that Peer.
func (self *Timer) Sample() {
	self.lock.RLock()
	peerId, peer, oldLatencies := self.randomPeer()
	self.lock.RUnlock()

	peerTime, latency, myTime := self.timeAndLatency(peer)

	self.lock.Lock()
	defer self.lock.Unlock()
	oldestLatencyIndex := 0
	if len(oldLatencies) > loglen {
		oldestLatencyIndex = len(oldLatencies) - loglen
	}
	newLatencies := append(oldLatencies[oldestLatencyIndex:], latency)
	self.peerLatencies[peerId] = newLatencies

	mean, deviation := newLatencies.stats()
	if math.Abs(float64(latency-mean)) < float64(deviation) {
		self.adjust(peerId, peerTime-myTime)
	}
}
func (self *Timer) hasState(s int32) bool {
	return atomic.LoadInt32(&self.state) == s
}
func (self *Timer) changeState(old, neu int32) bool {
	return atomic.CompareAndSwapInt32(&self.state, old, neu)
}
func (self *Timer) sleep() {
	err := self.Error()
	stability := self.Stability()
	if err == -1 || stability == -1 {
		time.Sleep(time.Second)
	} else {
		if err == 0 {
			err = 1
		}
		sleepyTime := ((time.Duration(stability) * time.Second) << 7) / time.Duration(err)
		if sleepyTime < time.Second {
			sleepyTime = time.Second
		}
		time.Sleep(sleepyTime)
	}
}

// Run will make this Timer regularly Sample. 
func (self *Timer) Run() {
	for self.hasState(started) {
		self.Sample()
		self.sleep()
	}
}

// Stop will permanently stop this Timer.
func (self *Timer) Stop() {
	self.changeState(started, stopped)
}

// Start will start a goroutine Running this Timer.
func (self *Timer) Start() {
	if self.changeState(created, started) {
		go self.Run()
	}
}
