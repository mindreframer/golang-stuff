package root

import (
	"../../libs"
	"../../models"
	"../../utils"
	"fmt"
	"io"
	"os"
	//"strconv"
	"strings"
	"time"
)

type RGalleryHandler struct {
	libs.RootHandler
}

func (self *RGalleryHandler) Get() {

	self.TplNames = "root/gallery.html"

	self.Data["MsgErr"], _ = self.GetSession("MsgErr").(string)
	self.DelSession("MsgErr")

	self.Data["MsgErr2"], _ = self.GetSession("MsgErr2").(string)
	self.DelSession("MsgErr2")

	self.Data["catpage"] = "about"
	self.Data["asidepage"] = "gallery"
	self.Data["images"] = models.GetAllFileByCtype(10)
	self.Data["images2"] = models.GetAllFileByCtype(20)
	self.Render()

}

func (self *RGalleryHandler) Post() {
	//Gallery上传状态
	file, handler, e := self.GetFile("uploadfile")
	file2, handler2, e2 := self.GetFile("uploadfile2")

	if handler != nil {

		defer file.Close()
		if e != nil {
			fmt.Println(e)
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
			fmt.Println(err)
			self.Data["MsgErr"] = "无法打开服务端文件存储路径！"
		} else {
			io.Copy(f, file)
			models.AddFile(10, path, "") //ctype設置為10,以免跟正常的已標記上傳文件混淆.

			input_file := "." + path
			output_file := "." + path
			output_size := "950x360"
			output_align := "center"
			utils.Thumbnail(input_file, output_file, output_size, output_align, "white")
			self.Data["MsgErr"] = "上传图片 " + handler.Filename + " 成功！"

		}
	}
	if handler2 != nil {

		defer file2.Close()
		if e2 != nil {
			fmt.Println(e2)
			self.Data["MsgErr2"] = "传输过程文件产生错误！"
		}

		ext := "." + strings.Split(handler2.Filename, ".")[1]
		filename := utils.MD5(time.Now().String()) + ext

		path := "/archives/upload/" + time.Now().Format("2006/01/02/")

		os.MkdirAll("."+path, 0644)
		path = path + filename
		f, err := os.OpenFile("."+path, os.O_WRONLY|os.O_CREATE, 0644)
		defer f.Close()
		if err != nil {
			fmt.Println(err)
			self.Data["MsgErr2"] = "无法打开服务端文件存储路径！"
		} else {
			io.Copy(f, file2)
			models.AddFile(20, path, "") //ctype設置為20,以免跟正常的已標記上傳文件混淆.

			input_file := "." + path
			output_file := "." + path
			output_size := "320x100"
			output_align := "center"
			utils.Thumbnail(input_file, output_file, output_size, output_align, "white")
			self.Data["MsgErr2"] = "上传图片 " + handler2.Filename + " 成功！"

		}
	} else {
		self.Data["MsgErr2"] = "你尚未选择图片文件！"
	}

	self.SetSession("MsgErr2", self.Data["MsgErr2"])
	self.Redirect("/root-about-gallery", 302)

}
