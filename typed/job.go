package typed

import (
	k8sType "github.com/fudali113/good-job/pkg/apis/goodjob/v1alpha1"
)

// Job 执行 Job 的配置
type Job = k8sType.JobSpec

// ExecConfig 执行 job 程序的配置
type ExecConfig = k8sType.ExecConfig

// ShardConfig 分片程序的配置
type ShardConfig = k8sType.ShardConfig

type JobStatus = k8sType.JobStatus

type Pipeline = k8sType.PipelineSpec
