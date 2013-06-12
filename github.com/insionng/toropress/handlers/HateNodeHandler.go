package handlers

import (
	"../libs"
	"../models"
	"../utils"
)

type HateNodeHandler struct {
	libs.BaseHandler
}

func (self *HateNodeHandler) Get() {
	//inputs := self.Input()
	//id, _ := strconv.Atoi(inputs.Get("id"))
	if utils.IsSpider(self.Ctx.Request.UserAgent()) != true {

		id, _ := self.GetInt(":nid")

		nd := models.GetNode(id)
		nd.Hotdown = nd.Hotdown + 1
		nd.Hotness = utils.Hotness(nd.Hotup, nd.Hotdown, nd.Created)

		models.SaveNode(nd)

		self.Ctx.WriteString("success")
		//self.Ctx.Redirect(302, "/")

	} else {
		self.Ctx.WriteString("R u spider?")
	}

}
