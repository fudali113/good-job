package controller

import (
	"fmt"
	"strconv"
	"encoding/json"


	"github.com/golang/glog"
	"github.com/fudali113/good-job/typed"
	"github.com/fudali113/good-job/pkg/apis/goodjob/v1alpha1"
	"k8s.io/api/core/v1"
	"github.com/fudali113/good-job/pkg/client/clientset/versioned"

	btachlisters "k8s.io/client-go/listers/batch/v1"
	goodlisters "github.com/fudali113/good-job/pkg/client/listers/goodjob/v1alpha1"
	batchinformers "k8s.io/client-go/informers/batch/v1"
	goodfinformers "github.com/fudali113/good-job/pkg/client/informers/externalversions/goodjob/v1alpha1"

	goodscheme "github.com/fudali113/good-job/pkg/client/clientset/versioned/scheme"

	corev1 "k8s.io/api/core/v1"
	typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"




	batchv1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"time"
	"k8s.io/apimachinery/pkg/api/errors"
	"github.com/fudali113/good-job/pkg/apis/goodjob"
	"bytes"
	"regexp"
)


// GoodJobController is the controller implementation for Foo resources
type GoodJobController struct {
	// kubeclientset is a standard kubernetes clientset
	kubeclientset kubernetes.Interface
	// sampleclientset is a clientset for our own API group
	goodclientset versioned.Interface

	jobsLister     btachlisters.JobLister
	jobsSynced     cache.InformerSynced
	goodJobsLister goodlisters.GoodJobLister
	goodJobsSynced cache.InformerSynced

	// workqueue is a rate limited work queue. This is used to queue work to be
	// processed instead of performing it as soon as a change happens. This
	// means we can ensure we only process a fixed amount of resources at a
	// time, and makes it easy to ensure we are never processing the same item
	// simultaneously in two different workers.
	workqueue workqueue.RateLimitingInterface
	// recorder is an event recorder for recording Event resources to the
	// Kubernetes API.
	recorder record.EventRecorder
}

const controllerAgentName = "goodjob-controller"

// NewGoodJobController returns a new sample controller
func NewGoodJobController(
	kubeclientset kubernetes.Interface,
	sampleclientset versioned.Interface,
	jobInformer batchinformers.JobInformer,
	goodInformer goodfinformers.GoodJobInformer) *GoodJobController {

	// Create event broadcaster
	// Add sample-controller types to the default Kubernetes Scheme so Events can be
	// logged for sample-controller types.
	goodscheme.AddToScheme(scheme.Scheme)
	glog.V(4).Info("Creating event broadcaster")
	eventBroadcaster := record.NewBroadcaster()
	eventBroadcaster.StartLogging(glog.Infof)
	eventBroadcaster.StartRecordingToSink(&typedcorev1.EventSinkImpl{Interface: kubeclientset.CoreV1().Events("")})
	recorder := eventBroadcaster.NewRecorder(scheme.Scheme, corev1.EventSource{Component: controllerAgentName})

	controller := &GoodJobController{
		kubeclientset:  kubeclientset,
		goodclientset:  sampleclientset,
		jobsLister:     jobInformer.Lister(),
		jobsSynced:     jobInformer.Informer().HasSynced,
		goodJobsLister: goodInformer.Lister(),
		goodJobsSynced: goodInformer.Informer().HasSynced,
		workqueue:      workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "GoodJobs"),
		recorder:       recorder,
	}



	glog.Info("Setting up event handlers")
	// Set up an event handler for when Foo resources change
	goodInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: controller.enqueueFoo,
		UpdateFunc: func(old, new interface{}) {
			controller.enqueueFoo(new)
		},
	})
	// Set up an event handler for when Deployment resources change. This
	// handler will lookup the owner of the given Deployment, and if it is
	// owned by a Foo resource will enqueue that Foo resource for
	// processing. This way, we don't need to implement custom logic for
	// handling Deployment resources. More info on this pattern:
	// https://github.com/kubernetes/community/blob/8cafef897a22026d42f5e5bb3f104febe7e29830/contributors/devel/controllers.md
	jobInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: controller.handleObject,
		UpdateFunc: func(old, new interface{}) {
			newJob := new.(*batchv1.Job)
			oldJob := old.(*batchv1.Job)
			if newJob.ResourceVersion == oldJob.ResourceVersion {
				// Periodic resync will send update events for all known Deployments.
				// Two different versions of the same Deployment will always have different RVs.
				return
			}
			// 判断是不是由 GoodJob 控制的
			if info, err := SwitchLabelInfo(newJob); err == nil {
				switch {
				case newJob.Status.Succeeded == 1 && oldJob.Status.Succeeded == 0:
					controller.jobSuccess(newJob, info)
				}
			}

			controller.handleObject(new)
		},
		DeleteFunc: controller.handleObject,
	})

	return controller
}

