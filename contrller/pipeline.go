package contrller

import (
	"github.com/emicklei/go-restful"
)

type PipelineResource struct{}

func (p PipelineResource) Register(container *restful.Container) {
	ws := new(restful.WebService)
	ws.
		Path("/pipelines").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	ws.Route(ws.POST("").To(p.Post))
	ws.Route(ws.GET("/{id}").To(p.GetOne))
	ws.Route(ws.PUT("/{id}").To(p.Update))
	ws.Route(ws.DELETE("/{id}").To(p.Delete))

	container.Add(ws)
}

func (p PipelineResource) Post(req *restful.Request, resp *restful.Response) {

}

func (p PipelineResource) GetOne(req *restful.Request, resp *restful.Response) {
	id := req.PathParameter("id")
	resp.WriteEntity(id)
}

func (p PipelineResource) Update(req *restful.Request, resp *restful.Response) {

}

func (p PipelineResource) Delete(req *restful.Request, resp *restful.Response) {

}
