package typed

import "k8s.io/client-go/kubernetes"

type K8sClient struct {
	Namespace string
	Client    *kubernetes.Clientset
}
