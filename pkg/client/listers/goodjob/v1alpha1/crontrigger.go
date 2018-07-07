/*
Copyright The Kubernetes Authors.

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

// Code generated by lister-gen. DO NOT EDIT.

package v1alpha1

import (
	v1alpha1 "github.com/fudali113/good-job/pkg/apis/goodjob/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// CronTriggerLister helps list CronTriggers.
type CronTriggerLister interface {
	// List lists all CronTriggers in the indexer.
	List(selector labels.Selector) (ret []*v1alpha1.CronTrigger, err error)
	// CronTriggers returns an object that can list and get CronTriggers.
	CronTriggers(namespace string) CronTriggerNamespaceLister
	CronTriggerListerExpansion
}

// cronTriggerLister implements the CronTriggerLister interface.
type cronTriggerLister struct {
	indexer cache.Indexer
}

// NewCronTriggerLister returns a new CronTriggerLister.
func NewCronTriggerLister(indexer cache.Indexer) CronTriggerLister {
	return &cronTriggerLister{indexer: indexer}
}

// List lists all CronTriggers in the indexer.
func (s *cronTriggerLister) List(selector labels.Selector) (ret []*v1alpha1.CronTrigger, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.CronTrigger))
	})
	return ret, err
}

// CronTriggers returns an object that can list and get CronTriggers.
func (s *cronTriggerLister) CronTriggers(namespace string) CronTriggerNamespaceLister {
	return cronTriggerNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// CronTriggerNamespaceLister helps list and get CronTriggers.
type CronTriggerNamespaceLister interface {
	// List lists all CronTriggers in the indexer for a given namespace.
	List(selector labels.Selector) (ret []*v1alpha1.CronTrigger, err error)
	// Get retrieves the CronTrigger from the indexer for a given namespace and name.
	Get(name string) (*v1alpha1.CronTrigger, error)
	CronTriggerNamespaceListerExpansion
}

// cronTriggerNamespaceLister implements the CronTriggerNamespaceLister
// interface.
type cronTriggerNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all CronTriggers in the indexer for a given namespace.
func (s cronTriggerNamespaceLister) List(selector labels.Selector) (ret []*v1alpha1.CronTrigger, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.CronTrigger))
	})
	return ret, err
}

// Get retrieves the CronTrigger from the indexer for a given namespace and name.
func (s cronTriggerNamespaceLister) Get(name string) (*v1alpha1.CronTrigger, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1alpha1.Resource("crontrigger"), name)
	}
	return obj.(*v1alpha1.CronTrigger), nil
}
