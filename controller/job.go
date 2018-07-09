package controller

import (
	"k8s.io/api/core/v1"

	batchv1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"fmt"
	"github.com/emicklei/go-restful/log"
)

func getJobLog(namespace, name string, opt *v1.PodLogOptions) ([]byte, error) {
	return clientset.CoreV1().Pods(namespace).GetLogs(name, opt).Do().Raw()
}

func findShardInfo(bytes []byte, matchShardInfoFunc func(string) bool) []string {
	return []string{}
}

func jobStatusUpdate(oldObj, newObj interface{}) {
	oldJob := oldObj.(*batchv1.Job)
	newJob := newObj.(*batchv1.Job)
	if newJob.ResourceVersion == oldJob.ResourceVersion {
		// Periodic resync will send update events for all known Deployments.
		// Two different versions of the same Deployment will always have different RVs.
		return
	}
	if _, ok := newJob.Labels[ShardIndexLabel]; ok /*&& newJob.Labels[ShardIndexLabel] == "-1" */{
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
			bytes, err := getJobLog(pod.Namespace, pod.Name, &v1.PodLogOptions{})
			if err != nil {
				log.Printf("job 执行成功获取 pod 的相关的日志出错， error: %s", err.Error())
			}
			log.Printf(string(bytes))
		}
	}
}
