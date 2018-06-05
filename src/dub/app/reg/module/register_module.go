package module

import (
	"dub/common"
	"dub/define"
	"dub/utils"
	"sync"

	json "github.com/json-iterator/go"
)

//注册模块 处理注册命令
type RegisterModule struct {
	log           *utils.Logger                            //日志对象
	servers       map[uint64]*define.ModelRegReqServerType //连接的服务器
	serverSession map[uint64]common.ISession
	mutex         sync.Mutex //锁
}

func (r *RegisterModule) OnMessage(id uint16, data []byte, ses common.ISession) bool {
	result := true

	switch id {
	case define.CmdSubRegServer_Register_Reg:
		result = r.registerServer(data, ses)
	case define.CmdSubRegServer_Register_Shut:
		result = r.registerShutDown(ses)
	case define.CmdSubRegServer_Register_Lobby_Proxy_Server:
		result = r.lobbyProxyServer(data, ses)
	}

	return result
}

//向大厅服务器发送代理微服务的地址信息
func (r *RegisterModule) lobbyProxyServer(data []byte, ses common.ISession) bool {
	req := &define.ModelRegReqLobbyProxyServer{}
	err := json.Unmarshal(data, req)
	if err != nil {
		r.log.Errorf("register_module.go lobbyProxyServer method json.Unmarshal(data, &reqRegister) err. %v\n", err)
		return true
	}

	for mk, mv := range r.servers {
		if mv.ServerName == req.InformServerName {
			r.serverSession[mk].Send(define.CmdRegServer_Register, define.CmdSubRegServer_Register_Lobby_Proxy_Server, data)
		}
	}
	return true
}

//关闭一个服务器
func (r *RegisterModule) registerShutDown(ses common.ISession) bool {
	req, ok := r.servers[ses.ID()]
	if ok {
		_, himToServer := r.relationServer(req)
		//通知所需要自己的服务器,已经下线
		if len(himToServer) > 0 {
			data, _ := json.Marshal(req)
			for i := 0; i < len(himToServer); i++ {
				himToServer[i].Send(define.CmdRegServer_Register, define.CmdSubRegServer_Register_Shut, data)
			}
		}
	}
	return true
}

//注册服务器
func (r *RegisterModule) registerServer(data []byte, ses common.ISession) bool {
	req := &define.ModelRegReqServerType{}
	res := &define.ModelRegResServerType{}

	err := json.Unmarshal(data, req)
	if err != nil {
		r.log.Errorf("register_module.go RegisterServer method json.Unmarshal(data, &reqRegister) err. %v\n", err)
		res.Err = 1
		reData, _ := json.Marshal(res)
		ses.Send(define.CmdRegServer_Register, define.CmdSubRegServer_Register_Reg, reData)
		return true
	}

	res.Err = 0
	reData, _ := json.Marshal(res)
	ses.Send(define.CmdRegServer_Register, define.CmdSubRegServer_Register_Reg, reData)
	r.mutex.Lock()

	toHimServer, himToServer := r.relationServer(req)

	//添加服务器
	r.serverSession[ses.ID()] = ses
	r.servers[ses.ID()] = req

	//给自己发需要的服务器
	if len(toHimServer) > 0 {
		r.log.Infoln("toHimServer > 0 gate server on line, server pre is ", req)
		for i := 0; i < len(toHimServer); i++ {
			r.log.Infoln("need server is ", toHimServer[i])
			reData, _ := json.Marshal(toHimServer[i])
			ses.Send(define.CmdRegServer_Register, define.CmdSubRegServer_Register_Reg_Inform, reData)
		}
	}

	//通知所需要自己的服务器
	if len(himToServer) > 0 {
		for i := 0; i < len(himToServer); i++ {
			himToServer[i].Send(define.CmdRegServer_Register, define.CmdSubRegServer_Register_Reg_Inform, data)
		}
	}
	r.mutex.Unlock()

	return true
}

//通知服务器
func (r *RegisterModule) relationServer(req *define.ModelRegReqServerType) ([]*define.ModelRegReqServerType, []common.ISession) {
	var (
		toHimServer []*define.ModelRegReqServerType //此服务器上线后需要那些服务器支持
		himToServer []common.ISession               //此服务器上线后要通知那些服务器
	)

	if len(r.servers) < 1 {
		return toHimServer, himToServer
	}

	for key, server := range r.servers {
		switch req.ServerType {
		case 0: //数据服务
			//通知所有逻辑服务器
			if server.ServerType == 1 {
				himToServer = append(himToServer, r.serverSession[key])
			}
		case 1: //逻辑服务
			//所有数据服务通知给此服务器，再把此服务器通知给应用服
			if server.ServerType == 0 {
				toHimServer = append(toHimServer, server)
			}
			if server.ServerType == 2 {
				himToServer = append(himToServer, r.serverSession[key])
			}
		case 2: //应用服务
			//所有的逻辑服务通知给应用服务，再把应用服务通知给网关服务
			if server.ServerType == 1 {
				toHimServer = append(toHimServer, server)
			}
			if server.ServerType == 3 {
				himToServer = append(himToServer, r.serverSession[key])
			}
		case 3: //网关服务
			//把所有的web应用服务器推送给网关服务
			if server.ServerType == 2 {
				toHimServer = append(toHimServer, server)
			}
		}
	}
	return toHimServer, himToServer
}

//实例化注册模块
func NewRegistModule() common.IMoudle {
	return &RegisterModule{
		servers:       make(map[uint64]*define.ModelRegReqServerType),
		serverSession: make(map[uint64]common.ISession),
		log:           utils.NewLogger(),
	}
}
