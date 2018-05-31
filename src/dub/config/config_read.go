package config

import (
	"dub/define"
)

//得到web的网关服务
func GetGateWebServerConfig(path string) (define.GateWebServerConfig, define.LogConfig, error) {
	var (
		gus_opt define.GateWebServerConfig
		log_opt define.LogConfig
	)
	c, err := ReadConfigFile(path)
	if err != nil {
		return gus_opt, log_opt, err
	}

	gus_opt.Addr, err = c.GetString("net", "addr")
	if err != nil {
		return gus_opt, log_opt, err
	}

	gus_opt.RegAddr, err = c.GetString("net", "regAddr")
	if err != nil {
		return gus_opt, log_opt, err
	}

	err = fullLog(c, &log_opt)
	if err != nil {
		return gus_opt, log_opt, err
	}

	return gus_opt, log_opt, nil
}

//得到web用户中心
func GetWebUserCenterServerConfig(path string) (define.WebUserCenterServerConfig, define.LogConfig, error) {
	var (
		wuc_opt define.WebUserCenterServerConfig
		log_opt define.LogConfig
	)

	c, err := ReadConfigFile(path)
	if err != nil {
		return wuc_opt, log_opt, err
	}

	wuc_opt.Addr, err = c.GetString("net", "addr")
	if err != nil {
		return wuc_opt, log_opt, err
	}
	wuc_opt.RegAddr, err = c.GetString("net", "regAddr")
	if err != nil {
		return wuc_opt, log_opt, err
	}
	wuc_opt.WebWiew, err = c.GetString("web", "webView")
	if err != nil {
		return wuc_opt, log_opt, err
	}
	wuc_opt.ProxyUrl, err = c.GetString("web", "proxyUrl")
	if err != nil {
		return wuc_opt, log_opt, err
	}
	wuc_opt.WebStaticUrl, err = c.GetString("web", "webStaticUrl")
	if err != nil {
		return wuc_opt, log_opt, err
	}
	wuc_opt.WebStaticPath, err = c.GetString("web", "webStaticPath")
	if err != nil {
		return wuc_opt, log_opt, err
	}
	wuc_opt.RunMode, err = c.GetString("web", "runMode")
	if err != nil {
		return wuc_opt, log_opt, err
	}

	err = fullLog(c, &log_opt)
	if err != nil {
		return wuc_opt, log_opt, err
	}

	return wuc_opt, log_opt, nil
}

//得到用户中心与其日志配置
func GetServiceUseServerConfig(path string) (define.ServiceUseServerConfig, define.LogConfig, error) {
	var (
		use_opt define.ServiceUseServerConfig
		log_opt define.LogConfig
	)

	c, err := ReadConfigFile(path)
	if err != nil {
		return use_opt, log_opt, err
	}

	use_opt.Addr, err = c.GetString("net", "addr")
	if err != nil {
		return use_opt, log_opt, err
	}

	use_opt.RegAddr, err = c.GetString("net", "regAddr")
	if err != nil {
		return use_opt, log_opt, err
	}

	err = fullLog(c, &log_opt)
	if err != nil {
		return use_opt, log_opt, err
	}

	return use_opt, log_opt, nil
}

//得到注册服务器与其日志的配置
func GetRegisterServerConfig(path string) (define.RegisterServerConfig, define.LogConfig, error) {
	var (
		reg_opt define.RegisterServerConfig
		log_opt define.LogConfig
	)

	c, err := ReadConfigFile(path)
	if err != nil {
		return reg_opt, log_opt, err
	}

	reg_opt.Addr, err = c.GetString("net", "addr")
	if err != nil {
		return reg_opt, log_opt, err
	}

	err = fullLog(c, &log_opt)
	if err != nil {
		return reg_opt, log_opt, err
	}

	return reg_opt, log_opt, nil
}

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

	db_opt.Addr, err = c.GetString("net", "addr")
	if err != nil {
		return db_opt, log_opt, err
	}

	db_opt.RegAddr, err = c.GetString("net", "regAddr")
	if err != nil {
		return db_opt, log_opt, err
	}

	err = fullLog(c, &log_opt)
	if err != nil {
		return db_opt, log_opt, err
	}

	return db_opt, log_opt, nil
}

//从配置文件里返回服务的日志对你
func fullLog(c *ConfigFile, log_opt *define.LogConfig) error {
	var err error
	log_opt.Level, err = c.GetString("log", "level")
	if err != nil {
		return err
	}
	log_opt.DeviceName, err = c.GetString("log", "diveName")
	if err != nil {
		return err
	}
	log_opt.MaxSize, err = c.GetInt("log", "maxSize")
	if err != nil {
		return err
	}
	return nil
}
