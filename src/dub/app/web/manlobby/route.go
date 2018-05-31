package manlobby

import (
	"github.com/astaxie/beego"
	"dub/app/web/manlobby/controller"
)

func RouteAdd() {
	manController := &controller.ManDefaultController{}
	beego.Router("/", manController, "get:Get")
	beego.Router("/", manController, "post:Post")
}
