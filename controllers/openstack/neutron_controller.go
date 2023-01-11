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

package openstack

import (
	"context"

	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/chartutil"
	"helm.sh/helm/v3/pkg/release"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/go-logr/logr"
	"github.com/operator-framework/helm-operator-plugins/pkg/hook"
	"github.com/operator-framework/helm-operator-plugins/pkg/reconciler"
	"github.com/operator-framework/helm-operator-plugins/pkg/values"
	openstackv1alpha1 "github.com/vexxhost/atmosphere-operator/apis/openstack/v1alpha1"
	"github.com/vexxhost/atmosphere-operator/pkg/endpoints"
	"github.com/vexxhost/atmosphere-operator/pkg/images"
)

// NeutronReconciler reconciles a Neutron object
type NeutronReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// TODO(mnaser): Tone down these RBAC rules
//+kubebuilder:rbac:groups=apps;batch;core;extensions,resources=deployments;daemonsets;endpoints;replicasets;services;statefulsets,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=apps;batch;core;extensions,resources=cronjobs;jobs,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=apps;batch;core;extensions,resources=pods;services;endpoints;persistentvolumeclaims;events;configmaps;secrets;serviceaccounts;namespaces,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=networking.k8s.io,resources=ingresses;networkpolicies,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=policy,resources=poddisruptionbudgets,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=rbac.authorization.k8s.io,resources=roles;rolebindings,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=storage.k8s.io,resources=storageclasses,verbs=get;list;watch
//+kubebuilder:rbac:groups=pxc.percona.com,resources=perconaxtradbclusters,verbs=get;list;watch
//+kubebuilder:rbac:groups=openstack.atmosphere.vexxhost.com,resources=neutrons,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=openstack.atmosphere.vexxhost.com,resources=neutrons/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=openstack.atmosphere.vexxhost.com,resources=neutrons/finalizers,verbs=update
//+kubebuilder:rbac:groups=infra.atmosphere.vexxhost.com,resources=rabbitmqclusters,verbs=get;list;watch;create;update;patch

