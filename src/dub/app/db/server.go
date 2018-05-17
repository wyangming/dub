package db

import (
	"dub/common"
	"dub/config"
	"dub/define"
	"dub/utils"
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
	opt, log_opt, err := config.GetDatabaseServerConfig(cfgPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	log := utils.NewLogger()
	log.SetConfig(log_opt)
	d.log = log

	conn := NewConnector()
	conn.CallBack = d.callBack
	d.regConn = conn

	//数据库代理
	dbProxy := NewDbProxy()
	dbProxy.Init(&opt)

	log.Infoln("db_server start")
}

//向注册服务器发注册命令
func (d *DbServer) reg() {
	err := d.regConn.Start(d.dbCfg.Register)
	if err != nil {
		d.log.Errorf("db app server.go reg method d.regConn.Start(d.dbCfg.Register) err %v\n", err)
		os.Exit(2)
	}

}

func (d *DbServer) callBack(mainId, subId uint16, data []byte) bool {

}

var dbServer *DbServer

func NewDbServer() *DbServer {
	if dbServer == nil {
		dbServer = new(DbServer)
	}
	return dbServer
}
