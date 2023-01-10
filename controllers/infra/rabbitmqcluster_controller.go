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

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/utils/pointer"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	rabbitmqv1beta1 "github.com/rabbitmq/cluster-operator/api/v1beta1"
	infrav1alpha1 "github.com/vexxhost/atmosphere-operator/apis/infra/v1alpha1"
	"github.com/vexxhost/atmosphere-operator/pkg/images"
)

// RabbitmqClusterReconciler reconciles a RabbitmqCluster object
type RabbitmqClusterReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=infra.atmosphere.vexxhost.com,resources=rabbitmqclusters,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=infra.atmosphere.vexxhost.com,resources=rabbitmqclusters/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=infra.atmosphere.vexxhost.com,resources=rabbitmqclusters/finalizers,verbs=update
//+kubebuilder:rbac:groups=rabbitmq.com,resources=rabbitmqclusters,verbs=get;list;watch;create;update;patch

func (r *RabbitmqClusterReconciler) createOrUpdateRabbitmqCluster(ctx context.Context, rabbitmq *infrav1alpha1.RabbitmqCluster) (*rabbitmqv1beta1.RabbitmqCluster, error) {
	cluster := &rabbitmqv1beta1.RabbitmqCluster{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("rabbitmq-%s", rabbitmq.Name),
			Namespace: rabbitmq.Namespace,
		},
	}

	_, err := ctrl.CreateOrUpdate(ctx, r.Client, cluster, func() error {
		cluster.Labels = rabbitmq.Labels

		image, err := images.GetImageReference("rabbitmq_server")
		if err != nil {
			return err
		}

		if rabbitmq.Spec.ImageRepository != "" {
			image, err = images.OverrideRegistry(image, rabbitmq.Spec.ImageRepository)
			if err != nil {
				return err
			}
		}

		cluster.Spec.Image = image.Remote()

		cluster.Spec.Affinity = &corev1.Affinity{
			NodeAffinity: &corev1.NodeAffinity{
				RequiredDuringSchedulingIgnoredDuringExecution: &corev1.NodeSelector{
					NodeSelectorTerms: []corev1.NodeSelectorTerm{
						{
							MatchExpressions: []corev1.NodeSelectorRequirement{
								{
									Key:      "openstack-control-plane",
									Operator: corev1.NodeSelectorOpIn,
									Values:   []string{"enabled"},
								},
							},
						},
					},
				},
			},
		}

		cluster.Spec.Rabbitmq.AdditionalConfig = "vm_memory_high_watermark.relative = 0.9\n"

		cluster.Spec.Resources = &corev1.ResourceRequirements{
			Requests: corev1.ResourceList{
				corev1.ResourceCPU:    resource.MustParse("500m"),
				corev1.ResourceMemory: resource.MustParse("1Gi"),
			},
			Limits: corev1.ResourceList{
				corev1.ResourceCPU:    resource.MustParse("1"),
				corev1.ResourceMemory: resource.MustParse("2Gi"),
			},
		}

		cluster.Spec.TerminationGracePeriodSeconds = pointer.Int64(15)

		return ctrl.SetControllerReference(rabbitmq, cluster, r.Scheme)
	})

	return cluster, err
}

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the RabbitmqCluster object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.13.0/pkg/reconcile
func (r *RabbitmqClusterReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	rabbitmq := &infrav1alpha1.RabbitmqCluster{}
	err := r.Get(ctx, req.NamespacedName, rabbitmq)
	if err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	cluster, err := r.createOrUpdateRabbitmqCluster(ctx, rabbitmq)
	if err != nil {
		return ctrl.Result{}, err
	}

	if cluster.Status.DefaultUser != nil {
		rabbitmq.Status.DefaultUser = *cluster.Status.DefaultUser
		err = r.Status().Update(ctx, rabbitmq)
	}

	return ctrl.Result{}, err
}

// SetupWithManager sets up the controller with the Manager.
func (r *RabbitmqClusterReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&infrav1alpha1.RabbitmqCluster{}).
		Complete(r)
}
