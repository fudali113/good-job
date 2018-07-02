package runtime

import (
	"github.com/fudali113/good-job"
	"k8s.io/api/batch/v1"
	"k8s.io/apimachinery/pkg/watch"
)

type JobRuntime struct {
	goodjob.Job
	Status int           `json:"status"`
	Logs   []interface{} `json:"logs"`
}

type PipelineRuntime struct {
	goodjob.Pipeline
	CurrentIndex int           `json:"currentIndex"`
	Status       int           `json:"status"`
	Logs         []interface{} `json:"logs"`
}

type Runtime interface {
	CreateJob(job goodjob.Job, exec goodjob.ExecConfig) (k8sJob *v1.Job, err error)
	WatchJob(name string) (watch watch.Interface, err error)
	CreateCronJob(resource, id, token string) error
}
