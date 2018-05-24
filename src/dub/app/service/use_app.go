package main

import (
	"flag"
	"dub/app/service/use"
)

func main() {
	config_path := flag.String("Config", "./config/service/use_server.cfg", "config file path!")
	flag.Parse()

	serviceUseServer := use.NewServiceUseServer()
	serviceUseServer.Init(*config_path)
}
