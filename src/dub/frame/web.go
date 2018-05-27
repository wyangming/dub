package frame

import "github.com/astaxie/beego"

type DubBaseController struct {
	ProxyBaseUrl string //web服务被代理的基础路径
	beego.Controller
}

func (d *DubBaseController) Prepare() {
	//设置web服务被代理的基础路径
	//暂时思路是在代理的服务器上使用request的header里设置一个信息来作为一个被代理的路径
	//以后直接打算使用tcp协议通知
	//d.Ctx.Request.Header.Get()
	d.ProxyBaseUrl = ""
}
