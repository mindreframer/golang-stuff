package handlers

import (
	"../libs"
	"../models"
	"strconv"
	"time"
)

type ModifyNodeHandler struct {
	libs.RootAuthHandler
}

func (self *ModifyNodeHandler) Get() {

	self.Layout = "layout.html"
	self.TplNames = "modify_node.html"
	self.Render()
}

func (self *ModifyNodeHandler) Post() {

	inputs := self.Input()
	cid, _ := strconv.Atoi(inputs.Get("categoryid"))
	nid, _ := strconv.Atoi(inputs.Get("nodeid"))

	nd_title := inputs.Get("title")
	nd_content := inputs.Get("content")
	if cid != 0 && nid != 0 && nd_title != "" && nd_content != "" {
		var nd models.Node
		nd.Id = int64(nid)
		nd.Pid = int64(cid)
		nd.Title = nd_title
		nd.Content = nd_content
		nd.Created = time.Now()
		models.SaveNode(nd)
		self.Ctx.Redirect(302, "/node/"+inputs.Get("nodeid"))
	} else {
		self.Ctx.Redirect(302, "/")
	}
}
