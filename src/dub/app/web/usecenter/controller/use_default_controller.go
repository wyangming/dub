package controller

import "dub/utils"

type UseDefaultController struct {
	UseBaseController
	log *utils.Logger
}

func (u *UseDefaultController) Log() *utils.Logger {
	if u.log == nil {
		u.log = utils.NewLogger()
	}
	return u.log
}
func (u *UseDefaultController) Get() {
	u.Log().Infof("all ready Get method do. the proxe url is %s\n", u.ProxyBaseUrl)
	u.TplName = "index.html"
}
