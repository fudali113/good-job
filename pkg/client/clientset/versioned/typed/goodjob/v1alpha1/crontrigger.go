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

// Code generated by client-gen. DO NOT EDIT.

package v1alpha1

import (
	v1alpha1 "github.com/fudali113/good-job/pkg/apis/goodjob/v1alpha1"
	scheme "github.com/fudali113/good-job/pkg/client/clientset/versioned/scheme"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// CronTriggersGetter has a method to return a CronTriggerInterface.
// A group's client should implement this interface.
type CronTriggersGetter interface {
	CronTriggers(namespace string) CronTriggerInterface
}

// CronTriggerInterface has methods to work with CronTrigger resources.
type CronTriggerInterface interface {
	Create(*v1alpha1.CronTrigger) (*v1alpha1.CronTrigger, error)
	Update(*v1alpha1.CronTrigger) (*v1alpha1.CronTrigger, error)
	UpdateStatus(*v1alpha1.CronTrigger) (*v1alpha1.CronTrigger, error)
	Delete(name string, options *v1.DeleteOptions) error
	DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error
	Get(name string, options v1.GetOptions) (*v1alpha1.CronTrigger, error)
	List(opts v1.ListOptions) (*v1alpha1.CronTriggerList, error)
	Watch(opts v1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.CronTrigger, err error)
	CronTriggerExpansion
}

// cronTriggers implements CronTriggerInterface
type cronTriggers struct {
	client rest.Interface
	ns     string
}

// newCronTriggers returns a CronTriggers
func newCronTriggers(c *GoodjobV1alpha1Client, namespace string) *cronTriggers {
	return &cronTriggers{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the cronTrigger, and returns the corresponding cronTrigger object, and an error if there is any.
func (c *cronTriggers) Get(name string, options v1.GetOptions) (result *v1alpha1.CronTrigger, err error) {
	result = &v1alpha1.CronTrigger{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("crontriggers").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of CronTriggers that match those selectors.
func (c *cronTriggers) List(opts v1.ListOptions) (result *v1alpha1.CronTriggerList, err error) {
	result = &v1alpha1.CronTriggerList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("crontriggers").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested cronTriggers.
func (c *cronTriggers) Watch(opts v1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("crontriggers").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch()
}

// Create takes the representation of a cronTrigger and creates it.  Returns the server's representation of the cronTrigger, and an error, if there is any.
func (c *cronTriggers) Create(cronTrigger *v1alpha1.CronTrigger) (result *v1alpha1.CronTrigger, err error) {
	result = &v1alpha1.CronTrigger{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("crontriggers").
		Body(cronTrigger).
		Do().
		Into(result)
	return
}

// Update takes the representation of a cronTrigger and updates it. Returns the server's representation of the cronTrigger, and an error, if there is any.
func (c *cronTriggers) Update(cronTrigger *v1alpha1.CronTrigger) (result *v1alpha1.CronTrigger, err error) {
	result = &v1alpha1.CronTrigger{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("crontriggers").
		Name(cronTrigger.Name).
		Body(cronTrigger).
		Do().
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().

func (c *cronTriggers) UpdateStatus(cronTrigger *v1alpha1.CronTrigger) (result *v1alpha1.CronTrigger, err error) {
	result = &v1alpha1.CronTrigger{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("crontriggers").
		Name(cronTrigger.Name).
		SubResource("status").
		Body(cronTrigger).
		Do().
		Into(result)
	return
}

// Delete takes name of the cronTrigger and deletes it. Returns an error if one occurs.
func (c *cronTriggers) Delete(name string, options *v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("crontriggers").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *cronTriggers) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("crontriggers").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched cronTrigger.
func (c *cronTriggers) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.CronTrigger, err error) {
	result = &v1alpha1.CronTrigger{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("crontriggers").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
