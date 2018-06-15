package managent

import (
	"dub/app/web/managent/controller"
	"dub/utils"

	"github.com/astaxie/beego"
)

func RouteAdd() {
	manController := &controller.AgentDefaultController{}
	beego.Router("/", manController, "get:Get")
}

func RegSessionGobStruct() {
	utils.RegSessionGobStruct()
}
