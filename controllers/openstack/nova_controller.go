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

// NovaReconciler reconciles a Nova object
type NovaReconciler struct {
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
//+kubebuilder:rbac:groups=openstack.atmosphere.vexxhost.com,resources=novas,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=openstack.atmosphere.vexxhost.com,resources=novas/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=openstack.atmosphere.vexxhost.com,resources=novas/finalizers,verbs=update
//+kubebuilder:rbac:groups=infra.atmosphere.vexxhost.com,resources=rabbitmqclusters,verbs=get;list;watch;create;update;patch

// SetupWithManager sets up the controller with the Manager.
func (r *NovaReconciler) SetupWithManager(mgr ctrl.Manager) error {
	chart, err := loader.Load("helm-charts/nova.tgz")
	if err != nil {
		return err
	}

	translator := values.TranslatorFunc(func(ctx context.Context, u *unstructured.Unstructured) (chartutil.Values, error) {
		nova := &openstackv1alpha1.Nova{}
		if err = runtime.DefaultUnstructuredConverter.FromUnstructured(u.Object, nova); err != nil {
			return nil, err
		}

		tags, err := images.GetImageTagsForOpenstackHelmChart(chart, nova.Spec.ImageRepository)
		if err != nil {
			return nil, err
		}

		keystoneRef := nova.Spec.KeystoneRef.WithNamespace(nova.Namespace)
		placementRef := nova.Spec.PlacementRef.WithNamespace(nova.Namespace)
		glanceRef := nova.Spec.GlanceRef.WithNamespace(nova.Namespace)
		neutronRef := nova.Spec.NeutronRef.WithNamespace(nova.Namespace)
		ironicRef := nova.Spec.IronicRef.WithNamespace(nova.Namespace)

		endpointConfig, err := endpoints.NewConfig(
			endpoints.WithNamespace(nova.Namespace),
			endpoints.WithKeystoneRef(ctx, r.Client, &keystoneRef),
			endpoints.WithPlacementRef(ctx, r.Client, &placementRef),
			endpoints.WithGlanceRef(ctx, r.Client, &glanceRef),
			endpoints.WithNeutronRef(ctx, r.Client, &neutronRef),
			endpoints.WithIronicRef(ctx, r.Client, &ironicRef),
			endpoints.WithNova(ctx, r.Client, nova),
		)
		if err != nil {
			return nil, err
		}

		endpoints, err := endpoints.ForChart(chart, endpointConfig)
		if err != nil {
			return nil, err
		}

		values := map[string]interface{}{
			"labels": map[string]interface{}{
				"agent": map[string]interface{}{
					"compute_ironic": map[string]interface{}{
						"node_selector_key":   "openstack-control-plane",
						"node_selector_value": "enabled",
					},
				},
			},
			"bootstrap": map[string]interface{}{
				"structured": map[string]interface{}{
					"flavors": map[string]interface{}{
						"enabled": false,
					},
				},
			},
			"images": map[string]interface{}{
				"tags": tags,
			},
			"endpoints": endpoints,
			"network": map[string]interface{}{
				"ssh": map[string]interface{}{
					"enabled":     true,
					"public_key":  endpointConfig.NovaSSHPublicKey,
					"private_key": endpointConfig.NovaSSHPrivateKey,
				},
			},
			"conf": map[string]interface{}{
				"paste": map[string]interface{}{
					"composite:openstack_compute_api_v21": map[string]interface{}{
						"keystone": "cors http_proxy_to_wsgi compute_req_id faultwrap sizelimit authtoken keystonecontext osapi_compute_app_v21",
					},
					"composite:openstack_compute_api_v21_legacy_v2_compatible": map[string]interface{}{
						"keystone": "cors http_proxy_to_wsgi compute_req_id faultwrap sizelimit authtoken keystonecontext legacy_v2_compatible osapi_compute_app_v21",
					},
				},
				"nova": map[string]interface{}{
					"DEFAULT": map[string]interface{}{
						"allow_resize_to_same_host":        true,
						"cpu_allocation_ratio":             float32(4.5),
						"ram_allocation_ratio":             float32(0.9),
						"disk_allocation_ratio":            float32(3.0),
						"resume_guests_state_on_host_boot": true,
						"osapi_compute_workers":            float32(8),
						"metadata_workers":                 float32(8),
					},
					"cache": map[string]interface{}{
						"backend": "oslo_cache.memcache_pool",
					},
					"cinder": map[string]interface{}{
						"catalog_info": "volumev3::internalURL",
					},
					"conductor": map[string]interface{}{
						"workers": float32(8),
					},
					"compute": map[string]interface{}{
						"consecutive_build_service_disable_threshold": float32(0),
					},
					"cors": map[string]interface{}{
						"allowed_origin": "*",
						"allow_headers":  "X-Auth-Token,X-OpenStack-Nova-API-Version",
					},
					"filter_scheduler": map[string]interface{}{
						"enabled_filters": strings.Join([]string{
							"AvailabilityZoneFilter",
							"ComputeFilter",
							"AggregateTypeAffinityFilter",
							"ComputeCapabilitiesFilter",
							"PciPassthroughFilter",
							"ImagePropertiesFilter",
							"ServerGroupAntiAffinityFilter",
							"ServerGroupAffinityFilter",
						}, ","),
						"image_properties_default_architecture": "x86_64",
						"max_instances_per_host":                float32(200),
					},
					"glance": map[string]interface{}{
						"enable_rbd_download": true,
					},
					"neutron": map[string]interface{}{
						"metadata_proxy_shared_secret": endpointConfig.NovaMetadataSecret,
					},
					"oslo_messaging_notifications": map[string]interface{}{
						"driver": "noop",
					},
					"scheduler": map[string]interface{}{
						"workers": float32(8),
					},
				},
				"nova_ironic": map[string]interface{}{
					"DEFAULT": map[string]interface{}{
						"log_config_append":  nil,
						"force_config_drive": true,
					},
				},
			},
			"pod": map[string]interface{}{
				"replicas": map[string]interface{}{
					"api_metadata": nova.Spec.Replicas,
					"osapi":        nova.Spec.Replicas,
					"conductor":    nova.Spec.Replicas,
					"scheduler":    nova.Spec.Replicas,
					"novncproxy":   nova.Spec.Replicas,
					"spiceproxy":   nova.Spec.Replicas,
				},
			},
			"manifests": map[string]interface{}{
				"ingress_metadata":           false,
				"ingress_novncproxy":         false,
				"ingress_osapi":              false,
				"service_ingress_metadata":   false,
				"service_ingress_novncproxy": false,
				"service_ingress_osapi":      false,
			},
		}

		overrides, err := nova.Spec.Overrides.GetAsMap()
		if err != nil {
			return nil, err
		}

		return chartutil.CoalesceTables(overrides, values), nil
	})

	postHook := hook.PostHookFunc(func(u *unstructured.Unstructured, release release.Release, _ logr.Logger) error {
		nova := &openstackv1alpha1.Nova{}
		if err = runtime.DefaultUnstructuredConverter.FromUnstructured(u.Object, nova); err != nil {
			return err
		}

		ingress := &networkingv1.Ingress{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "nova-api",
				Namespace: u.GetNamespace(),
			},
		}
		_, err := ctrl.CreateOrUpdate(context.Background(), r.Client, ingress, func() error {
			GenerateIngress(ingress, &nova.Spec.Ingress, endpoints.GetPortFromChart(chart, "compute", "api"))
			return ctrl.SetControllerReference(nova, ingress, r.Scheme)
		})

		if err != nil {
			return err
		}

		vncIngress := &networkingv1.Ingress{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "nova-novncproxy",
				Namespace: u.GetNamespace(),
			},
		}
		_, err = ctrl.CreateOrUpdate(context.Background(), r.Client, vncIngress, func() error {
			GenerateIngress(vncIngress, &nova.Spec.VncIngress, endpoints.GetPortFromChart(chart, "compute_novnc_proxy", "novnc_proxy"))
			return ctrl.SetControllerReference(nova, vncIngress, r.Scheme)
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
			Kind:    "Nova",
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
