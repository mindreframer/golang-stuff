package bench

import (
	"flag"
	"fmt"
	"github.com/zond/god/client"
	"github.com/zond/god/murmur"
	"math/rand"
	"net"
	"net/rpc"
	"sync"
	"sync/atomic"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

const (
	stopped = iota
	started
)

type SpinResult struct {
	Nodes int
	Keys  int
}

type PrepareCommand struct {
	Addr  string
	Range [2]int64
}

type SpinCommand struct {
	Addr   string
	MaxKey int64
}

type Nothing struct{}

type Slave struct {
	maxRps  int64
	currRps int64
	maxKey  int64
	req     int64
	start   time.Time
	addr    string
	state   int32
	client  *client.Conn
	wg      *sync.WaitGroup
}

func (self *Slave) switchState(expected, wanted int32) bool {
	return atomic.CompareAndSwapInt32(&self.state, expected, wanted)
}
func (self *Slave) hasState(s int32) bool {
	return atomic.LoadInt32(&self.state) == s
}

func (self *Slave) spinner() {
	var kv []byte
	for self.hasState(started) {
		kv = murmur.HashInt64(rand.Int63n(self.maxKey))
		self.client.Put(kv, kv)
		atomic.AddInt64(&self.req, 1)
	}
}

func (self *Slave) run() {
	freebies := 2
	peaked := false
	var curr int64
	for self.hasState(started) {
		curr = atomic.LoadInt64(&self.req) / ((time.Now().UnixNano() - self.start.UnixNano()) / int64(time.Second))
		atomic.StoreInt64(&self.currRps, curr)
		if self.maxRps == 0 || freebies > 0 || curr > self.maxRps {
			fmt.Println("Spinning up one more loader, curr:", curr, "max:", self.maxRps)
			go self.spinner()
			if curr < self.maxRps {
				freebies--
			}
			if self.maxRps == 0 || curr > self.maxRps {
				self.maxRps = curr
			}
		} else if !peaked {
			fmt.Println("Peaked at", self.maxRps)
			self.wg.Done()
			peaked = true
		}
		self.req = 0
		self.start = time.Now()
		time.Sleep(time.Second)
	}
}

func (self *Slave) Prepare(command PrepareCommand, x *Nothing) error {
	if self.hasState(stopped) {
		fmt.Printf("Preparing %+v\n", command)
		self.client = client.MustConn(command.Addr)
		var kv []byte
		for i := command.Range[0]; i < command.Range[1]; i++ {
			kv = murmur.HashInt64(i)
			self.client.Put(kv, kv)
			if i%1000 == 0 {
				fmt.Print(".")
			}
		}
		fmt.Println("done!")
	}
	return nil
}

func (self *Slave) Current(x Nothing, rps *int64) error {
	if self.hasState(started) {
		self.wg.Wait()
		*rps = atomic.LoadInt64(&self.currRps)
		return nil
	}
	return fmt.Errorf("%v is not started", self)
}

func (self *Slave) Stop(x Nothing, y *Nothing) error {
	if self.hasState(started) {
		self.wg.Wait()
		self.switchState(started, stopped)
	}
	return nil
}

func (self *Slave) Wait(x Nothing, y *Nothing) error {
	if self.hasState(started) {
		self.wg.Wait()
		return nil
	}
	return fmt.Errorf("%v is not started", self)
}

func (self *Slave) Spin(command SpinCommand, result *SpinResult) error {
	if self.switchState(stopped, started) {
		fmt.Println("Spinning on ", command.Addr)
		self.wg = new(sync.WaitGroup)
		self.wg.Add(1)
		self.maxRps = 0
		self.client = client.MustConn(command.Addr)
		self.maxKey = command.MaxKey
		go self.run()
	} else {
		fmt.Println("Already started on", self.addr)
	}
	*result = SpinResult{
		Nodes: len(self.client.Nodes()),
		Keys:  self.client.Size(),
	}
	return nil
}

func RunSlave() {
	var ip = flag.String("ip", "127.0.0.1", "IP address to listen to")
	var port = flag.Int("port", 19191, "Port to connect to")
	flag.Parse()
	var err error
	var addr *net.TCPAddr
	if addr, err = net.ResolveTCPAddr("tcp", fmt.Sprintf("%v:%v", *ip, *port)); err != nil {
		panic(err)
	}
	var listener *net.TCPListener
	if listener, err = net.ListenTCP("tcp", addr); err != nil {
		panic(err)
	}
	rpc.Register(&Slave{})
	rpc.Accept(listener)
}
