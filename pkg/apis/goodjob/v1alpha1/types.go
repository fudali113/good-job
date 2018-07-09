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
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// GoodJob is a specification for a GoodJob resource
type GoodJob struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   GoodJobSpec   `json:"spec,omitempty" protobuf:"bytes,,opt,name=spec"`
	Status GoodJobStatus `json:"status,omitempty" protobuf:"bytes,,opt,name=status"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// GoodJobList is a list of GoodJob resources
type GoodJobList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []GoodJob `json:"items"`
}

// FooSpec is the spec for a Foo resource
type GoodJobSpec struct {
	// 储存执行的程序
	Template v1.PodTemplateSpec `json:"template" protobuf:"bytes,6,opt,name=template"`
	// 分片的配置
	Shard GoodJobShard `json:"shard" protobuf:"bytes,,opt,name=shard"`
	// 指定并行度
	Parallel int `json:"parallel" protobuf:"bytes,,opt,name=parallel"`
}

// GoodJobShard 分片程序的配置
type GoodJobShard struct {
	// 分片的类型
	Type string 				`json:"type" protobuf:"bytes,,opt,name=type"`
	// 执行分片程序的配置
	Template v1.PodTemplateSpec `json:"template" protobuf:"bytes,6,opt,name=template"`
	// 匹配日志的正则表达式
	MatchPattern string 		`json:"matchPattern" protobuf:"bytes,6,opt,name=matchPattern"`
	// 手动设置分片
	Shards []string 			`json:"shards" protobuf:"bytes,,opt,name=shards"`
}

// FooStatus is the status for a Foo resource
type GoodJobStatus struct {
	Status        int               `json:"status" protobuf:"bytes,,opt,name=status"`
	Pipeline      string            `json:"pipeline" protobuf:"bytes,,opt,name=pipeline"`
	Shards        []string          `json:"shards" protobuf:"bytes,,opt,name=shards"`
	ShardStatuses map[string]string `json:"shardStatuses" protobuf:"bytes,,opt,name=shardStatuses"`
	Logs          []string          `json:"logs" protobuf:"bytes,,opt,name=logs"`
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

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// PipelineList is a list of Pipeline resources
type PipelineList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []Pipeline `json:"items"`
}

// PipelineSpec is the spec for a Pipeline resource
type PipelineSpec struct {
	Name string               `json:"name"`
	Jobs map[string][]GoodJob `json:"jobs"`
}

// PipelineStatus is the status for a Pipeline resource
type PipelineStatus struct {
	NowJob       string            `json:"nowJob"`
	AdditionInfo map[string]string `json:"addition_info"`
	Logs         []string          `json:"logs"`
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

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// CronTriggerList is a list of CornTrigger resources
type CronTriggerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []Pipeline `json:"items"`
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
	JobId        string
	RunTotal     int
	SuccessTotal int
}
