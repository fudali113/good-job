package controller

import (
	batchv1 "k8s.io/api/batch/v1"
	"k8s.io/api/core/v1"
	"github.com/fudali113/good-job/pkg/apis/goodjob/v1alpha1"
	"github.com/fudali113/good-job/typed"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"fmt"
	"github.com/fudali113/good-job/pkg/apis/goodjob"
	"github.com/golang/glog"
	"encoding/json"
)

func addGoodJob(obj interface{}) {

}

func updateGoodJob(oldObj, newObj interface{})  {
	info, _ := json.Marshal(oldObj)
	glog.Info(string(info))
	info, _ = json.Marshal(newObj)
	glog.Info(string(info))
	oldGoodjob := oldObj.(*v1alpha1.Job)
	newGoodjob := newObj.(*v1alpha1.Job)
	if newGoodjob.ResourceVersion == oldGoodjob.ResourceVersion {
		// Periodic resync will send update events for all known Deployments.
		// Two different versions of the same Deployment will always have different RVs.
		return
	}
	oldStatus := oldGoodjob.Status.Status
	newStatus := newGoodjob.Status.Status
	if oldStatus <= newStatus {
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
			_, err := clientset.GoodjobV1alpha1().Jobs(newGoodjob.Namespace).UpdateStatus(goodjob)
			if err != nil {
				glog.Errorf("更新 GoodJob 状态失败, error: %s", err.Error())
			}
		case "exec":
			goodjob := newGoodjob.DeepCopy()
			job := newJob(goodjob.Spec.Shard.Template, goodjob, true)
			_, err := clientset.BatchV1().Jobs(newGoodjob.Namespace).Create(&job)
			if err != nil {
				glog.Errorf("创建 Job 失败, error: %s", err.Error())
				goodjob.Status.Status = typed.ShardFail
			} else {
				goodjob.Status.Status = typed.Sharding
			}
			clientset.GoodjobV1alpha1().Jobs(newGoodjob.Namespace).UpdateStatus(goodjob)
		}
	case typed.ShardSuccess:

	}
}

func newJob(podTemplate v1.PodTemplateSpec, job *v1alpha1.Job, shard bool) batchv1.Job {
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
