package handlers

import (
	"../libs"
	"../models"
	"strconv"
	"time"
)

type ModifyCategoryHandler struct {
	libs.RootAuthHandler
}

func (self *ModifyCategoryHandler) Get() {
	self.TplNames = "modify_category.html"
	self.Layout = "layout.html"

	self.Render()
}

func (self *ModifyCategoryHandler) Post() {
	inputs := self.Input()
	cid, _ := strconv.Atoi(inputs.Get("categoryid"))

	cat_title := inputs.Get("title")
	cat_content := inputs.Get("content")
	if cid != 0 && cat_title != "" && cat_content != "" {
		var cat models.Category
		cat.Id = int64(cid)
		cat.Title = cat_title
		cat.Content = cat_content
		cat.Created = time.Now()
		models.SaveCategory(cat)
		self.Ctx.Redirect(302, "/category/"+inputs.Get("categoryid"))
	} else {
		self.Ctx.Redirect(302, "/")
	}
}
