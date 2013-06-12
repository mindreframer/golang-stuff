package handlers

import (
	"../libs"
	"../models"
	"../utils"
)

type LikeTopicHandler struct {
	libs.BaseHandler
}

func (self *LikeTopicHandler) Get() {
	//inputs := self.Input()
	//id, _ := strconv.Atoi(inputs.Get("id"))
	if utils.IsSpider(self.Ctx.Request.UserAgent()) != true {

		id, _ := self.GetInt(":tid")

		tp := models.GetTopic(id)
		tp.Hotup = tp.Hotup + 1
		tp.Hotness = utils.Hotness(tp.Hotup, tp.Hotdown, tp.Created)

		models.SaveTopic(tp)

		self.Ctx.WriteString("success")

	} else {
		self.Ctx.WriteString("R u spider?")
	}

}
