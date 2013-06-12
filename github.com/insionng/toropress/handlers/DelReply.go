package handlers

import (
	"../libs"
	"../models"
)

type DeleteReplyHandler struct {
	libs.RootAuthHandler
}

func (self *DeleteReplyHandler) Get() {
	rid, _ := self.GetInt(":rid")
	models.DelReply(rid)
	self.Ctx.Redirect(302, "/")
}
