package v1alpha1

import (
	"encoding/json"

	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/types"
)

type NamespacedName struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace,omitempty"`
}

func (n NamespacedName) WithNamespace(namespace string) NamespacedName {
	if n.Namespace == "" {
		n.Namespace = namespace
	}

	return n
}

func (n NamespacedName) NativeNamespacedName() types.NamespacedName {
	return types.NamespacedName{
		Name:      n.Name,
		Namespace: n.Namespace,
	}
}

type IngressConfig struct {
	Host        string                  `json:"host"`
	ClassName   string                  `json:"className"`
	Labels      map[string]string       `json:"labels,omitempty"`
	Annotations map[string]string       `json:"annotations,omitempty"`
	TLS         networkingv1.IngressTLS `json:"tls,omitempty"`
}

// +kubebuilder:validation:Schemaless
// +kubebuilder:pruning:PreserveUnknownFields
// +kubebuilder:validation:Type=object
type HelmOverrides struct {
	json.RawMessage `json:"-"`
}

func (h *HelmOverrides) GetAsMap() (map[string]interface{}, error) {
	val := make(map[string]interface{})

	if h.RawMessage != nil {
		err := json.Unmarshal(h.RawMessage, &val)
		if err != nil {
			return nil, err
		}
	}

	return val, nil
}
