package controller

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"fmt"
	"strconv"
)

func SwitchLabelInfo(object v1.Object) (info labelInfo, err error) {
	labels := object.GetLabels()
	annotations := object.GetAnnotations()
	goodJob, ok := labels[GoodJobNameLabel]
	if !ok {
		err = fmt.Errorf("%s filed must in labels", GoodJobNameLabel)
		return
	}
	info.shardIndex, err = strconv.Atoi(labels[ShardIndexLabel])
	if err != nil {
		err = fmt.Errorf("%s filed must in labels and must be int, error: %s", ShardIndexLabel, err.Error())
		return
	}
	info.goodJob = goodJob
	info.pipeline = labels[PipelineNameLabel]
	info.shard = annotations[ShardLabel]
	info.shardMatchPattern = annotations[ShardMatchPatternLabel]
	return
}

func duplicateCheck(datas []string) error {
	for i := 0; i<len(datas); i++ {
		masterValue := datas[i]
		for j := i + 1; j<len(datas); j++ {
			value := datas[j]
			if masterValue == value {
				return fmt.Errorf("存在重复的数据，master: %d, branch: %d", i, j)
			}
		}
	}
	return nil
}