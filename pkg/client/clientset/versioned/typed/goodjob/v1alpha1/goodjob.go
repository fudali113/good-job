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

// GoodJobsGetter has a method to return a GoodJobInterface.
// A group's client should implement this interface.
type GoodJobsGetter interface {
	GoodJobs(namespace string) GoodJobInterface
}

// GoodJobInterface has methods to work with GoodJob resources.
type GoodJobInterface interface {
	Create(*v1alpha1.GoodJob) (*v1alpha1.GoodJob, error)
	Update(*v1alpha1.GoodJob) (*v1alpha1.GoodJob, error)
	UpdateStatus(*v1alpha1.GoodJob) (*v1alpha1.GoodJob, error)
	Delete(name string, options *v1.DeleteOptions) error
	DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error
	Get(name string, options v1.GetOptions) (*v1alpha1.GoodJob, error)
	List(opts v1.ListOptions) (*v1alpha1.GoodJobList, error)
	Watch(opts v1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.GoodJob, err error)
	GoodJobExpansion
}

// goodJobs implements GoodJobInterface
type goodJobs struct {
	client rest.Interface
	ns     string
}

// newGoodJobs returns a GoodJobs
func newGoodJobs(c *GoodjobV1alpha1Client, namespace string) *goodJobs {
	return &goodJobs{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the goodJob, and returns the corresponding goodJob object, and an error if there is any.
func (c *goodJobs) Get(name string, options v1.GetOptions) (result *v1alpha1.GoodJob, err error) {
	result = &v1alpha1.GoodJob{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("goodjobs").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of GoodJobs that match those selectors.
func (c *goodJobs) List(opts v1.ListOptions) (result *v1alpha1.GoodJobList, err error) {
	result = &v1alpha1.GoodJobList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("goodjobs").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested goodJobs.
func (c *goodJobs) Watch(opts v1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("goodjobs").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch()
}

// Create takes the representation of a goodJob and creates it.  Returns the server's representation of the goodJob, and an error, if there is any.
func (c *goodJobs) Create(goodJob *v1alpha1.GoodJob) (result *v1alpha1.GoodJob, err error) {
	result = &v1alpha1.GoodJob{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("goodjobs").
		Body(goodJob).
		Do().
		Into(result)
	return
}

// Update takes the representation of a goodJob and updates it. Returns the server's representation of the goodJob, and an error, if there is any.
func (c *goodJobs) Update(goodJob *v1alpha1.GoodJob) (result *v1alpha1.GoodJob, err error) {
	result = &v1alpha1.GoodJob{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("goodjobs").
		Name(goodJob.Name).
		Body(goodJob).
		Do().
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().

func (c *goodJobs) UpdateStatus(goodJob *v1alpha1.GoodJob) (result *v1alpha1.GoodJob, err error) {
	result = &v1alpha1.GoodJob{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("goodjobs").
		Name(goodJob.Name).
		SubResource("status").
		Body(goodJob).
		Do().
		Into(result)
	return
}

// Delete takes name of the goodJob and deletes it. Returns an error if one occurs.
func (c *goodJobs) Delete(name string, options *v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("goodjobs").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *goodJobs) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("goodjobs").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched goodJob.
func (c *goodJobs) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.GoodJob, err error) {
	result = &v1alpha1.GoodJob{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("goodjobs").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
