package root

import (
	"../../libs"
	"../../models"
	"../../utils"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

type RIntroductionHandler struct {
	libs.RootHandler
}

func (self *RIntroductionHandler) Get() {
	self.Data["catpage"] = "about"

	var tid int64 = 0
	for _, v := range models.GetAllTopicByCid(1, 0, 1, 10, "id") {
		if v.Id > 0 {
			tid = v.Id
		}
	}

	self.Data["topic"] = models.GetTopic(tid)

	//self.Data["nodes"] = models.GetAllNodeByCid(1, 0, 0, "id")
	self.Data["MsgErr"], _ = self.GetSession("MsgErr").(string)
	self.DelSession("MsgErr")

	//簡介設置
	self.Data["asidepage"] = "root-about-introduction"
	self.TplNames = "root/about_introduction.html"

	self.Render()
}

func (self *RIntroductionHandler) Post() {
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

	if content == "" {
		msg = msg + "内容不能为空！"
	}

	self.Data["MsgErr"] = msg

	if msg == "" {

		//新增内容POST状态
		file, handler, e := self.GetFile("image")
		path := ""

		if handler == nil {
			var tid int64 = 0
			for _, v := range models.GetAllTopicByCid(1, 0, 1, 10, "id") {
				if v.Id > 0 {
					tid = v.Id
				}
			}

			if m := models.GetTopic(tid); m.Attachment != "" {
				path = m.Attachment
			} else {

				self.Data["MsgErr"] = "你还没有选择封面！"
			}
		}

		if handler != nil {
			if e != nil {
				self.Data["MsgErr"] = "传输过程文件产生错误！"

			} else {

				ext := "." + strings.Split(handler.Filename, ".")[1]
				filename := utils.MD5(time.Now().String()) + ext

				path = "/archives/upload/" + time.Now().Format("2006/01/02/")

				os.MkdirAll("."+path, 0644)
				path = path + filename
				f, err := os.OpenFile("."+path, os.O_WRONLY|os.O_CREATE, 0644)
				defer f.Close()
				if err != nil {
					self.Data["MsgErr"] = "无法打开服务端文件存储路径！"

				} else {
					//拷貝成功之後執行刪除舊文件
					if _, err := io.Copy(f, file); err == nil {
						var tid int64 = 0
						for _, v := range models.GetAllTopicByCid(1, 0, 1, 10, "id") {
							if v.Id > 0 {
								tid = v.Id
							}
						}

						if m := models.GetTopic(tid); m.Attachment != "" {
							if err := os.Remove("." + m.Attachment); err != nil {
								fmt.Println("Remove Old Image", err)
							}
						}

					}
				}
			}

		}

		if self.Data["MsgErr"] == "" && path != "" && ftitle != "" && content != "" {
			//开始添加内容
			input_file := "." + path
			output_file := "." + path
			output_size := "233x139"
			output_align := "center"
			if e := utils.Thumbnail(input_file, output_file, output_size, output_align, "white"); e != nil {
				fmt.Println("Thumbnail", e)
			}

			if e := models.SetTopic(1, cid, nodeid, uid, 10, ftitle, content, stitle, path); e != nil {
				self.Data["MsgErr"] = "添加“" + ftitle + "”失败，无法写入数据库！"
			} else {
				self.Data["MsgErr"] = "添加“" + ftitle + "”成功，你可以继续添加其他内容！"
			}

		}
	}

	self.SetSession("MsgErr", self.Data["MsgErr"])
	self.Redirect(self.Ctx.Request.RequestURI, 302)

}
