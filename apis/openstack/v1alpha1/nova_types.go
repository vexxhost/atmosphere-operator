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

// NovaSpec defines the desired state of Nova
type NovaSpec struct {
	// +kubebuilder:default=1
	Replicas          int32          `json:"replicas"`
	RegionName        string         `json:"regionName"`
	Ingress           IngressConfig  `json:"ingress"`
	VncIngress        IngressConfig  `json:"vncIngress"`
	KeystoneRef       NamespacedName `json:"keystoneRef"`
	PlacementRef      NamespacedName `json:"placementRef"`
	GlanceRef         NamespacedName `json:"glanceRef"`
	NeutronRef        NamespacedName `json:"neutronRef"`
	IronicRef         NamespacedName `json:"ironicRef"`
	SecretsRef        NamespacedName `json:"secretsRef"`
	RabbitmqReference NamespacedName `json:"rabbitmqRef"`
	DatabaseReference NamespacedName `json:"databaseRef"`
	ImageRepository   string         `json:"imageRepository,omitempty"`
	Overrides         HelmOverrides  `json:"overrides,omitempty"`
}

// NovaStatus defines the observed state of Nova
type NovaStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Nova is the Schema for the nova API
type Nova struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   NovaSpec   `json:"spec,omitempty"`
	Status NovaStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// NovaList contains a list of Nova
type NovaList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Nova `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Nova{}, &NovaList{})
}
