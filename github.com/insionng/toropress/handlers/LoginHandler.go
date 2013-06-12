package handlers

import (
	"../libs"
	"../models"
	"../utils"
)

type LoginHandler struct {
	libs.BaseHandler
}

func (self *LoginHandler) Get() {

	sess_username, _ := self.GetSession("username").(string)
	//如果未登录
	if sess_username == "" {
		self.TplNames = "login.html"
		self.Render()
	} else { //如果已登录
		self.Ctx.Redirect(302, "/")
	}

}

func (self *LoginHandler) Post() {
	username := self.GetString("username")
	password := self.GetString("password")

	if username != "" && password != "" {

		if userInfo := models.GetUserByNickname(username); userInfo.Password != "" {

			if utils.Validate_password(userInfo.Password, password) {

				//登录成功设置session
				self.SetSession("userid", userInfo.Id)
				self.SetSession("username", userInfo.Nickname)
				self.SetSession("userrole", userInfo.Role)
				self.SetSession("useremail", userInfo.Email)

				self.Ctx.Redirect(302, "/")
			} else {

				self.Ctx.Redirect(302, "/login")
			}
		} else {

			self.Ctx.Redirect(302, "/login")
		}
	} else {

		self.Ctx.Redirect(302, "/login")
	}
}
