package root

import (
	"../../libs"
	"../../models"
	"../../utils"
)

type RSettingHandler struct {
	libs.RootHandler
}

func (self *RSettingHandler) Get() {
	self.Data["catpage"] = "setting"
	self.Data["MsgErr"], _ = self.GetSession("MsgErr").(string)
	self.DelSession("MsgErr")
	switch {
	case utils.Rex(self.Ctx.Request.RequestURI, "^/root-setting-setroot-del/([0-9]+)$"):
		rootid, _ := self.GetInt(":rid")
		if e := models.DelUser(rootid); e != nil {
			self.Data["MsgErr"] = "删除管理员失败！"
		} else {
			self.Data["MsgErr"] = "删除管理员成功！"
		}

		self.TplNames = "root/setting_setroot.html"

	case utils.Rex(self.Ctx.Request.RequestURI, "^/root-setting-setroot-edit/([0-9]+)$"):
		self.Data["asidepage"] = "root_setting_setroot_edit"
		rootid, _ := self.GetInt(":rid")
		self.Data["root"] = models.GetUser(rootid)
		self.TplNames = "root/setting_setroot.html"

	case self.Ctx.Request.RequestURI == "/root-setting-setroot":
		self.Data["asidepage"] = "root_setting_setroot"
		self.Data["roots"] = models.GetAllUserByRole(-1000)
		self.TplNames = "root/setting_setroot.html"

	case self.Ctx.Request.RequestURI == "/root-setting-password":
		self.Data["asidepage"] = "root_setting_password"
		self.TplNames = "root/setting_password.html"

	case self.Ctx.Request.RequestURI == "/root-setting":
		self.Data["asidepage"] = "root_setting"
		self.TplNames = "root/setting.html"
	}

	self.Render()
}

func (self *RSettingHandler) Post() {
	switch {
	case self.Ctx.Request.RequestURI == "/root-setting-setroot":
		//设置管理员
		newroot := self.GetString("newroot")
		realname := self.GetString("realname")
		curpassword := self.GetString("curpassword")
		newpassword := self.GetString("newpassword")
		repassword := self.GetString("repassword")
		if newroot != "" && realname != "" && curpassword != "" && repassword != "" && newpassword == repassword {
			sess_username, _ := self.GetSession("username").(string)
			usr := models.GetUserByNickname(sess_username)

			if utils.Validate_password(usr.Password, curpassword) {

				if e := models.AddUser("", newroot, realname, utils.Encrypt_password(newpassword, nil), -1000); e != nil {
					self.Data["MsgErr"] = "添加新管理员“" + newroot + "”失败！"

				} else {
					self.Data["MsgErr"] = "添加新管理员“" + newroot + "”成功！"
				}

			} else {

				self.Data["MsgErr"] = "当前密码不正确！"
			}
		} else {
			msg := ""
			if curpassword == "" {
				msg = msg + "当前管理员密码不能为空！"
			}

			if newpassword == "" {
				msg = msg + "新增管理员密码不能为空！"
			}

			if repassword == "" {
				msg = msg + "新增管理员确认密码不能为空！"
			}

			if newpassword != repassword {
				msg = msg + "两次输入的新增管理员密码不一致！"
			}

			self.Data["MsgErr"] = msg
		}

	case self.Ctx.Request.RequestURI == "/root-setting-password":
		//密码修改
		oldpassword := self.GetString("oldpassword")
		newpassword := self.GetString("newpassword")
		repassword := self.GetString("repassword")

		if oldpassword != "" && repassword != "" && newpassword == repassword {
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
			msg := ""
			if oldpassword == "" {
				msg = msg + "原密码不能为空！"
			}

			if newpassword == "" {
				msg = msg + "新密码不能为空！"
			}

			if repassword == "" {
				msg = msg + "请输入确认密码！"
			}

			if newpassword != repassword {
				msg = msg + "两次输入的新密码不一致！"
			}

			self.Data["MsgErr"] = msg
		}

	case self.Ctx.Request.RequestURI == "/root-setting":
		//常规设置 POST
		title := self.GetString("title")
		title_en := self.GetString("title_en")
		keywords := self.GetString("keywords")
		description := self.GetString("description")
		company := self.GetString("company")
		copyright := self.GetString("copyright")
		site_email := self.GetString("site_email")
		tweibo := self.GetString("tweibo")
		sweibo := self.GetString("sweibo")
		statistics := self.GetString("statistics")

		models.SetKV("title", title)
		models.SetKV("title_en", title_en)
		models.SetKV("keywords", keywords)
		models.SetKV("description", description)

		models.SetKV("company", company)
		models.SetKV("copyright", copyright)
		models.SetKV("site_email", site_email)

		models.SetKV("tweibo", tweibo)
		models.SetKV("sweibo", sweibo)

		models.SetKV("statistics", statistics)

	}

	self.SetSession("MsgErr", self.Data["MsgErr"])
	self.Redirect(self.Ctx.Request.RequestURI, 302)
}
