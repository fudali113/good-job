package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/emicklei/go-restful"
	"github.com/fudali113/good-job/typed"
)

var wsContainer = restful.NewContainer()

func init() {
	register(wsContainer, JobResource{})
	register(wsContainer, PipelineResource{})
}

func Start(config typed.ServerConfig) {
	port := fmt.Sprintf(":%d", config.Port)
	log.Printf("start listening on localhost %s", port)
	server := &http.Server{Addr: port, Handler: wsContainer}
	log.Fatal(server.ListenAndServe())
}

func register(container *restful.Container, register RouteRegister) {
	register.Register(container)
}

// RouteRegister 帮助我们注册相关的路由
type RouteRegister interface {
	Register(container *restful.Container)
}
