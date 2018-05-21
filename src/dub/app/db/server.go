package db

import (
	"dub/common"
	"dub/config"
	"dub/define"
	"dub/frame"
	"dub/utils"
	"fmt"
	json "github.com/json-iterator/go"
	"os"
)

type DbServer struct {
	dbCfg   define.DatabaseServerConfig //数据库配置
	logCfg  define.LogConfig            //日志配置
	regConn common.IConnector           //与注册服务器连接
	log     *utils.Logger
}

func (d *DbServer) Init(cfgPath string) {
	//读取配置
	var err error
	d.dbCfg, d.logCfg, err = config.GetDatabaseServerConfig(cfgPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	//日志信息
	log := utils.NewLogger()
	log.SetConfig(d.logCfg)
	d.log = log

	//连接注册服务器
	conn := frame.NewConnector()
	conn.CallBack = d.connCallBack
	d.regConn = conn
	err = d.regConn.Start(d.dbCfg.RegAddr)
	if err != nil {
		d.log.Errorf("server.go reg method d.regConn.Start(d.dbCfg.Register) err %v\n", err)
		os.Exit(2)
	}
	d.reg()

	//数据库代理
	dbProxy := NewDbProxy()
	dbProxy.Init(&d.dbCfg)

	//加载rpc服务

	log.Infoln("db_server start")
}

//向注册服务器发注册命令
func (d *DbServer) reg() {
	//如果有问题重新发送命令
	for {
		serverInfo := &define.ModelRegReqServerType{
			Addr:       d.dbCfg.Addr,
			ServerName: "dbServer",
		}
		data, err := json.Marshal(serverInfo)
		if err != nil {
			d.log.Errorf("server.go reg method json.Marshal(serverInfo) err %v\n", err)
			continue
		}

		//发送注册命令
		err = d.regConn.Send(define.CmdRegServer_Register, define.CmdSubRegServer_Register_Reg, data)
		if err != nil {
			d.log.Errorf("server.go reg method d.regConn.Send(define.CmdRegServer_Register %d, define.CmdSubRegServer_Register_Reg %d, data) err %v\n", define.CmdRegServer_Register, define.CmdSubRegServer_Register_Reg, err)
			continue
		}
		break
	}
}

//
func (d *DbServer) connCallBack(mainId, subId uint16, data []byte) bool {
	if mainId == 1 && subId == 1 {
		res := &define.ModelRegResServerType{}
		err := json.Unmarshal(data, res)
		if err != nil {
			d.log.Errorf("server.go connCallBack(%d, %d, data []byte) json.Unmarshal(data,res) err %v\n", mainId, subId, err)
			return true
		}

		if res.Err == 0 {
			d.log.Infoln("db server reg success")
		} else {
			d.log.Infoln("db server reg fail")
		}
	}

	return true
}

var dbServer *DbServer

func NewDbServer() *DbServer {
	if dbServer == nil {
		dbServer = new(DbServer)
	}
	return dbServer
}
