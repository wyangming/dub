package controller

import ()

type ManDefaultController struct {
	UseBaseController
}

func (m *ManDefaultController) Get() {
	m.TplName = "index.html"
}
