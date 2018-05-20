package main

import (
	"dub/app/reg"
	"flag"
	"fmt"
)

func main() {
	config_path := flag.String("Config", "./config/reg_server.cfg", "config file path!")
	flag.Parse()
	fmt.Println(config_path)
	regServer := reg.NewDbServer()
	regServer.Init(*config_path)

}
