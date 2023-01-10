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

	"k8s.io/apimachinery/pkg/runtime"
	clientretry "k8s.io/client-go/util/retry"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"github.com/go-logr/logr"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	pxcv1 "github.com/percona/percona-xtradb-cluster-operator/pkg/apis/pxc/v1"
	infrav1alpha1 "github.com/vexxhost/atmosphere-operator/apis/infra/v1alpha1"
	"github.com/vexxhost/atmosphere-operator/pkg/builders"
)

const DatabaseDeletionFinalizer = "deletion.finalizers.databases.infra.atmosphere.vexxhost.com"

// DatabaseReconciler reconciles a Database object
type DatabaseReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// addFinalizerIfNeeded adds a deletion finalizer if the Database does not have one yet and is not marked for deletion
func (r *DatabaseReconciler) addFinalizerIfNeeded(ctx context.Context, database *infrav1alpha1.Database) error {
	if database.ObjectMeta.DeletionTimestamp.IsZero() && !controllerutil.ContainsFinalizer(database, DatabaseDeletionFinalizer) {
		controllerutil.AddFinalizer(database, DatabaseDeletionFinalizer)
		if err := r.Client.Update(ctx, database); err != nil {
			return err
		}
	}
	return nil
}

func (r *DatabaseReconciler) removeFinalizer(ctx context.Context, database *infrav1alpha1.Database) error {
	controllerutil.RemoveFinalizer(database, DatabaseDeletionFinalizer)
	return r.Client.Update(ctx, database)
}

func (r *DatabaseReconciler) prepareForDeletion(ctx context.Context, database *infrav1alpha1.Database) error {
	if controllerutil.ContainsFinalizer(database, DatabaseDeletionFinalizer) {
		if err := clientretry.RetryOnConflict(clientretry.DefaultRetry, func() error {
			pxc := &pxcv1.PerconaXtraDBCluster{
				ObjectMeta: metav1.ObjectMeta{
					Name:      database.Name,
					Namespace: database.Namespace,
				},
			}
			if err := r.Client.Delete(ctx, pxc); client.IgnoreNotFound(err) != nil {
				return fmt.Errorf("cannot delete PerconaXtradbCluster: %w", err)
			}
			return nil
		}); err != nil {
			ctrl.LoggerFrom(ctx).Error(err, "Database deletion")
		}

		if err := r.removeFinalizer(ctx, database); err != nil {
			ctrl.LoggerFrom(ctx).Error(err, "Failed to remove finalizer for deletion")
			return err
		}
	}
	return nil
}

func (r *DatabaseReconciler) setReconcileSuccess(ctx context.Context, database *infrav1alpha1.Database, condition corev1.ConditionStatus, reason, msg string) {
	// TODO(oleks): Add status field in database type

	// database.Status.SetCondition(status.ReconcileSuccess, condition, reason, msg)
	// if writerErr := r.Status().Update(ctx, database); writerErr != nil {
	// 	ctrl.LoggerFrom(ctx).Error(writerErr, "Failed to update Custom Resource status",
	// 		"namespace", database.Namespace,
	// 		"name", database.Name)
	// }
}

// createOrUpdate updates current rabbitmqCluster with operator defaults from the Reconciler
// it handles RabbitMQ image, imagePullSecrets, and user updater image
func (r *DatabaseReconciler) createOrUpdate(ctx context.Context, database *infrav1alpha1.Database) (ctrl.Result, error) {

	var operationResult controllerutil.OperationResult
	pxc := &pxcv1.PerconaXtraDBCluster{
		ObjectMeta: metav1.ObjectMeta{
			Name:      database.Name,
			Namespace: database.Namespace,
		},
	}

	if err != nil {
		r.setReconcileSuccess(ctx, database, corev1.ConditionFalse, "Error", err.Error())
		return ctrl.Result{}, err
	}

	err = clientretry.RetryOnConflict(clientretry.DefaultRetry, func() error {
		var apiError error
		operationResult, apiError = controllerutil.CreateOrUpdate(ctx, r.Client, pxc, func() error {
			pxc, err := builders.defaultPxc(
				pxc,
			).overrideImageRegistry(
				database.Spec.ImageRepository,
			).build()
			if err != nil {
				return fmt.Errorf("failed to build pxc object: %w", err)
			}
			if err = controllerutil.SetControllerReference(database, pxc, r.Scheme); err != nil {
				return fmt.Errorf("failed setting controller reference: %w", err)
			}
		})
		return apiError
	})
	return ctrl.Result{}, err
}

//+kubebuilder:rbac:groups=infra.atmosphere.vexxhost.com,resources=databases,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=infra.atmosphere.vexxhost.com,resources=databases/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=infra.atmosphere.vexxhost.com,resources=databases/finalizers,verbs=update

func (r *DatabaseReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("Database", req.NamespacedName)

	instance := &infrav1alpha1.Database{}
	err := r.Get(ctx, req.NamespacedName, instance)
	if err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Check if the resource has been marked for deletion
	if !instance.ObjectMeta.DeletionTimestamp.IsZero() {
		logger.Info("Deleting")
		return ctrl.Result{}, r.prepareForDeletion(ctx, instance)
	}

	// Ensure the resource has a deletion marker
	if err := r.addFinalizerIfNeeded(ctx, database); err != nil {
		return ctrl.Result{}, err
	}

	return r.createOrUpdate(ctx, req, instance)
}

// SetupWithManager sets up the controller with the Manager.
func (r *DatabaseReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&infrav1alpha1.Database{}).
		Complete(r)
}
