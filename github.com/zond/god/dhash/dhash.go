package dhash

import (
	"bytes"
	"fmt"
	"github.com/zond/god/common"
	"github.com/zond/god/discord"
	"github.com/zond/god/murmur"
	"github.com/zond/god/radix"
	"github.com/zond/god/timenet"
	"sync"
	"sync/atomic"
	"time"
)

// SyncListener is a function listening for sync events where one dhash.Node has pushed items to, and pulled items from, another dhash.Node.
type SyncListener func(source, dest common.Remote, pulled, pushed int) (keep bool)

// CleanListener is a function listening for clean events where one dhash.Node has cleaned items from itself and pushed items to another dhash.Node.
type CleanListener func(source, dest common.Remote, cleaned, pushed int) (keep bool)

// MigrateListener is a function listening for migrate events where one dhash.Node has migrated from one position to another.
type MigrateListener func(dhash *Node, source, destination []byte) (keep bool)

// CommListener is a function listening to generic communications between two dhash.Nodes.
type CommListener func(comm Comm) (keep bool)

type commListenerContainer struct {
	channel  chan Comm
	listener CommListener
	node     *Node
}

func (self *commListenerContainer) run() {
	for self.listener(<-self.channel) {
	}
	self.node.removeCommListener(self)
}

// Comm contains metadata about one communication between two dhash.Nodes.
type Comm struct {
	Source      common.Remote
	Destination common.Remote
	Key         []byte
	SubKey      []byte
	Type        string
}

const (
	syncInterval      = time.Second
	migrateHysteresis = 1.5
	migrateWaitFactor = 2
)

const (
	created = iota
	started
	stopped
)

// Node is a node in the database. It contains a discord.Node containing routing and rpc functionality, 
// a timenet.Timer containing time synchronization functionality and a radix.Tree containing the actual data.
type Node struct {
	lastSync         int64
	lastMigrate      int64
	lastReroute      int64
	state            int32
	lock             *sync.RWMutex
	syncListeners    []SyncListener
	cleanListeners   []CleanListener
	migrateListeners []MigrateListener
	commListeners    map[*commListenerContainer]bool
	nCommListeners   int32
	node             *discord.Node
	timer            *timenet.Timer
	tree             *radix.Tree
}

func NewNode(listenAddr, broadcastAddr string) *Node {
	return NewNodeDir(listenAddr, broadcastAddr, broadcastAddr)
}

// NewNode will return a dhash.Node publishing itself on the given address.
func NewNodeDir(listenAddr, broadcastAddr, dir string) (result *Node) {
	result = &Node{
		node:          discord.NewNode(listenAddr, broadcastAddr),
		lock:          new(sync.RWMutex),
		commListeners: make(map[*commListenerContainer]bool),
		state:         created,
	}
	result.node.AddCommListener(func(source, dest common.Remote, typ string) bool {
		if result.hasState(started) {
			if result.hasCommListeners() {
				result.triggerCommListeners(Comm{
					Source:      source,
					Destination: dest,
					Type:        typ,
				})
			}
		}
		return !result.hasState(stopped)
	})
	result.AddChangeListener(func(r *common.Ring) bool {
		atomic.StoreInt64(&result.lastReroute, time.Now().UnixNano())
		return true
	})
	result.timer = timenet.NewTimer((*dhashPeerProducer)(result))
	result.tree = radix.NewTreeTimer(result.timer)
	if dir != "" {
		result.tree.Log(dir).Restore()
	}
	result.node.Export("Timenet", (*timerServer)(result.timer))
	result.node.Export("DHash", (*dhashServer)(result))
	result.node.Export("HashTree", (*hashTreeServer)(result))
	return
}
func (self *Node) AddCommListener(l CommListener) {
	self.lock.Lock()
	defer self.lock.Unlock()
	newListener := &commListenerContainer{
		listener: l,
		channel:  make(chan Comm),
		node:     self,
	}
	atomic.AddInt32(&self.nCommListeners, 1)
	self.commListeners[newListener] = true
	go newListener.run()
}
func (self *Node) removeCommListener(lc *commListenerContainer) {
	self.lock.Lock()
	defer self.lock.Unlock()
	if _, ok := self.commListeners[lc]; ok {
		atomic.AddInt32(&self.nCommListeners, -1)
	}
	delete(self.commListeners, lc)
	close(lc.channel)
}
func (self *Node) hasCommListeners() bool {
	return atomic.LoadInt32(&self.nCommListeners) > 0
}
func (self *Node) triggerCommListeners(comm Comm) {
	self.lock.RLock()
	for lc, _ := range self.commListeners {
		self.lock.RUnlock()
		select {
		case lc.channel <- comm:
		default:
		}
		self.lock.RLock()
	}
	self.lock.RUnlock()
}
func (self *Node) AddCleanListener(l CleanListener) {
	self.lock.Lock()
	defer self.lock.Unlock()
	self.cleanListeners = append(self.cleanListeners, l)
}
func (self *Node) AddMigrateListener(l MigrateListener) {
	self.lock.Lock()
	defer self.lock.Unlock()
	self.migrateListeners = append(self.migrateListeners, l)
}
func (self *Node) AddSyncListener(l SyncListener) {
	self.lock.Lock()
	defer self.lock.Unlock()
	self.syncListeners = append(self.syncListeners, l)
}
func (self *Node) hasState(s int32) bool {
	return atomic.LoadInt32(&self.state) == s
}
func (self *Node) changeState(old, neu int32) bool {
	return atomic.CompareAndSwapInt32(&self.state, old, neu)
}
func (self *Node) GetListenAddr() string {
	return self.node.GetListenAddr()
}
func (self *Node) GetBroadcastAddr() string {
	return self.node.GetBroadcastAddr()
}
func (self *Node) AddChangeListener(f common.RingChangeListener) {
	self.node.AddChangeListener(f)
}

