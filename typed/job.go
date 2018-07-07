package typed

import (
	k8sType "github.com/fudali113/good-job/pkg/apis/goodjob/v1alpha1"
)

// GoodJob 执行 GoodJob 的配置
type Job = k8sType.GoodJobSpec

// GoodJobShard 分片程序的配置
type ShardConfig = k8sType.GoodJobShard

type JobStatus = k8sType.GoodJobStatus

type Pipeline = k8sType.PipelineSpec

type CronTrigger = k8sType.CronTriggerSpec
