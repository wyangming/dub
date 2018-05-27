package main

import (
	"flag"
	"dub/app/web/usecenter"
)

func main() {
	config_path := flag.String("Config", "./config/web/web_center_use.cfg", "config file path!")
	flag.Parse()

	webUseCenter := usecenter.NewWebUseCenterServer()
	webUseCenter.Init(*config_path)
}
