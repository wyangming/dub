package manuse

import (
	"dub/app/web/manuse/controller"

	"github.com/astaxie/beego"
)

func RouteAdd() {
	manController := &controller.ManDefaultController{}
	beego.Router("/", manController, "get:Get")
}
