package manlobby

import (
	"dub/app/web/manlobby/controller"
	"dub/define"
	"encoding/gob"
	"github.com/astaxie/beego"
)

func RouteAdd() {
	manController := &controller.ManDefaultController{}
	beego.Router("/", manController, "get:Get")
	beego.Router("/", manController, "post:Post")
	beego.Router("/ajaxAuth", manController, "post:AjaxAuth")
}
func RegSessionGobStruct() {
	gob.Register(define.RpcSecUseResLoginByLoginName{})
}
