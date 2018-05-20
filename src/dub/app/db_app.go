package main

import (
	"dub/app/db"
	"flag"
	"sync"
)

func main() {
	config_path := flag.String("Config", "./config/db_server.cfg", "config file path!")
	flag.Parse()

	dbServer := db.NewDbServer()
	dbServer.Init(*config_path)

	waite := &sync.WaitGroup{}
	waite.Add(1)
	waite.Wait()
}