func (c *GoodJobController) jobSuccess(job *batchv1.Job, info labelInfo)  {
	if info.shardIndex == 0 {
		pods, err := clientset.CoreV1().Pods(job.Namespace).List(metav1.ListOptions{
			LabelSelector: fmt.Sprintf("job-name=%s", job.Name),
		})
		if err != nil {
			glog.Errorf("job 执行成功获取相关的 pods 出错， error: %s", err.Error())
		}
		if len(pods.Items) == 0 {
			glog.Errorf("Job %s 相关 Pod 已经被删除", job.Name)
		}
		pod := pods.Items[0]
		// FIXME
		// 直接读取的话可能会消耗比较大的内存，建议所有 Shard 运行时不要打印太多的日志
		byteSlice, err := getJobLog(pod.Namespace, pod.Name, &v1.PodLogOptions{})
		if err != nil {
			glog.Errorf("job 执行成功获取 pod 的相关的日志出错， error: %s", err.Error())
		}
		logs := bytes.Split(byteSlice, []byte("\n"))

		// 找到用户设置的匹配参数
		shardMatchPattern := func(shardMatchPattern string) string {
			if shardMatchPattern == "" {
				shardMatchPattern = "GOOD_JOB_SHARDS[\\s\\S\\w\\W]+"
			}
			return shardMatchPattern
		}(info.shardMatchPattern)

		// 找到日志中的 Shard 信息
		shards := findShardInfo(logs, func(byteArray []byte) bool {
			match, err :=  regexp.Match(shardMatchPattern, byteArray)
			if err != nil {
				glog.Errorf("匹配 shard 数据出错; error: %s", err.Error())
				return false
			}
			return match
		})

		goodJob, err := clientset.GoodjobV1alpha1().GoodJobs(job.Namespace).Get(info.goodJob, metav1.GetOptions{})
		if err != nil {
			glog.Errorf("shards 日志获取完成，寻找对应的 GoodJob 失败, name: %s; error: %s", info.goodJob, err.Error())
		}

		goodJob = goodJob.DeepCopy()
		goodJob.Status.Status = typed.ShardSuccess
		goodJob.Status.Shards = shards

		_, err = clientset.GoodjobV1alpha1().GoodJobs(job.Namespace).Update(goodJob)
		if err != nil {
			glog.Errorf("更新 GoodJob 状态失败，error: %s", err.Error())
		}

		glog.Infof("shard 成功, shards is %v", shards)
	}
}

// Run will set up the event handlers for types we are interested in, as well
// as syncing informer caches and starting workers. It will block until stopCh
// is closed, at which point it will shutdown the workqueue and wait for
// workers to finish processing their current work items.
func (c *GoodJobController) Run(threadiness int, stopCh <-chan struct{}) error {
	defer runtime.HandleCrash()
	defer c.workqueue.ShutDown()

	// Start the informer factories to begin populating the informer caches
	glog.Info("Starting Foo controller")

	// Wait for the caches to be synced before starting workers
	glog.Info("Waiting for informer caches to sync")
	if ok := cache.WaitForCacheSync(stopCh, c.jobsSynced, c.goodJobsSynced); !ok {
		return fmt.Errorf("failed to wait for caches to sync")
	}

	glog.Info("Starting GoodJobController workers")
	// Launch two workers to process Foo resources
	for i := 0; i < threadiness; i++ {
		go wait.Until(c.runWorker, time.Second, stopCh)
	}

	glog.Info("Started GoodJobController workers")
	<-stopCh
	glog.Info("Shutting down GoodJobController workers")

	return nil
}

