package main

import (
	"config"
	"flag"
	"fmt"
	"os"
	"utils"
	"app/db"
)

func main() {
	//初始化日志
	log := utils.NewLogger()

	config_path := flag.String("Config", "./config/db_server.cfg", "config file path!")
	flag.Parse()
	opt, log_opt, err := config.GetDatabaseServerConfig(*config_path)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	log.SetConfig(log_opt)

	log.Infoln("db_server will be start")
	log.Infoln(opt)
	dbProxy := db.NewDbProxy()
	dbProxy.Init(&opt)
}
