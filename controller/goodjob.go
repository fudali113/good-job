package controller

import (
	"encoding/json"
	"fmt"
	"github.com/fudali113/good-job/pkg/apis/goodjob"
	"github.com/fudali113/good-job/pkg/apis/goodjob/v1alpha1"
	"github.com/fudali113/good-job/typed"
	batchv1 "k8s.io/api/batch/v1"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"log"
)

func addGoodJob(obj interface{}) {
	info, _ := json.Marshal(obj)
	log.Printf("create ---> " + string(info))
}

func updateGoodJob(oldObj, newObj interface{}) {
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
	log.Printf("new ---> " + string(info))

	oldStatus := oldGoodjob.Status.Status
	newStatus := newGoodjob.Status.Status
	if oldStatus >= newStatus {
		return
	}
	switch newStatus {
	case typed.Begin:
		goodjob := newGoodjob.DeepCopy()
		shard := goodjob.Spec.Shard
		switch shard.Type {
		case "config":
			goodjob.Status.Status = typed.ShardSuccess
			goodjob.Status.Shards = shard.Shards
		case "exec":
			goodjob := newGoodjob.DeepCopy()
			job := newJob(goodjob.Spec.Shard.Template, goodjob, "")
			_, err := clientset.BatchV1().Jobs(newGoodjob.Namespace).Create(job)
			if err != nil {
				log.Printf("创建 GoodJob 失败, error: %s", err.Error())
				goodjob.Status.Status = typed.ShardFail
			} else {
				goodjob.Status.Status = typed.Sharding
			}
		}
		goodjob.Status.ShardStatuses = map[string]string{}
		goodjob.Status.Logs = []string{}
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
	case typed.ShardSuccess:
		dealShardSuccess(newGoodjob)
	}
}

func dealShardSuccess(goodJob *v1alpha1.GoodJob) error {
	for _, shard := range goodJob.Status.Shards {
		goodjob := goodJob.DeepCopy()
		podTemplate := goodjob.Spec.Template
		containers := podTemplate.Spec.Containers
		for i := range containers {
			c := containers[i]
			c.Env = append(c.Env, v1.EnvVar{
				Name:  "GOOD_JOB_SHARD",
				Value: shard,
			})
			containers[i] = c
		}

		job := newJob(podTemplate, goodjob, shard)

		// 在创建 job 前检查是否状态为可重新执行
		nowGoodJob, err := clientset.GoodjobV1alpha1().GoodJobs(goodjob.Namespace).Get(goodjob.Name, metav1.GetOptions{})
		if err != nil {
			log.Printf(err.Error())
		}
		if _, ok := nowGoodJob.Status.ShardStatuses[shard]; ok {
			continue
		}

		_, err = clientset.BatchV1().Jobs(goodJob.Namespace).Create(job)
		if err != nil {
			log.Printf(err.Error())
		}
		goodjob.Status.ShardStatuses[shard] = "create"
		clientset.GoodjobV1alpha1().GoodJobs(goodjob.Namespace).Update(goodjob)
	}
	return nil
}

func newJob(podTemplate v1.PodTemplateSpec, job *v1alpha1.GoodJob, shard string) *batchv1.Job {
	action := func(shard string) string {
		if shard == "" {
			return "shard"
		}
		return "run"
	}(shard)
	labels := map[string]string{
		goodjob.GroupName + "/goodJobName":  job.Name,
		goodjob.GroupName + "/pipelineName": job.Status.Pipeline,
		goodjob.GroupName + "/action":       action,
		goodjob.GroupName + "/shard":        shard,
	}
	return &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("good-job-%s-%s-%s", job.Name, action, shard),
			Namespace: job.Namespace,
			Labels:    labels,
		},
		Spec: batchv1.JobSpec{
			Template: podTemplate,
		},
	}
}
