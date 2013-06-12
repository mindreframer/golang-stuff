package root

import (
	"../../libs"
	"../../models"
	"../../utils"
)

type RChangePasswordHandler struct {
	libs.RootHandler
}

func (self *RChangePasswordHandler) Get() {
	self.Data["MsgErr"], _ = self.GetSession("msgerr").(string)
	self.DelSession("MsgErr")
	self.TplNames = "root/change_password.html"
	self.Render()
}

func (self *RChangePasswordHandler) Post() {
	inputs := self.Input()
	oldpassword := inputs.Get("oldpassword")
	newpassword := inputs.Get("newpassword")
	renewpassword := inputs.Get("renewpassword")

	if oldpassword != "" && newpassword != "" && newpassword == renewpassword {

		sess_username, _ := self.GetSession("username").(string)
		usr := models.GetUserByNickname(sess_username)

		if utils.Validate_password(usr.Password, oldpassword) {
			usr.Password = utils.Encrypt_password(newpassword, nil)
			if e := models.SaveUser(usr); e != nil {
				self.Data["MsgErr"] = "更新密码失败！"

			} else {
				self.Data["MsgErr"] = "更新密码成功！"
			}

		}

	} else {
		switch {
		case oldpassword == "":
			self.Data["MsgErr"] = "原密码为空！"
		case newpassword == "":
			self.Data["MsgErr"] = "新密码为空！"
		case newpassword != renewpassword:
			self.Data["MsgErr"] = "两次输入的新密码不一致！"
		default:
			self.Data["MsgErr"] = "提交的信息有误！"
		}
	}

	self.SetSession("msgerr", self.Data["MsgErr"])
	self.Ctx.Redirect(302, "/root/change_password")

}
