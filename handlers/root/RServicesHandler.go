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

type RServicesHandler struct {
	libs.RootHandler
}

func (self *RServicesHandler) Get() {
	var cid int64 = 3
	self.Data["catpage"] = "services"
	self.Data["nodes"] = models.GetAllNodeByCid(cid, 0, 0, 0,"id")
	self.Data["MsgErr"], _ = self.GetSession("MsgErr").(string)
	self.DelSession("MsgErr")

	self.Data["topics"] = models.GetAllTopicByCid(cid, 0, 0, 0, "id")
	switch {
	case utils.Rex(self.Ctx.Request.RequestURI, "^/root-services-del/([0-9]+)$"):
		self.Data["asidepage"] = "root_services_list"
		self.TplNames = "root/services_list.html"

	case self.Ctx.Request.RequestURI == "/root-services-new-node":
		//新建內容分类
		self.Data["asidepage"] = "root_services_new_node"
		self.TplNames = "root/services_new_node.html"

	case self.Ctx.Request.RequestURI == "/root-services-node-list":
		//內容分类列表
		self.Data["asidepage"] = "root_services_node"
		self.TplNames = "root/services_node.html"

	case self.Ctx.Request.RequestURI == "/root-services-topic-list":
		//內容列表
		self.Data["asidepage"] = "root-services-topic-list"
		self.TplNames = "root/services_topic_list.html"

	case self.Ctx.Request.RequestURI == "/root-services":
		//设置內容
		self.Data["asidepage"] = "root_services"
		self.TplNames = "root/services.html"

	}

	self.Render()

}

func (self *RServicesHandler) Post() {
	title := self.GetString("title")
	content := self.GetString("content")
	nodeid, _ := self.GetInt("nodeid")
	var cid int64 = 3
	uid, _ := self.GetSession("userid").(int64)
	sess_username, _ := self.GetSession("username").(string)

	msg := ""
	if title == "" {
		msg = msg + "內容标题不能为空！"
	}
	if nodeid == 0 {
		msg = msg + "必须选择一个上级分类！"
	}
	if content == "" {
		msg = msg + "內容内容不能为空！"
	}

	self.Data["MsgErr"] = msg

	if msg == "" {
		switch {
		case utils.Rex(self.Ctx.Request.RequestURI, "^/root-services-edit/([0-9]+)$"):
			//编辑POST状态
			tid, _ := self.GetInt(":tid")
			file, handler, e := self.GetFile("image")
			defer file.Close()
			if handler == nil {

				if tid != 0 {
					item := models.GetTopic(tid)
					if item.Attachment != "" {
						self.Data["file_location"] = item.Attachment
					} else {
						self.Data["MsgErr"] = "你还没有选择內容封面！"
						self.Data["file_location"] = ""
					}
				} else {
					self.Data["MsgErr"] = "你编辑的內容不存在！"
				}
			}

			if nodeid != 0 && title != "" && content != "" && tid != 0 {
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
						output_size := "211x134"
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
					if e := models.SetTopic(tid, cid, nodeid, uid, 1, title, content, sess_username, self.Data["file_location"].(string)); e != nil {
						self.Data["MsgErr"] = "你提交的修改保存失败，无法写入数据库！"
					} else {
						self.Data["MsgErr"] = "你提交的修改已保存成功！"
					}

				} else {
					self.Data["MsgErr"] = "你提交的內容缺少封面！"
				}
			}
		case self.Ctx.Request.RequestURI == "/root-services":
			//新增內容POST状态
			file, handler, e := self.GetFile("image")
			switch {
			case handler == nil:
				self.Data["MsgErr"] = "你还没有选择內容封面！"
			case handler != nil && nodeid != 0 && title != "" && content != "":
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
					output_size := "211x134"
					output_align := "center"
					utils.Thumbnail(input_file, output_file, output_size, output_align, "white")
					if e := models.SetTopic(0, cid, nodeid, uid, 1, title, content, sess_username, path); e != nil {
						self.Data["MsgErr"] = "添加內容“" + title + "”失败，无法写入数据库！"
					} else {
						self.Data["MsgErr"] = "添加內容“" + title + "”成功，你可以继续添加其他內容！"
					}

				}
			}
		}
	} else {

		switch {
		case self.Ctx.Request.RequestURI == "/root-services-new-node":
			//新建內容分类
			title := self.GetString("title")
			content := self.GetString("content")
			if title != "" {
				models.AddNode(title, content, 3, -1000)
				self.Data["MsgErr"] = "內容分类“" + title + "”创建成功！"
			} else {
				self.Data["MsgErr"] = "內容分类标题不能为空！"
			}
		}

	}

	self.SetSession("MsgErr", self.Data["MsgErr"])
	self.Redirect(self.Ctx.Request.RequestURI, 302)

}
