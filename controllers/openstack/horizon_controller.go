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
	"fmt"

	"github.com/lithammer/dedent"
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

// HorizonReconciler reconciles a Horizon object
type HorizonReconciler struct {
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
//+kubebuilder:rbac:groups=openstack.atmosphere.vexxhost.com,resources=horizons,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=openstack.atmosphere.vexxhost.com,resources=horizons/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=openstack.atmosphere.vexxhost.com,resources=horizons/finalizers,verbs=update
//+kubebuilder:rbac:groups=infra.atmosphere.vexxhost.com,resources=rabbitmqclusters,verbs=get;list;watch;create;update;patch

// SetupWithManager sets up the controller with the Manager.
func (r *HorizonReconciler) SetupWithManager(mgr ctrl.Manager) error {
	chart, err := loader.Load("helm-charts/horizon.tgz")
	if err != nil {
		return err
	}

	translator := values.TranslatorFunc(func(ctx context.Context, u *unstructured.Unstructured) (chartutil.Values, error) {
		horizon := &openstackv1alpha1.Horizon{}
		if err = runtime.DefaultUnstructuredConverter.FromUnstructured(u.Object, horizon); err != nil {
			return nil, err
		}

		tags, err := images.GetImageTagsForOpenstackHelmChart(chart, horizon.Spec.ImageRepository)
		if err != nil {
			return nil, err
		}

		keystoneRef := horizon.Spec.KeystoneRef.WithNamespace(horizon.Namespace)

		endpointConfig, err := endpoints.NewConfig(
			endpoints.WithNamespace(horizon.Namespace),
			endpoints.WithKeystoneRef(ctx, r.Client, &keystoneRef),
			endpoints.WithHorizon(ctx, r.Client, horizon),
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
				"horizon": map[string]interface{}{
					"local_settings": map[string]interface{}{
						"config": map[string]interface{}{
							"secure_proxy_ssl_header":            "True",
							"horizon_images_upload_mode":         "direct",
							"openstack_enable_password_retrieve": "True",
							"raw": map[string]interface{}{
								"WEBSSO_KEYSTONE_URL": fmt.Sprintf("https://%s/v3", endpointConfig.KeystoneHost),
							},
						},
					},
					"local_settings_d": map[string]interface{}{
						"_50_monasca_ui_settings": dedent.Dedent(`
							from django.conf import settings
							from django.utils.translation import ugettext_lazy as _
							
							# Service group names (global across all projects):
							MONITORING_SERVICES_GROUPS = [
									{"name": _("OpenStack Services"), "groupBy": "service"},
									{"name": _("Servers"), "groupBy": "hostname"},
							]
							
							# Services being monitored
							MONITORING_SERVICES = getattr(
									settings, "MONITORING_SERVICES_GROUPS", MONITORING_SERVICES_GROUPS
							)
							
							MONITORING_SERVICE_VERSION = getattr(settings, "MONITORING_SERVICE_VERSION", "2_0")
							MONITORING_SERVICE_TYPE = getattr(settings, "MONITORING_SERVICE_TYPE", "monitoring")
							MONITORING_ENDPOINT_TYPE = getattr(
									# NOTE(trebskit) # will default to OPENSTACK_ENDPOINT_TYPE
									settings,
									"MONITORING_ENDPOINT_TYPE",
									None,
							)
							
							# Grafana button titles/file names (global across all projects):
							# GRAFANA_LINKS = [{"raw": True, "path": "monasca-dashboard", "title": "Sub page1"}]
							GRAFANA_LINKS = []
							DASHBOARDS = getattr(settings, "GRAFANA_LINKS", GRAFANA_LINKS)
							
							GRAFANA_URL = {"regionOne": "/grafana"}
							
							SHOW_GRAFANA_HOME = getattr(settings, "SHOW_GRAFANA_HOME", True)
							
							ENABLE_LOG_MANAGEMENT_BUTTON = getattr(settings, "ENABLE_LOG_MANAGEMENT_BUTTON", False)
							ENABLE_EVENT_MANAGEMENT_BUTTON = getattr(
									settings, "ENABLE_EVENT_MANAGEMENT_BUTTON", False
							)
							
							KIBANA_POLICY_RULE = getattr(settings, "KIBANA_POLICY_RULE", "monitoring:kibana_access")
							KIBANA_POLICY_SCOPE = getattr(settings, "KIBANA_POLICY_SCOPE", "monitoring")
							KIBANA_HOST = getattr(settings, "KIBANA_HOST", "http://192.168.10.6:5601/")
							
							OPENSTACK_SSL_NO_VERIFY = getattr(settings, "OPENSTACK_SSL_NO_VERIFY", False)
							OPENSTACK_SSL_CACERT = getattr(settings, "OPENSTACK_SSL_CACERT", None)
							
							POLICY_FILES = getattr(settings, "POLICY_FILES", {})
							POLICY_FILES.update(
									{
											"monitoring": "monitoring_policy.json",
									}
							)  # noqa
							setattr(settings, "POLICY_FILES", POLICY_FILES)
						`),
					},
					"extra_panels": []string{"designatedashboard", "heat_dashboard", "ironic_ui", "magnum_ui", "monitoring", "neutron_vpnaas_dashboard", "octavia_dashboard", "senlin_dashboard"},
					"policy": map[string]interface{}{
						"monitoring": map[string]interface{}{
							"default":                  "@",
							"monasca_user_role":        "role:monasca-user",
							"monitoring:monitoring":    "rule:monasca_user_role",
							"monitoring:kibana_access": "rule:monasca_user_role",
						},
					},
				},
			},
			"pod": map[string]interface{}{
				"replicas": map[string]interface{}{
					"api": horizon.Spec.Replicas,
				},
			},
			"manifests": map[string]interface{}{
				"ingress_api":         false,
				"service_ingress_api": false,
			},
		}

		overrides, err := horizon.Spec.Overrides.GetAsMap()
		if err != nil {
			return nil, err
		}

		return chartutil.CoalesceTables(overrides, values), nil
	})

	postHook := hook.PostHookFunc(func(u *unstructured.Unstructured, release release.Release, _ logr.Logger) error {
		horizon := &openstackv1alpha1.Horizon{}
		if err = runtime.DefaultUnstructuredConverter.FromUnstructured(u.Object, horizon); err != nil {
			return err
		}

		ingress := &networkingv1.Ingress{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "horizon-int",
				Namespace: u.GetNamespace(),
			},
		}
		_, err := ctrl.CreateOrUpdate(context.Background(), r.Client, ingress, func() error {
			GenerateIngress(ingress, &horizon.Spec.Ingress, endpoints.GetPortFromChart(chart, "dashboard", "web"))
			return ctrl.SetControllerReference(horizon, ingress, r.Scheme)
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
			Kind:    "Horizon",
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
