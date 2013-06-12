package dhash

import (
	"github.com/zond/god/timenet"
)

type dhashPeerProducer Node

func (self *dhashPeerProducer) Peers() (result map[string]timenet.Peer) {
	result = make(map[string]timenet.Peer)
	for _, node := range (*Node)(self).node.GetNodes() {
		result[node.Addr] = (remotePeer)(node)
	}
	return
}
