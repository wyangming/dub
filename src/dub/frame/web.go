package frame

import (
	"dub/define"
	"dub/utils"
	"github.com/astaxie/beego"
)

type DubBaseController struct {
	ProxyBaseUrl string //web服务被代理的基础路径
	beego.Controller
	logger *utils.Logger
}

func (d *DubBaseController) Prepare() {
	//设置web服务被代理的基础路径
	//	//暂时思路是在代理的服务器上使用request的header里设置一个信息来作为一个被代理的路径
	proxy_url := d.Ctx.Request.Header.Get(define.Gate_String_Web_Proxy)
	if proxy_url == "/" {
		proxy_url = ""
	}
	d.Data["baseUrl"] = proxy_url
}

func (d *DubBaseController) Log() *utils.Logger {
	if d.logger == nil {
		d.logger = utils.NewLogger()
	}
	return d.logger
}
