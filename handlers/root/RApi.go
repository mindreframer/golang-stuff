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

const (
	outtimesz = "<h1 style=\"text-align:center;\">请你登录后操作，当前已超时！</h1>"
)

//RAPI 控制器 没有引入权限验证基类，注意手动验证用户权限，以免造成人为漏洞
//这是为了实现在用户超时的情况下返回一个正常提示，而不是因为没有权限而直接跳转~
type RApi struct {
	libs.BaseHandler
}

func (self *RApi) Get() {

	switch {
	case utils.Rex(self.Ctx.Request.RequestURI, "^/root-about-gallery-edit/([0-9]+)$"):
		//# Gallery编辑GET状态
		if sess_role, _ := self.GetSession("userrole").(int64); sess_role != -1000 {
			self.Ctx.WriteString(outtimesz)
		} else {

			editmid, _ := self.GetInt(":editmid")
			img := models.GetFile(editmid)
			self.Data["img"] = img
			self.TplNames = "root/gallery_editurl.html"
			self.Render()
		}

	case utils.Rex(self.Ctx.Request.RequestURI, "^/root-about-gallery-del/([0-9]+)$"):
		//# Gallery删除GET状态
		if sess_role, _ := self.GetSession("userrole").(int64); sess_role != -1000 {
			self.Ctx.Redirect(302, "/root-login")
		} else {
			if mid, _ := self.GetInt(":delmid"); mid != 0 {
				if e := models.DelFile(mid); e != nil {
					self.Data["MsgErr"] = "删除图片文件失败！"
				} else {
					self.Data["MsgErr"] = "成功删除图片文件！"
				}

			} else {
				self.Data["MsgErr"] = "错误对象！"
			}
			self.SetSession("MsgErr", self.Data["MsgErr"])
			self.Ctx.Redirect(302, "/root-gallery")
		}
	case utils.Rex(self.Ctx.Request.RequestURI, "^/root-about-node-del/([0-9]+)$"):
		//#ABOUT 节点删除 GET状态
		if sess_role, _ := self.GetSession("userrole").(int64); sess_role != -1000 {
			self.Ctx.Redirect(302, "/root-login")
		} else {
			if mid, _ := self.GetInt(":nid"); mid != 0 {
				if e := models.DelNode(mid); e != nil {
					self.Data["MsgErr"] = "删除节点失败！"
				} else {
					self.Data["MsgErr"] = "成功删除节点！"
				}

			} else {
				self.Data["MsgErr"] = "错误节点！"
			}
			self.SetSession("MsgErr", self.Data["MsgErr"])
			self.Redirect("/root-about-node-list", 302)
		}
	case utils.Rex(self.Ctx.Request.RequestURI, "^/root-contact-node-del/([0-9]+)$"):
		//#contact 节点删除 GET状态
		if sess_role, _ := self.GetSession("userrole").(int64); sess_role != -1000 {
			self.Ctx.Redirect(302, "/root-login")
		} else {
			if mid, _ := self.GetInt(":nid"); mid != 0 {
				if e := models.DelNode(mid); e != nil {
					self.Data["MsgErr"] = "删除节点失败！"
				} else {
					self.Data["MsgErr"] = "成功删除节点！"
				}
				self.SetSession("MsgErr", self.Data["MsgErr"])
				self.Redirect("/root-contact-node-list", 302)

			}
		}
	case utils.Rex(self.Ctx.Request.RequestURI, "^/root-about-topic-del/([0-9]+)$"):
		//#ABOUT 内容删除GET状态
		if sess_role, _ := self.GetSession("userrole").(int64); sess_role != -1000 {
			self.Ctx.Redirect(302, "/root-login")
		} else {
			if mid, _ := self.GetInt(":tid"); mid != 0 {
				if e := models.DelTopic(mid); e != nil {
					self.Data["MsgErr"] = "删除内容失败！"
				} else {
					self.Data["MsgErr"] = "成功删除内容！"
				}
				self.SetSession("MsgErr", self.Data["MsgErr"])
				self.Redirect("/root-about-topic-list", 302)

			}
		}
	case utils.Rex(self.Ctx.Request.RequestURI, "^/root-crafts-topic-del/([0-9]+)$"):
		// crafts 内容删除GET状态
		if sess_role, _ := self.GetSession("userrole").(int64); sess_role != -1000 {
			self.Redirect("/root-login", 302)
		} else {
			if mid, _ := self.GetInt(":tid"); mid != 0 {
				if e := models.DelTopic(mid); e != nil {
					self.Data["MsgErr"] = "删除内容失败！"
				} else {
					self.Data["MsgErr"] = "成功删除内容！"
				}
				self.SetSession("MsgErr", self.Data["MsgErr"])
				self.Redirect("/root-crafts-topic-list", 302)

			}
		}
	case utils.Rex(self.Ctx.Request.RequestURI, "^/root-services-topic-del/([0-9]+)$"):
		// services 内容删除GET状态
		if sess_role, _ := self.GetSession("userrole").(int64); sess_role != -1000 {
			self.Redirect("/root-login", 302)
		} else {
			if mid, _ := self.GetInt(":tid"); mid != 0 {
				if e := models.DelTopic(mid); e != nil {
					self.Data["MsgErr"] = "删除内容失败！"
				} else {
					self.Data["MsgErr"] = "成功删除内容！"
				}
				self.SetSession("MsgErr", self.Data["MsgErr"])
				self.Redirect("/root-services-topic-list", 302)

			}
		}
	case utils.Rex(self.Ctx.Request.RequestURI, "^/root-about-topic-edit/([0-9]+)$"):
		//ABOUT 内容编辑GET状态
		if sess_role, _ := self.GetSession("userrole").(int64); sess_role != -1000 {
			self.Redirect("/root-login", 302)
		} else {
			self.Data["MsgErr"], _ = self.GetSession("MsgErr").(string)
			self.DelSession("MsgErr")
			tid, _ := self.GetInt(":tid")
			tid_handler := models.GetTopic(tid)
			self.Data["catpage"] = "about"
			self.Data["asidepage"] = "root_about_topic_edit"
			self.Data["topic"] = tid_handler

			self.TplNames = "root/about.html"
			self.Render()
		}
	case utils.Rex(self.Ctx.Request.RequestURI, "^/root-services-topic-edit/([0-9]+)$"):
		// services 内容编辑 GET状态
		if sess_role, _ := self.GetSession("userrole").(int64); sess_role != -1000 {
			self.Redirect("/root-login", 302)
		} else {
			var cid int64 = 3
			self.Data["MsgErr"], _ = self.GetSession("MsgErr").(string)
			self.DelSession("MsgErr")

			tid, _ := self.GetInt(":tid")
			tid_handler := models.GetTopic(tid)
			self.Data["catpage"] = "services"
			self.Data["asidepage"] = "root-services-topic-edit"
			self.Data["topic"] = tid_handler
			self.Data["inode"] = models.GetNode(tid_handler.Nid)
			self.Data["nodes"] = models.GetAllNodeByCid(cid, 0, 0,0, "id")
			self.TplNames = "root/published.html"
			self.Render()
		}
	case utils.Rex(self.Ctx.Request.RequestURI, "^/root-contact-topic-edit/([0-9]+)$"):
		//CONTACT 内容编辑GET状态
		if sess_role, _ := self.GetSession("userrole").(int64); sess_role != -1000 {
			self.Redirect("/root-login", 302)
		} else {
			var cid int64 = 6
			self.Data["MsgErr"], _ = self.GetSession("MsgErr").(string)
			self.DelSession("MsgErr")
			tid, _ := self.GetInt(":tid")
			tid_handler := models.GetTopic(tid)
			self.Data["asidepage"] = "root-contact-topic-edit"
			self.Data["topic"] = tid_handler
			self.Data["inode"] = models.GetNode(tid_handler.Nid)
			self.Data["catpage"] = "contact"
			self.Data["nodes"] = models.GetAllNodeByCid(cid, 0, 0,0, "id")
			self.TplNames = "root/contact.html"
			self.Render()
		}
	case utils.Rex(self.Ctx.Request.RequestURI, "^/root-crafts-topic-edit/([0-9]+)$"):
		//crafts 内容编辑GET状态
		if sess_role, _ := self.GetSession("userrole").(int64); sess_role != -1000 {
			self.Redirect("/root-login", 302)
		} else {
			var cid int64 = 4
			self.Data["MsgErr"], _ = self.GetSession("MsgErr").(string)
			self.DelSession("MsgErr")

			tid, _ := self.GetInt(":tid")
			tid_handler := models.GetTopic(tid)
			self.Data["catpage"] = "crafts"
			self.Data["asidepage"] = "root-crafts-topic-edit"
			self.Data["topic"] = tid_handler
			self.Data["inode"] = models.GetNode(tid_handler.Nid)
			self.Data["nodes"] = models.GetAllNodeByCid(cid, 0, 0,0, "id")
			self.TplNames = "root/published.html"
			self.Render()
		}
	case utils.Rex(self.Ctx.Request.RequestURI, "^/root-information-topic-del/([0-9]+)$"):
		//information 内容删除GET状态
		if sess_role, _ := self.GetSession("userrole").(int64); sess_role != -1000 {
			self.Ctx.Redirect(302, "/root-login")
		} else {
			if mid, _ := self.GetInt(":tid"); mid != 0 {
				if e := models.DelTopic(mid); e != nil {
					self.Data["MsgErr"] = "删除内容失败！"
				} else {
					self.Data["MsgErr"] = "成功删除内容！"
				}
				self.SetSession("MsgErr", self.Data["MsgErr"])
				self.Redirect("/root-information-topic-list", 302)

			}
		}
	case utils.Rex(self.Ctx.Request.RequestURI, "^/root-contact-topic-del/([0-9]+)$"):
		//contact 内容删除GET状态
		if sess_role, _ := self.GetSession("userrole").(int64); sess_role != -1000 {
			self.Ctx.Redirect(302, "/root-login")
		} else {
			if mid, _ := self.GetInt(":tid"); mid != 0 {
				if e := models.DelTopic(mid); e != nil {
					self.Data["MsgErr"] = "删除内容失败！"
				} else {
					self.Data["MsgErr"] = "成功删除内容！"
				}
				self.SetSession("MsgErr", self.Data["MsgErr"])
				self.Redirect("/root-contact-topic-list", 302)

			}
		}
	case utils.Rex(self.Ctx.Request.RequestURI, "^/root-about-node-edit/([0-9]+)$"):
		// about 节点编辑GET状态
		if sess_role, _ := self.GetSession("userrole").(int64); sess_role != -1000 {
			self.Ctx.Redirect(302, "/root-login")
		} else {
			nid, _ := self.GetInt(":nid")
			//var cid int64 = 6
			self.Data["MsgErr"], _ = self.GetSession("MsgErr").(string)
			self.DelSession("MsgErr")

			self.Data["catpage"] = "about"
			self.Data["asidepage"] = "root_about_node_edit"
			self.Data["node"] = models.GetNode(nid)
			self.TplNames = "root/published_node.html"
			self.Render()
		}

	case utils.Rex(self.Ctx.Request.RequestURI, "^/root-contact-node-edit/([0-9]+)$"):
		//contact 节点编辑GET状态
		if sess_role, _ := self.GetSession("userrole").(int64); sess_role != -1000 {
			self.Ctx.Redirect(302, "/root-login")
		} else {
			nid, _ := self.GetInt(":nid")
			//var cid int64 = 6
			self.Data["MsgErr"], _ = self.GetSession("MsgErr").(string)
			self.DelSession("MsgErr")

			self.Data["catpage"] = "contact"
			self.Data["asidepage"] = "root_contact_node_edit"
			self.Data["node"] = models.GetNode(nid)
			self.TplNames = "root/published_node.html"
			self.Render()
		}

	case utils.Rex(self.Ctx.Request.RequestURI, "^/root-services-node-edit/([0-9]+)$"):
		//services 节点编辑GET状态
		if sess_role, _ := self.GetSession("userrole").(int64); sess_role != -1000 {
			self.Ctx.Redirect(302, "/root-login")
		} else {
			nid, _ := self.GetInt(":nid")
			//var cid int64 = 6
			self.Data["MsgErr"], _ = self.GetSession("MsgErr").(string)
			self.DelSession("MsgErr")

			self.Data["catpage"] = "services"
			self.Data["asidepage"] = "root_services_node_edit"
			self.Data["node"] = models.GetNode(nid)
			self.TplNames = "root/published_node.html"
			self.Render()
		}
	case utils.Rex(self.Ctx.Request.RequestURI, "^/root-information-node-edit/([0-9]+)$"):
		//information 节点编辑GET状态
		if sess_role, _ := self.GetSession("userrole").(int64); sess_role != -1000 {
			self.Ctx.Redirect(302, "/root-login")
		} else {
			nid, _ := self.GetInt(":nid")
			self.Data["MsgErr"], _ = self.GetSession("MsgErr").(string)
			self.DelSession("MsgErr")

			self.Data["catpage"] = "information"
			self.Data["asidepage"] = "root-information-node-edit"
			self.Data["node"] = models.GetNode(nid)
			self.TplNames = "root/published_node.html"
			self.Render()
		}
	case utils.Rex(self.Ctx.Request.RequestURI, "^/root-crafts-node-edit/([0-9]+)$"):
		//crafts节点编辑GET状态
		if sess_role, _ := self.GetSession("userrole").(int64); sess_role != -1000 {
			self.Ctx.Redirect(302, "/root-login")
		} else {
			nid, _ := self.GetInt(":nid")
			self.Data["MsgErr"], _ = self.GetSession("MsgErr").(string)
			self.DelSession("MsgErr")

			self.Data["catpage"] = "crafts"
			self.Data["asidepage"] = "root-crafts-node-edit"
			self.Data["node"] = models.GetNode(nid)
			self.TplNames = "root/published_node.html"
			self.Render()
		}
	case utils.Rex(self.Ctx.Request.RequestURI, "^/root-information-topic-edit/([0-9]+)$"):
		//information 内容编辑GET状态
		if sess_role, _ := self.GetSession("userrole").(int64); sess_role != -1000 {
			self.Ctx.Redirect(302, "/root-login")
		} else {
			var cid int64 = 5
			self.Data["MsgErr"], _ = self.GetSession("MsgErr").(string)
			self.DelSession("MsgErr")

			tid, _ := self.GetInt(":tid")
			tid_handler := models.GetTopic(tid)
			self.Data["catpage"] = "information"
			self.Data["asidepage"] = "root-information-topic-edit"
			self.Data["topic"] = tid_handler
			self.Data["inode"] = models.GetNode(tid_handler.Nid)
			self.Data["nodes"] = models.GetAllNodeByCid(cid, 0, 0, 0,"id")
			self.TplNames = "root/published.html"
			self.Render()
		}
	case utils.Rex(self.Ctx.Request.RequestURI, "^/root-services-node-del/([0-9]+)$"):
		//services 节点删除 GET状态
		if sess_role, _ := self.GetSession("userrole").(int64); sess_role != -1000 {
			self.Ctx.Redirect(302, "/root-login")
		} else {
			if mid, _ := self.GetInt(":nid"); mid != 0 {
				if e := models.DelNode(mid); e != nil {
					self.Data["MsgErr"] = "删除节点失败！"
				} else {
					self.Data["MsgErr"] = "成功删除节点！"
				}
				self.SetSession("MsgErr", self.Data["MsgErr"])
				self.Redirect("/root-services-node-list", 302)

			}
		}
	case utils.Rex(self.Ctx.Request.RequestURI, "^/root-information-node-del/([0-9]+)$"):
		//#INFORMATION 节点删除 GET状态
		if sess_role, _ := self.GetSession("userrole").(int64); sess_role != -1000 {
			self.Ctx.Redirect(302, "/root-login")
		} else {
			if mid, _ := self.GetInt(":nid"); mid != 0 {
				if e := models.DelNode(mid); e != nil {
					self.Data["MsgErr"] = "删除节点失败！"
				} else {
					self.Data["MsgErr"] = "成功删除节点！"
				}
				self.SetSession("MsgErr", self.Data["MsgErr"])
				self.Redirect("/root-information-node-list", 302)

			}
		}
	case utils.Rex(self.Ctx.Request.RequestURI, "^/root-crafts-node-del/([0-9]+)$"):
		//#crafts 删除节点 GET状态
		if sess_role, _ := self.GetSession("userrole").(int64); sess_role != -1000 {
			self.Ctx.Redirect(302, "/root-login")
		} else {
			if mid, _ := self.GetInt(":nid"); mid != 0 {
				if e := models.DelNode(mid); e != nil {
					self.Data["MsgErr"] = "删除节点失败！"
				} else {
					self.Data["MsgErr"] = "成功删除节点！"
				}
				self.SetSession("MsgErr", self.Data["MsgErr"])
				self.Redirect("/root-crafts-node-list", 302)

			}
		}
	}

}

