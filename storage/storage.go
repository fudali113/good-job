package storage

import (
	"github.com/fudali113/good-job/typed"
)

type JobStorage interface {
	Get(name string) (typed.Job, error)
	Create(job typed.Job) error
	Update(job typed.Job) error
	Delete(name string) error
}

type PipelineStorage interface {
	Get(name string) (typed.Pipeline, error)
	Create(job typed.Pipeline) error
	Update(job typed.Pipeline) error
	Delete(name string) error
}

type Storage struct {
	JobStorage      JobStorage
	PipelineStorage PipelineStorage
}

// CreateStorage 创建一个 Storage
func CreateStoage() *Storage {
	return &Storage{
		JobStorage:      nil,
		PipelineStorage: nil,
	}
}
