package controller

import "github.com/fudali113/good-job/pkg/apis/goodjob"

const (
	GoodJobNameLabel       = goodjob.GroupName + "/goodJobName"
	PipelineNameLabel      = goodjob.GroupName + "/pipelineName"
	ShardLabel             = goodjob.GroupName + "/shard"
	ShardIndexLabel        = goodjob.GroupName + "/shardIndex"
	ShardMatchPatternLabel = goodjob.GroupName + "/shardMatchPattern"
)

type labelInfo struct {
	pipeline string
	goodJob string
	shard string
	shardIndex int
	shardMatchPattern string
}

