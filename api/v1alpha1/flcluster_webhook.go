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

package v1alpha1

import (
	"errors"
	"k8s.io/apimachinery/pkg/runtime"
	"regexp"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// log is for logging in this package.
var flclusterlog = logf.Log.WithName("flcluster-resource")

func (f *FLCluster) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(f).
		Complete()
}

//+kubebuilder:webhook:path=/mutate-kfl-aut-tech-v1alpha1-flcluster,mutating=true,failurePolicy=fail,sideEffects=None,groups=kfl.aut.tech,resources=flclusters,verbs=create;update,versions=v1alpha1,name=mflcluster.kb.io,admissionReviewVersions=v1

var _ webhook.Defaulter = &FLCluster{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (f *FLCluster) Default() {
	flclusterlog.Info("default", "name", f.Name)
	if f.Spec.EdgeServer != nil {
		v := "True"
		f.Spec.EdgeBased = &v
	} else if f.Spec.EdgeBased == nil {
		v := "False"
		f.Spec.EdgeBased = &v
	}
	if f.Spec.Splitting == nil {
		v := "none_splitting"
		f.Spec.Splitting = &v
	} else if *f.Spec.Splitting != "none_splitting" {
		v := "True"
		f.Spec.Offload = &v
	} else if f.Spec.Offload == nil {
		v := "False"
		f.Spec.Offload = &v
	}
	if f.Spec.Aggegation == nil {
		v := "fed_avg"
		f.Spec.Aggegation = &v
	}
	if f.Spec.Clustering == nil {
		v := "none_clustering"
		f.Spec.Clustering = &v
	}
	if f.Spec.ModelName == nil {
		v := "vgg"
		f.Spec.ModelName = &v
	}
	if f.Spec.Dataset == nil {
		v := "cifar10"
		f.Spec.Dataset = &v
	}
	if f.Spec.Index == nil {
		v := "0"
		f.Spec.Index = &v
	}
}

//+kubebuilder:webhook:path=/validate-kfl-aut-tech-v1alpha1-flcluster,mutating=false,failurePolicy=fail,sideEffects=None,groups=kfl.aut.tech,resources=flclusters,verbs=create;update,versions=v1alpha1,name=vflcluster.kb.io,admissionReviewVersions=v1

var _ webhook.Validator = &FLCluster{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (f *FLCluster) ValidateCreate() (admission.Warnings, error) {
	flclusterlog.Info("validate create", "name", f.Name)

	if &f.Spec.CentralServer == nil {
		return nil, errors.New("empty central server")
	}
	if &f.Spec.EdgeClient == nil || len(f.Spec.EdgeClient) == 0 {
		return nil, errors.New("empty edge client")
	}
	e1 := validateDevice(&f.Spec.CentralServer)
	if e1 != nil {
		return nil, errors.New("central server: " + e1.Error())
	}
	for i := 0; i < len(f.Spec.EdgeClient); i++ {
		e2 := validateDevice(&f.Spec.EdgeClient[i])
		if e2 != nil {
			return nil, errors.New("edge client: " + e2.Error())
		}
	}
	if f.Spec.EdgeServer != nil && len(f.Spec.EdgeServer) != 0 {
		for i := 0; i < len(f.Spec.EdgeServer); i++ {
			e3 := validateDevice(f.Spec.EdgeServer[i])
			if e3 != nil {
				return nil, errors.New("edge server: " + e3.Error())
			}
		}
	}
	pattern := `^(True|False)$`

	compile := regexp.MustCompile(pattern)
	if f.Spec.EdgeBased != nil {
		if !compile.MatchString(*f.Spec.EdgeBased) {
			return nil, errors.New("invalid new edgeBased, must be True or False")
		}
	}
	if f.Spec.Offload != nil {
		if !compile.MatchString(*f.Spec.Offload) {
			return nil, errors.New("invalid new offload, must be True or False")
		}
	}
	return nil, nil
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (f *FLCluster) ValidateUpdate(old runtime.Object) (admission.Warnings, error) {
	flclusterlog.Info("validate update", "name", f.Name)
	oldCluster := old.(*FLCluster)
	if &f.Spec.CentralServer == nil {
		return nil, errors.New("empty central server")
	}
	if &f.Spec.EdgeClient == nil {
		return nil, errors.New("empty edge client")
	}
	e1 := validateDevice(&f.Spec.CentralServer)
	if e1 != nil {
		return nil, errors.New("central server: " + e1.Error())
	}
	if !validateResourceUpdate(&f.Spec.CentralServer.Resources, &oldCluster.Spec.CentralServer.Resources) {
		return nil, errors.New("invalid new central server resource")
	}
	for i := 0; i < len(f.Spec.EdgeClient); i++ {
		e2 := validateDevice(&f.Spec.EdgeClient[i])
		if e2 != nil {
			return nil, errors.New("edge client: " + e2.Error())
		}
		if !validateResourceUpdate(&f.Spec.EdgeClient[i].Resources, &oldCluster.Spec.EdgeClient[i].Resources) {
			return nil, errors.New("invalid new edge client resource")
		}
	}
	if f.Spec.EdgeServer != nil && len(f.Spec.EdgeServer) != 0 {
		for i := 0; i < len(f.Spec.EdgeServer); i++ {
			e3 := validateDevice(f.Spec.EdgeServer[i])
			if e3 != nil {
				return nil, errors.New("edge server: " + e3.Error())
			}
			if !validateResourceUpdate(&f.Spec.EdgeServer[i].Resources, &oldCluster.Spec.EdgeServer[i].Resources) {
				return nil, errors.New("invalid new edge server resource")
			}
		}
	}
	pattern := `^(True|False)$`

	compile := regexp.MustCompile(pattern)
	if f.Spec.EdgeBased != nil {
		if !compile.MatchString(*f.Spec.EdgeBased) {
			return nil, errors.New("invalid edgeBased, must be True or False")
		}
	}
	if f.Spec.Offload != nil {
		if !compile.MatchString(*f.Spec.Offload) {
			return nil, errors.New("invalid offload, must be True or False")
		}
	}
	return nil, nil
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (f *FLCluster) ValidateDelete() (admission.Warnings, error) {
	flclusterlog.Info("validate delete", "name", f.Name)

	return nil, nil
}
func validateDevice(d *Device) error {

	if &d.Replica == nil || d.Replica != 1 {
		return errors.New("invalid replica")
	}

	if !validateResource(&d.Resources) {
		return errors.New("invalid resource")
	}
	return nil
}
func validateResource(r *Resources) bool {
	pattern := `^([+-]?[0-9.]+)([eEinumkKMGTP]*[-+]?[0-9]*)$`

	compile := regexp.MustCompile(pattern)

	if r.Cpu != "" && !compile.MatchString(r.Cpu) {
		return false
	}
	if r.Memory != "" && !compile.MatchString(r.Memory) {
		return false
	}
	//if !compile.MatchString(r.Storage) {
	//	return false
	//}
	return true
}

func validateResourceUpdate(r *Resources, o *Resources) bool {
	return true
}
