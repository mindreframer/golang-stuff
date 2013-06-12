package root

import (
	"../../libs"
	"../../models"
	"../../utils"
)

type RLoginHandler struct {
	libs.BaseHandler
}

func (self *RLoginHandler) Get() {
	browser := ""
	switch {
	case utils.Rex(self.Ctx.Request.UserAgent(), "MSIE 5.0"):
		browser = "MSIE 5.0"
	case utils.Rex(self.Ctx.Request.UserAgent(), "MSIE 5.5"):
		browser = "MSIE 5.5"
	}

	if browser != "" {
		self.Data["BrowserMsgErr"] = "你的浏览器使用的当前内核版本为“" + browser + "”，该版本太古老导致部分功能无法正常使用，建议安装IE8或IE9+以上浏览器,或使用浏览器的非IE内核模式，或选择使用当前最新的主流浏览器IE8、IE9+、Opera10+、Firefox 3.6+、Safari4+和Google Chrome浏览器。"
		self.TplNames = "root/login.html"
		self.Render()
	} else {

		sess_userrole, _ := self.GetSession("userrole").(int64)
		//如果未登录root
		if sess_userrole != -1000 {
			self.TplNames = "root/login.html"
			self.Render()
		} else { //如果已登录root
			self.Ctx.Redirect(302, "/root")
		}
	}

}

func (self *RLoginHandler) Post() {

	username := self.GetString("username")
	password := self.GetString("password")
	if username != "" && password != "" {

		if userInfo := models.GetUserByNickname(username); userInfo.Nickname != "" {

			if utils.Validate_password(userInfo.Password, password) && userInfo.Role == -1000 {

				//登录成功设置session
				self.SetSession("userid", userInfo.Id)
				self.SetSession("username", userInfo.Nickname)
				self.SetSession("userrole", userInfo.Role)
				self.SetSession("useremail", userInfo.Email)

			}
		}

	}
	self.Ctx.Redirect(302, "/root")
}
