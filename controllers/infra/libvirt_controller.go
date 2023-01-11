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

package infra

import (
	"context"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/operator-framework/helm-operator-plugins/pkg/reconciler"
	"github.com/operator-framework/helm-operator-plugins/pkg/values"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/chartutil"

	"github.com/vexxhost/atmosphere-operator/apis/infra/v1alpha1"
	"github.com/vexxhost/atmosphere-operator/pkg/images"
)

// LibvirtReconciler reconciles a Libvirt object
type LibvirtReconciler struct {
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
//+kubebuilder:rbac:groups=infra.atmosphere.vexxhost.com,resources=libvirts,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=infra.atmosphere.vexxhost.com,resources=libvirts/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=infra.atmosphere.vexxhost.com,resources=libvirts/finalizers,verbs=update

// SetupWithManager sets up the controller with the Manager.
func (r *LibvirtReconciler) SetupWithManager(mgr ctrl.Manager) error {
	chart, err := loader.Load("helm-charts/libvirt.tgz")
	if err != nil {
		return err
	}

	translator := values.TranslatorFunc(func(ctx context.Context, u *unstructured.Unstructured) (chartutil.Values, error) {
		libvirt := &v1alpha1.Libvirt{}
		if err = runtime.DefaultUnstructuredConverter.FromUnstructured(u.Object, libvirt); err != nil {
			return nil, err
		}

		tags, err := images.GetImageTagsForOpenstackHelmChart(chart, libvirt.Spec.ImageRepository)
		if err != nil {
			return nil, err
		}

		values := map[string]interface{}{
			"images": map[string]interface{}{
				"tags": tags,
			},
			"conf": map[string]interface{}{
				"libvirt": map[string]interface{}{
					// TODO(mnaser): This allows live migration but we should really use TLS.
					"listen_addr": "0.0.0.0",
				},
			},
		}

		overrides, err := libvirt.Spec.Overrides.GetAsMap()
		if err != nil {
			return nil, err
		}

		return chartutil.CoalesceTables(overrides, values), nil
	})

	reconciler, err := reconciler.New(
		reconciler.WithChart(*chart),
		reconciler.WithClient(r.Client),
		reconciler.WithValueTranslator(translator),
		reconciler.SkipPrimaryGVKSchemeRegistration(true),
		reconciler.WithGroupVersionKind(schema.GroupVersionKind{
			Group:   v1alpha1.GroupVersion.Group,
			Version: v1alpha1.GroupVersion.Version,
			Kind:    "Libvirt",
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
