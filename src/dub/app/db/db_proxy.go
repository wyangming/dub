package db

import (
	"database/sql"
	"dub/utils"
	"dub/define"
)

const (
	DB_Accounts = iota //基础账号库
)

type DbProxy struct {
	dbMap map[int]*sql.DB
	log   *utils.Logger
}

//初始化数据库代理的信息
func (d *DbProxy) Init(cfg *define.DatabaseServerConfig) error {
	d.log.Infof("db_server %s db_server start ...\n", cfg.Ip)

	d.dbMap = make(map[int]*sql.DB)
	return nil
}

var _db_instance *DbProxy

//得到数据库代理信息，如果第一次调用则需要再次调用
// Init方法进行初始化
func NewDbProxy() *DbProxy {
	if _db_instance == nil {
		_db_instance = new(DbProxy)
		_db_instance.log = utils.NewLogger()
	}
	return _db_instance
}
