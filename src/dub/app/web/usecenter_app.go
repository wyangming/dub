package main

import "flag"

func main()  {
	//beego.Run(":81")
	config_path := flag.String("Config", "./config/web/use_center.cfg", "config file path!")
	flag.Parse()
	println(config_path)
}