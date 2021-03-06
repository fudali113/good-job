package main

import (
	"flag"

	"github.com/golang/glog"
	"github.com/fudali113/good-job/api"
	"github.com/fudali113/good-job/controller"
	"github.com/fudali113/good-job/pkg/signals"
	"github.com/fudali113/good-job/typed"
)

var port int

func main() {
	defer glog.Flush()

	flag.Parse()
	config := typed.RunConfig{
		Server: typed.ServerConfig{
			Port: port,
		},
		Runtime: typed.RuntimeConfig{},
	}

	stop := signals.SetupSignalHandler()
	go controller.Start(config.Runtime, stop)
	go api.Start(config.Server, stop)
	<-stop

	glog.Info("程序收到停止信号终止运行")
}

func init() {
	flag.IntVar(&port, "port", 3333, "启动应用的 server 所监听的端口")
}