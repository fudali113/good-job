package typed

import (
	"github.com/fudali113/good-job/pkg/client/clientset/versioned"
	"k8s.io/client-go/kubernetes"
)

type Clientset struct {
	*kubernetes.Clientset
	*GoodJobClientset
	Namespace string
}

type GoodJobClientset = versioned.Clientset
