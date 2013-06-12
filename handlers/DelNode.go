package handlers

import (
	"../libs"
	"../models"
)

type NodeDeleteHandler struct {
	libs.RootAuthHandler
}

func (self *NodeDeleteHandler) Get() {
	nid, _ := self.GetInt(":nid")
	models.DelNode(nid)
	self.Ctx.Redirect(302, "/")
}
