package discord

import (
	"github.com/zond/god/common"
)

type nodeServer Node

func (self *nodeServer) Notify(caller common.Remote, predecessor *common.Remote) error {
	*predecessor = (*Node)(self).Notify(caller)
	return nil
}
func (self *nodeServer) Nodes(x int, nodes *common.Remotes) error {
	*nodes = (*Node)(self).GetNodes()
	return nil
}
func (self *nodeServer) Ping(ping PingPack, remote *common.Remote) error {
	*remote = (*Node)(self).Ping(ping)
	return nil
}
func (self *nodeServer) GetPredecessor(x int, predecessor *common.Remote) error {
	*predecessor = (*Node)(self).GetPredecessor()
	return nil
}
func (self *nodeServer) GetSuccessorFor(key []byte, successor *common.Remote) error {
	*successor = (*Node)(self).GetSuccessorFor(key)
	return nil
}
