package usecenter

import (
	"github.com/astaxie/beego"
	"dub/app/web/usecenter/controller"
)

func RouteAdd() {
	useController := &controller.UseDefaultController{}
	beego.Router("/", useController, "get:Get")
}
