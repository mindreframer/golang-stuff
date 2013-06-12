package handlers

import (
	"../libs"
	"../models"
	"../utils"
	"strconv"
)

type ViewHandler struct {
	libs.BaseHandler
}

func (self *ViewHandler) Get() {
	tid, _ := self.GetInt(":tid")
	tid_handler := models.GetTopic(tid)

	self.TplNames = "view.html"
	self.Layout = "layout.html"

	if tid_handler.Id > 0 {

		tid_handler.Views = tid_handler.Views + 1
		models.UpdateTopic(tid, tid_handler)

		self.Data["article"] = tid_handler
		self.Data["replys"] = models.GetReplyByPid(tid, 0, 0, "id")

		tps := models.GetAllTopicByCid(tid_handler.Cid, 0, 0, 0, "asc")

		if tps != nil && tid != 0 {

			for i, v := range tps {

				if v.Id == tid {
					prev := i - 1
					next := i + 1

					for i, v := range tps {
						if prev == i {
							self.Data["previd"] = v.Id
							self.Data["prev"] = v.Title
						}
						if next == i {
							self.Data["nextid"] = v.Id
							self.Data["next"] = v.Title
						}
					}
				}
			}
		}

		if sess_userrole, _ := self.GetSession("userrole").(int64); sess_userrole == -1000 {
			self.Render()
		} else {
			tid_path := strconv.Itoa(int(tid_handler.Cid)) + "/" + strconv.Itoa(int(tid_handler.Nid)) + "/"
			tid_name := strconv.Itoa(int(tid_handler.Id)) + ".html"
			rs, _ := self.RenderString()
			utils.Writefile("./archives/"+tid_path, tid_name, rs)
			self.Redirect("/archives/"+tid_path+tid_name, 301)
		}
	} else {
		self.Redirect("/", 302)
	}

}
