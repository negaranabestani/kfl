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

package controller

import (
	"context"
	"fmt"
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/api/errors"
	"strings"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	kflv1alpha1 "github.com/negaranabestani/kfl/api/v1alpha1"
)

// FLClusterReconciler reconciles a FLCluster object
type FLClusterReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=kfl.aut.tech,resources=flclusters,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=kfl.aut.tech,resources=flclusters/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=kfl.aut.tech,resources=flclusters/finalizers,verbs=update
//+kubebuilder:rbac:groups=hdfs.aut.tech,resources=hdfsclusters/finalizers,verbs=update
//+kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups="",resources=persistentvolumeclaims,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups="",resources=persistentvolume,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups="",resources=services,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the FLCluster object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.15.0/pkg/reconcile
func (r *FLClusterReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	var flc kflv1alpha1.FLCluster
	err := r.Get(ctx, req.NamespacedName, &flc)
	if err != nil {
		if errors.IsNotFound(err) {
			logger.Info("Object not found, it could have been deleted")
			return ctrl.Result{}, nil
		}
		logger.Info("Error occurred during fetching the object")
		return ctrl.Result{}, err
	}

	requestArray := strings.Split(fmt.Sprint(req), "/")
	requestName := requestArray[1]

	if requestName == flc.Name {
		err = r.createOrUpdateComponents(ctx, &flc, logger)
		if err != nil {
			logger.Info("Error occurred during create Or Update Components")
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{}, nil
}

func (r *FLClusterReconciler) createOrUpdateComponents(ctx context.Context, flc *kflv1alpha1.FLCluster, logger logr.Logger) error {
	err := r.createOrUpdateCentralServer(ctx, flc)
	if err != nil {
		logger.Info("Error occurred during createOrUpdateCentralServer")
		return err
	}

	if flc.Spec.EdgeServer != nil {
		err1 := r.createOrUpdateEdgeServer(ctx, flc)
		if err1 != nil {
			logger.Info("Error occurred during createOrUpdateEdgeServer")
			return err1
		}
	}

	err2 := r.createOrUpdateEdgeClient(ctx, flc)
	if err2 != nil {
		logger.Info("Error occurred during createOrUpdateEdgeClient")
		return err2
	}

	return nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *FLClusterReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&kflv1alpha1.FLCluster{}).
		Complete(r)
}
