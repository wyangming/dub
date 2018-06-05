package controller

import (
	"dub/utils"
)

type ManDefaultController struct {
	UseBaseController
	logger *utils.Logger
}

func (m *ManDefaultController) Get() {
	m.TplName = "index.html"
}
