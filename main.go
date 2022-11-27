package main

import (
	"flag"
	"year-end/config"
	"year-end/internal"
	"year-end/model"
)

var port = flag.String("port", "8080", "http port")

const (
	uname = ""
	psd   = ""
)

func main() {
	// init
	//flag.Parse()
	cfg := &config.Config{}
	cfg.Init("./config/config.yaml")
	model.Start()
	// run engine - > route
	internal.StartHTTP()
	
	// run server
	//internal.StartHTTP(*port)
}
