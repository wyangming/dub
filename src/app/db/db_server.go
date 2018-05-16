package db
import (
	"config"
	"flag"
	"fmt"
	"os"
	"utils"
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

	log.Infoln("test log")
	log.Infoln(opt)
}