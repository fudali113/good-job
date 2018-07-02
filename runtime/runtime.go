package runtime

import (
	"context"
	"github.com/fudali113/good-job/storage"
)

type JobRuntime struct {
	storage.Job
	Status int           `json:"status"`
	Logs   []interface{} `json:"logs"`
}

type PipelineRuntime struct {
	storage.Pipeline
	CurrentIndex int           `json:"currentIndex"`
	Status       int           `json:"status"`
	Logs         []interface{} `json:"logs"`
}

type Runtime interface {
	RunJob(job storage.Job) (context.CancelFunc, error)
	RunPipeline(pipeline storage.Pipeline) (context.CancelFunc, error)
	CreateCron(resource, id, token string) error
}
