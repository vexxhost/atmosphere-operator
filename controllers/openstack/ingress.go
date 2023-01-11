package openstack

import (
	"fmt"

	openstackv1alpha1 "github.com/vexxhost/atmosphere-operator/apis/openstack/v1alpha1"
	"golang.org/x/exp/maps"
	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/utils/pointer"
)

var pathPrefix networkingv1.PathType = networkingv1.PathTypePrefix

func GenerateIngress(ingress *networkingv1.Ingress, config *openstackv1alpha1.IngressConfig, port int32) {
	ingress.Labels = config.Labels
	ingress.Annotations = map[string]string{
		"nginx.ingress.kubernetes.io/proxy-body-size":         "0",
		"nginx.ingress.kubernetes.io/proxy-request-buffering": "off",
	}

	maps.Copy(ingress.Annotations, config.Annotations)
	fmt.Println(config.Host)

	ingress.Spec.IngressClassName = pointer.String(config.ClassName)
	ingress.Spec.Rules = []networkingv1.IngressRule{
		{
			Host: config.Host,
			IngressRuleValue: networkingv1.IngressRuleValue{
				HTTP: &networkingv1.HTTPIngressRuleValue{
					Paths: []networkingv1.HTTPIngressPath{
						{
							Path:     "/",
							PathType: &pathPrefix,
							Backend: networkingv1.IngressBackend{
								Service: &networkingv1.IngressServiceBackend{
									Name: ingress.GetName(),
									Port: networkingv1.ServiceBackendPort{
										Number: port,
									},
								},
							},
						},
					},
				},
			},
		},
	}

	secretName := ingress.GetName() + "-certs"
	if config.TLS.SecretName != "" {
		secretName = config.TLS.SecretName
	}

	ingress.Spec.TLS = []networkingv1.IngressTLS{
		{
			SecretName: secretName,
			Hosts:      []string{config.Host},
		},
	}
}
