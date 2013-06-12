package handlers

import (
	"../libs"
	"../models"
	"../utils"
)

type LikeNodeHandler struct {
	libs.BaseHandler
}

func (self *LikeNodeHandler) Get() {
	//inputs := self.Input()
	//id, _ := strconv.Atoi(inputs.Get("id"))
	if utils.IsSpider(self.Ctx.Request.UserAgent()) != true {

		id, _ := self.GetInt(":nid")

		nd := models.GetNode(id)
		nd.Hotup = nd.Hotup + 1
		nd.Hotness = utils.Hotness(nd.Hotup, nd.Hotdown, nd.Created)

		models.SaveNode(nd)

		self.Ctx.WriteString("success")
		//self.Ctx.Redirect(302, "/")

	} else {
		self.Ctx.WriteString("R u spider?")
	}

}
