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
	openstackv1alpha "github.com/vexxhost/atmosphere-operator/apis/openstack/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// OpenvswitchSpec defines the desired state of Openvswitch
type OpenvswitchSpec struct {
	ImageRepository string                         `json:"imageRepository,omitempty"`
	Overrides       openstackv1alpha.HelmOverrides `json:"overrides,omitempty"`
}

// OpenvswitchStatus defines the observed state of Openvswitch
type OpenvswitchStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Openvswitch is the Schema for the openvswitches API
type Openvswitch struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   OpenvswitchSpec   `json:"spec,omitempty"`
	Status OpenvswitchStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// OpenvswitchList contains a list of Openvswitch
type OpenvswitchList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Openvswitch `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Openvswitch{}, &OpenvswitchList{})
}
