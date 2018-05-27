package web

import (
	"dub/define"
	"dub/common"
	"dub/utils"
	"dub/config"
	"fmt"
	"os"
	"dub/frame"
	json "github.com/json-iterator/go"
)

type GateWebServer struct {
	gwsCfg   define.GateWebServerConfig //web网关服务配置
	logCfg   define.LogConfig           //日志配置
	regConn  common.IConnector          //与注册服务器连接
	log      *utils.Logger              //日志对象
	proxyUrl map[string]string          //代理的路径与服务器映射
}

func (g *GateWebServer) Init(cfgPath string) {
	//读取配置
	var err error
	g.gwsCfg, g.logCfg, err = config.GetGateWebServerConfig(cfgPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	//日志配置
	log := utils.NewLogger()
	log.SetConfig(g.logCfg)
	g.log = log

	//连接主服务器
	conn := frame.NewConnector()
	conn.CallBack = g.RegServiceCallBack
	g.regConn = conn
	err = g.regConn.Start(g.gwsCfg.RegAddr)
	if err != nil {
		log.Errorf("server.go Init method g.regCon.Start err %v\n", err)
		os.Exit(2)
	}
	g.Reg()
	log.Infof("reg %s server is %s", g.logCfg.DeviceName, g.gwsCfg.Addr)

	g.proxyUrl = make(map[string]string)
	//TODO:自己应有的业务
}

func (g *GateWebServer) Reg() {
	serverInfo := &define.ModelRegReqServerType{
		Addr:       g.gwsCfg.Addr,
		ServerName: define.ServerNameGate_WebServer,
		ServerType: 3,
	}

	for {
		data, err := json.Marshal(serverInfo)
		if err != nil {
			g.log.Errorf("server.go reg method json.Marshal(serverInfo) err %v\n", err)
			continue
		}

		//发送注册命令
		err = g.regConn.Send(define.CmdRegServer_Register, define.CmdSubRegServer_Register_Reg, data)
		if err != nil {
			g.log.Errorf("server.go reg method g.regConn.Send(define.CmdRegServer_Register %d, define.CmdSubRegServer_Register_Reg %d, data) err %v\n", define.CmdRegServer_Register, define.CmdSubRegServer_Register_Reg, err)
			continue
		}
		break
	}
}

//注册服务回调函数
func (g *GateWebServer) RegServiceCallBack(mainId, subId uint16, data []byte) bool {
	if mainId == define.CmdRegServer_Register {
		switch subId {
		case define.CmdSubRegServer_Register_Reg_Inform:
			//下层服务上线通知
			res := &define.ModelRegReqServerType{}
			err := json.Unmarshal(data, res)
			if err != nil {
				g.log.Errorf("server.go RegServiceCallBack(%d, %d, data []byte) json.Unmarshal(data,res) err %v\n", mainId, subId, err)
				return true
			}
			g.log.Infof("Reg_Inform is %+v	\n", res)

			//判断是否是逻辑服务
			//web网关接入web应用服务
			if res.ServerType == 2 {
				protocol := "tcp"
				tmp_rpc := utils.NewRpcProxy(protocol, res.Addr)
				err := tmp_rpc.Start()
				if err != nil {
					g.log.Errorf("server.go RegServiceCallBack method s.dbRpc.Start method err %v\n", err)
					return true
				}

				//配置rpc服务
				switch res.ServerName {
				case define.ServerNameWeb_UseCenterServer:
					if len(res.ProxyUrl) < 0 {
						res.ProxyUrl = "/"
					}
					g.proxyUrl[res.ProxyUrl] = res.Addr
				}
			}
		case define.CmdSubRegServer_Register_Reg:
			res := &define.ModelRegResServerType{}
			err := json.Unmarshal(data, res)
			if err != nil {
				g.log.Errorf("server.go RegServiceCallBack(%d, %d, data []byte) json.Unmarshal(data,res) err %v\n", mainId, subId, err)
				return true
			}

			if res.Err == 0 {
				g.log.Infoln("service use server reg success")
			} else {
				g.log.Infoln("service use server reg fail")
			}
		}
	}
	return true
}

var gateWebServer *GateWebServer

func NewGateWebServer() *GateWebServer {
	if gateWebServer == nil {
		gateWebServer = new(GateWebServer)
	}
	return gateWebServer
}
