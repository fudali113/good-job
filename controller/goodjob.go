package controller

import (
	batchv1 "k8s.io/api/batch/v1"
	"k8s.io/api/core/v1"
	"github.com/fudali113/good-job/pkg/apis/goodjob/v1alpha1"
	"github.com/fudali113/good-job/typed"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"fmt"
	"github.com/fudali113/good-job/pkg/apis/goodjob"
	"encoding/json"
	"log"
)

func addGoodJob(obj interface{}) {
	info, _ := json.Marshal(obj)
	log.Printf("create ---> " + string(info))
}

func updateGoodJob(oldObj, newObj interface{})  {
	oldGoodjob := oldObj.(*v1alpha1.GoodJob)
	newGoodjob := newObj.(*v1alpha1.GoodJob)
	if newGoodjob.ResourceVersion == oldGoodjob.ResourceVersion {
		// Periodic resync will send update events for all known Deployments.
		// Two different versions of the same Deployment will always have different RVs.
		return
	}

	info, _ := json.Marshal(oldObj)
	log.Printf("old ---> " + string(info))
	info, _ = json.Marshal(newObj)
	log.Printf("new ---> " +string(info))

	oldStatus := oldGoodjob.Status.Status
	newStatus := newGoodjob.Status.Status
	if oldStatus >= newStatus {
		return
	}
	switch newStatus {
	case typed.Begin:
		shard := oldGoodjob.Spec.Shard
		switch shard.Type {
		case "config":
			goodjob := newGoodjob.DeepCopy()
			goodjob.Status.Status = typed.ShardSuccess
			goodjob.Status.Shards = shard.Shards
			// FIXME
			// 本来应该使用 UpdateStatus 方法进行更新状态
			// 1.10 才开始支持 CRD 定义 subresource, 如果在
			// https://v1-10.docs.kubernetes.io/docs/tasks/access-kubernetes-api/extend-api-custom-resource-definitions/#status-subresource
			// 文档对比
			// https://v1-10.docs.kubernetes.io/docs/concepts/extend-kubernetes/api-extension/custom-resources/#advanced-features-and-flexibility
			// https://v1-9.docs.kubernetes.io/docs/concepts/api-extension/custom-resources/#advanced-features-and-flexibility
			_, err := clientset.GoodjobV1alpha1().GoodJobs(newGoodjob.Namespace).Update(goodjob)
			if err != nil {
				log.Printf("update GoodJob status error, error: %s", err.Error())
			}
		case "exec":
			goodjob := newGoodjob.DeepCopy()
			job := newJob(goodjob.Spec.Shard.Template, goodjob, true)
			_, err := clientset.BatchV1().Jobs(newGoodjob.Namespace).Create(&job)
			if err != nil {
				log.Printf("创建 GoodJob 失败, error: %s", err.Error())
				goodjob.Status.Status = typed.ShardFail
			} else {
				goodjob.Status.Status = typed.Sharding
			}
			clientset.GoodjobV1alpha1().GoodJobs(newGoodjob.Namespace).UpdateStatus(goodjob)
		}
	case typed.ShardSuccess:

	}
}

func newJob(podTemplate v1.PodTemplateSpec, job *v1alpha1.GoodJob, shard bool) batchv1.Job {
	action := func(shard bool) string {
		if shard {
			return "shard"
		}
		return "run"
	}(shard)
	labels := map[string]string{
		goodjob.GroupName + "/goodjobName": job.Name,
		goodjob.GroupName + "/pipelineName": job.Status.Pipeline,
		goodjob.GroupName + "/action": action,
	}
	return batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:fmt.Sprintf("good-job-%s-%s", job.Name, action),
			Namespace:job.Namespace,
			Labels:labels,
		},
		Spec: batchv1.JobSpec{
				Template: podTemplate,
		},
	}
}
