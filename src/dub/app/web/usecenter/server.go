package usecenter

import (
	"dub/define"
	"dub/common"
	"dub/utils"
)

const (
	ServiceUseRpc = iota
)

type WebUseCenterServer struct {
	wucCfg     define.WebUserCenterServerConfig //用户的微服务
	logCfg     define.LogConfig                 //日志配置
	regCon     common.IConnector                //与注册服务器连接
	log        *utils.Logger                    //日志对象
	serviceRpc map[uint8]*utils.RpcProxy        //服务ipc集合
}