// processNextWorkItem will read a single work item off the workqueue and
// attempt to process it, by calling the syncHandler.
func (c *GoodJobController) processNextWorkItem() bool {
	obj, shutdown := c.workqueue.Get()

	if shutdown {
		return false
	}

	// We wrap this block in a func so we can defer c.workqueue.Done.
	err := func(obj interface{}) error {
		// We call Done here so the workqueue knows we have finished
		// processing this item. We also must remember to call Forget if we
		// do not want this work item being re-queued. For example, we do
		// not call Forget if a transient error occurs, instead the item is
		// put back on the workqueue and attempted again after a back-off
		// period.
		defer c.workqueue.Done(obj)
		var key string
		var ok bool
		// We expect strings to come off the workqueue. These are of the
		// form namespace/name. We do this as the delayed nature of the
		// workqueue means the items in the informer cache may actually be
		// more up to date that when the item was initially put onto the
		// workqueue.
		if key, ok = obj.(string); !ok {
			// As the item in the workqueue is actually invalid, we call
			// Forget here else we'd go into a loop of attempting to
			// process a work item that is invalid.
			c.workqueue.Forget(obj)
			runtime.HandleError(fmt.Errorf("expected string in workqueue but got %#v", obj))
			return nil
		}
		// Run the syncHandler, passing it the namespace/name string of the
		// Foo resource to be synced.
		if err := c.syncHandler(key); err != nil {
			return fmt.Errorf("error syncing '%s': %s", key, err.Error())
		}
		// Finally, if no error occurs we Forget this item so it does not
		// get queued again until another change happens.
		c.workqueue.Forget(obj)
		glog.Infof("Successfully synced '%s'", key)
		return nil
	}(obj)

	if err != nil {
		runtime.HandleError(err)
		return true
	}

	return true
}


// syncHandler compares the actual state with the desired, and attempts to
// converge the two. It then updates the Status block of the Foo resource
// with the current status of the resource.
func (c *GoodJobController) syncHandler(key string) error {
	// Convert the namespace/name string into a distinct namespace and name
	namespace, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		runtime.HandleError(fmt.Errorf("invalid resource key: %s", key))
		return nil
	}

	// Get the Foo resource with this namespace/name
	goodJob, err := c.goodJobsLister.GoodJobs(namespace).Get(name)
	if err != nil {
		// The Foo resource may no longer exist, in which case we stop
		// processing.
		if errors.IsNotFound(err) {
			runtime.HandleError(fmt.Errorf("goodJob '%s' in work queue no longer exists", key))
			return nil
		}

		return err
	}

	// Get the deployment with the name specified in Foo.spec
	deployment, err := c.jobsLister.Jobs(goodJob.Namespace).Get(deploymentName)
	// If the resource doesn't exist, we'll create it
	if errors.IsNotFound(err) {
		deployment, err = c.kubeclientset.BatchV1().Jobs(goodJob.Namespace).Create(newDeployment(goodJob))
	}

	// If an error occurs during Get/Create, we'll requeue the item so we can
	// attempt processing again later. This could have been caused by a
	// temporary network failure, or any other transient reason.
	if err != nil {
		return err
	}

	// If the Deployment is not controlled by this Foo resource, we should log
	// a warning to the event recorder and ret
	if !metav1.IsControlledBy(deployment, goodJob) {
		msg := fmt.Sprintf(MessageResourceExists, deployment.Name)
		c.recorder.Event(goodJob, corev1.EventTypeWarning, ErrResourceExists, msg)
		return fmt.Errorf(msg)
	}

	// If this number of the replicas on the Foo resource is specified, and the
	// number does not equal the current desired replicas on the Deployment, we
	// should update the Deployment resource.
	if goodJob.Spec.Replicas != nil && *goodJob.Spec.Replicas != *deployment.Spec.Replicas {
		glog.V(4).Infof("Foo %s replicas: %d, deployment replicas: %d", name, *goodJob.Spec.Replicas, *deployment.Spec.Replicas)
		deployment, err = c.kubeclientset.AppsV1().Deployments(goodJob.Namespace).Update(newDeployment(goodJob))
	}

	// If an error occurs during Update, we'll requeue the item so we can
	// attempt processing again later. THis could have been caused by a
	// temporary network failure, or any other transient reason.
	if err != nil {
		return err
	}

	// Finally, we update the status block of the Foo resource to reflect the
	// current state of the world
	err = c.updateFooStatus(goodJob, deployment)
	if err != nil {
		return err
	}

	c.recorder.Event(goodJob, corev1.EventTypeNormal, SuccessSynced, MessageResourceSynced)
	return nil
}

// runWorker is a long-running function that will continually call the
// processNextWorkItem function in order to read and process a message on the
// workqueue.
func (c *GoodJobController) runWorker() {
	for c.processNextWorkItem() {
	}
}


// enqueueFoo takes a Foo resource and converts it into a namespace/name
// string which is then put onto the work queue. This method should *not* be
// passed resources of any type other than Foo.
func (c *GoodJobController) enqueueFoo(obj interface{}) {
	var key string
	var err error
	if key, err = cache.MetaNamespaceKeyFunc(obj); err != nil {
		runtime.HandleError(err)
		return
	}
	c.workqueue.AddRateLimited(key)
}