// SetupWithManager sets up the controller with the Manager.
func (r *NeutronReconciler) SetupWithManager(mgr ctrl.Manager) error {
	chart, err := loader.Load("helm-charts/neutron.tgz")
	if err != nil {
		return err
	}

	translator := values.TranslatorFunc(func(ctx context.Context, u *unstructured.Unstructured) (chartutil.Values, error) {
		neutron := &openstackv1alpha1.Neutron{}
		if err = runtime.DefaultUnstructuredConverter.FromUnstructured(u.Object, neutron); err != nil {
			return nil, err
		}

		tags, err := images.GetImageTagsForOpenstackHelmChart(chart, neutron.Spec.ImageRepository)
		if err != nil {
			return nil, err
		}

		keystoneRef := neutron.Spec.KeystoneRef.WithNamespace(neutron.Namespace)
		novaRef := neutron.Spec.NovaRef.WithNamespace(neutron.Namespace)
		octaviaRef := neutron.Spec.OctaviaRef.WithNamespace(neutron.Namespace)
		designateRef := neutron.Spec.DesignateRef.WithNamespace(neutron.Namespace)
		ironicRef := neutron.Spec.IronicRef.WithNamespace(neutron.Namespace)

		endpointConfig, err := endpoints.NewConfig(
			endpoints.WithNamespace(neutron.Namespace),
			endpoints.WithKeystoneRef(ctx, r.Client, &keystoneRef),
			endpoints.WithNovaRef(ctx, r.Client, &novaRef),
			endpoints.WithOctaviaRef(ctx, r.Client, &octaviaRef),
			endpoints.WithDesignateRef(ctx, r.Client, &designateRef),
			endpoints.WithIronicRef(ctx, r.Client, &ironicRef),
			endpoints.WithNeutron(ctx, r.Client, neutron),
		)
		if err != nil {
			return nil, err
		}

		endpoints, err := endpoints.ForChart(chart, endpointConfig)
		if err != nil {
			return nil, err
		}

		values := map[string]interface{}{
			"images": map[string]interface{}{
				"tags": tags,
			},
			"endpoints": endpoints,
			"conf": map[string]interface{}{
				"paste": map[string]interface{}{
					"composite:neutronapi_v2_0": map[string]interface{}{
						"keystone": "cors http_proxy_to_wsgi request_id catch_errors authtoken keystonecontext extensions neutronapiapp_v2_0",
					},
				},
				"neutron": map[string]interface{}{
					"DEFAULT": map[string]interface{}{
						"api_workers":             8,
						"rpc_workers":             8,
						"dhcp_agents_per_network": 3,
						"log_config_append":       nil,
						"service_plugins":         "qos,router,segments,trunk,vpnaas",
						"external_dns_driver":     "designate",
					},
					"cors": map[string]interface{}{
						"allowed_origin": "*",
					},
					"nova": map[string]interface{}{
						"live_migration_events": true,
					},
					"oslo_messaging_notifications": map[string]interface{}{
						"driver": "noop",
					},
					"service_providers": map[string]interface{}{
						"service_provider": "VPN:strongswan:neutron_vpnaas.services.vpn.service_drivers.ipsec.IPsecVPNDriver:default",
					},
				},
				"dhcp_agent": map[string]interface{}{
					"DEFAULT": map[string]interface{}{
						"dnsmasq_dns_servers":      endpointConfig.CoreDNSClusterIP,
						"enable_isolated_metadata": true,
					},
				},
				"l3_agent": map[string]interface{}{
					"AGENT": map[string]interface{}{
						"extensions": "vpnaas",
					},
					"vpnagent": map[string]interface{}{
						"vpn_device_driver": "neutron_vpnaas.services.vpn.device_drivers.strongswan_ipsec.StrongSwanDriver",
					},
				},
				"metadata_agent": map[string]interface{}{
					"DEFAULT": map[string]interface{}{
						"nova_metadata_port":           float32(8775),
						"metadata_proxy_shared_secret": endpointConfig.NovaMetadataSecret,
					},
				},
				"plugins": map[string]interface{}{
					"ml2_conf": map[string]interface{}{
						"ml2": map[string]interface{}{
							"extension_drivers": "dns_domain_ports,port_security,qos",
							"type_drivers":      "flat,gre,vlan,vxlan",
						},
						"ml2_type_gre": map[string]interface{}{
							"tunnel_id_ranges": "1:1000",
						},
						"ml2_type_vlan": map[string]interface{}{
							"network_vlan_ranges": "external:1:4094",
						},
					},
				},
			},
			"pod": map[string]interface{}{
				"replicas": map[string]interface{}{
					"api": neutron.Spec.Replicas,
				},
			},
			"manifests": map[string]interface{}{
				"ingress_server":         false,
				"service_ingress_server": false,
			},
		}

		overrides, err := neutron.Spec.Overrides.GetAsMap()
		if err != nil {
			return nil, err
		}

		return chartutil.CoalesceTables(values, overrides), nil
	})

	postHook := hook.PostHookFunc(func(u *unstructured.Unstructured, release release.Release, _ logr.Logger) error {
		neutron := &openstackv1alpha1.Neutron{}
		if err = runtime.DefaultUnstructuredConverter.FromUnstructured(u.Object, neutron); err != nil {
			return err
		}

		ingress := &networkingv1.Ingress{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "neutron-server",
				Namespace: u.GetNamespace(),
			},
		}
		_, err := ctrl.CreateOrUpdate(context.Background(), r.Client, ingress, func() error {
			GenerateIngress(ingress, &neutron.Spec.Ingress, endpoints.GetPortFromChart(chart, "network", "api"))
			return ctrl.SetControllerReference(neutron, ingress, r.Scheme)
		})

		return err
	})

	reconciler, err := reconciler.New(
		reconciler.WithChart(*chart),
		reconciler.WithClient(r.Client),
		reconciler.WithPostHook(postHook),
		reconciler.WithValueTranslator(translator),
		reconciler.SkipPrimaryGVKSchemeRegistration(true),
		reconciler.WithGroupVersionKind(schema.GroupVersionKind{
			Group:   openstackv1alpha1.GroupVersion.Group,
			Version: openstackv1alpha1.GroupVersion.Version,
			Kind:    "Neutron",
		}),
	)
	if err != nil {
		return err
	}

	if err := reconciler.SetupWithManager(mgr); err != nil {
		return err
	}

	return nil
}
