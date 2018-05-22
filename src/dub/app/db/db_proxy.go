package db

import (
	"database/sql"
	"dub/define"
	"dub/utils"
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
	d.log.Infof("db_server %s db_server start ...\n", cfg.Addr)

	d.dbMap = make(map[int]*sql.DB)
	return nil
}

//得到一个数据库代理
func (d *DbProxy) Get(dbtype int) *sql.DB {
	sqlDb, ok := d.dbMap[dbtype]
	if ok {
		return sqlDb
	}
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
