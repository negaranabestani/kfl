//go:build !ignore_autogenerated
// +build !ignore_autogenerated

/*
Copyright 2023.

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

// Code generated by controller-gen. DO NOT EDIT.

package v1alpha1

import (
	"k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Device) DeepCopyInto(out *Device) {
	*out = *in
	out.Resources = in.Resources
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Device.
func (in *Device) DeepCopy() *Device {
	if in == nil {
		return nil
	}
	out := new(Device)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FLCluster) DeepCopyInto(out *FLCluster) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.Spec = in.Spec
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FLCluster.
func (in *FLCluster) DeepCopy() *FLCluster {
	if in == nil {
		return nil
	}
	out := new(FLCluster)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *FLCluster) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FLClusterList) DeepCopyInto(out *FLClusterList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]FLCluster, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FLClusterList.
func (in *FLClusterList) DeepCopy() *FLClusterList {
	if in == nil {
		return nil
	}
	out := new(FLClusterList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *FLClusterList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FLClusterSpec) DeepCopyInto(out *FLClusterSpec) {
	*out = *in
	out.CentralServer = in.CentralServer
	out.EdgeServer = in.EdgeServer
	out.EdgeClient = in.EdgeClient
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FLClusterSpec.
func (in *FLClusterSpec) DeepCopy() *FLClusterSpec {
	if in == nil {
		return nil
	}
	out := new(FLClusterSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FLClusterStatus) DeepCopyInto(out *FLClusterStatus) {
	*out = *in
	if in.LocalTrainings != nil {
		in, out := &in.LocalTrainings, &out.LocalTrainings
		*out = make([]LocalTrainingData, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FLClusterStatus.
func (in *FLClusterStatus) DeepCopy() *FLClusterStatus {
	if in == nil {
		return nil
	}
	out := new(FLClusterStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *LocalTrainingData) DeepCopyInto(out *LocalTrainingData) {
	*out = *in
	out.EdgeClient = in.EdgeClient
	out.EdgeServer = in.EdgeServer
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LocalTrainingData.
func (in *LocalTrainingData) DeepCopy() *LocalTrainingData {
	if in == nil {
		return nil
	}
	out := new(LocalTrainingData)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Resources) DeepCopyInto(out *Resources) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Resources.
func (in *Resources) DeepCopy() *Resources {
	if in == nil {
		return nil
	}
	out := new(Resources)
	in.DeepCopyInto(out)
	return out
}
