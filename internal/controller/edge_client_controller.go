package controller

import (
	"context"
	"github.com/negaranabestani/kfl/api/v1alpha1"
	utils "github.com/negaranabestani/kfl/internal/util"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"reflect"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	EdgeClient                  = "edge-client"
	EdgeClientSelectorApp       = "edge-client"
	EdgeClientImage             = "something"
	EdgeClientContainerPort     = 8080
	EdgeClientContainerPortName = "httpbin"
	EdgeClientServicePort       = 8080
)

func (r *FLClusterReconciler) createOrUpdateEdgeClient(ctx context.Context, cluster *v1alpha1.FLCluster) error {
	desiredDep, er2 := r.desiredEdgeClientDeployment(cluster)
	desiredService, er1 := r.desiredEdgeClientService(cluster)
	if er1 != nil {
		return er1
	}
	if er2 != nil {
		return er2
	}
	existingDep := &appsv1.Deployment{}
	err3 := r.Get(ctx, client.ObjectKeyFromObject(desiredDep), existingDep)
	if err3 != nil && !errors.IsNotFound(err3) {
		return err3
	}
	if errors.IsNotFound(err3) {
		if err := r.Create(ctx, desiredDep); err != nil {
			return err
		}
	}
	if !reflect.DeepEqual(existingDep, desiredDep) {
		existingDep = desiredDep
		if err4 := r.Update(ctx, existingDep); err4 != nil {
			return err4
		}
	}

	existingSer := &corev1.Service{}
	err5 := r.Get(ctx, client.ObjectKeyFromObject(desiredService), existingSer)
	if err5 != nil && !errors.IsNotFound(err5) {
		return err5
	}
	if errors.IsNotFound(err5) {
		if err := r.Create(ctx, desiredService); err != nil {
			return err
		}
	}
	if !reflect.DeepEqual(existingSer, desiredService) {
		existingSer = desiredService
		if err4 := r.Update(ctx, existingSer); err4 != nil {
			return err4
		}
	}

	err6 := r.Status().Update(ctx, cluster)
	if err6 != nil {
		return err6
	}
func (r *FLClusterReconciler) createOrUpdateEdgeClient(ctx context.Context, cluster *v1alpha1.FLCluster) error {
	//TODO implement and add required input params
	return nil
}

func (r *FLClusterReconciler) deleteEdgeClient(ctx context.Context, cluster v1alpha1.FLCluster) error {
	//TODO implement and add required input params
	return nil
}
