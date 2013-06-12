package main

import (
	"./handlers"
	"./handlers/root"
	"./models"
	//"github.com/astaxie/beego"
	"github.com/insionng/torgo"
	//"./torgo"
)

func main() {
	models.CreateDb()
	torgo.SetStaticPath("/static", "./static")
	torgo.SetStaticPath("/archives", "./archives")

	torgo.Router("/", &handlers.MainHandler{})
	torgo.Router("/category/:cid([0-9]+)", &handlers.MainHandler{})
	torgo.Router("/search", &handlers.SearchHandler{})

	torgo.Router("/node/:nid([0-9]+)", &handlers.NodeHandler{})
	torgo.Router("/view/:tid([0-9]+)", &handlers.ViewHandler{})

	torgo.Router("/register", &handlers.RegHandler{})
	torgo.Router("/login", &handlers.LoginHandler{})
	torgo.Router("/logout", &handlers.LogoutHandler{})

	torgo.Router("/like/topic/:tid([0-9]+)", &handlers.LikeTopicHandler{})
	torgo.Router("/hate/topic/:tid([0-9]+)", &handlers.HateTopicHandler{})

	torgo.Router("/like/node/:nid([0-9]+)", &handlers.LikeNodeHandler{})
	torgo.Router("/hate/node/:nid([0-9]+)", &handlers.HateNodeHandler{})

	torgo.Router("/new/category", &handlers.NewCategoryHandler{})
	torgo.Router("/new/node", &handlers.NewNodeHandler{})
	torgo.Router("/new/topic", &handlers.NewTopicHandler{})
	torgo.Router("/new/reply/:tid([0-9]+)", &handlers.NewReplyHandler{})

	torgo.Router("/modify/category", &handlers.ModifyCategoryHandler{})
	torgo.Router("/modify/node", &handlers.ModifyNodeHandler{})

	torgo.Router("/topic/delete/:tid([0-9]+)", &handlers.TopicDeleteHandler{})
	torgo.Router("/topic/edit/:tid([0-9]+)", &handlers.TopicEditHandler{})

	torgo.Router("/node/delete/:nid([0-9]+)", &handlers.NodeDeleteHandler{})
	torgo.Router("/node/edit/:nid([0-9]+)", &handlers.NodeEditHandler{})

	torgo.Router("/delete/reply/:rid([0-9]+)", &handlers.DeleteReplyHandler{})

	//root routes
	torgo.Router("/root", &root.RMainHandler{})
	torgo.Router("/root-login", &root.RLoginHandler{})
	torgo.Router("/root/account", &root.RAccountHandler{})


	torgo.SessionOn = true
	torgo.Run()
}
