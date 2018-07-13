package controller

import (
	"os"
	"time"
	"path/filepath"
	"encoding/json"

	"github.com/golang/glog"
	"github.com/fudali113/good-job/pkg/client/clientset/versioned"
	"github.com/fudali113/good-job/typed"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/cache"

	kubeinformers "k8s.io/client-go/informers"
	goodinformers "github.com/fudali113/good-job/pkg/client/informers/externalversions"
)

var clientset *typed.Clientset

// Start 根据 Config 运行 controller
func Start(config typed.RuntimeConfig, stopCh <-chan struct{}) {

	goodInformersFactory := goodinformers.NewSharedInformerFactoryWithOptions(
		clientset.GoodJobClientset,
		1*time.Second,
		goodinformers.WithNamespace("good-job"))
	kubeInformersFactory := kubeinformers.NewSharedInformerFactoryWithOptions(
		clientset.Clientset,
		1*time.Second,
		kubeinformers.WithNamespace("good-job"))

	googjobInformer := goodInformersFactory.Goodjob().V1alpha1().GoodJobs()

	googjobInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    addGoodJob,
		UpdateFunc: updateGoodJob,
		DeleteFunc: func(obj interface{}) {
			info, _ := json.Marshal(obj)
			glog.Info(string(info))
		},
	})

	jobInformer := kubeInformersFactory.Batch().V1().Jobs()

	jobInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			info, _ := json.Marshal(obj)
			glog.V(1).Info(string(info))
		},
		UpdateFunc: jobStatusUpdate,
		DeleteFunc: func(obj interface{}) {
			info, _ := json.Marshal(obj)
			glog.Info(string(info))
		},
	})

	goodInformersFactory.Start(stopCh)
	kubeInformersFactory.Start(stopCh)


	goodJobController := NewGoodJobController(
		clientset.Clientset,
		clientset.GoodJobClientset,
		kubeInformersFactory.Batch().V1().Jobs(),
		goodInformersFactory.Goodjob().V1alpha1().GoodJobs())

	if err := goodJobController.Run(2, stopCh); err != nil {
		glog.Fatalf("Error running goodJobController: %s", err.Error())
	}

}

func init() {
	clientset = CreateClientset("")
}

func CreateClientset(token string) *typed.Clientset {
	clientset, err := CreateOriginClientset(CreteConfig(""))
	if err != nil {
		panic(err)
	}
	goodJobClientset, err := CreateGoodJobClientset(CreteConfig(""))
	if err != nil {
		panic(err)
	}
	return &typed.Clientset{
		Clientset:        clientset,
		GoodJobClientset: goodJobClientset,
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
	config, err := rest.InClusterConfig()
	if err == nil {
		return config, err
	}
	config, err = clientcmd.BuildConfigFromFlags("", kubeConfigPath())
	if err != nil {
		return nil, err
	}
	if token != "" {
		config.BearerToken = token
	}
	return config, err
}

func kubeConfigPath() string {
	return filepath.Join(homeDir(), ".kube", "config")
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}
