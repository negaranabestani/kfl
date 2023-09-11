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
	edgeServer                  = "edge-server"
	edgeServerSelectorApp       = "edge-server"
	edgeServerImage             = "kennethreitz/httpbin"
	edgeServerContainerPort     = 8080
	edgeServerContainerPortName = "httpbin"
	edgeServerServicePort       = 8080
)

func (r *FLClusterReconciler) createOrUpdateEdgeServer(ctx context.Context, cluster *v1alpha1.FLCluster) error {
	desiredDep, er2 := r.edgeServerDesiredDeployment(cluster)
	desiredService, er1 := r.edgeServerDesiredService(cluster)
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
	return nil
}

func (r *FLClusterReconciler) deleteEdgeServer(ctx context.Context, cluster v1alpha1.FLCluster) error {
	//TODO implement and add required input params
	return nil
}

func (r *FLClusterReconciler) edgeServerDesiredDeployment(cluster *v1alpha1.FLCluster) (*appsv1.Deployment, error) {
	resources, _ := utils.ResourceRequirements(cluster.Spec.EdgeServer.Resources)
	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cluster.Name + "-" + edgeServer,
			Namespace: cluster.Namespace,
			Labels: map[string]string{
				"cluster": cluster.Name,
				"app":     edgeServerSelectorApp,
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: utils.Int32ptr(cluster.Spec.EdgeServer.Replica),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"cluster": cluster.Name,
					"app":     edgeServerSelectorApp,
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"cluster": cluster.Name,
						"app":     edgeServerSelectorApp,
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  cluster.Name + "-" + edgeServer,
							Image: edgeServerImage,
							Ports: []corev1.ContainerPort{
								{
									Name:          edgeServerContainerPortName,
									ContainerPort: edgeServerContainerPort,
								},
							},
							Resources: *resources,
						},
					},
				},
			},
		},
	}
	if err := ctrl.SetControllerReference(cluster, dep, r.Scheme); err != nil {
		return dep, err
	}
	return dep, nil
}

func (r *FLClusterReconciler) edgeServerDesiredService(cluster *v1alpha1.FLCluster) (*corev1.Service, error) {
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: cluster.Namespace,
			Name:      cluster.Name + "-" + edgeServer,
			Labels: map[string]string{
				"cluster": cluster.Name,
				"app":     edgeServerSelectorApp,
			},
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Name:       "default",
					Port:       int32(edgeServerServicePort),
					TargetPort: intstr.FromString("default"),
				},
			},
			Selector: map[string]string{
				"cluster": cluster.Name,
				"app":     edgeServerSelectorApp,
			},
		},
	}

	if err := ctrl.SetControllerReference(cluster, service, r.Scheme); err != nil {
		return service, err
	}

	return service, nil
}
