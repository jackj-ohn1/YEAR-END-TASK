package main

import (
	"context"
	"flag"
	"year-end/crawler"
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
	ctx, _ := context.WithCancel(context.Background())
	model.Start(ctx)
	
	// run engine - > route
	crawler.Crawler(ctx, uname, psd)
	
	// run server
	//internal.StartHTTP(*port)
}
