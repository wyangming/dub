package use

import (
	"dub/define"
	"dub/utils"
	"dub/common"
	"dub/config"
	"fmt"
	"os"
	"dub/frame"
	json "github.com/json-iterator/go"
	"dub/app/service/use/secrpc"
	"net/rpc"
	"net"
	"net/http"
)

type ServiceUseServer struct {
	useCfg      define.ServiceUseServerConfig //用户服务配置
	logCfg      define.LogConfig              //日志配置
	regConn     common.IConnector             //与注册服务器连接
	log         *utils.Logger                 //日志对象
	dbRpc       *utils.RpcProxy               //数据服务的rpc
	dbRpcOnline bool                          //数据服务是否上线
}

func (s *ServiceUseServer) Init(cfgPath string) {
	//读取配置
	var err error
	s.useCfg, s.logCfg, err = config.GetServiceUseServerConfig(cfgPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	//日志信息
	log := utils.NewLogger()
	log.SetConfig(s.logCfg)
	s.log = log

	//连接注册服务器
	conn := frame.NewConnector()
	conn.CallBack = s.RegServiceCallBack
	s.regConn = conn
	err = s.regConn.Start(s.useCfg.RegAddr)
	if err != nil {
		s.log.Errorf("server.go Init method s.regConn.Start err %v\n", err)
		os.Exit(2)
	}
	s.Reg()
	log.Infof("reg %s server is %s", s.logCfg.DeviceName, s.useCfg.Addr)

	//加载自己的rpc服务
	secRpc := secrpc.NewSecRpc(s.log, s.dbRpc)
	rpc.Register(secRpc)

	rpc.HandleHTTP()
	l, err := net.Listen("tcp", s.useCfg.Addr)
	if err != nil {
		s.log.Errorf("rpc listen err:%v\n", err)
		os.Exit(2)
	}
	http.Serve(l, nil)
	s.log.Infoln("service useService server start")
}

func (s *ServiceUseServer) Reg() {
	serverInfo := &define.ModelRegReqServerType{
		Addr:       s.useCfg.Addr,
		ServerName: define.ServerNameService_UseServer,
		ServerType: 1,
	}

	//如果有问题重新发送命令
	for {
		data, err := json.Marshal(serverInfo)
		if err != nil {
			s.log.Errorf("server.go reg method json.Marshal(serverInfo) err %v\n", err)
			continue
		}

		//发送注册命令
		err = s.regConn.Send(define.CmdRegServer_Register, define.CmdSubRegServer_Register_Reg, data)
		if err != nil {
			s.log.Errorf("server.go reg method d.regConn.Send(define.CmdRegServer_Register %d, define.CmdSubRegServer_Register_Reg %d, data) err %v\n", define.CmdRegServer_Register, define.CmdSubRegServer_Register_Reg, err)
			continue
		}
		break
	}
}

//注册服务器回调函数
func (s *ServiceUseServer) RegServiceCallBack(mainId, subId uint16, data []byte) bool {
	if mainId == define.CmdRegServer_Register {
		switch subId {
		case define.CmdSubRegServer_Register_Reg_Inform:
			//下层服务器上线
			res := &define.ModelRegReqServerType{}
			err := json.Unmarshal(data, res)
			if err != nil {
				s.log.Errorf("server.go RegServiceCallBack(%d, %d, data []byte) json.Unmarshal(data,res) err %v\n", mainId, subId, err)
				return true
			}

			//判断是否是数据服务
			if res.ServerType == 0 {
				protocol := "tcp"
				s.dbRpc = utils.NewRpcProxy(protocol, res.Addr)
				err := s.dbRpc.Start()
				if err != nil {
					s.log.Errorf("server.go RegServiceCallBack method  s.dbRpc.Start method err %v\n", err)
					return true
				}
				s.log.Infof("link %s server rpc.\n", res.ServerName)
				s.dbRpcOnline = s.dbRpc.Status()
			}
		case define.CmdSubRegServer_Register_Reg:
			res := &define.ModelRegResServerType{}
			err := json.Unmarshal(data, res)
			if err != nil {
				s.log.Errorf("server.go RegServiceCallBack(%d, %d, data []byte) json.Unmarshal(data,res) err %v\n", mainId, subId, err)
				return true
			}

			if res.Err == 0 {
				s.log.Infoln("service use server reg success")
			} else {
				s.log.Infoln("service use server reg fail")
			}
		}
	}
	return true
}

var serviceUseServer *ServiceUseServer

func NewServiceUseServer() *ServiceUseServer {
	if serviceUseServer == nil {
		serviceUseServer = &ServiceUseServer{
			dbRpcOnline: false,
		}
	}
	return serviceUseServer
}
