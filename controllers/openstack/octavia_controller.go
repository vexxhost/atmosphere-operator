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
	"strings"

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

// OctaviaReconciler reconciles a Octavia object
type OctaviaReconciler struct {
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
//+kubebuilder:rbac:groups=openstack.atmosphere.vexxhost.com,resources=octavias,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=openstack.atmosphere.vexxhost.com,resources=octavias/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=openstack.atmosphere.vexxhost.com,resources=octavias/finalizers,verbs=update
//+kubebuilder:rbac:groups=infra.atmosphere.vexxhost.com,resources=rabbitmqclusters,verbs=get;list;watch;create;update;patch

// SetupWithManager sets up the controller with the Manager.
func (r *OctaviaReconciler) SetupWithManager(mgr ctrl.Manager) error {
	chart, err := loader.Load("helm-charts/octavia.tgz")
	if err != nil {
		return err
	}

	translator := values.TranslatorFunc(func(ctx context.Context, u *unstructured.Unstructured) (chartutil.Values, error) {
		octavia := &openstackv1alpha1.Octavia{}
		if err = runtime.DefaultUnstructuredConverter.FromUnstructured(u.Object, octavia); err != nil {
			return nil, err
		}

		tags, err := images.GetImageTagsForOpenstackHelmChart(chart, octavia.Spec.ImageRepository)
		if err != nil {
			return nil, err
		}

		keystoneRef := octavia.Spec.KeystoneRef.WithNamespace(octavia.Namespace)
		neutronRef := octavia.Spec.NeutronRef.WithNamespace(octavia.Namespace)

		endpointConfig, err := endpoints.NewConfig(
			endpoints.WithNamespace(octavia.Namespace),
			endpoints.WithKeystoneRef(ctx, r.Client, &keystoneRef),
			endpoints.WithNeutronRef(ctx, r.Client, &neutronRef),
			endpoints.WithOctavia(ctx, r.Client, octavia),
		)
		if err != nil {
			return nil, err
		}

		endpoints, err := endpoints.ForChart(chart, endpointConfig)
		if err != nil {
			return nil, err
		}

		// 		volumes:
		// 			- name: octavia-server-ca
		// 				secret:
		// 					secretName: octavia-server-ca
		// 			- name: octavia-client-certs
		// 				secret:
		// 					secretName: octavia-client-certs

		volumes := []interface{}{
			map[string]interface{}{
				"name": "octavia-server-ca",
				"secret": map[string]interface{}{
					"secretName": octavia.Spec.AmphoraConfig.ServerCertificateAuthority.Name,
				},
			},
			map[string]interface{}{
				"name": "octavia-client-certs",
				"secret": map[string]interface{}{
					"secretName": octavia.Spec.AmphoraConfig.ClientCertificate.Name,
				},
			},
		}

		volumeMounts := []interface{}{
			map[string]interface{}{
				"name":      "octavia-server-ca",
				"mountPath": "/etc/octavia/certs/server",
			},
			map[string]interface{}{
				"name":      "octavia-client-certs",
				"mountPath": "/etc/octavia/certs/client",
			},
		}

		values := map[string]interface{}{
			"images": map[string]interface{}{
				"tags": tags,
			},
			"endpoints": endpoints,
			"conf": map[string]interface{}{
				"octavia": map[string]interface{}{
					"DEFAULT": map[string]interface{}{
						"log_config_append": nil,
					},
					"certificates": map[string]interface{}{
						"ca_certificate":            "/etc/octavia/certs/server/ca.crt",
						"ca_private_key":            "/etc/octavia/certs/server/tls.key",
						"ca_private_key_passphrase": "",
						"endpoint_type":             "internalURL",
					},
					"cinder": map[string]interface{}{
						"endpoint_type": "internalURL",
					},
					"controller_worker": map[string]interface{}{
						"workers":               float32(4),
						"client_ca":             "/etc/octavia/certs/client/ca.crt",
						"amp_boot_network_list": octavia.Spec.AmphoraConfig.Network,
						"amp_flavor_id":         octavia.Spec.AmphoraConfig.Flavor,
						"amp_image_owner_id":    octavia.Spec.AmphoraConfig.ImageOwner,
						"amp_secgroup_list":     octavia.Spec.AmphoraConfig.SecurityGroup,
						"amp_ssh_key_name":      octavia.Spec.AmphoraConfig.SSHKeyName,
					},
					"glance": map[string]interface{}{
						"endpoint_type": "internalURL",
					},
					"haproxy_amphora": map[string]interface{}{
						"client_cert": "/etc/octavia/certs/client/tls-combined.pem",
						"server_ca":   "/etc/octavia/certs/server/ca.crt",
					},
					"health_manager": map[string]interface{}{
						"controller_ip_port_list": strings.Join(octavia.Spec.HealthManagers, ","),
						"heartbeat_key":           endpointConfig.OctaviaHeartbeatKey,
					},
					"oslo_messaging_notifications": map[string]interface{}{
						"driver": "noop",
					},
					"neutron": map[string]interface{}{
						"endpoint_type": "internalURL",
					},
					"nova": map[string]interface{}{
						"endpoint_type": "internalURL",
					},
					"service_auth": map[string]interface{}{
						"endpoint_type": "internalURL",
					},
				},
			},
			"pod": map[string]interface{}{
				"mounts": map[string]interface{}{
					"octavia_api": map[string]interface{}{
						"octavia_api": map[string]interface{}{
							"volumeMounts": volumeMounts,
							"volumes":      volumes,
						},
					},
					"octavia_worker": map[string]interface{}{
						"octavia_worker": map[string]interface{}{
							"volumeMounts": volumeMounts,
							"volumes":      volumes,
						},
					},
					"octavia_housekeeping": map[string]interface{}{
						"octavia_housekeeping": map[string]interface{}{
							"volumeMounts": volumeMounts,
							"volumes":      volumes,
						},
					},
					"octavia_health_manager": map[string]interface{}{
						"octavia_health_manager": map[string]interface{}{
							"volumeMounts": volumeMounts,
							"volumes":      volumes,
						},
					},
				},
				"replicas": map[string]interface{}{
					"api":          octavia.Spec.Replicas,
					"worker":       octavia.Spec.Replicas,
					"housekeeping": octavia.Spec.Replicas,
				},
			},
			"manifests": map[string]interface{}{
				"ingress_api":         false,
				"service_ingress_api": false,
			},
		}

		overrides, err := octavia.Spec.Overrides.GetAsMap()
		if err != nil {
			return nil, err
		}

		return chartutil.CoalesceTables(values, overrides), nil
	})

	postHook := hook.PostHookFunc(func(u *unstructured.Unstructured, release release.Release, _ logr.Logger) error {
		octavia := &openstackv1alpha1.Octavia{}
		if err = runtime.DefaultUnstructuredConverter.FromUnstructured(u.Object, octavia); err != nil {
			return err
		}

		ingress := &networkingv1.Ingress{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "octavia-api",
				Namespace: u.GetNamespace(),
			},
		}
		_, err := ctrl.CreateOrUpdate(context.Background(), r.Client, ingress, func() error {
			GenerateIngress(ingress, &octavia.Spec.Ingress, endpoints.GetPortFromChart(chart, "load_balancer", "api"))
			return ctrl.SetControllerReference(octavia, ingress, r.Scheme)
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
			Kind:    "Octavia",
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
