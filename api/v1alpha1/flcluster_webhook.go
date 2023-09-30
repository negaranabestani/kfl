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
		f.Spec.EdgeBased = "True"
	} else if &f.Spec.EdgeBased == nil {
		f.Spec.EdgeBased = "False"
	}
	if f.Spec.Splitting != "none_splitting" {
		f.Spec.Offload = "True"
	} else if &f.Spec.Splitting == nil {
		f.Spec.Splitting = "none_splitting"
	} else if &f.Spec.Offload == nil {
		f.Spec.Offload = "False"
	}
	if &f.Spec.Aggegation == nil {
		f.Spec.Aggegation = "fed_avg"
	}
	if &f.Spec.Clustering == nil {
		f.Spec.Clustering = "none_clustering"
	}
	if &f.Spec.ModelName == nil {
		f.Spec.ModelName = "vgg"
	}
	if &f.Spec.Dataset == nil {
		f.Spec.Dataset = "cifar10"
	}
	if &f.Spec.Index == nil {
		f.Spec.Index = "0"
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
	if &f.Spec.EdgeClient == nil {
		return nil, errors.New("empty edge client")
	}
	e1 := validateDevice(&f.Spec.CentralServer)
	if e1 != nil {
		return nil, errors.New("central server: " + e1.Error())
	}
	e2 := validateDevice(&f.Spec.EdgeClient)
	if e2 != nil {
		return nil, errors.New("edge client: " + e2.Error())
	}
	if f.Spec.EdgeServer != nil {
		e3 := validateDevice(f.Spec.EdgeServer)
		if e3 != nil {
			return nil, errors.New("edge server: " + e3.Error())
		}
	}
	pattern := `^(True|False)$`

	compile := regexp.MustCompile(pattern)
	if &f.Spec.EdgeBased != nil {
		if !compile.MatchString(f.Spec.EdgeBased) {
			return nil, errors.New("invalid new edgeBased, must be True or False")
		}
	}
	if &f.Spec.Offload != nil {
		if !compile.MatchString(f.Spec.Offload) {
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
	e2 := validateDevice(&f.Spec.EdgeClient)
	if e2 != nil {
		return nil, errors.New("edge client: " + e2.Error())
	}
	if !validateResourceUpdate(&f.Spec.EdgeClient.Resources, &oldCluster.Spec.EdgeClient.Resources) {
		return nil, errors.New("invalid new edge client resource")
	}
	if f.Spec.EdgeServer != nil {
		e3 := validateDevice(f.Spec.EdgeServer)
		if e3 != nil {
			return nil, errors.New("edge server: " + e3.Error())
		}
		if !validateResourceUpdate(&f.Spec.EdgeServer.Resources, &oldCluster.Spec.EdgeServer.Resources) {
			return nil, errors.New("invalid new edge server resource")
		}
	}
	pattern := `^(True|False)$`

	compile := regexp.MustCompile(pattern)
	if &f.Spec.EdgeBased != nil {
		if !compile.MatchString(f.Spec.EdgeBased) {
			return nil, errors.New("invalid new edgeBased, must be True or False")
		}
	}
	if &f.Spec.Offload != nil {
		if !compile.MatchString(f.Spec.Offload) {
			return nil, errors.New("invalid new offload, must be True or False")
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
