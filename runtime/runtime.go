package runtime

import (
	"github.com/fudali113/good-job/typed"
	"k8s.io/api/batch/v1"
	"k8s.io/apimachinery/pkg/watch"
)

type JobRuntime struct {
	typed.Job
	Status int           `json:"status"`
	Logs   []interface{} `json:"logs"`
}

type PipelineRuntime struct {
	typed.Pipeline
	CurrentIndex int           `json:"currentIndex"`
	Status       int           `json:"status"`
	Logs         []interface{} `json:"logs"`
}

type Runtime interface {
	CreateJob(job typed.Job, exec typed.ExecConfig) (k8sJob *v1.Job, err error)
	WatchJob(name string) (watch watch.Interface, err error)
	CreateCronJob(resource, id, token string) error
}

func Start(config typed.RuntimeConfig)  {
	
}
