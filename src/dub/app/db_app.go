package main

import (
	"dub/app/db"
)

func main() {
	dbServer := db.NewDbServer()
	dbServer.Init()
}
