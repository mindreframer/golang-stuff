package discord

import (
	"bytes"
	"fmt"
	"github.com/zond/god/common"
	"github.com/zond/god/murmur"
	"net"
	"net/rpc"
	"sync"
	"sync/atomic"
	"time"
)

// CommListener is a function listening for generic communication between two Nodes.
type CommListener func(source, dest common.Remote, typ string) bool

// PingPack contains the sender and a hash of its discord ring, to let the receiver compare to its current ring.
type PingPack struct {
	Caller   common.Remote
	RingHash []byte
}

const (
	created = iota
	started
	stopped
)

// Node is a node in a chord like cluster.
//
// Like chord networks, it is a ring of nodes ordered by a position metric. Unlike chord, every node has every other node in its routing table.
// This allows stable networks to route with a constant time complexity.
type Node struct {
	ring          *common.Ring
	position      []byte
	listenAddr    string
	broadcastAddr string
	listener      *net.TCPListener
	metaLock      *sync.RWMutex
	routeLock     *sync.Mutex
	state         int32
	exports       map[string]interface{}
	commListeners []CommListener
}

func NewNode(listenAddr, broadcastAddr string) (result *Node) {
	return &Node{
		ring:          common.NewRing(),
		position:      make([]byte, murmur.Size),
		listenAddr:    listenAddr,
		broadcastAddr: broadcastAddr,
		exports:       make(map[string]interface{}),
		metaLock:      new(sync.RWMutex),
		routeLock:     new(sync.Mutex),
		state:         created,
	}
}

// Export will export the given api on a net/rpc server running on this Node.
func (self *Node) Export(name string, api interface{}) error {
	if self.hasState(created) {
		self.metaLock.Lock()
		defer self.metaLock.Unlock()
		self.exports[name] = api
		return nil
	}
	return fmt.Errorf("%v can only export when in state 'created'")
}
func (self *Node) AddCommListener(f CommListener) {
	self.metaLock.Lock()
	defer self.metaLock.Unlock()
	self.commListeners = append(self.commListeners, f)
}
func (self *Node) triggerCommListeners(source, dest common.Remote, typ string) {
	self.metaLock.RLock()
	newListeners := make([]CommListener, 0, len(self.commListeners))
	for _, l := range self.commListeners {
		self.metaLock.RUnlock()
		if l(source, dest, typ) {
			newListeners = append(newListeners, l)
		}
		self.metaLock.RLock()
	}
	self.metaLock.RUnlock()
	self.metaLock.Lock()
	defer self.metaLock.Unlock()
	self.commListeners = newListeners
}
func (self *Node) AddChangeListener(f common.RingChangeListener) {
	self.ring.AddChangeListener(f)
}
func (self *Node) SetPosition(position []byte) *Node {
	self.metaLock.Lock()
	self.position = make([]byte, len(position))
	copy(self.position, position)
	self.metaLock.Unlock()
	self.routeLock.Lock()
	defer self.routeLock.Unlock()
	self.ring.Add(self.Remote())
	return self
}

// GetNodes will return remotes to all Nodes in the ring.
func (self *Node) GetNodes() (result common.Remotes) {
	return self.ring.Nodes()
}

// Redundancy will return the current maximum redundancy in the ring.
func (self *Node) Redundancy() int {
	return self.ring.Redundancy()
}

// CountNodes returns the number of Nodes in the ring.
func (self *Node) CountNodes() int {
	return self.ring.Size()
}
func (self *Node) GetPosition() (result []byte) {
	self.metaLock.RLock()
	defer self.metaLock.RUnlock()
	result = make([]byte, len(self.position))
	copy(result, self.position)
	return
}
func (self *Node) GetListenAddr() string {
	self.metaLock.RLock()
	defer self.metaLock.RUnlock()
	return self.listenAddr
}
func (self *Node) GetBroadcastAddr() string {
	self.metaLock.RLock()
	defer self.metaLock.RUnlock()
	return self.broadcastAddr
}
func (self *Node) String() string {
	return fmt.Sprintf("<%v@%v>", common.HexEncode(self.GetPosition()), self.GetBroadcastAddr())
}

