package config

import (
	"define"
)

//得到数据服务器的配置与日志配置
func GetDatabaseServerConfig(path string) (define.DatabaseServerConfig, define.LogConfig, error) {
	var (
		db_opt  define.DatabaseServerConfig
		log_opt define.LogConfig
	)

	c, err := ReadConfigFile(path)
	if err != nil {
		return db_opt, log_opt, err
	}

	db_opt.Ip, err = c.GetString("net", "ip")
	if err != nil {
		return db_opt, log_opt, err
	}

	log_opt.Level, err = c.GetString("log", "level")
	if err != nil {
		return db_opt, log_opt, err
	}
	log_opt.DeviceName, err = c.GetString("log", "diveName")
	if err != nil {
		return db_opt, log_opt, err
	}
	log_opt.MaxSize, err = c.GetInt("log", "maxSize")
	if err != nil {
		return db_opt, log_opt, err
	}

	return db_opt, log_opt, nil
}
