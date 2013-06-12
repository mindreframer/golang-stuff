package dhash

import (
	"bytes"
	"fmt"
	"github.com/zond/god/common"
	"os"
	"runtime"
	"sort"
	"testing"
	"time"
)

type dhashAry []*Node

func (self dhashAry) Less(i, j int) bool {
	return self[i].node.Remote().Less(self[j].node.Remote())
}
func (self dhashAry) Swap(i, j int) {
	self[i], self[j] = self[j], self[i]
}
func (self dhashAry) Len() int {
	return len(self)
}

func countHaving(t *testing.T, dhashes []*Node, key, value []byte) (result int) {
	for _, d := range dhashes {
		if foundValue, _, existed := d.tree.Get(key); existed && bytes.Compare(foundValue, value) == 0 {
			result++
		}
	}
	return
}

func testStartup(t *testing.T, n, port int) (dhashes []*Node) {
	for i := 0; i < n; i++ {
		os.RemoveAll(fmt.Sprintf("127.0.0.1:%v", port+i*2))
	}
	dhashes = make([]*Node, n)
	for i := 0; i < n; i++ {
		dhashes[i] = NewNode(fmt.Sprintf("127.0.0.1:%v", port+i*2), fmt.Sprintf("127.0.0.1:%v", port+i*2))
		dhashes[i].MustStart()
	}
	for i := 1; i < n; i++ {
		dhashes[i].MustJoin(fmt.Sprintf("127.0.0.1:%v", port))
	}
	common.AssertWithin(t, func() (string, bool) {
		routes := make(map[string]bool)
		for _, dhash := range dhashes {
			routes[dhash.node.GetNodes().Describe()] = true
		}
		return fmt.Sprint(routes), len(routes) == 1
	}, time.Second*10)
	return
}

func testSync(t *testing.T, dhashes []*Node) {
	dhashes[0].tree.Put([]byte{3}, []byte{0}, 1)
	common.AssertWithin(t, func() (string, bool) {
		having := countHaving(t, dhashes, []byte{3}, []byte{0})
		return fmt.Sprint(having), having == common.Redundancy
	}, time.Second*10)
}

func testClean(t *testing.T, dhashes []*Node) {
	for _, n := range dhashes {
		n.tree.Put([]byte{1}, []byte{1}, 1)
	}
	common.AssertWithin(t, func() (string, bool) {
		having := countHaving(t, dhashes, []byte{1}, []byte{1})
		return fmt.Sprint(having), having == common.Redundancy
	}, time.Second*20)
}

func testPut(t *testing.T, dhashes []*Node) {
	for index, n := range dhashes {
		n.Put(common.Item{Key: []byte{byte(index + 100)}, Timestamp: 1, Value: []byte{byte(index + 100)}})
	}
	common.AssertWithin(t, func() (string, bool) {
		haves := make(map[int]bool)
		for index, _ := range dhashes {
			count := countHaving(t, dhashes, []byte{byte(index + 100)}, []byte{byte(index + 100)})
			haves[count] = true
		}
		return fmt.Sprint(haves), len(haves) == 1 && haves[common.Redundancy] == true
	}, time.Second*10)
}

func testMigrate(t *testing.T, dhashes []*Node) {
	for _, d := range dhashes {
		d.Clear()
	}
	var item common.Item
	for i := 0; i < 1000; i++ {
		item.Key = []byte(fmt.Sprint(i))
		item.Value = []byte(fmt.Sprint(i))
		item.Timestamp = 1
		dhashes[0].Put(item)
	}
	common.AssertWithin(t, func() (string, bool) {
		sum := 0
		status := new(bytes.Buffer)
		ordered := dhashAry(dhashes)
		sort.Sort(ordered)
		lastOwned := ordered[len(ordered)-1].Owned()
		ok := true
		for _, d := range ordered {
			sum += d.Owned()
			fmt.Fprintf(status, "%v %v %v\n", d.node.GetBroadcastAddr(), common.HexEncode(d.node.GetPosition()), d.Owned())
			if float64(lastOwned)/float64(d.Owned()) > migrateHysteresis {
				ok = false
			}
			if d.Owned() == 0 {
				ok = false
			}
			lastOwned = d.Owned()
		}
		return string(status.Bytes()), ok && sum == 1000
	}, time.Second*100)
}

func stopServers(servers []*Node) {
	for _, d := range servers {
		d.Stop()
	}
}

func TestDHash(t *testing.T) {
	runtime.GOMAXPROCS(runtime.NumCPU())
	dhashes := testStartup(t, 6, 10191)
	testSync(t, dhashes)
	testClean(t, dhashes)
	testPut(t, dhashes)
	testMigrate(t, dhashes)
}
