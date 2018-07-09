package controller

import (
	"k8s.io/api/core/v1"

	batchv1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"bytes"
	"fmt"
	"github.com/emicklei/go-restful/log"
	"github.com/fudali113/good-job/typed"
	"regexp"
)

func getJobLog(namespace, name string, opt *v1.PodLogOptions) ([]byte, error) {
	return clientset.CoreV1().Pods(namespace).GetLogs(name, opt).Do().Raw()
}

func findShardInfo(bytes [][]byte, matchShardInfoFunc func([]byte) bool) []string {
	shards := make([]string, 0, 8)
	for i := range bytes {
		if matchShardInfoFunc(bytes[i]) {
			shards = append(shards, string(bytes[i]))
		}
	}
	return shards
}

func jobStatusUpdate(oldObj, newObj interface{}) {
	oldJob := oldObj.(*batchv1.Job)
	newJob := newObj.(*batchv1.Job)
	if newJob.ResourceVersion == oldJob.ResourceVersion {
		// Periodic resync will send update events for all known Deployments.
		// Two different versions of the same Deployment will always have different RVs.
		return
	}
	if _, ok := newJob.Labels[ShardIndexLabel]; ok && newJob.Labels[ShardIndexLabel] == "0" {
		// 状态扭转的时候触发
		if newJob.Status.Succeeded == 1 && oldJob.Status.Succeeded == 0 {
			pods, err := clientset.CoreV1().Pods(newJob.Namespace).List(metav1.ListOptions{
				LabelSelector: fmt.Sprintf("job-name=%s", newJob.Name),
			})
			if err != nil {
				log.Printf("job 执行成功获取相关的 pods 出错， error: %s", err.Error())
			}
			if len(pods.Items) == 0 {
				log.Printf("不可思议，pods 数量为 0")
			}
			pod := pods.Items[0]
			// FIXME
			// 直接读取的话可能会消耗比较大的内存，建议所有 Shard 运行时不要打印太多的日志
			byteSlice, err := getJobLog(pod.Namespace, pod.Name, &v1.PodLogOptions{})
			if err != nil {
				log.Printf("job 执行成功获取 pod 的相关的日志出错， error: %s", err.Error())
			}
			logs := bytes.Split(byteSlice, []byte("\n"))

			// 找到用户设置的匹配参数
			shardMatchPattern := func(annotations map[string]string) string {
				shardMatchPattern := annotations[ShardMatchPatternLabel]
				if shardMatchPattern == "" {
					shardMatchPattern = "GOOD_JOB_SHARDS[\\s\\S\\w\\W]+"
				}
				return shardMatchPattern
			}(newJob.Annotations)

			// 找到日志中的 Shard 信息
			shards := findShardInfo(logs, func(byteArray []byte) bool {
				match, err :=  regexp.Match(shardMatchPattern, byteArray)
				if err != nil {
					log.Printf("匹配 shard 数据出错; error: %s", err.Error())
					return false
				}
				return match
			})

			goodJobName := newJob.Labels[GoodJobNameLabel]

			goodJob, err := clientset.GoodjobV1alpha1().GoodJobs(newJob.Namespace).Get(goodJobName, metav1.GetOptions{})
			if err != nil {
				log.Printf("shards 日志获取完成，寻找对应的 GoodJob 失败, name: %s; error: %s", goodJobName, err.Error())
			}

			goodJob = goodJob.DeepCopy()
			goodJob.Status.Status = typed.ShardSuccess
			goodJob.Status.Shards = shards

			_, err = clientset.GoodjobV1alpha1().GoodJobs(newJob.Namespace).Update(goodJob)
			if err != nil {
				log.Printf("更新 GoodJob 状态失败，error: %s", err.Error())
			}

			log.Printf("shard 成功, shards is %v", shards)

		}
	}
}
