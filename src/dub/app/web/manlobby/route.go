package manlobby

import (
	"dub/app/web/manlobby/controller"
	"dub/utils"
	"github.com/astaxie/beego"
)

func RouteAdd() {
	manController := &controller.ManDefaultController{}
	beego.Router("/", manController, "get:Get")
	beego.Router("/", manController, "post:Post")
	beego.Router("/index", manController, "*:Index")
}
func RegSessionGobStruct() {
	utils.RegSessionGobStruct()
}

func AddFunMap() {
	beego.AddFuncMap("c2int", utils.WebTemplateC2Int)
}
