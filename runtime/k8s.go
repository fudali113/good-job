package runtime

import (
	"bytes"
	"encoding/json"
	"fmt"
	"text/template"

	"github.com/fudali113/good-job/typed"
	"gopkg.in/yaml.v2"
	"k8s.io/api/batch/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var jobTemplate = `
apiVersion: batch/v1
kind: Job
metadata:
 labels:
   pipeline: {{.pipeline}}
   job: {{.job}}
   type: {{.type}}
 name: {{.name}}
spec:
 template:
   metadata:
     labels:
       pipeline: {{.pipeline}}
       job: {{.job}}
	   type: {{.type}}
     name: {{.name}}
   spec:
     containers:
     - image: {{.config.Image}}
       command: {{toJson .config.Cmd}}
       args: {{toJson .config.Args}}
	   env: {{toJson .config.Env}}
`

type K8sRuntime struct {
	typed.K8sClient
	jobTemplate *template.Template
}

// CreateK8sRuntime 生成一个 k8s 的运行时
func CreateK8sRuntime(namespace string) (k8sRuntime *K8sRuntime, err error) {
	k8sClient := typed.K8sClient{
		Namespace: namespace,
	}
	config, err := clientcmd.BuildConfigFromFlags("", "")
	if err != nil {
		return
	}
	k8sClient.Client, err = kubernetes.NewForConfig(config)
	if err != nil {
		return
	}
	jobTemplate, err := createTemplate()
	k8sRuntime = &K8sRuntime{
		jobTemplate: jobTemplate,
		K8sClient:   k8sClient,
	}
	return
}

func createTemplate() (tpl *template.Template, err error) {
	tpl, err = template.New("jobTemplate").Funcs(template.FuncMap{
		"toJson": func(v interface{}) string {
			bytes, err := json.Marshal(v)
			if err != nil {
				panic(err)
			}
			return string(bytes)
		},
	}).Parse(jobTemplate)
	if err != nil {
		return
	}
	return
}

func (k K8sRuntime) CreateJob(metaData map[string]interface{},
	exec typed.ExecConfig) (k8sJob *v1.Job, err error) {
	buffer := bytes.NewBufferString("")
	metaData["config"] = exec
	k.jobTemplate.Execute(buffer, metaData)
	yaml.Unmarshal(buffer.Bytes(), k8sJob)
	return k.Client.BatchV1().Jobs(k.Namespace).Create(k8sJob)
}

func (k K8sRuntime) WatchJob(name string) (watch watch.Interface, err error) {
	return k.Client.BatchV1().Jobs(k.Namespace).Watch(metav1.ListOptions{FieldSelector: fmt.Sprintf("metadata.name=%s", name)})
}
