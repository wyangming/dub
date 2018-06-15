package manlobby

import (
	"dub/common"
	"dub/config"
	"dub/define"
	"dub/frame"
	"dub/secrec"
	"dub/utils"
	"fmt"
	"os"

	"github.com/astaxie/beego"
	json "github.com/json-iterator/go"
)

type WebManLobbyCenterServer struct {
	wucCfg  define.WebManLobbyCenterServerConfig //大厅的微服务配置
	logCfg  define.LogConfig                     //日志配置
	regConn common.IConnector                    //与注册服务器连接
	log     *utils.Logger                        //日志对象
}

func (w *WebManLobbyCenterServer) Init(cfgPath string) {
	//读取配置
	var err error
	w.wucCfg, w.logCfg, err = config.GetWebManLobbyCenterServerConfig(cfgPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	//日志配置
	log := utils.NewLogger()
	log.SetConfig(w.logCfg)
	w.log = log

	//连接主服务器
	conn := frame.NewConnector()
	conn.CallBack = w.RegServiceCallBack
	w.regConn = conn
	err = w.regConn.Start(w.wucCfg.RegAddr)
	if err != nil {
		log.Errorf("server.go Init method w.regCon.Start err %v\n", err)
		os.Exit(2)
	}
	w.Reg()
	log.Infof("reg %s server is %s", w.logCfg.DeviceName, w.wucCfg.Addr)

	//启动自己的http服务器
	//设定基础的配置
	beego.BConfig.WebConfig.ViewsPath = w.wucCfg.WebWiew
	beego.SetStaticPath(w.wucCfg.WebStaticUrl, w.wucCfg.WebStaticPath)
	beego.BConfig.WebConfig.Session.SessionOn = true
	//如果想让session共享把所有的SessionProvider与SessionProviderConfig配置为一个地方，前提是在同一台服务器上
	//内存方式实现不了共享
	beego.BConfig.WebConfig.Session.SessionProvider = w.wucCfg.SessionProvider
	beego.BConfig.WebConfig.Session.SessionProviderConfig = w.wucCfg.SessionProviderConfig
	if w.wucCfg.SessionProvider != "memory" {
		RegSessionGobStruct()
	}

	beego.BConfig.RunMode = w.wucCfg.RunMode

	//日志设定
	beego.SetLogger("file", fmt.Sprintf(`{"filename":"./log/%s.log"}`, w.logCfg.DeviceName))
	//路由设定
	RouteAdd()
	//添加自定义函数
	AddFunMap()
	//启动
	beego.Run(w.wucCfg.Addr)
}

//注册服务器
func (w *WebManLobbyCenterServer) Reg() {
	serverInfo := &define.ModelRegReqServerType{
		Addr:       w.wucCfg.Addr,
		ServerName: define.ServerNameWeb_ManLobbyServer,
		ServerType: 2,
		ProxyUrl:   w.wucCfg.ProxyUrl,
	}

	for {
		data, err := json.Marshal(serverInfo)
		if err != nil {
			w.log.Errorf("server.go reg method json.Marshal(serverInfo) err %v\n", err)
			continue
		}

		//发送注册命令
		err = w.regConn.Send(define.CmdRegServer_Register, define.CmdSubRegServer_Register_Reg, data)
		if err != nil {
			w.log.Errorf("server.go reg method w.regConn.Send(define.CmdRegServer_Register %d, define.CmdSubRegServer_Register_Reg %d, data) err %v\n", define.CmdRegServer_Register, define.CmdSubRegServer_Register_Reg, err)
			continue
		}
		break
	}
}

//注册服务回调函数
func (w *WebManLobbyCenterServer) RegServiceCallBack(mainId, subId uint16, data []byte) bool {
	if mainId == define.CmdRegServer_Register {
		switch subId {
		case define.CmdSubRegServer_Register_Lobby_Proxy_Server:
			//理新gate代理微服务的信息
			req := &define.ModelRegReqLobbyProxyServer{}
			err := json.Unmarshal(data, req)
			if err != nil {
				w.log.Errorf("server.go RegServiceCallBack method json.Unmarshal(data, &reqRegister) err. %v\n", err)
				return true
			}

			if len(req.ProxyUrls) > 0 && len(req.ProxyUrls) == len(req.ServerNames) {
				secrec.UpdateProxyInfo(req.ProxyUrls, req.ServerNames)
			}
		case define.CmdSubRegServer_Register_Reg_Inform:
			//下层服务上线通知
			res := &define.ModelRegReqServerType{}
			err := json.Unmarshal(data, res)
			if err != nil {
				w.log.Errorf("server.go RegServiceCallBack(%d, %d, data []byte) json.Unmarshal(data,res) err %v\n", mainId, subId, err)
				return true
			}
			w.log.Infof("Reg_Inform is %+v	\n", res)

			//判断是否是逻辑服务
			if res.ServerType == 1 {
				protocol := "tcp"
				tmp_rpc := utils.NewRpcProxy(protocol, res.Addr)
				err := tmp_rpc.Start()
				if err != nil {
					w.log.Errorf("server.go RegServiceCallBack method s.dbRpc.Start method err %v\n", err)
					return true
				}

				//配置rpc服务
				switch res.ServerName {
				case define.ServerNameService_UseServer:
					secrec.AddSecRpc(secrec.ConstServiceUseRpc, tmp_rpc)
				}
			}
		case define.CmdSubRegServer_Register_Reg:
			res := &define.ModelRegResServerType{}
			err := json.Unmarshal(data, res)
			if err != nil {
				w.log.Errorf("server.go RegServiceCallBack(%d, %d, data []byte) json.Unmarshal(data,res) err %v\n", mainId, subId, err)
				return true
			}

			if res.Err == 0 {
				w.log.Infoln("service use server reg success")
			} else {
				w.log.Infoln("service use server reg fail")
			}
		}
	}
	return true
}

var webUseCenter *WebManLobbyCenterServer

func NewWebManLobbyCenterServer() *WebManLobbyCenterServer {
	if webUseCenter == nil {
		webUseCenter = new(WebManLobbyCenterServer)
	}
	return webUseCenter
}
