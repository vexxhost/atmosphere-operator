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
	"github.com/vexxhost/atmosphere-operator/apis/openstack/v1alpha1"
	openstackv1alpha1 "github.com/vexxhost/atmosphere-operator/apis/openstack/v1alpha1"
	"github.com/vexxhost/atmosphere-operator/pkg/endpoints"
	"github.com/vexxhost/atmosphere-operator/pkg/images"
)

// KeystoneReconciler reconciles a Keystone object
type KeystoneReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// 1.6733330759414449e+09  ERROR   Reconciler error        {"controller": "keystone-controller", "object": {"name":"keystone","namespace":"openstack"}, "namespace": "openstack", "name": "keystone", "reconcileID": "6f86e0e8-c6de-445a-912f-228e07d7229b", "error": "roles.rbac.authorization.k8s.io \"keystone-openstack-keystone-credential-rotate\" is forbidden: user \"system:serviceaccount:atmosphere-system:atmosphere-controller-manager\" (groups=[\"system:serviceaccounts\" \"system:serviceaccounts:atmosphere-system\" \"system:authenticated\"]) is attempting to grant RBAC permissions not currently held:\n{APIGroups:[\"extensions\"], Resources:[\"jobs\"], Verbs:[\"get\" \"list\"]}\n{APIGroups:[\"extensions\"], Resources:[\"pods\"], Verbs:[\"get\" \"list\"]}"}

// TODO(mnaser): Tone down these RBAC rules
//+kubebuilder:rbac:groups=apps;batch;core;extensions,resources=deployments;daemonsets;endpoints;replicasets;services;statefulsets,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=apps;batch;core;extensions,resources=cronjobs;jobs,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=apps;batch;core;extensions,resources=pods;services;endpoints;persistentvolumeclaims;events;configmaps;secrets;serviceaccounts;namespaces,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=networking.k8s.io,resources=ingresses;networkpolicies,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=policy,resources=poddisruptionbudgets,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=rbac.authorization.k8s.io,resources=roles;rolebindings,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=storage.k8s.io,resources=storageclasses,verbs=get;list;watch
//+kubebuilder:rbac:groups=pxc.percona.com,resources=perconaxtradbclusters,verbs=get;list;watch
//+kubebuilder:rbac:groups=openstack.atmosphere.vexxhost.com,resources=keystones,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=openstack.atmosphere.vexxhost.com,resources=keystones/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=openstack.atmosphere.vexxhost.com,resources=keystones/finalizers,verbs=update
//+kubebuilder:rbac:groups=infra.atmosphere.vexxhost.com,resources=rabbitmqclusters,verbs=get;list;watch;create;update;patch

// SetupWithManager sets up the controller with the Manager.
func (r *KeystoneReconciler) SetupWithManager(mgr ctrl.Manager) error {
	chart, err := loader.Load("helm-charts/keystone.tgz")
	if err != nil {
		return err
	}

	translator := values.TranslatorFunc(func(ctx context.Context, u *unstructured.Unstructured) (chartutil.Values, error) {
		keystone := &openstackv1alpha1.Keystone{}
		if err = runtime.DefaultUnstructuredConverter.FromUnstructured(u.Object, keystone); err != nil {
			return nil, err
		}

		tags, err := images.GetImageTagsForOpenstackHelmChart(chart, keystone.Spec.ImageRepository)
		if err != nil {
			return nil, err
		}

		endpointConfig, err := endpoints.NewConfig(
			endpoints.WithNamespace(keystone.Namespace),
			endpoints.WithKeystone(ctx, r.Client, keystone),
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
				"keystone": map[string]interface{}{
					"DEFAULT": map[string]interface{}{
						"log_config_append": nil,
					},
					"auth": map[string]interface{}{
						"methods": "password,token,openid,application_credential",
					},
					"cors": map[string]interface{}{
						"allowed_origins": "*",
					},
					"federation": map[string]interface{}{
						"assertion_prefix":    "OIDC-",
						"remote_id_attribute": "OIDC-iss",
						// TODO(mnaser): Look-up using endpoints
						// "trusted_dashboard": "https://{{ openstack_helm_endpoints_horizon_api_host }}/auth/websso/"
					},
					"identity": map[string]interface{}{
						"domain_configurations_from_database": true,
					},
					"oslo_messaging_notifications": map[string]interface{}{
						"driver": "noop",
					},
				},
			},
			"pod": map[string]interface{}{
				"replicas": map[string]interface{}{
					"api": keystone.Spec.Replicas,
				},
			},
			"manifests": map[string]interface{}{
				"job_credential_cleanup": false,
				"ingress_api":            false,
				"service_ingress_api":    false,
			},
		}

		overrides, err := keystone.Spec.Overrides.GetAsMap()
		if err != nil {
			return nil, err
		}
		fmt.Println(overrides)

		return chartutil.CoalesceTables(values, overrides), nil
	})

	preHook := hook.PreHookFunc(func(u *unstructured.Unstructured, _ chartutil.Values, _ logr.Logger) error {
		keystone := &openstackv1alpha1.Keystone{}
		if err = runtime.DefaultUnstructuredConverter.FromUnstructured(u.Object, keystone); err != nil {
			return err
		}

		// TODO: Wait for PXC to be ready?

		return err
	})

	postHook := hook.PostHookFunc(func(u *unstructured.Unstructured, release release.Release, _ logr.Logger) error {
		keystone := &openstackv1alpha1.Keystone{}
		if err = runtime.DefaultUnstructuredConverter.FromUnstructured(u.Object, keystone); err != nil {
			return err
		}

		ingress := &networkingv1.Ingress{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "keystone-api",
				Namespace: u.GetNamespace(),
			},
		}
		_, err := ctrl.CreateOrUpdate(context.Background(), r.Client, ingress, func() error {
			GenerateIngress(ingress, &keystone.Spec.Ingress, 5000)
			return ctrl.SetControllerReference(keystone, ingress, r.Scheme)
		})

		return err
	})

	reconciler, err := reconciler.New(
		reconciler.WithChart(*chart),
		reconciler.WithClient(r.Client),
		reconciler.WithPreHook(preHook),
		reconciler.WithPostHook(postHook),
		reconciler.WithValueTranslator(translator),
		reconciler.SkipPrimaryGVKSchemeRegistration(true),
		reconciler.WithGroupVersionKind(schema.GroupVersionKind{
			Group:   v1alpha1.GroupVersion.Group,
			Version: v1alpha1.GroupVersion.Version,
			Kind:    "Keystone",
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