// Describe returns a humanly readable string describing the broadcast address, position and ring of this Node.
func (self *Node) Describe() string {
	self.metaLock.RLock()
	buffer := bytes.NewBufferString(fmt.Sprintf("%v@%v\n", common.HexEncode(self.position), self.broadcastAddr))
	self.metaLock.RUnlock()
	fmt.Fprint(buffer, self.ring.Describe())
	return string(buffer.Bytes())
}

func (self *Node) hasState(s int32) bool {
	return atomic.LoadInt32(&self.state) == s
}
func (self *Node) changeState(old, neu int32) bool {
	return atomic.CompareAndSwapInt32(&self.state, old, neu)
}
func (self *Node) getListener() *net.TCPListener {
	self.metaLock.RLock()
	defer self.metaLock.RUnlock()
	return self.listener
}
func (self *Node) setListener(l *net.TCPListener) {
	self.metaLock.Lock()
	defer self.metaLock.Unlock()
	self.listener = l
}

// Remote returns a remote to this Node.
func (self *Node) Remote() common.Remote {
	return common.Remote{self.GetPosition(), self.GetBroadcastAddr()}
}

// Stop will shut down this Node permanently.
func (self *Node) Stop() {
	if self.changeState(started, stopped) {
		self.getListener().Close()
	}
}
func (self *Node) MustStart() {
	if err := self.Start(); err != nil {
		panic(err)
	}
}

// Start will spin up this Node, export all its api interfaces and start its notify and ping jobs.
func (self *Node) Start() (err error) {
	if !self.changeState(created, started) {
		return fmt.Errorf("%v can only be started when in state 'created'", self)
	}
	if self.listenAddr == "" {
		return fmt.Errorf("%v needs to have an address to listen at", self)
	}
	var addr *net.TCPAddr
	if addr, err = net.ResolveTCPAddr("tcp", self.listenAddr); err != nil {
		return
	}
	var listener *net.TCPListener
	if listener, err = net.ListenTCP("tcp", addr); err != nil {
		return
	}
	self.setListener(listener)
	server := rpc.NewServer()
	if err = server.RegisterName("Discord", (*nodeServer)(self)); err != nil {
		return
	}
	for name, api := range self.exports {
		if err = server.RegisterName(name, api); err != nil {
			return
		}
	}
	self.ring.Add(self.Remote())
	go server.Accept(self.getListener())
	go self.notifyPeriodically()
	go self.pingPeriodically()
	return
}
func (self *Node) notifyPeriodically() {
	for self.hasState(started) {
		self.notifySuccessor()
		time.Sleep(common.PingInterval)
	}
}
func (self *Node) pingPeriodically() {
	for self.hasState(started) {
		self.pingPredecessor()
		time.Sleep(common.PingInterval)
	}
}

// RingHash returns a hash of the discord ring of this Node.
func (self *Node) RingHash() []byte {
	return self.ring.Hash()
}

// Ping will compare the hash of this Node with the one in the received PingPack, and request the entire routing ring from the sender if they are not equal.
func (self *Node) Ping(ping PingPack) (me common.Remote) {
	me = self.Remote()
	if bytes.Compare(ping.RingHash, self.ring.Hash()) != 0 {
		var newNodes common.Remotes
		if err := ping.Caller.Call("Discord.Nodes", 0, &newNodes); err != nil {
			self.RemoveNode(ping.Caller)
		} else {
			self.routeLock.Lock()
			defer self.routeLock.Unlock()
			pred := self.ring.Predecessor(me)
			self.ring.SetNodes(newNodes)
			self.ring.Add(me)
			self.ring.Add(pred)
			self.ring.Clean(pred, me)
		}
	}
	return
}
func (self *Node) pingPredecessor() {
	pred := self.GetPredecessor()
	ping := PingPack{
		RingHash: self.ring.Hash(),
		Caller:   self.Remote(),
	}
	var newPred common.Remote
	op := "Discord.Ping"
	self.triggerCommListeners(self.Remote(), pred, op)
	if err := pred.Call(op, ping, &newPred); err != nil {
		self.RemoveNode(pred)
	} else {
		self.routeLock.Lock()
		defer self.routeLock.Unlock()
		self.ring.Add(newPred)
	}
}

// Nodes will return remotes for all Nodes in the ring.
func (self *Node) Nodes() common.Remotes {
	return self.ring.Nodes()
}

