package controller

import "dub/utils"

type ManDefaultController struct {
	UseBaseController
	logger *utils.Logger
}

func (m *ManDefaultController) log() *utils.Logger {
	if m.logger == nil {
		m.logger = utils.NewLogger()
	}
	return m.logger
}

func (m *ManDefaultController) Post() {
	m.TplName = "index.html"
}

func (m *ManDefaultController) Get() {
	m.TplName = "index.html"
}
