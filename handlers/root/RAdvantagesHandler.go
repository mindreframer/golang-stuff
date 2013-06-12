package root

import (
	"../../libs"
	"../../models"
	"../../utils"
	"io"
	"os"
	"strings"
	"time"
)

type RAdvantagesHandler struct {
	libs.RootHandler
}

func (self *RAdvantagesHandler) Get() {
	var cid int64 = 2 //優勢属于第二个分类
	self.Data["MsgErr"], _ = self.GetSession("MsgErr").(string)
	self.DelSession("MsgErr")
	self.Data["catpage"] = "advantages"

	switch {

	case utils.Rex(self.Ctx.Request.RequestURI, "^/root-advantages-del/([0-9]+)$"):
		//删除GET状态 删除內容
		tid, _ := self.GetInt(":tid")

		if e := models.DelTopic(tid); e != nil {
			self.Data["MsgErr"] = "删除內容失败！"
		} else {

			self.Data["MsgErr"] = "删除內容成功！"
		}

		self.SetSession("MsgErr", self.Data["MsgErr"])
		self.Redirect("/root-advantages-list", 302)

	case utils.Rex(self.Ctx.Request.RequestURI, "^/root-advantages-edit/([0-9]+)$"):
		//编辑GET状态
		self.Data["asidepage"] = "root_advantages_edit"
		tid, _ := self.GetInt(":tid")
		self.Data["topic"] = models.GetTopic(tid)
		self.TplNames = "root/advantages.html"
		self.Render()

	case self.Ctx.Request.RequestURI == "/root-advantages-list":
		//優勢列表
		self.Data["asidepage"] = "root-advantages-list"
		self.Data["topics"] = models.GetAllTopicByCid(cid, 0, 0, 0, "asc")
		self.TplNames = "root/advantages_list.html"
		self.Render()

	default:
		//设置優勢
		self.Data["asidepage"] = "root_advantages"
		self.TplNames = "root/advantages.html"
		self.Render()
	}
}

func (self *RAdvantagesHandler) Post() {
	var cid int64 = 2

	title := self.GetString("title")
	content := self.GetString("content")
	nodeid, _ := self.GetInt("nodeid")
	uid, _ := self.GetSession("userid").(int64)
	uname, _ := self.GetSession("username").(string)

	tid, _ := self.GetInt(":tid")
	file, handler, e := self.GetFile("image")

	msg := ""
	if title == "" {
		msg = msg + "標題不能为空！"
	}

	if content == "" {
		msg = msg + "内容不能为空！"
	}

	self.Data["MsgErr"] = msg
	if msg == "" {
		switch {
		case utils.Rex(self.Ctx.Request.RequestURI, "^/root-advantages-edit/([0-9]+)$"):
			//编辑POST状态
			if handler == nil {

				if tid != 0 {
					advantages := models.GetTopic(tid)
					if advantages.Attachment != "" {
						self.Data["file_location"] = advantages.Attachment
					} else {
						self.Data["MsgErr"] = "你还没有选择內容封面圖片！"
						self.Data["file_location"] = ""
					}
				} else {
					self.Data["MsgErr"] = "你编辑的內容不存在！"
				}
			}
			if title != "" && content != "" && tid != 0 {
				if handler != nil {
					if e != nil {
						self.Data["MsgErr"] = "传输图片文件过程中产生错误！"
					}

					ext := "." + strings.Split(handler.Filename, ".")[1]
					filename := utils.MD5(time.Now().String()) + ext

					path := "/archives/upload/" + time.Now().Format("2006/01/02/")

					os.MkdirAll("."+path, 0644)
					path = path + filename
					f, err := os.OpenFile("."+path, os.O_WRONLY|os.O_CREATE, 0644)
					defer f.Close()
					if err != nil {
						self.Data["MsgErr"] = "无法打开服务端文件存储路径！"
					} else {
						io.Copy(f, file)
						input_file := "." + path
						output_file := "." + path
						output_size := "248x171"
						output_align := "center"
						utils.Thumbnail(input_file, output_file, output_size, output_align, "white")

						//若文件存在则删除，不存在就当忽略处理
						if self.Data["file_location"] != nil {
							if utils.Exist("." + self.Data["file_location"].(string)) {
								if err := os.Remove("." + self.Data["file_location"].(string)); err != nil {
									self.Data["MsgErr"] = "删除旧形象图片错误！"
								}
							}
						}
						self.Data["file_location"] = path

					}
				}
				if self.Data["file_location"] != "" {
					//保存编辑
					if e := models.SetTopic(tid, cid, nodeid, uid, 0, title, content, uname, self.Data["file_location"].(string)); e != nil {
						self.Data["MsgErr"] = "修改“" + title + "”失败，无法写入数据库！"
					} else {
						self.Data["MsgErr"] = "修改“" + title + "”成功，你可以继续修改其他内容！"
					}

				} else {
					self.Data["MsgErr"] = "你提交的內容缺少封面圖片！"
				}
			}
		case self.Ctx.Request.RequestURI == "/root-advantages":
			//新增內容POST状态
			if handler == nil {
				self.Data["MsgErr"] = "你还没有选择內容封面圖片！"
			}

			if handler != nil && title != "" && content != "" {
				//开始添加內容
				if e != nil {
					self.Data["MsgErr"] = "传输过程文件产生错误！"
				}

				ext := "." + strings.Split(handler.Filename, ".")[1]
				filename := utils.MD5(time.Now().String()) + ext

				path := "/archives/upload/" + time.Now().Format("2006/01/02/")

				os.MkdirAll("."+path, 0644)
				path = path + filename
				f, err := os.OpenFile("."+path, os.O_WRONLY|os.O_CREATE, 0644)
				defer f.Close()
				if err != nil {
					self.Data["MsgErr"] = "无法打开服务端文件存储路径！"
				} else {
					io.Copy(f, file)
					input_file := "." + path
					output_file := "." + path
					output_size := "248x171"
					output_align := "center"
					utils.Thumbnail(input_file, output_file, output_size, output_align, "white")
					if e := models.SetTopic(0, cid, nodeid, uid, 0, title, content, uname, path); e != nil {
						self.Data["MsgErr"] = "添加“" + title + "”失败，无法写入数据库！"
					} else {
						self.Data["MsgErr"] = "添加“" + title + "”成功，你可以继续添加其他案例！"
					}

				}
			}
		}
	}

	self.SetSession("MsgErr", self.Data["MsgErr"])
	self.Redirect(self.Ctx.Request.RequestURI, 302)

}
