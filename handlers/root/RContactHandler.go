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

type RContactHandler struct {
	libs.RootHandler
}

func (self *RContactHandler) Get() {
	var cid int64 = 5
	self.Data["catpage"] = "contact"
	self.Data["topics"] = models.GetAllTopicByCid(cid, 0, 0, 0, "id")
	self.Data["nodes"] = models.GetAllNodeByCid(cid, 0, 0,0, "id")
	self.Data["MsgErr"], _ = self.GetSession("MsgErr").(string)
	self.DelSession("MsgErr")

	switch {
	case self.Ctx.Request.RequestURI == "/root-contact":
		self.Data["asidepage"] = "root_contact"
		self.TplNames = "root/contact.html"
	case self.Ctx.Request.RequestURI == "/root-contact-node-list":
		self.Data["asidepage"] = "root_contact_node"
		self.TplNames = "root/contact_node.html"
	case self.Ctx.Request.RequestURI == "/root-contact-new-node":
		self.Data["asidepage"] = "root_contact_new_node"
		self.TplNames = "root/contact_new_node.html"

	case self.Ctx.Request.RequestURI == "/root-contact-topic-list":
		self.Data["asidepage"] = "root_contact_topic_list"
		self.TplNames = "root/contact_topic_list.html"
	}

	self.Render()
}

func (self *RContactHandler) Post() {
	var cid int64 = 5
	title := self.GetString("title")
	content := self.GetString("content")
	nodeid, _ := self.GetInt("nodeid")
	uid, _ := self.GetSession("userid").(int64)
	sess_username, _ := self.GetSession("username").(string)

	msg := ""
	if title == "" {
		msg = msg + "内容标题不能为空！"
	}
	if nodeid == 0 {
		msg = msg + "必须选择一个上级分类！"
	}
	if content == "" {
		msg = msg + "内容内容不能为空！"
	}

	self.Data["MsgErr"] = msg
	if self.Ctx.Request.RequestURI == "/root-contact-new-node" {

		//新建内容分类
		if title != "" {
			if e := models.AddNode(title, content, cid, 1); e != nil {

				self.Data["MsgErr"] = "内容分类无法保存到数据库！"
			} else {

				self.Data["MsgErr"] = "内容分类成功保存到数据库！"
			}

			if content == "" {
				self.Data["MsgErr"] = ""
			}
			self.Data["MsgErr"] = "内容分类“" + title + "”创建成功！"
		} else {
			self.Data["MsgErr"] = "内容分类标题不能为空！"
		}
	} else {

		if msg == "" {
			switch {
			case utils.Rex(self.Ctx.Request.RequestURI, "^/root-contact-edit/([0-9]+)$"):
				//编辑POST状态
				tid, _ := self.GetInt(":tid")
				file, handler, e := self.GetFile("image")
				defer file.Close()
				if handler == nil {

					if tid != 0 {
						icase := models.GetTopic(tid)
						if icase.Attachment != "" {
							self.Data["file_location"] = icase.Attachment
						} else {
							self.Data["MsgErr"] = "你还没有选择内容封面！"
							self.Data["file_location"] = ""
						}
					} else {
						self.Data["MsgErr"] = "你编辑的内容不存在！"
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
							output_size := "214x335"
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
						//保存编辑
						if e := models.SetTopic(tid, cid, nodeid, uid, 0, title, content, sess_username, self.Data["file_location"].(string)); e != nil {
							self.Data["MsgErr"] = "你提交的修改保存失败，无法写入数据库！"
						} else {
							self.Data["MsgErr"] = "你提交的修改已保存成功！"
						}

					} else {
						self.Data["MsgErr"] = "你提交的内容缺少封面！"
					}
				}
			case self.Ctx.Request.RequestURI == "/root-contact":
				//新增内容POST状态
				file, handler, e := self.GetFile("image")
				switch { /*
					case handler == nil:
						self.Data["MsgErr"] = "你还没有选择封面！"*/
				case /* handler != nil &&*/ nodeid != 0 && title != "" && content != "":
					//开始添加内容
					if e != nil {
						self.Data["MsgErr"] = "传输过程文件产生错误！"
					}
					if handler != nil {

						ext := "." + strings.Split(handler.Filename, ".")[1]
						filename := utils.MD5(time.Now().String()) + ext

						path := "/archives/upload/" + time.Now().Format("2006/01/02/")

						os.MkdirAll("."+path, 0644)
						path = path + filename
						self.Data["file_localtion"] = path

					}

					uid, _ := self.GetSession("userid").(int64)
					username, _ := self.GetSession("username").(string)
					path := ""
					if self.Data["file_localtion"] != nil {
						path := self.Data["file_localtion"].(string)
						f, err := os.OpenFile("."+path, os.O_WRONLY|os.O_CREATE, 0644)
						defer f.Close()

						if err != nil {
							self.Data["MsgErr"] = "无法打开服务端文件存储路径！"
						} else {

							io.Copy(f, file)
							input_file := "." + path
							output_file := "." + path
							output_size := "288x180"
							output_align := "center"
							utils.Thumbnail(input_file, output_file, output_size, output_align, "black")
						}
					}
					if e := models.SetTopic(0, cid, nodeid, uid, 0, title, content, username, path); e != nil {
						self.Data["MsgErr"] = "添加“" + title + "”失败，无法写入数据库！"
					} else {
						self.Data["MsgErr"] = "添加“" + title + "”成功，你可以继续添加其他内容！"
					}

				}
			}
		}
	}
	self.SetSession("MsgErr", self.Data["MsgErr"])
	self.Redirect(self.Ctx.Request.RequestURI, 302)

}
