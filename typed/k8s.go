package typed

import (
	"github.com/fudali113/good-job/pkg/client/clientset/versioned"
)

type K8sClient struct {
	Namespace string
	Client    *versioned.Clientset
}
