package storage

import (
	"github.com/fudali113/good-job/typed"
	"github.com/fudali113/good-job/pkg/apis/goodjob/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type K8sJobStorage struct {
	*typed.Clientset
}

func (k K8sJobStorage) Get(name string) (jobSpec typed.Job, err error) {
	job, err := k.
		GoodjobV1alpha1().
		Jobs(k.Namespace).
		Get(name, metav1.GetOptions{})
	jobSpec = job.Spec
	return
}

func (k K8sJobStorage) Create(jobSpec typed.Job) (err error) {
	_, err = k.
		GoodjobV1alpha1().
		Jobs(k.Namespace).
		Create(createJob(jobSpec, map[string]string{}))
	return
}

func (k K8sJobStorage) Update(jobSpec typed.Job) (err error) {
	_,err = k.
		GoodjobV1alpha1().
		Jobs(k.Namespace).
		Update(createJob(jobSpec, map[string]string{}))
	return
}

func (k K8sJobStorage) Delete(name string) (err error) {
	return k.
		GoodjobV1alpha1().
		Jobs(k.Namespace).
		Delete(name, &metav1.DeleteOptions{})
}

func createJob(jobSpec typed.Job, labels map[string]string) *v1alpha1.Job {
	return &v1alpha1.Job{
		TypeMeta: metav1.TypeMeta{
			Kind: typed.JobCrdName,
			APIVersion: v1alpha1.SchemeGroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:jobSpec.Name,
			Labels: labels,
		},
		Spec:jobSpec,
	}
}