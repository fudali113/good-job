package main

import (
	"flag"

	"github.com/fudali113/good-job/contrller"
	"github.com/fudali113/good-job/typed"
)

var port int

func main() {
	flag.Parse()
	config := typed.RunConfig{
		Server: typed.Server{
			Port: port,
		},
	}
	contrller.Start(config)
}

func init() {
	flag.IntVar(&port, "port", 3333, "启动应用的 server 所监听的端口")
}
