package v1alpha1

import (
	"encoding/json"

	runtime "k8s.io/apimachinery/pkg/runtime"
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
	Host        string            `json:"host"`
	ClassName   string            `json:"className"`
	Labels      map[string]string `json:"labels,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty"`
}

type HelmOverrides map[string]runtime.RawExtension

func (h *HelmOverrides) GetAsMap() (map[string]interface{}, error) {
	m := make(map[string]interface{})

	for k, v := range *h {
		val := make(map[string]interface{})
		err := json.Unmarshal(v.Raw, &val)
		if err != nil {
			return nil, err
		}

		m[k] = val
	}

	return m, nil
}
