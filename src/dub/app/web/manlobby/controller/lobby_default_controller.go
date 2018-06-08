package controller

import (
	"dub/define"
	"dub/secrec"
	"fmt"
)

type ManDefaultController struct {
	LobbyBaseController
}

func (m *ManDefaultController) Post() {
	m.TplName = "index.html"
	m.Data["result"] = false

	//非post请求直接跳转到登录页面不做处理
	if "POST" != m.Ctx.Request.Method {
		return
	}

	//登录参数处理
	login_phone := m.GetString("phone_number", "")
	pwd := m.GetString("pwd", "")
	if len(login_phone) < 1 || len(pwd) < 1 {
		m.Data["msg"] = "参数错误"
		return
	}

	//处理登录逻辑
	use_rpc := secrec.GetRpc(secrec.ConstServiceUseRpc)
	if use_rpc == nil {
		m.Data["msg"] = "服务器错误"
		m.Log().Errorf("not find use_rpc info\n")
		return
	}

	arg := define.RpcSecUseReqLoginByLoginName{
		LoginName: login_phone,
		LoginPwd:  pwd,
	}
	reply := define.RpcSecUseResLoginByLoginName{}

	err := use_rpc.Call("SecUseRpc.LoginByLoginName", &arg, &reply)
	if err != nil {
		m.Data["msg"] = "服务器错误"
		m.Log().Errorf("use_default_controller.go Post method s.dbRpc.Call SecUseRpc.LoginByLoginName err. %v\n", err)
		return
	}

	if reply.Err == 1 {
		m.Data["msg"] = "没有找到用户"
		return
	} else if reply.Err == 2 {
		m.Data["msg"] = "服务器错误"
		return
	} else if reply.Err == 3 {
		m.Data["msg"] = "用户名密码不正确"
		return
	}

	//替换所有微服务的url
	if len(reply.Auths) > 0 {
		for i := 0; i < len(reply.Auths); i++ {
			proxy_url := secrec.GetProxyByServerName(reply.Auths[i].AuthMicroServerName)
			if len(proxy_url) > 0 && proxy_url != "/" {
				reply.Auths[i].AuthUrl = fmt.Sprintf("%s%s", proxy_url, reply.Auths[i].AuthUrl)
			}
		}
	}

	m.SetSession(define.Web_Session_Name_Man_Login_Use, reply)

	//权限信息放到前台
	m.Data["LoginInfo"] = reply

	m.Data["result"] = true
	m.Data["msg"] = "登录成功"
	m.TplName = "main.html"
}

func (m *ManDefaultController) Get() {
	m.TplName = "index.html"
}

//得到当前用户的权限
func (m *ManDefaultController) AjaxAuth() {
	result := make(map[string]interface{})
	result["result"] = false

	ses_val := m.GetSession(define.Web_Session_Name_Man_Login_Use)
	if ses_val == nil {
		result["msg"] = "no login"
		m.Data["json"] = result
		m.ServeJSON()
		return
	}

	if login_info, login_ok := ses_val.(define.RpcSecUseResLoginByLoginName); !login_ok {
		result["msg"] = "no login"
	} else {
		result["result"] = true
		result["auths"] = login_info.Auths
	}
	m.Data["json"] = result

	m.ServeJSON()
}
