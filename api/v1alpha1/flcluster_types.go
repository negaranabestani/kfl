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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// FLClusterSpec defines the desired state of FLCluster
type FLClusterSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	Server     Device `json:"server"`
	EdgeServer Device `json:"edge-server"`
	EdgeClient Device `json:"edge-client"`
	Dataset    string `json:"dataset"`
	ModelName  string `json:"model-name"`
}


type Device struct {
	Replica int    `json:"replica"`
	Memory  string `json:"memory"`
	Cpu     string `json:"cpu"`
}

type LocalTrainingData struct {
	EdgeClient        Device `json:"edge-client"`
	EdgeServer        Device `json:"edge-server"`
	LocalTrainingTime int64  `json:"local-training-time"`
	LocalRounds       int    `json:"local-rounds"`
}

// FLClusterStatus defines the observed state of FLCluster
type FLClusterStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	GlobalTrainingTime int64               `json:"global-training-time"`
	GlobalRounds       int                 `json:"global-rounds"`
	GlobalAccuracy     float64             `json:"global-accuracy"`
	LocalTrainings     []LocalTrainingData `json:"local-trainings"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// FLCluster is the Schema for the flclusters API
type FLCluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   FLClusterSpec   `json:"spec,omitempty"`
	Status FLClusterStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// FLClusterList contains a list of FLCluster
type FLClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []FLCluster `json:"items"`
}

func init() {
	SchemeBuilder.Register(&FLCluster{}, &FLClusterList{})
}
