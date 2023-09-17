package util

import (
	"github.com/negaranabestani/kfl/api/v1alpha1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

func Int32ptr(value int32) *int32 {
	return &value
}
func ResourceRequirements(resources v1alpha1.Resources) (*v1.ResourceRequirements, error) {
	var err error

	req := v1.ResourceRequirements{
		Requests: v1.ResourceList{},
		Limits:   v1.ResourceList{},
	}

	if resources.Cpu != "" {
		req.Requests[v1.ResourceCPU], err = resource.ParseQuantity(resources.Cpu)
		if err != nil {
			return nil, err
		}
		req.Limits[v1.ResourceCPU] = req.Requests[v1.ResourceCPU]
	}

	if resources.Memory != "" {
		req.Requests[v1.ResourceMemory], err = resource.ParseQuantity(resources.Memory)
		if err != nil {
			return nil, err
		}
		req.Limits[v1.ResourceMemory] = req.Requests[v1.ResourceMemory]
	}

	return &req, nil
}