// handleObject will take any resource implementing metav1.Object and attempt
// to find the Foo resource that 'owns' it. It does this by looking at the
// objects metadata.ownerReferences field for an appropriate OwnerReference.
// It then enqueues that Foo resource to be processed. If the object does not
// have an appropriate OwnerReference, it will simply be skipped.
func (c *GoodJobController) handleObject(obj interface{}) {
	var object metav1.Object
	var ok bool
	if object, ok = obj.(metav1.Object); !ok {
		tombstone, ok := obj.(cache.DeletedFinalStateUnknown)
		if !ok {
			runtime.HandleError(fmt.Errorf("error decoding object, invalid type"))
			return
		}
		object, ok = tombstone.Obj.(metav1.Object)
		if !ok {
			runtime.HandleError(fmt.Errorf("error decoding object tombstone, invalid type"))
			return
		}
		glog.V(4).Infof("Recovered deleted object '%s' from tombstone", object.GetName())
	}
	glog.V(4).Infof("Processing object: %s", object.GetName())
	if ownerRef := metav1.GetControllerOf(object); ownerRef != nil {
		// If this object is not owned by a Foo, we should not do anything more
		// with it.
		if ownerRef.Kind != "Foo" {
			return
		}

		foo, err := c.goodJobsLister.GoodJobs(object.GetNamespace()).Get(ownerRef.Name)
		if err != nil {
			glog.V(4).Infof("ignoring orphaned object '%s' of foo '%s'", object.GetSelfLink(), ownerRef.Name)
			return
		}

		c.enqueueFoo(foo)
		return
	}
}


func addGoodJob(obj interface{}) {
	info, _ := json.Marshal(obj)
	glog.V(1).Info("create ---> " + string(info))
}

func updateGoodJob(oldObj, newObj interface{}) {
	oldGoodjob := oldObj.(*v1alpha1.GoodJob)
	newGoodjob := newObj.(*v1alpha1.GoodJob)
	if newGoodjob.ResourceVersion == oldGoodjob.ResourceVersion {
		// Periodic resync will send update events for all known Deployments.
		// Two different versions of the same Deployment will always have different RVs.
		return
	}

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
			job := newJob(goodjob.Spec.Shard.Template, goodjob, "", 0)
			_, err := clientset.BatchV1().Jobs(newGoodjob.Namespace).Create(job)
			if err != nil {
				glog.Errorf("创建 Shard Job 失败, error: %s", err.Error())
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
			glog.Errorf("update GoodJob status error, error: %s", err.Error())
		}
	case typed.ShardSuccess:
		dealShardSuccess(newGoodjob)
	}
}


func dealShardSuccess(goodJob *v1alpha1.GoodJob) error {
	for i, shard := range goodJob.Status.Shards {
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

		job := newJob(podTemplate, goodjob, shard, i+1)

		// 在创建 job 前检查是否状态为可重新执行
		nowGoodJob, err := clientset.GoodjobV1alpha1().GoodJobs(goodjob.Namespace).Get(goodjob.Name, metav1.GetOptions{})
		if err != nil {
			glog.Errorf("获取 GoodJob 失败， error: %s", err.Error())
		}
		if _, ok := nowGoodJob.Status.ShardStatuses[shard]; ok {
			continue
		}

		_, err = clientset.BatchV1().Jobs(goodJob.Namespace).Create(job)
		if err != nil {
			glog.Errorf("创建 Job 失败, error: %s", err.Error())
		}
		goodjob.Status.ShardStatuses[shard] = "create"
		clientset.GoodjobV1alpha1().GoodJobs(goodjob.Namespace).Update(goodjob)
	}
	return nil
}

func newJob(podTemplate v1.PodTemplateSpec, goodJob *v1alpha1.GoodJob, shard string, shardIndex int) *batchv1.Job {
	labels := map[string]string{
		GoodJobNameLabel:  goodJob.Name,
		PipelineNameLabel: goodJob.Status.Pipeline,
		ShardIndexLabel:   strconv.Itoa(shardIndex),
	}
	annotations := map[string]string{
		ShardLabel: shard,
	}
	if shardIndex == 0 {
		annotations[ShardMatchPatternLabel] = goodJob.Spec.Shard.MatchPattern
	}
	return &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:        fmt.Sprintf("good-goodJob-%s-%s", goodJob.Name, fmt.Sprintf("shard-%d", shardIndex)),
			Namespace:   goodJob.Namespace,
			Labels:      labels,
			Annotations: annotations,
			OwnerReferences:[]metav1.OwnerReference{
				*metav1.NewControllerRef(goodJob, v1alpha1.SchemeGroupVersion.WithKind(goodjob.GoodJob)),
			},
		},
		Spec: batchv1.JobSpec{
			Template: podTemplate,
		},
	}
}