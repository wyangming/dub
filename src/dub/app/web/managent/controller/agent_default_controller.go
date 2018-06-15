package controller

import ()

type AgentDefaultController struct {
	UseBaseController
}

func (m *AgentDefaultController) Get() {
	m.TplName = "index.html"
}

//select * from dubregion where regionName regexp '^[0-9][0-9]0000'
