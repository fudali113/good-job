package contrller

import (
	"github.com/emicklei/go-restful"
	"github.com/fudali113/good-job/storage"
)

type JobResource struct {
}

func (j JobResource) Register(container *restful.Container) {
	ws := new(restful.WebService)
	ws.
		Path("/v1/jobs").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	ws.Route(ws.GET("").To(j.Get))
	ws.Route(ws.POST("").To(j.Post))
	ws.Route(ws.GET("/{id}").To(j.GetOne))
	ws.Route(ws.PUT("/{id}").To(j.Update))
	ws.Route(ws.DELETE("/{id}").To(j.Delete))

	container.Add(ws)
}

func (j JobResource) Get(req *restful.Request, resp *restful.Response) {
}

func (j JobResource) Post(req *restful.Request, resp *restful.Response) {

}

func (j JobResource) GetOne(req *restful.Request, resp *restful.Response) {
	id := req.PathParameter("id")
	storage := storage.CreateStoage()
	job, _ := storage.JobStorage.Get(id)
	resp.WriteEntity(job)
}

func (j JobResource) Update(req *restful.Request, resp *restful.Response) {

}

func (j JobResource) Delete(req *restful.Request, resp *restful.Response) {

}
