package handlers

import (
	"../libs"
	"../models"
	"fmt"
	"net/url"
)

type NewReplyHandler struct {
	libs.BaseHandler
}

func (self *NewReplyHandler) Post() {
	tid, _ := self.GetInt("comment_parent")
	sess_userid, _ := self.GetSession("userid").(int64)

	gmt, _ := self.Ctx.Request.Cookie("gmt")

	if gmt != nil {

		if gmtstr, err := url.QueryUnescape(gmt.Value); err == nil {
			fmt.Println("Reply on tid:", tid, gmtstr)

			author := self.GetString("author")
			email := self.GetString("email")
			website := self.GetString("website")

			rc := self.GetString("comment")

			if author != "" && email != "" && tid != 0 && rc != "" {
				if err := models.AddReply(tid, sess_userid, rc, author, email, website); err != nil {
					fmt.Println(err)
				}
				self.Ctx.Redirect(302, "/view/"+self.GetString("comment_parent"))
			} else {
				self.Ctx.Redirect(302, "/")
			}
		} else {
			self.Ctx.Redirect(302, "/")
		}

	} else {
		self.Ctx.Redirect(302, "/")
	}

}
