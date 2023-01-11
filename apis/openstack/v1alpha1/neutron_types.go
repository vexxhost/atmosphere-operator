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

// NeutronSpec defines the desired state of Neutron
type NeutronSpec struct {
	// +kubebuilder:default=1
	Replicas          int32          `json:"replicas"`
	RegionName        string         `json:"regionName"`
	Ingress           IngressConfig  `json:"ingress"`
	KeystoneRef       NamespacedName `json:"keystoneRef"`
	NovaRef           NamespacedName `json:"novaRef"`
	OctaviaRef        NamespacedName `json:"octaviaRef"`
	DesignateRef      NamespacedName `json:"designateRef"`
	IronicRef         NamespacedName `json:"ironicRef"`
	CoreDNSRef        NamespacedName `json:"coreDnsRef"`
	SecretsRef        NamespacedName `json:"secretsRef"`
	RabbitmqReference NamespacedName `json:"rabbitmqRef"`
	DatabaseReference NamespacedName `json:"databaseRef"`
	ImageRepository   string         `json:"imageRepository,omitempty"`
	Overrides         HelmOverrides  `json:"overrides,omitempty"`
}

// NeutronStatus defines the observed state of Neutron
type NeutronStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Neutron is the Schema for the neutrons API
type Neutron struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   NeutronSpec   `json:"spec,omitempty"`
	Status NeutronStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// NeutronList contains a list of Neutron
type NeutronList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Neutron `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Neutron{}, &NeutronList{})
}
