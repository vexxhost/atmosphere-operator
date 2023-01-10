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
	"fmt"

	monitoringv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	corev1 "k8s.io/api/core/v1"
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
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/chartutil"
	"helm.sh/helm/v3/pkg/release"

	"github.com/vexxhost/atmosphere-operator/apis/infra/v1alpha1"
	"github.com/vexxhost/atmosphere-operator/pkg/images"
	"github.com/vexxhost/atmosphere-operator/pkg/monitoring"
)

// MemcachedReconciler reconciles a Memcached object
type MemcachedReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=core,resources=secrets,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=services,verbs=get;create;update;delete
//+kubebuilder:rbac:groups=infra.atmosphere.vexxhost.com,resources=memcacheds,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=infra.atmosphere.vexxhost.com,resources=memcacheds/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=infra.atmosphere.vexxhost.com,resources=memcacheds/finalizers,verbs=update
//+kubebuilder:rbac:groups=monitoring.coreos.com,resources=servicemonitors,verbs=get;create;update
//+kubebuilder:rbac:groups=monitoring.coreos.com,resources=prometheusrules,verbs=get;create;update

func (r *MemcachedReconciler) createOrUpdateService(memcached *v1alpha1.Memcached) (*corev1.Service, error) {
	service := &corev1.Service{
		ObjectMeta: ctrl.ObjectMeta{
			Name:      fmt.Sprintf("%s-metrics", memcached.Name),
			Namespace: memcached.Namespace,
		},
	}

	_, err := ctrl.CreateOrUpdate(context.TODO(), r.Client, service, func() error {
		service.Labels = map[string]string{"application": memcached.Name, "component": "server"}
		service.Spec.Selector = map[string]string{"application": memcached.Name, "component": "server"}
		service.Spec.Ports = []corev1.ServicePort{
			{
				Name:     "metrics",
				Port:     9150,
				Protocol: "TCP",
			},
		}

		return ctrl.SetControllerReference(memcached, service, r.Scheme)
	})

	if err != nil {
		return nil, err
	}

	return service, nil
}

func (r *MemcachedReconciler) createOrUpdateServiceMonitor(memcached *v1alpha1.Memcached, service *corev1.Service) (*monitoringv1.ServiceMonitor, error) {
	serviceMonitor := &monitoringv1.ServiceMonitor{
		ObjectMeta: metav1.ObjectMeta{
			Name:      memcached.Name,
			Namespace: memcached.Namespace,
		},
	}

	_, err := ctrl.CreateOrUpdate(context.Background(), r.Client, serviceMonitor, func() error {
		serviceMonitor.Labels = map[string]string{
			"app.kubernetes.io/managed-by": "atmosphere-operator",
			"app.kubernetes.io/instance":   memcached.Name,
			"app.kubernetes.io/name":       "memcached",
		}
		serviceMonitor.Spec = monitoringv1.ServiceMonitorSpec{
			Endpoints: []monitoringv1.Endpoint{
				{
					Path: "/metrics",
					Port: "metrics",
					RelabelConfigs: []*monitoringv1.RelabelConfig{
						{
							Action: "replace",
							SourceLabels: []monitoringv1.LabelName{
								"__meta_kubernetes_pod_name",
							},
							TargetLabel: "instance",
						},
						{
							Action: "labeldrop",
							Regex:  "^(container|endpoint|namespace|pod|service)$",
						},
					},
				},
			},
			JobLabel: "application",
			NamespaceSelector: monitoringv1.NamespaceSelector{
				MatchNames: []string{
					memcached.Namespace,
				},
			},
			Selector: metav1.LabelSelector{
				MatchLabels: service.Labels,
			},
		}

		return ctrl.SetControllerReference(memcached, serviceMonitor, r.Scheme)
	})

	if err != nil {
		return nil, err
	}

	return serviceMonitor, nil
}

func (r *MemcachedReconciler) createOrUpdatePrometheusRule(memcached *v1alpha1.Memcached, serviceMonitor *monitoringv1.ServiceMonitor) (*monitoringv1.PrometheusRule, error) {
	prometheusRule := &monitoringv1.PrometheusRule{
		ObjectMeta: metav1.ObjectMeta{
			Name:      memcached.Name,
			Namespace: memcached.Namespace,
		},
	}

	groups, err := monitoring.GetMemcachedPrometheusRules()
	if err != nil {
		return nil, err
	}

	_, err = ctrl.CreateOrUpdate(context.Background(), r.Client, prometheusRule, func() error {
		prometheusRule.Labels = serviceMonitor.Labels
		prometheusRule.Spec = monitoringv1.PrometheusRuleSpec{
			Groups: groups,
		}

		return ctrl.SetControllerReference(memcached, prometheusRule, r.Scheme)
	})

	if err != nil {
		return nil, err
	}

	return prometheusRule, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *MemcachedReconciler) SetupWithManager(mgr ctrl.Manager) error {
	chart, err := loader.Load("helm-charts/memcached.tgz")
	if err != nil {
		return err
	}

	translator := values.TranslatorFunc(func(ctx context.Context, u *unstructured.Unstructured) (chartutil.Values, error) {
		memcached := &v1alpha1.Memcached{}
		if err = runtime.DefaultUnstructuredConverter.FromUnstructured(u.Object, memcached); err != nil {
			return nil, err
		}

		tags, err := images.GetImageTagsForOpenstackHelmChart(chart, memcached.Spec.ImageRepository)
		if err != nil {
			return nil, err
		}

		return chartutil.Values{
			"images": map[string]interface{}{
				"tags": tags,
			},
			"monitoring": map[string]interface{}{
				"prometheus": map[string]interface{}{
					"enabled": true,
				},
			},
		}, nil
	})

	postHook := hook.PostHookFunc(func(u *unstructured.Unstructured, _ release.Release, _ logr.Logger) error {
		memcached := &v1alpha1.Memcached{}
		if err = runtime.DefaultUnstructuredConverter.FromUnstructured(u.Object, memcached); err != nil {
			return err
		}

		service, err := r.createOrUpdateService(memcached)
		if err != nil {
			return err
		}

		serviceMonitor, err := r.createOrUpdateServiceMonitor(memcached, service)
		if err != nil {
			return err
		}

		_, err = r.createOrUpdatePrometheusRule(memcached, serviceMonitor)
		if err != nil {
			return err
		}

		return nil
	})

	reconciler, err := reconciler.New(
		reconciler.WithChart(*chart),
		reconciler.WithClient(r.Client),
		reconciler.WithValueTranslator(translator),
		reconciler.WithPostHook(postHook),
		reconciler.SkipPrimaryGVKSchemeRegistration(true),
		reconciler.WithGroupVersionKind(schema.GroupVersionKind{
			Group:   v1alpha1.GroupVersion.Group,
			Version: v1alpha1.GroupVersion.Version,
			Kind:    "Memcached",
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
