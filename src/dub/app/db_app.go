package main

import (
	"dub/app/db"
	"flag"
)

func main() {
	config_path := flag.String("Config", "./config/db_server.cfg", "config file path!")
	flag.Parse()

	dbServer := db.NewDbServer()
	dbServer.Init(*config_path)
}
