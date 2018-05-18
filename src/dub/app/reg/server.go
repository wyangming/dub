package reg

import (
	"dub/define"
	"dub/config"
	"os"
	"fmt"
	"dub/utils"
)

type RegServer struct {
	logCfg define.LogConfig            //日志配置
	regCfg define.RegisterServerConfig //注册服务器配置
	log     *utils.Logger
}

func (r *RegServer) Init(cfgPath string)  {
	//读取配置
	var err error
	r.regCfg,r.logCfg,err=config.GetRegisterServerConfig(cfgPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	log := utils.NewLogger()
	log.SetConfig(r.logCfg)
	r.log=log
}

var regServer *RegServer

func NewDbServer() *RegServer {
	if regServer == nil {
		regServer = new(RegServer)
	}
	return regServer
}