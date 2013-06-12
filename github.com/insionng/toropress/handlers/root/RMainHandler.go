package root

import (
	"../../libs"
	//"../../models"
	//"../../utils"
)

type RMainHandler struct {
	libs.RootHandler
}

func (self *RMainHandler) Get() {
	/*
		tpcc_today, tpcc_week, tpcc_month := models.TopicCount()
		cs_categorys, cs_nodes, cs_topics, cs_menbers := models.Counts()

		page := 1
		limit := 5
		pages, pageout, beginnum, endnum, offset := utils.Pages(cs_topics, page, limit)

		self.Data["topic_pagesbar"] = utils.Pagesbar("", cs_topics, pages, pageout, beginnum, endnum, 2)

		self.Data["topic_hotness"] = models.GetAllTopic(offset, limit, "hotness")

		self.Data["tpcc_today"] = tpcc_today
		self.Data["tpcc_week"] = tpcc_week
		self.Data["tpcc_month"] = tpcc_month
		self.Data["cs_categorys"] = cs_categorys
		self.Data["cs_nodes"] = cs_nodes
		self.Data["cs_topics"] = cs_topics
		self.Data["cs_menbers"] = cs_menbers
	*/
	self.Data["catpage"] = "default"
	self.TplNames = "root/index.html"
	self.Render()

}
