package openstack

import (
	openstackv1alpha1 "github.com/vexxhost/atmosphere-operator/apis/openstack/v1alpha1"
	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/utils/pointer"
)

var pathPrefix networkingv1.PathType = networkingv1.PathTypePrefix

func GenerateIngress(ingress *networkingv1.Ingress, config *openstackv1alpha1.IngressConfig, port int32) {
	ingress.Labels = config.Labels
	ingress.Annotations = config.Annotations

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

	ingress.Spec.TLS = []networkingv1.IngressTLS{
		{
			SecretName: ingress.GetName() + "-certs",
			Hosts:      []string{config.Host},
		},
	}
}
