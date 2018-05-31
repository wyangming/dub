package main

import (
	"flag"
	"dub/app/web/manlobby"
)

//后台管理大厅
func main() {
	config_path := flag.String("Config", "./config/web/web_center_man_lobby.cfg", "config file path!")
	flag.Parse()

	webUseCenter := manlobby.NewWebUseCenterServer()
	webUseCenter.Init(*config_path)
}
