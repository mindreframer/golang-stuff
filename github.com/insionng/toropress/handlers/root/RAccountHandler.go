package root

import (
	"../../libs"
	"../../models"
	"../../utils"
)

type RAccountHandler struct {
	libs.RootHandler
}

func (self *RAccountHandler) Get() {
	self.Data["MsgErr"], _ = self.GetSession("MsgErr").(string)
	self.DelSession("MsgErr")
	self.TplNames = "root/account.html"
	self.Render()
}

func (self *RAccountHandler) Post() {
	inputs := self.Input()
	//nickname := inputs.Get("nickname")
	realname := inputs.Get("realname")
	email := inputs.Get("email")
	mobile := inputs.Get("mobile")
	company := inputs.Get("company")
	address := inputs.Get("address")
	uid, _ := self.GetSession("userid").(int64)

	if utils.CheckEmail(email) {
		ur := models.GetUser(uid)
		ur.Email = email
		ur.Mobile = mobile
		ur.Company = company
		ur.Address = address
		ur.Realname = realname

		if e := models.UpdateUser(int(uid), ur); e != nil {
			self.Data["MsgErr"] = "更新资料失败！"
		} else {
			self.Data["MsgErr"] = "更新资料成功！"
			self.SetSession("useremail", email)
		}
	} else {
		self.Data["MsgErr"] = "Email地址有误！"
	}

	self.SetSession("MsgErr", self.Data["MsgErr"])
	self.Ctx.Redirect(302, "/root-account")
}
