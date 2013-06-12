package libs

import (
	"../models"
	"../utils"
	"github.com/insionng/torgo"
	//"github.com/astaxie/beego"
	//"../torgo"
	"runtime"
	"time"
)

var (
	sess_username string
	sess_uid      int64
	sess_role     int64
	sess_email    string

	bc *torgo.BeeCache
)

type BaseHandler struct {
	torgo.Handler
	//beego.Controller
}

type AuthHandler struct {
	BaseHandler
}

type RootAuthHandler struct {
	BaseHandler
}

type RootHandler struct {
	BaseHandler
}

func init() {
	bc = torgo.NewBeeCache()
	bc.Every = 259200 //該單位為秒，0為不過期，259200 三天,604800 即一個星期清空一次緩存
	bc.Start()
}

//用户等级划分：正数是普通用户，负数是管理员各种等级划分，为0则尚未注册
func (self *BaseHandler) Prepare() {
	sess_username, _ = self.GetSession("username").(string)
	sess_uid, _ = self.GetSession("userid").(int64)
	sess_role, _ = self.GetSession("userrole").(int64)
	sess_email, _ = self.GetSession("useremail").(string)

	if sess_role == 0 {
		self.Data["Userid"] = 0
		self.Data["Username"] = ""
		self.Data["Userrole"] = 0
		self.Data["Useremail"] = ""
	} else {
		self.Data["Userid"] = sess_uid
		self.Data["Username"] = sess_username
		self.Data["Userrole"] = sess_role
		self.Data["Useremail"] = sess_email
	}
	self.Data["categorys"] = models.GetAllCategory()
	self.Data["nodes"] = models.GetAllNode()
	self.Data["topics_5s"] = models.GetAllTopic(0, 5, "id")
	self.Data["topics_10s"] = models.GetAllTopic(0, 10, "id")
	self.Data["nodes_10s"] = models.GetAllNodeByCid(0, 0, 10, 0, "id")
	self.Data["replys_5s"] = models.GetReplyByPid(0, 0, 5, "id")
	self.Data["replys_10s"] = models.GetReplyByPid(0, 0, 10, "id")

	self.Data["author"] = models.GetKV("author")
	self.Data["title"] = models.GetKV("title")
	self.Data["title_en"] = models.GetKV("title_en")
	self.Data["keywords"] = models.GetKV("keywords")
	self.Data["description"] = models.GetKV("description")

	self.Data["company"] = models.GetKV("company")
	self.Data["copyright"] = models.GetKV("copyright")
	self.Data["site_email"] = models.GetKV("site_email")

	self.Data["tweibo"] = models.GetKV("tweibo")
	self.Data["sweibo"] = models.GetKV("sweibo")
	self.Data["timenow"] = time.Now()
	self.Data["statistics"] = models.GetKV("statistics")

}

//会员或管理员前台权限认证
func (self *AuthHandler) Prepare() {
	self.BaseHandler.Prepare()

	if sess_role == 0 {
		self.Ctx.Redirect(302, "/login")
	}
}

//管理员前台权限认证
func (self *RootAuthHandler) Prepare() {
	self.BaseHandler.Prepare()
	if sess_role != -1000 {
		self.Ctx.Redirect(302, "/login")
	}
}

//管理员后台后台认证
func (self *RootHandler) Prepare() {
	self.BaseHandler.Prepare()

	if !utils.IsSpider(self.Ctx.Request.UserAgent()) {
		if sess_role != -1000 {
			self.Ctx.Redirect(302, "/root-login")
		} else {
			self.Data["remoteproto"] = self.Ctx.Request.Proto
			self.Data["remotehost"] = self.Ctx.Request.Host
			self.Data["remoteos"] = runtime.GOOS
			self.Data["remotearch"] = runtime.GOARCH
			self.Data["remotecpus"] = runtime.NumCPU()
			self.Data["golangver"] = runtime.Version()
		}
	} else {
		self.Ctx.Redirect(302, "/")
	}
}

func (self *BaseHandler) Render() (err error) {

	var ivalue []byte
	ck, _ := self.Ctx.Request.Cookie("lang")
	lang := ""

	if ck != nil {
		lang = ck.Value
	} else {
		lang = "normal"
	}

	if self.GetString("lang") != "" {

		if self.GetString("lang") == "normal" {
			lang = "normal"
		}

		if self.GetString("lang") == "cn" {
			lang = "zh-cn"
		}

		if self.GetString("lang") == "hk" {
			lang = "zh-hk"
		}

	}

	self.Ctx.SetCookie("lang", lang, "", "", 0)
	self.Data["lang"] = lang

	rb, e := self.RenderBytes()
	rs := string(rb)
	ikey := utils.MD5(rs + lang)
	if bc.IsExist(ikey) {
		ivalue = bc.Get(ikey).([]byte)
	} else {

		if lang == "normal" {
			ivalue = rb
		} else {
			ivalue = utils.Convzh(rs, lang)
		}

		bc.Put(ikey, ivalue, 259200)

	}

	return self.RenderCore(ivalue, e)

}

func (self *RootHandler) Render() (err error) {
	return self.BaseHandler.Render()
}
