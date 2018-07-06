package controller

import (
	"time"

	"github.com/fudali113/good-job/typed"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"github.com/fudali113/good-job/pkg/client/clientset/versioned"

	kubeinformers "k8s.io/client-go/informers"
	informers "github.com/fudali113/good-job/pkg/client/informers/externalversions"
	"k8s.io/client-go/tools/cache"
)

var clientset typed.Clientset

// Start 根据 Config 运行 controller
func Start(config typed.RuntimeConfig, stop <- chan struct{})  {

	goodInformers := informers.NewSharedInformerFactory(clientset.GoodJobClientset, 1 * time.Second)
	kubeInformers := kubeinformers.NewSharedInformerFactory(clientset.Clientset, 1 * time.Second)

	googjobInformer := goodInformers.Goodjob().V1alpha1().Jobs()

	googjobInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {

		},
		UpdateFunc: func(oldObj, newObj interface{}) {

		},
		DeleteFunc: func(obj interface{}) {

		},
	})

	jobInformer := kubeInformers.Batch().V1().Jobs()

	jobInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {

		},
		UpdateFunc: func(oldObj, newObj interface{}) {

		},
		DeleteFunc: func(obj interface{}) {

		},
	})

	goodInformers.Start(stop)
	kubeInformers.Start(stop)

}

func CreateClientset(token string) *typed.Clientset {
	return &typed.Clientset{
		Clientset: nil,
		GoodJobClientset: nil,
	}
}

func CreateOriginClientset(config *rest.Config, err error) (*kubernetes.Clientset, error) {
	if err != nil {
		return nil, err
	}
	return kubernetes.NewForConfig(config)
}

func CreateGoodJobClientset(config *rest.Config, err error) (*versioned.Clientset, error) {
	if err != nil {
		return nil, err
	}
	return versioned.NewForConfig(config)
}

func CreteConfig(token string) (*rest.Config, error) {
	config, err := clientcmd.BuildConfigFromFlags("", "")
	if err != nil {
		return nil, err
	}
	if token != "" {
		config.BearerToken = token
	}
	return config, err
}

