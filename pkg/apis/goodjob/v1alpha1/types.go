/*
Copyright 2017 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Job is a specification for a Job resource
type Job struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   JobSpec   `json:"spec"`
	Status JobStatus `json:"status"`
}

// FooSpec is the spec for a Foo resource
type JobSpec struct {
	// Job 的名字
	Name string `json:"name"`
	// 储存分片后的数据
	Shards []string `json:"shards"`
	// 储存执行的程序
	Exec ExecConfig `json:"exec"`
	// 分片的配置
	Shard ShardConfig `json:"shard"`
	// 指定并行度
	Parallel int `json:"parallel"`
	// 该 job 所在的 pipeline
	Pipeline string `json:"pipeline"`
}

// ExecConfig 执行 job 程序的配置
type ExecConfig struct {
	// 镜像
	Image string `json:"image"`
	// 启动命令
	Cmd []string `json:"cmd"`
	// 启动参数
	Args []string `json:"args"`
	// 环境变量
	Env []string `json:"env"`
}

// ShardConfig 分片程序的配置
type ShardConfig struct {
	// 分片的类型
	Type string `json:"type"`
	// 执行分片程序的配置
	Exec ExecConfig `json:"exec"`
	// 手动设置分片
	Shards []string `json:"shards"`
}

// FooStatus is the status for a Foo resource
type JobStatus struct {
	AdditionInfo map[string]string `json:"addition_info"`
	Pipeline     string            `json:"pipeline"`
	Logs         []string          `json:"logs"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// JobList is a list of Job resources
type JobList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []Job `json:"items"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Pipeline is a specification for a Pipeline resource
type Pipeline struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PipelineSpec   `json:"spec"`
	Status PipelineStatus `json:"status"`
}

// PipelineSpec is the spec for a Pipeline resource
type PipelineSpec struct {
	Name string           `json:"name"`
	Jobs map[string][]Job `json:"jobs"`
}

// PipelineStatus is the status for a Pipeline resource
type PipelineStatus struct {
	NowJob       string            `json:"nowJob"`
	AdditionInfo map[string]string `json:"addition_info"`
	Logs         []string          `json:"logs"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// PipelineList is a list of Pipeline resources
type PipelineList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []Pipeline `json:"items"`
}


// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// CronTrigger is a specification for a CornTrigger resource
type CronTrigger struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CronTriggerSpec   `json:"spec"`
	Status CronTriggerStatus `json:"status"`
}

type CronTriggerSpec struct {
	// 对应 k8s CornJob 的 Scheduler 字段
	Scheduler string `json:"scheduler"`
	// 出发资源的 Type ，可能为 pipelines 或者 jobs
	Type string `json:"type"`
	// 出发任务的 id
	Id string `json:"id"`
}

type CronTriggerStatus struct {
	JobId string
	RunTotal int
	SuccessTotal int
}


// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// CronTriggerList is a list of CornTrigger resources
type CronTriggerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []Pipeline `json:"items"`
}

