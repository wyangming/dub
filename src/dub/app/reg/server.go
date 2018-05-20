package reg

import (
	"dub/app/reg/module"
	"dub/common"
	"dub/config"
	"dub/define"
	"dub/frame"
	"dub/utils"
	"fmt"
	"os"
)

type RegServer struct {
	logCfg  define.LogConfig            //日志配置
	regCfg  define.RegisterServerConfig //注册服务器配置
	log     *utils.Logger               //日志对象
	ln      *frame.ListenerTcp          //监听的服务
	modules map[uint16]common.IMoudle   //所有模块
}

func (r *RegServer) Init(cfgPath string) {
	//读取配置
	var err error
	r.regCfg, r.logCfg, err = config.GetRegisterServerConfig(cfgPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	log := utils.NewLogger()
	log.SetConfig(r.logCfg)
	r.log = log

	log.Infoln("reg_server start")
	r.ln = frame.NewListener(r.regCfg.Addr)
	r.ln.CallBack = r.lCallBack
	r.ln.OnShutdown = r.lnConnShutDown
	err = r.ln.Start()
	if err != nil {
		log.Errorf("service.go Init method frame.NewListener(r.regCfg.Addr) err. %v\n", err)
		os.Exit(2)
	}

	//添加模块
	r.modules[define.CmdRegServer_Register] = module.NewRegistModule()
	r.log.Infof("register server load register module success.\n")
}

//当监听关闭连接时的事件
func (r *RegServer) lnConnShutDown(ses common.ISession) bool {
	return r.lCallBack(define.CmdRegServer_Register, define.CmdSubRegServer_Register_Shut, nil, ses)
}

//服务回调函数
func (r *RegServer) lCallBack(mainId uint16, subId uint16, data []byte, ses common.ISession) bool {
	module, ok := r.modules[mainId]
	if ok {
		return module.OnMessage(subId, data, ses)
	}

	return true
}

var regServer *RegServer

//返回注册服务器
func NewDbServer() *RegServer {
	if regServer == nil {
		regServer = new(RegServer)
	}
	return regServer
}
