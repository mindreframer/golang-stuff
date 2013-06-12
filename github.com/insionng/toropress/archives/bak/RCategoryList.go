package root

import (
	"../../libs"
	//"../../models"
	//"../../utils"

)

type RCategoryListHandler struct {
	libs.RootHandler
}

func (self *RCategoryListHandler) Get() {
	self.TplNames = "root/category_list.html"
	self.Render()
}
