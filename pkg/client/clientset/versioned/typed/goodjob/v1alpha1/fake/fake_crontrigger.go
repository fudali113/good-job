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

package fake

import (
	v1alpha1 "github.com/fudali113/good-job/pkg/apis/goodjob/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeCronTriggers implements CronTriggerInterface
type FakeCronTriggers struct {
	Fake *FakeGoodjobV1alpha1
	ns   string
}

var crontriggersResource = schema.GroupVersionResource{Group: "goodjob.k8s.io", Version: "v1alpha1", Resource: "crontriggers"}

var crontriggersKind = schema.GroupVersionKind{Group: "goodjob.k8s.io", Version: "v1alpha1", Kind: "CronTrigger"}

// Get takes name of the cronTrigger, and returns the corresponding cronTrigger object, and an error if there is any.
func (c *FakeCronTriggers) Get(name string, options v1.GetOptions) (result *v1alpha1.CronTrigger, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(crontriggersResource, c.ns, name), &v1alpha1.CronTrigger{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.CronTrigger), err
}

// List takes label and field selectors, and returns the list of CronTriggers that match those selectors.
func (c *FakeCronTriggers) List(opts v1.ListOptions) (result *v1alpha1.CronTriggerList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(crontriggersResource, crontriggersKind, c.ns, opts), &v1alpha1.CronTriggerList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.CronTriggerList{ListMeta: obj.(*v1alpha1.CronTriggerList).ListMeta}
	for _, item := range obj.(*v1alpha1.CronTriggerList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested cronTriggers.
func (c *FakeCronTriggers) Watch(opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(crontriggersResource, c.ns, opts))

}

// Create takes the representation of a cronTrigger and creates it.  Returns the server's representation of the cronTrigger, and an error, if there is any.
func (c *FakeCronTriggers) Create(cronTrigger *v1alpha1.CronTrigger) (result *v1alpha1.CronTrigger, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(crontriggersResource, c.ns, cronTrigger), &v1alpha1.CronTrigger{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.CronTrigger), err
}

// Update takes the representation of a cronTrigger and updates it. Returns the server's representation of the cronTrigger, and an error, if there is any.
func (c *FakeCronTriggers) Update(cronTrigger *v1alpha1.CronTrigger) (result *v1alpha1.CronTrigger, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(crontriggersResource, c.ns, cronTrigger), &v1alpha1.CronTrigger{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.CronTrigger), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeCronTriggers) UpdateStatus(cronTrigger *v1alpha1.CronTrigger) (*v1alpha1.CronTrigger, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(crontriggersResource, "status", c.ns, cronTrigger), &v1alpha1.CronTrigger{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.CronTrigger), err
}

// Delete takes name of the cronTrigger and deletes it. Returns an error if one occurs.
func (c *FakeCronTriggers) Delete(name string, options *v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(crontriggersResource, c.ns, name), &v1alpha1.CronTrigger{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeCronTriggers) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(crontriggersResource, c.ns, listOptions)

	_, err := c.Fake.Invokes(action, &v1alpha1.CronTriggerList{})
	return err
}

// Patch applies the patch and returns the patched cronTrigger.
func (c *FakeCronTriggers) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.CronTrigger, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(crontriggersResource, c.ns, name, data, subresources...), &v1alpha1.CronTrigger{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.CronTrigger), err
}
