package discord

import (
	"fmt"
	"github.com/zond/god/common"
	"testing"
	"time"
)

func TestStartup(t *testing.T) {
	firstPort := 9191
	var nodes []*Node
	n := 10
	for i := 0; i < n; i++ {
		nodes = append(nodes, NewNode(fmt.Sprintf("%v:%v", "127.0.0.1", firstPort+i), fmt.Sprintf("%v:%v", "127.0.0.1", firstPort+i)))
	}
	for i := 0; i < n; i++ {
		nodes[i].MustStart()
	}
	for i := 1; i < n; i++ {
		nodes[i].MustJoin(nodes[0].GetBroadcastAddr())
	}
	common.AssertWithin(t, func() (string, bool) {
		routes := make(map[string]bool)
		for i := 0; i < n; i++ {
			routes[nodes[i].Nodes().Describe()] = true
		}
		return fmt.Sprint(routes), len(routes) == 1 && nodes[0].ring.Size() > 0
	}, time.Second*30)
}
