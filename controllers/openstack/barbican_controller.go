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
	openstackv1alpha1 "github.com/vexxhost/atmosphere-operator/apis/openstack/v1alpha1"
	"github.com/vexxhost/atmosphere-operator/pkg/endpoints"
	"github.com/vexxhost/atmosphere-operator/pkg/images"
)

// BarbicanReconciler reconciles a Barbican object
type BarbicanReconciler struct {
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
//+kubebuilder:rbac:groups=openstack.atmosphere.vexxhost.com,resources=barbicans,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=openstack.atmosphere.vexxhost.com,resources=barbicans/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=openstack.atmosphere.vexxhost.com,resources=barbicans/finalizers,verbs=update
//+kubebuilder:rbac:groups=infra.atmosphere.vexxhost.com,resources=rabbitmqclusters,verbs=get;list;watch;create;update;patch

// SetupWithManager sets up the controller with the Manager.
func (r *BarbicanReconciler) SetupWithManager(mgr ctrl.Manager) error {
	chart, err := loader.Load("helm-charts/barbican.tgz")
	if err != nil {
		return err
	}

	translator := values.TranslatorFunc(func(ctx context.Context, u *unstructured.Unstructured) (chartutil.Values, error) {
		barbican := &openstackv1alpha1.Barbican{}
		if err = runtime.DefaultUnstructuredConverter.FromUnstructured(u.Object, barbican); err != nil {
			return nil, err
		}

		tags, err := images.GetImageTagsForOpenstackHelmChart(chart, barbican.Spec.ImageRepository)
		if err != nil {
			return nil, err
		}

		keystoneRef := barbican.Spec.KeystoneRef.WithNamespace(barbican.Namespace)

		endpointConfig, err := endpoints.NewConfig(
			endpoints.WithNamespace(barbican.Namespace),
			endpoints.WithKeystoneRef(ctx, r.Client, &keystoneRef),
			endpoints.WithBarbican(ctx, r.Client, barbican),
		)
		if err != nil {
			return nil, err
		}

		endpoints, err := endpoints.ForChart(chart, endpointConfig)
		if err != nil {
			return nil, err
		}

		if endpointConfig.BarbicanKeyEncryptionKey == "" {
			return nil, fmt.Errorf("barbican key encryption key is empty")
		}

		values := map[string]interface{}{
			"images": map[string]interface{}{
				"tags": tags,
			},
			"endpoints": endpoints,
			"conf": map[string]interface{}{
				"barbican": map[string]interface{}{
					"DEFAULT": map[string]interface{}{
						"log_config_append": nil,
					},
					"oslo_messaging_notifications": map[string]interface{}{
						"driver": "noop",
					},
					"simple_crypto_plugin": map[string]interface{}{
						"kek": endpointConfig.BarbicanKeyEncryptionKey,
					},
				},
				"simple_crypto_kek_rewrap": map[string]interface{}{
					"old_kek": endpointConfig.BarbicanKeyEncryptionKey,
				},
			},
			"pod": map[string]interface{}{
				"replicas": map[string]interface{}{
					"api": barbican.Spec.Replicas,
				},
			},
			"manifests": map[string]interface{}{
				"ingress_api":         false,
				"service_ingress_api": false,
			},
		}

		overrides, err := barbican.Spec.Overrides.GetAsMap()
		if err != nil {
			return nil, err
		}

		return chartutil.CoalesceTables(values, overrides), nil
	})

	postHook := hook.PostHookFunc(func(u *unstructured.Unstructured, release release.Release, _ logr.Logger) error {
		barbican := &openstackv1alpha1.Barbican{}
		if err = runtime.DefaultUnstructuredConverter.FromUnstructured(u.Object, barbican); err != nil {
			return err
		}

		ingress := &networkingv1.Ingress{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "barbican-api",
				Namespace: u.GetNamespace(),
			},
		}
		_, err := ctrl.CreateOrUpdate(context.Background(), r.Client, ingress, func() error {
			GenerateIngress(ingress, &barbican.Spec.Ingress, endpoints.GetPortFromChart(chart, "key_manager", "api"))
			return ctrl.SetControllerReference(barbican, ingress, r.Scheme)
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
			Kind:    "Barbican",
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
