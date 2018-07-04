package main

import (
	"flag"
	"log"

	"github.com/fudali113/good-job/typed"
	"github.com/fudali113/good-job/api"
	"github.com/fudali113/good-job/pkg/signals"
	"github.com/fudali113/good-job/controller"
)

var port int

func main() {
	flag.Parse()
	config := typed.RunConfig{
		Server: typed.ServerConfig{
			Port: port,
		},
		Runtime:typed.RuntimeConfig{

		},
	}
	go controller.Start(config.Runtime)
	go api.Start(config.Server)
	stop := signals.SetupSignalHandler()
	<- stop
	log.Printf("程序收到停止信号终止运行")
}

func init() {
	flag.IntVar(&port, "port", 3333, "启动应用的 server 所监听的端口")
}
