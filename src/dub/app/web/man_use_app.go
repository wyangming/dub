package main

import (
	"dub/app/web/manuse"
	"flag"
)

//用户中心
func main() {
	config_path := flag.String("Config", "./config/web/web_center_man_use.cfg", "config file path!")
	flag.Parse()

	webUseCenter := manuse.NewWebManUseCenterServer()
	webUseCenter.Init(*config_path)
}