// Notify will add the caller to the ring of this Node.
func (self *Node) Notify(caller common.Remote) common.Remote {
	self.routeLock.Lock()
	defer self.routeLock.Unlock()
	self.ring.Add(caller)
	return self.GetPredecessor()
}
func (self *Node) notifySuccessor() {
	succ := self.GetSuccessor()
	var otherPred common.Remote
	op := "Discord.Notify"
	selfRemote := self.Remote()
	self.triggerCommListeners(selfRemote, succ, op)
	if err := succ.Call(op, selfRemote, &otherPred); err != nil {
		self.RemoveNode(succ)
	} else {
		if otherPred.Addr != self.GetBroadcastAddr() {
			self.routeLock.Lock()
			defer self.routeLock.Unlock()
			self.ring.Add(otherPred)
		}
	}
}
func (self *Node) MustJoin(addr string) {
	if err := self.Join(addr); err != nil {
		panic(err)
	}
}

// Join will fetch the routing ring of the Node at addr, pick a location on an empty spot in the received ring and notify the other Node of our joining.
func (self *Node) Join(addr string) (err error) {
	var newNodes common.Remotes
	if err = common.Switch.Call(addr, "Discord.Nodes", 0, &newNodes); err != nil {
		return
	}
	if bytes.Compare(self.GetPosition(), make([]byte, murmur.Size)) == 0 {
		self.SetPosition(common.NewRingNodes(newNodes).GetSlot())
	}
	self.routeLock.Lock()
	self.ring.SetNodes(newNodes)
	self.routeLock.Unlock()
	var x common.Remote
	if err = common.Switch.Call(addr, "Discord.Notify", self.Remote(), &x); err != nil {
		return
	}
	return
}

// RemoveNode will remove the provided remote from our routing ring.
func (self *Node) RemoveNode(remote common.Remote) {
	if remote.Addr == self.GetBroadcastAddr() {
		panic(fmt.Errorf("%v is trying to remove itself from the routing!", self))
	}
	self.routeLock.Lock()
	defer self.routeLock.Unlock()
	self.ring.Remove(remote)
}

// GetPredecessor will return our predecessor on the ring.
func (self *Node) GetPredecessor() common.Remote {
	return self.GetPredecessorForRemote(self.Remote())
}

// GetPredecessorForRemote will return the predecessor for the provided remote.
func (self *Node) GetPredecessorForRemote(r common.Remote) common.Remote {
	return self.ring.Predecessor(r)
}

// GetPredecessorFor will return the predecessor for the provided key.
func (self *Node) GetPredecessorFor(key []byte) common.Remote {
	pred, _, _ := self.ring.Remotes(key)
	return *pred
}

// HasNode will return true if there is a Node on the ring with the given pos.
func (self *Node) HasNode(pos []byte) bool {
	if _, match, _ := self.ring.Remotes(pos); match != nil {
		return true
	}
	return false
}

// GetSuccessor will return our successor on the ring.
func (self *Node) GetSuccessor() common.Remote {
	return self.GetSuccessorForRemote(self.Remote())
}

// GetSuccessorFor will return the successor for the provided remote.
func (self *Node) GetSuccessorForRemote(r common.Remote) common.Remote {
	return self.ring.Successor(r)
}

// GetSuccessorFor will return the successor for the provided key.
// If the successor is not this Node, it will assert that the provided key is between the found successor and the predecessor it claims to have.
func (self *Node) GetSuccessorFor(key []byte) common.Remote {
	// Guess according to our route cache
	predecessor, match, successor := self.ring.Remotes(key)
	if match != nil {
		predecessor = match
	}
	// If we consider ourselves successors, just return us
	if successor.Addr != self.GetBroadcastAddr() {
		// Double check by asking the successor we found what predecessor it has
		if err := successor.Call("Discord.GetPredecessor", 0, predecessor); err != nil {
			self.RemoveNode(*successor)
			return self.GetSuccessorFor(key)
		}
		// If the key we are looking for is between them, just return the successor
		if !common.BetweenIE(key, predecessor.Pos, successor.Pos) {
			// Otherwise, ask the predecessor we actually found about who is the successor of the key
			if err := predecessor.Call("Discord.GetSuccessorFor", key, successor); err != nil {
				self.RemoveNode(*predecessor)
				return self.GetSuccessorFor(key)
			}
		}
	}
	return *successor
}
