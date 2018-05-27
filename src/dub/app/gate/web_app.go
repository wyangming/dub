package main

import (
	"flag"
	"dub/app/gate/web"
)

func main() {
	config_path := flag.String("Config", "./config/gate/gate_web.cfg", "config file path!")
	flag.Parse()

	gateWeb := web.NewGateWebServer()
	gateWeb.Init(*config_path)
}
