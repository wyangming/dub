package manuse

import (
	"dub/app/web/manuse/controller"
	"dub/define"
	"encoding/gob"

	"github.com/astaxie/beego"
)

func RouteAdd() {
	manController := &controller.ManDefaultController{}
	beego.Router("/", manController, "get:Get")
}

func RegSessionGobStruct() {
	gob.Register(define.RpcSecUseResLoginByLoginName{})
}
