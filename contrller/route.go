package contrller

import (
	"github.com/emicklei/go-restful"
)

func init() {
	wsContainer := restful.NewContainer()
	register(wsContainer, JobResource{})
	register(wsContainer, PipelineResource{})
}

func register(container *restful.Container, register RouteRegister) {
	register.Register(container)
}

// RouteRegister 帮助我们注册相关的路由
type RouteRegister interface {
	Register(container *restful.Container)
}