// Stop will shut down this dhash.Node, including its discord.Node and timenet.Timer,  permanently.
func (self *Node) Stop() {
	if self.changeState(started, stopped) {
		self.node.Stop()
		self.timer.Stop()
	}
}

// Start will spin up this dhash.Node, including its discord.Node and timenet.Timer.
// It will also start the sync, clean and migrate jobs.
func (self *Node) Start() (err error) {
	if !self.changeState(created, started) {
		return fmt.Errorf("%v can only be started when in state 'created'", self)
	}
	if err = self.node.Start(); err != nil {
		return
	}
	self.timer.Start()
	go self.syncPeriodically()
	go self.cleanPeriodically()
	go self.migratePeriodically()
	self.startJson()
	return
}
func (self *Node) triggerSyncListeners(source, dest common.Remote, pulled, pushed int) {
	self.lock.RLock()
	newListeners := make([]SyncListener, 0, len(self.syncListeners))
	for _, l := range self.syncListeners {
		self.lock.RUnlock()
		if l(source, dest, pulled, pushed) {
			newListeners = append(newListeners, l)
		}
		self.lock.RLock()
	}
	self.lock.RUnlock()
	self.lock.Lock()
	defer self.lock.Unlock()
	self.syncListeners = newListeners
}
func (self *Node) sync() {
	var pulled int
	var pushed int
	selfRemote := self.node.Remote()
	nextSuccessor := self.node.GetSuccessor()
	for i := 0; i < self.node.Redundancy()-1; i++ {
		myPos := self.node.GetPosition()
		remoteHash := remoteHashTree{
			source:      selfRemote,
			destination: nextSuccessor,
			node:        self,
		}
		pushed = radix.NewSync(self.tree, remoteHash).From(self.node.GetPredecessor().Pos).To(myPos).Run().PutCount()
		pulled = radix.NewSync(remoteHash, self.tree).From(self.node.GetPredecessor().Pos).To(myPos).Run().PutCount()
		if pushed != 0 || pulled != 0 {
			self.triggerSyncListeners(selfRemote, nextSuccessor, pulled, pushed)
		}
		nextSuccessor = self.node.GetSuccessorForRemote(nextSuccessor)
	}
}
func (self *Node) syncPeriodically() {
	for self.hasState(started) {
		self.sync()
		time.Sleep(syncInterval)
	}
}
func (self *Node) cleanPeriodically() {
	for self.hasState(started) {
		self.clean()
		time.Sleep(syncInterval)
	}
}
func (self *Node) triggerMigrateListeners(oldPos, newPos []byte) {
	self.lock.RLock()
	newListeners := make([]MigrateListener, 0, len(self.migrateListeners))
	for _, l := range self.migrateListeners {
		self.lock.RUnlock()
		if l(self, oldPos, newPos) {
			newListeners = append(newListeners, l)
		}
		self.lock.RLock()
	}
	self.lock.RUnlock()
	self.lock.Lock()
	defer self.lock.Unlock()
	self.migrateListeners = newListeners
}
func (self *Node) changePosition(newPos []byte) {
	for len(newPos) < murmur.Size {
		newPos = append(newPos, 0)
	}
	oldPos := self.node.GetPosition()
	if bytes.Compare(newPos, oldPos) != 0 {
		self.node.SetPosition(newPos)
		atomic.StoreInt64(&self.lastMigrate, time.Now().UnixNano())
		self.triggerMigrateListeners(oldPos, newPos)
	}
}
func (self *Node) isLeader() bool {
	return bytes.Compare(self.node.GetPredecessor().Pos, self.node.GetPosition()) > 0
}
func (self *Node) migratePeriodically() {
	for self.hasState(started) {
		self.migrate()
		time.Sleep(syncInterval)
	}
}
func (self *Node) migrate() {
	lastAllowedChange := time.Now().Add(-1 * migrateWaitFactor * syncInterval).UnixNano()
	if lastAllowedChange > common.Max64(atomic.LoadInt64(&self.lastSync), atomic.LoadInt64(&self.lastReroute), atomic.LoadInt64(&self.lastMigrate)) {
		var succSize int
		succ := self.node.GetSuccessor()
		if err := succ.Call("DHash.Owned", 0, &succSize); err != nil {
			self.node.RemoveNode(succ)
		} else {
			mySize := self.Owned()
			if mySize > 10 && float64(mySize) > float64(succSize)*migrateHysteresis {
				wantedDelta := (mySize - succSize) / 2
				var existed bool
				var wantedPos []byte
				pred := self.node.GetPredecessor()
				if bytes.Compare(pred.Pos, self.node.GetPosition()) < 1 {
					if wantedPos, existed = self.tree.NextMarkerIndex(self.tree.RealSizeBetween(nil, self.node.GetPosition(), true, false) - wantedDelta); !existed {
						return
					}
				} else {
					ownedAfterNil := self.tree.RealSizeBetween(nil, succ.Pos, true, false)
					if ownedAfterNil > wantedDelta {
						if wantedPos, existed = self.tree.NextMarkerIndex(ownedAfterNil - wantedDelta); !existed {
							return
						}
					} else {
						if wantedPos, existed = self.tree.NextMarkerIndex(self.tree.RealSize() + ownedAfterNil - wantedDelta); !existed {
							return
						}
					}
				}
				if common.BetweenIE(wantedPos, self.node.GetPredecessor().Pos, self.node.GetPosition()) {
					self.changePosition(wantedPos)
				}
			}
		}
	}
}
func (self *Node) circularNext(key []byte) (nextKey []byte, existed bool) {
	if nextKey, existed = self.tree.NextMarker(key); existed {
		return
	}
	nextKey = make([]byte, murmur.Size)
	if _, _, existed = self.tree.Get(nextKey); existed {
		return
	}
	nextKey, existed = self.tree.NextMarker(nextKey)
	return
}
func (self *Node) owners(key []byte) (owners common.Remotes, isOwner bool) {
	owners = append(owners, self.node.GetSuccessorFor(key))
	if owners[0].Addr == self.node.GetBroadcastAddr() {
		isOwner = true
	}
	for i := 1; i < self.node.Redundancy(); i++ {
		owners = append(owners, self.node.GetSuccessorForRemote(owners[i-1]))
		if owners[i].Addr == self.node.GetBroadcastAddr() {
			isOwner = true
		}
	}
	return
}
func (self *Node) triggerCleanListeners(source, dest common.Remote, cleaned, pushed int) {
	self.lock.RLock()
	newListeners := make([]CleanListener, 0, len(self.cleanListeners))
	for _, l := range self.cleanListeners {
		self.lock.RUnlock()
		if l(source, dest, cleaned, pushed) {
			newListeners = append(newListeners, l)
		}
		self.lock.RLock()
	}
	self.lock.RUnlock()
	self.lock.Lock()
	defer self.lock.Unlock()
	self.cleanListeners = newListeners
}
func (self *Node) clean() {
	selfRemote := self.node.Remote()
	var cleaned int
	var pushed int
	if nextKey, existed := self.circularNext(self.node.GetPosition()); existed {
		if owners, isOwner := self.owners(nextKey); !isOwner {
			var sync *radix.Sync
			for index, owner := range owners {
				sync = radix.NewSync(self.tree, remoteHashTree{
					source:      selfRemote,
					destination: owner,
					node:        self,
				}).From(nextKey).To(owners[0].Pos)
				if index == len(owners)-2 {
					sync.Destroy()
				}
				sync.Run()
				cleaned = sync.DelCount()
				pushed = sync.PutCount()
				if cleaned != 0 || pushed != 0 {
					self.triggerCleanListeners(selfRemote, owner, cleaned, pushed)
				}
			}
		}
	}
}
func (self *Node) MustStart() *Node {
	if err := self.Start(); err != nil {
		panic(err)
	}
	return self
}
func (self *Node) MustJoin(addr string) {
	self.timer.Conform(remotePeer(common.Remote{Addr: addr}))
	self.node.MustJoin(addr)
}
func (self *Node) Time() time.Time {
	return time.Unix(0, self.timer.ContinuousTime())
}

// Owned returns the number of items, including tombstones, that this node has responsibility for.
func (self *Node) Owned() int {
	pred := self.node.GetPredecessor()
	me := self.node.Remote()
	cmp := bytes.Compare(pred.Pos, me.Pos)
	if cmp < 0 {
		return self.tree.RealSizeBetween(pred.Pos, me.Pos, true, false)
	} else if cmp > 0 {
		return self.tree.RealSizeBetween(pred.Pos, nil, true, false) + self.tree.RealSizeBetween(nil, me.Pos, true, false)
	}
	if pred.Less(me) {
		return 0
	}
	return self.tree.RealSize()
}
