package main

import (
	"dub/app/web/manlobby"
	"flag"
)

//后台管理大厅
func main() {
	config_path := flag.String("Config", "./config/web/web_center_man_lobby.cfg", "config file path!")
	flag.Parse()

	webUseCenter := manlobby.NewWebManLobbyCenterServer()
	webUseCenter.Init(*config_path)
}
