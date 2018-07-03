package storage

import (
	"github.com/fudali113/good-job/typed"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Job struct {
	metav1.TypeMeta `json:",inline"`
	// Standard object's metadata.
	// More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	Spec typed.Job `json:"spec,omitempty"`

	Status typed.JobStatus `json:"status,omitempty"`
}

type K8sJobStorage struct {
	typed.K8sClient
}

func (k K8sJobStorage) Get(name string) (typed.Job, error) {
	job := Job{}
	k.Client.RESTClient().Get().Namespace(k.Namespace).Resource("GoodJob").Name(name).Do().Into(&job)
}
