package main

import (
	"dub/app/reg"
	"flag"
)

func main() {
	config_path := flag.String("Config", "./config/reg_server.cfg", "config file path!")
	flag.Parse()
	regServer := reg.NewDbServer()
	regServer.Init(*config_path)

}
