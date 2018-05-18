package main

import (
	"flag"
	"fmt"
)

func main()  {
	config_path := flag.String("Config", "./config/reg_server.cfg", "config file path!")
	flag.Parse()
	fmt.Println(config_path)
}