package handlers

import (
	"../libs"
	"../models"
)

type NewTopicHandler struct {
	libs.AuthHandler
}

func (self *NewTopicHandler) Get() {
	self.TplNames = "topic_new.html"
	self.Layout = "layout.html"
	self.Data["nodes"] = models.GetAllNode()
	self.Render()
}

func (self *NewTopicHandler) Post() {
	nid, _ := self.GetInt("nodeid")
	cid := models.GetNode(nid).Pid
	uid, _ := self.GetSession("userid").(int64)
	tid_title := self.GetString("title")
	tid_content := self.GetString("content")
	if tid_title != "" && tid_content != "" {
		models.AddTopic(self.GetString("title"), self.GetString("content"), cid, nid, uid)
		self.Ctx.Redirect(302, "/node/"+self.GetString("nodeid"))
	} else {
		self.Ctx.Redirect(302, "/")
	}
}
