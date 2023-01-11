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

// LibvirtSpec defines the desired state of Libvirt
type LibvirtSpec struct {
	ImageRepository string                         `json:"imageRepository,omitempty"`
	Overrides       openstackv1alpha.HelmOverrides `json:"overrides,omitempty"`
}

// LibvirtStatus defines the observed state of Libvirt
type LibvirtStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Libvirt is the Schema for the libvirts API
type Libvirt struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   LibvirtSpec   `json:"spec,omitempty"`
	Status LibvirtStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// LibvirtList contains a list of Libvirt
type LibvirtList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Libvirt `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Libvirt{}, &LibvirtList{})
}
