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

type RAboutHandler struct {
	libs.RootHandler
}

func (self *RAboutHandler) Get() {
	self.Data["catpage"] = "about"
	tid, _ := self.GetInt(":tid")
	self.Data["topic"] = models.GetTopic(tid)
	self.Data["topics"] = models.GetAllTopicByCid(1, 0, 0, 1, "id")
	self.Data["nodes"] = models.GetAllNodeByCid(1, 0, 0,0, "id")
	self.Data["MsgErr"], _ = self.GetSession("MsgErr").(string)
	self.DelSession("MsgErr")

	switch {
	case self.Ctx.Request.RequestURI == "/root-about":
		//发布内容
		self.Data["asidepage"] = "root_about"
		self.TplNames = "root/about.html"

		self.Render()
	case self.Ctx.Request.RequestURI == "/root-about-topic-list":
		//内容列表
		self.Data["asidepage"] = "root-about-topic-list"
		self.TplNames = "root/about_topic_list.html"

		self.Render()
	case self.Ctx.Request.RequestURI == "/root-about-new-node":
		//创建分类
		self.Data["asidepage"] = "root_about_new_node"
		self.TplNames = "root/about_new_node.html"

		self.Render()
	case self.Ctx.Request.RequestURI == "/root-about-node-list":
		//分类列表
		self.Data["asidepage"] = "root_about_node"
		self.TplNames = "root/about_node.html"

		self.Render()
	}
}

func (self *RAboutHandler) Post() {
	ftitle := self.GetString("ftitle")
	stitle := self.GetString("stitle")
	content := self.GetString("content")
	var nodeid int64 = 1
	var cid int64 = 1
	uid, _ := self.GetSession("userid").(int64)

	msg := ""
	if ftitle == "" {
		msg = msg + "标题不能为空！"
	}
	if nodeid == 0 {
		msg = msg + "必须选择一个上级分类！"
	}
	if content == "" {
		msg = msg + "内容不能为空！"
	}

	self.Data["MsgErr"] = msg

	if msg == "" {
		switch {
		case utils.Rex(self.Ctx.Request.RequestURI, "^/root-about-edit/([0-9]+)$"):
			//编辑POST状态
			tid, _ := self.GetInt(":tid")
			file, handler, e := self.GetFile("image")
			defer file.Close()
			if handler == nil {

				if tid != 0 {
					icontent := models.GetTopic(tid)
					if icontent.Attachment != "" {
						self.Data["file_location"] = icontent.Attachment
					} else {
						self.Data["MsgErr"] = "你还没有选择封面！"
						self.Data["file_location"] = ""
					}
				} else {
					self.Data["MsgErr"] = "你编辑的内容不存在！"
				}
			}
			if nodeid != 0 && ftitle != "" && content != "" && tid != 0 {
				if handler != nil {
					if e != nil {
						self.Data["MsgErr"] = "传输文件过程中产生错误！"
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
						output_size := "232x135"
						output_align := "center"
						utils.Thumbnail(input_file, output_file, output_size, output_align, "black")

						//若文件存在则删除，不存在就当忽略处理
						if self.Data["file_location"] != nil {
							if utils.Exist("." + self.Data["file_location"].(string)) {
								if err := os.Remove("." + self.Data["file_location"].(string)); err != nil {
									self.Data["MsgErr"] = "删除旧文件错误！"
								}
							}
						}
						self.Data["file_location"] = path

					}
				}
				if self.Data["file_location"] != "" {
					if e := models.SetTopic(tid, cid, nodeid, uid, 1, ftitle, content, stitle, self.Data["file_location"].(string)); e != nil {
						self.Data["MsgErr"] = "你提交的修改保存失败，无法写入数据库！"
					} else {
						self.Data["MsgErr"] = "你提交的修改已保存成功！"
					}

				} else {
					self.Data["MsgErr"] = "你提交的内容缺少封面！"
				}
			}
		case self.Ctx.Request.RequestURI == "/root-about":
			//新增内容POST状态
			file, handler, e := self.GetFile("image")
			switch {
			case handler == nil:
				self.Data["MsgErr"] = "你还没有选择封面！"
			case handler != nil && nodeid != 0 && ftitle != "" && content != "":
				//开始添加内容
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
					output_size := "232x135"
					output_align := "center"
					utils.Thumbnail(input_file, output_file, output_size, output_align, "black")
					if e := models.SetTopic(0, cid, nodeid, uid, 1, ftitle, content, stitle, path); e != nil {
						self.Data["MsgErr"] = "添加“" + ftitle + "”失败，无法写入数据库！"
					} else {
						self.Data["MsgErr"] = "添加“" + ftitle + "”成功，你可以继续添加其他内容！"
					}

				}
			}
		}
	} else {

		switch {
		case self.Ctx.Request.RequestURI == "/root-about-new-node":
			//新建分类
			title := self.GetString("title")
			content := self.GetString("content")
			if title != "" {
				models.AddNode(title, content, cid, -1000)
				self.Data["MsgErr"] = "分类“" + title + "”创建成功！"
			} else {
				self.Data["MsgErr"] = "分类标题不能为空！"
			}
		}

	}

	self.SetSession("MsgErr", self.Data["MsgErr"])
	self.Redirect(self.Ctx.Request.RequestURI, 302)

}