func (self *RApi) Post() {
	switch {
	case utils.Rex(self.Ctx.Request.RequestURI, "^/root-about-gallery-edit/([0-9]+)$"):
		//# Gallery 圖片编辑 POST状态 ,主要是爲了設置url
		if sess_role, _ := self.GetSession("userrole").(int64); sess_role != -1000 {
			self.Ctx.WriteString(outtimesz)
		} else {
			inputs := self.Input()
			url := inputs.Get("url")
			//圖片ID
			editmid, _ := self.GetInt(":editmid")
			img := models.GetFile(editmid)

			if e := models.SetFile(editmid, img.Pid, img.Ctype, img.Filename, img.Content, img.Hash, img.Location, url, img.Size); e != nil {
				self.Data["MsgErr"] = "设置图片链接失败！"
			} else {
				self.Data["MsgErr"] = "设置图片链接成功！"
			}

			self.TplNames = "root/gallery_editurl.html"
			self.Render()
		}
	case utils.Rex(self.Ctx.Request.RequestURI, "^/root-about-topic-edit/([0-9]+)$"):
		//ABOUT内容编辑 POST状态
		if sess_role, _ := self.GetSession("userrole").(int64); sess_role != -1000 {
			self.Ctx.Redirect(302, "/root-login")
		} else {
			var cid, nid int64 = 1, 1
			file_location := ""
			ftitle := self.GetString("ftitle")
			stitle := self.GetString("stitle")
			content := self.GetString("content")
			uid, _ := self.GetSession("userid").(int64)
			tid, _ := self.GetInt(":tid")

			file, handler, e := self.GetFile("image")

			if handler != nil && e == nil {

				if tid != 0 {
					tp := models.GetTopic(tid)
					if tp.Attachment != "" {
						file_location = tp.Attachment
					}
				} else {
					self.Data["MsgErr"] = "你编辑的内容不存在！"
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
					defer file.Close()
					input_file := "." + path
					output_file := "." + path
					output_size := "232x135"
					output_align := "center"
					utils.Thumbnail(input_file, output_file, output_size, output_align, "white")

					//若文件存在则删除，不存在就当忽略处理
					if file_location != "" {
						if utils.Exist("." + file_location) {
							if err := os.Remove("." + file_location); err != nil {
								self.Data["MsgErr"] = "删除旧文件错误！"
							}
						}
					}
					file_location = path

				}
			}

			if ftitle != "" && nid != 0 && content != "" {
				//保存编辑
				if e := models.SetTopic(tid, cid, nid, uid, 1, ftitle, content, stitle, file_location); e != nil {
					self.Data["MsgErr"] = "你提交的修改保存失败，无法写入数据库！"

				} else {
					self.Data["MsgErr"] = "你提交的修改已保存成功！"
				}
			} else {
				//下面三个为基本诉求
				msg := ""
				if ftitle == "" {
					msg = msg + "标题不能为空！"
				}
				if nid == 0 {
					msg = msg + "分类不正确！"
				}
				if content == "" {
					msg = msg + "内容不能为空！"
				}
				self.Data["MsgErr"] = msg
			}
			self.SetSession("MsgErr", self.Data["MsgErr"])
			self.Redirect(self.Ctx.Request.RequestURI, 302)
		}
	case utils.Rex(self.Ctx.Request.RequestURI, "^/root-about-node-edit/([0-9]+)$"):
		/// about 节点编辑 POST状态

		if sess_role, _ := self.GetSession("userrole").(int64); sess_role != -1000 {
			self.Ctx.Redirect(302, "/root-login")
		} else {
			var cid int64 = 1
			title := self.GetString("title")
			content := self.GetString("content")
			nid, _ := self.GetInt(":nid")
			//编辑NODE
			if title != "" && nid != 0 {
				if e := models.SetNode(nid, title, content, cid, 1); e != nil {

					self.Data["MsgErr"] = "内容分类无法保存到数据库！"
				} else {

					self.Data["MsgErr"] = "内容分类成功保存到数据库！"
				}

				if content == "" {
					self.Data["MsgErr"] = ""
				}
				self.Data["MsgErr"] = "内容分类“" + title + "”保存成功！"
			} else {
				self.Data["MsgErr"] = "内容分类标题不能为空！"
			}

			self.SetSession("MsgErr", self.Data["MsgErr"])
			self.Redirect(self.Ctx.Request.RequestURI, 302)
		}
	case utils.Rex(self.Ctx.Request.RequestURI, "^/root-services-node-edit/([0-9]+)$"):
		/// services 节点编辑 POST状态

		if sess_role, _ := self.GetSession("userrole").(int64); sess_role != -1000 {
			self.Ctx.Redirect(302, "/root-login")
		} else {
			var cid int64 = 3
			title := self.GetString("title")
			content := self.GetString("content")
			nid, _ := self.GetInt(":nid")
			//编辑NODE
			if title != "" && nid != 0 {
				if e := models.SetNode(nid, title, content, cid, 1); e != nil {

					self.Data["MsgErr"] = "内容分类无法保存到数据库！"
				} else {

					self.Data["MsgErr"] = "内容分类成功保存到数据库！"
				}

				if content == "" {
					self.Data["MsgErr"] = ""
				}
				self.Data["MsgErr"] = "内容分类“" + title + "”保存成功！"
			} else {
				self.Data["MsgErr"] = "内容分类标题不能为空！"
			}

			self.SetSession("MsgErr", self.Data["MsgErr"])
			self.Redirect(self.Ctx.Request.RequestURI, 302)
		}
	case utils.Rex(self.Ctx.Request.RequestURI, "^/root-services-topic-edit/([0-9]+)$"):
		// services 内容编辑 POST状态
		if sess_role, _ := self.GetSession("userrole").(int64); sess_role != -1000 {
			self.Ctx.Redirect(302, "/root-login")
		} else {
			var cid int64 = 3
			file_location := ""
			ftitle := self.GetString("ftitle")
			stitle := self.GetString("stitle")
			content := self.GetString("content")
			uid, _ := self.GetSession("userid").(int64)
			tid, _ := self.GetInt(":tid")
			nid, _ := self.GetInt("nodeid")
			file, handler, e := self.GetFile("image")

			if handler != nil && e == nil {

				if tid != 0 {
					tp := models.GetTopic(tid)
					if tp.Attachment != "" {
						file_location = tp.Attachment
					}
				} else {
					self.Data["MsgErr"] = "你编辑的内容不存在！"
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
					defer file.Close()
					input_file := "." + path
					output_file := "." + path
					output_size := "211x134"
					output_align := "center"
					utils.Thumbnail(input_file, output_file, output_size, output_align, "white")

					//若文件存在则删除，不存在就当忽略处理
					if file_location != "" {
						if utils.Exist("." + file_location) {
							if err := os.Remove("." + file_location); err != nil {
								self.Data["MsgErr"] = "删除旧文件错误！"
							}
						}
					}
					file_location = path

				}
			}

			if ftitle != "" && nid != 0 && content != "" {
				//保存编辑
				if e := models.SetTopic(tid, cid, nid, uid, 1, ftitle, content, stitle, file_location); e != nil {
					self.Data["MsgErr"] = "你提交的修改保存失败，无法写入数据库！"

				} else {
					self.Data["MsgErr"] = "你提交的修改已保存成功！"
				}
			} else {
				//下面三个为基本诉求
				msg := ""
				if ftitle == "" {
					msg = msg + "标题不能为空！"
				}
				if nid == 0 {
					msg = msg + "分类不正确！"
				}
				if content == "" {
					msg = msg + "内容不能为空！"
				}
				self.Data["MsgErr"] = msg
			}
			self.SetSession("MsgErr", self.Data["MsgErr"])
			self.Redirect(self.Ctx.Request.RequestURI, 302)
		}
	case utils.Rex(self.Ctx.Request.RequestURI, "^/root-crafts-node-edit/([0-9]+)$"):
		/// crafts 节点编辑 POST状态

		if sess_role, _ := self.GetSession("userrole").(int64); sess_role != -1000 {
			self.Ctx.Redirect(302, "/root-login")
		} else {
			var cid int64 = 4
			title := self.GetString("title")
			content := self.GetString("content")
			nid, _ := self.GetInt(":nid")
			//编辑NODE
			if title != "" && nid != 0 {
				if e := models.SetNode(nid, title, content, cid, 1); e != nil {

					self.Data["MsgErr"] = "内容分类无法保存到数据库！"
				} else {

					self.Data["MsgErr"] = "内容分类成功保存到数据库！"
				}

				if content == "" {
					self.Data["MsgErr"] = ""
				}
				self.Data["MsgErr"] = "内容分类“" + title + "”保存成功！"
			} else {
				self.Data["MsgErr"] = "内容分类标题不能为空！"
			}

			self.SetSession("MsgErr", self.Data["MsgErr"])
			self.Redirect(self.Ctx.Request.RequestURI, 302)
		}
	case utils.Rex(self.Ctx.Request.RequestURI, "^/root-crafts-topic-edit/([0-9]+)$"):
		// crafts 内容编辑 POST状态
		if sess_role, _ := self.GetSession("userrole").(int64); sess_role != -1000 {
			self.Ctx.Redirect(302, "/root-login")
		} else {
			var cid int64 = 4
			file_location := ""
			ftitle := self.GetString("ftitle")
			stitle := self.GetString("stitle")
			content := self.GetString("content")
			uid, _ := self.GetSession("userid").(int64)
			tid, _ := self.GetInt(":tid")
			nid, _ := self.GetInt("nodeid")
			file, handler, e := self.GetFile("image")

			if handler != nil && e == nil {

				if tid != 0 {
					tp := models.GetTopic(tid)
					if tp.Attachment != "" {
						file_location = tp.Attachment
					}
				} else {
					self.Data["MsgErr"] = "你编辑的内容不存在！"
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
					defer file.Close()
					input_file := "." + path
					output_file := "." + path
					output_size := "534"
					output_align := "center"
					utils.Thumbnail(input_file, output_file, output_size, output_align, "black")

					//若文件存在则删除，不存在就当忽略处理
					if file_location != "" {
						if utils.Exist("." + file_location) {
							if err := os.Remove("." + file_location); err != nil {
								self.Data["MsgErr"] = "删除旧文件错误！"
							}
						}
					}
					file_location = path

				}
			}

			if ftitle != "" && nid != 0 && content != "" {
				//保存编辑
				if e := models.SetTopic(tid, cid, nid, uid, 0, ftitle, content, stitle, file_location); e != nil {
					self.Data["MsgErr"] = "你提交的修改保存失败，无法写入数据库！"

				} else {
					self.Data["MsgErr"] = "你提交的修改已保存成功！"
				}
			} else {
				//下面三个为基本诉求
				msg := ""
				if ftitle == "" {
					msg = msg + "标题不能为空！"
				}
				if nid == 0 {
					msg = msg + "分类不正确！"
				}
				if content == "" {
					msg = msg + "内容不能为空！"
				}
				self.Data["MsgErr"] = msg
			}
			self.SetSession("MsgErr", self.Data["MsgErr"])
			self.Redirect(self.Ctx.Request.RequestURI, 302)
		}
	case utils.Rex(self.Ctx.Request.RequestURI, "^/root-information-topic-edit/([0-9]+)$"):
		// information 内容编辑 POST状态
		if sess_role, _ := self.GetSession("userrole").(int64); sess_role != -1000 {
			self.Ctx.Redirect(302, "/root-login")
		} else {
			var cid int64 = 5
			file_location := ""
			ftitle := self.GetString("ftitle")
			stitle := self.GetString("stitle")
			content := self.GetString("content")
			uid, _ := self.GetSession("userid").(int64)
			tid, _ := self.GetInt(":tid")
			nid, _ := self.GetInt("nodeid")
			file, handler, e := self.GetFile("image")

			if handler != nil && e == nil {

				if tid != 0 {
					tp := models.GetTopic(tid)
					if tp.Attachment != "" {
						file_location = tp.Attachment
					}
				} else {
					self.Data["MsgErr"] = "你编辑的内容不存在！"
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
					defer file.Close()
					input_file := "." + path
					output_file := "." + path
					output_size := "534"
					output_align := "center"
					utils.Thumbnail(input_file, output_file, output_size, output_align, "black")

					//若文件存在则删除，不存在就当忽略处理
					if file_location != "" {
						if utils.Exist("." + file_location) {
							if err := os.Remove("." + file_location); err != nil {
								self.Data["MsgErr"] = "删除旧文件错误！"
							}
						}
					}
					file_location = path

				}
			}

			if ftitle != "" && nid != 0 && content != "" {
				//保存编辑
				if e := models.SetTopic(tid, cid, nid, uid, 0, ftitle, content, stitle, file_location); e != nil {
					self.Data["MsgErr"] = "你提交的修改保存失败，无法写入数据库！"

				} else {
					self.Data["MsgErr"] = "你提交的修改已保存成功！"
				}
			} else {
				//下面三个为基本诉求
				msg := ""
				if ftitle == "" {
					msg = msg + "标题不能为空！"
				}
				if nid == 0 {
					msg = msg + "分类不正确！"
				}
				if content == "" {
					msg = msg + "内容不能为空！"
				}
				self.Data["MsgErr"] = msg
			}
			self.SetSession("MsgErr", self.Data["MsgErr"])
			self.Redirect(self.Ctx.Request.RequestURI, 302)
		}
	case utils.Rex(self.Ctx.Request.RequestURI, "^/root-information-node-edit/([0-9]+)$"):
		// information 节点编辑 POST状态

		if sess_role, _ := self.GetSession("userrole").(int64); sess_role != -1000 {
			self.Ctx.Redirect(302, "/root-login")
		} else {
			var cid int64 = 5
			title := self.GetString("title")
			content := self.GetString("content")
			nid, _ := self.GetInt(":nid")
			//编辑NODE
			if title != "" && nid != 0 {
				if e := models.SetNode(nid, title, content, cid, 1); e != nil {

					self.Data["MsgErr"] = "内容分类无法保存到数据库！"
				} else {

					self.Data["MsgErr"] = "内容分类成功保存到数据库！"
				}

				if content == "" {
					self.Data["MsgErr"] = ""
				}
				self.Data["MsgErr"] = "内容分类“" + title + "”保存成功！"
			} else {
				self.Data["MsgErr"] = "内容分类标题不能为空！"
			}

			self.SetSession("MsgErr", self.Data["MsgErr"])
			self.Redirect(self.Ctx.Request.RequestURI, 302)
		}
	case utils.Rex(self.Ctx.Request.RequestURI, "^/root-contact-node-edit/([0-9]+)$"):
		//contact 节点编辑 POST状态

		if sess_role, _ := self.GetSession("userrole").(int64); sess_role != -1000 {
			self.Ctx.Redirect(302, "/root-login")
		} else {
			var cid int64 = 6
			title := self.GetString("title")
			content := self.GetString("content")
			nid, _ := self.GetInt(":nid")
			//编辑NODE
			if title != "" && nid != 0 {
				if e := models.SetNode(nid, title, content, cid, 1); e != nil {

					self.Data["MsgErr"] = "内容分类无法保存到数据库！"
				} else {

					self.Data["MsgErr"] = "内容分类成功保存到数据库！"
				}

				if content == "" {
					self.Data["MsgErr"] = ""
				}
				self.Data["MsgErr"] = "内容分类“" + title + "”保存成功！"
			} else {
				self.Data["MsgErr"] = "内容分类标题不能为空！"
			}

			self.SetSession("MsgErr", self.Data["MsgErr"])
			self.Redirect(self.Ctx.Request.RequestURI, 302)
		}
	case utils.Rex(self.Ctx.Request.RequestURI, "^/root-contact-topic-edit/([0-9]+)$"):
		//CONTACT内容编辑 POST状态
		if sess_role, _ := self.GetSession("userrole").(int64); sess_role != -1000 {
			self.Ctx.Redirect(302, "/root-login")
		} else {
			var cid int64 = 6
			file_location := ""
			ftitle := self.GetString("title")
			content := self.GetString("content")
			uid, _ := self.GetSession("userid").(int64)
			tid, _ := self.GetInt(":tid")
			nid, _ := self.GetInt("nodeid")

			if ftitle != "" && nid != 0 && content != "" {
				//保存编辑
				if e := models.SetTopic(tid, cid, nid, uid, 0, ftitle, content, "", file_location); e != nil {
					self.Data["MsgErr"] = "你提交的修改保存失败，无法写入数据库！"

				} else {
					self.Data["MsgErr"] = "你提交的修改已保存成功！"
				}
			} else {
				//下面三个为基本诉求
				msg := ""
				if ftitle == "" {
					msg = msg + "标题不能为空！"
				}
				if nid == 0 {
					msg = msg + "分类不正确！"
				}
				if content == "" {
					msg = msg + "内容不能为空！"
				}
				self.Data["MsgErr"] = msg
			}
			self.SetSession("MsgErr", self.Data["MsgErr"])
			self.Redirect(self.Ctx.Request.RequestURI, 302)
		}
	}
}
