package controller

import (
	"context"
	"github.com/negaranabestani/kfl/api/v1alpha1"
	utils "github.com/negaranabestani/kfl/internal/util"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	ctrl "sigs.k8s.io/controller-runtime"
)

const (
	CentralServerSelectorApp       = "central-server"
	CentralServerImage             = "kennethreitz/httpbin"
	CentralServerContainerPort     = 8080
	CentralServerContainerPortName = "httpbin"
	centralServerServicePort       = 8080
)

func (r *FLClusterReconciler) createOrUpdateCentralServer(ctx context.Context, cluster v1alpha1.FLCluster) error {
	//TODO implement
	return nil
}

func (r *FLClusterReconciler) deleteCentralServer(ctx context.Context, cluster v1alpha1.FLCluster) error {
	//TODO implement
	return nil
}

func (r *FLClusterReconciler) centralServerDesiredDeployment(cluster *v1alpha1.FLCluster) (*appsv1.Deployment, error) {
	resources, _ := utils.ResourceRequirements(cluster.Spec.CentralServer.Resources)
	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cluster.Name + "-central-server-deployment",
			Namespace: cluster.Namespace,
			Labels: map[string]string{
				"cluster": cluster.Name,
				"app":     CentralServerSelectorApp,
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: utils.Int32ptr(cluster.Spec.CentralServer.Replica),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"cluster": cluster.Name,
					"app":     CentralServerSelectorApp,
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"cluster": cluster.Name,
						"app":     CentralServerSelectorApp,
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  cluster.Name + "-central-server-container",
							Image: CentralServerImage,
							Ports: []corev1.ContainerPort{
								{
									Name:          CentralServerContainerPortName,
									ContainerPort: CentralServerContainerPort,
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

func (r *FLClusterReconciler) centralServerDesiredService(cluster *v1alpha1.FLCluster) (*corev1.Service, error) {
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: cluster.Namespace,
			Name:      cluster.Name + "-central-server-service",
			Labels: map[string]string{
				"cluster": cluster.Name,
				"app":     CentralServerSelectorApp,
			},
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Name:       "default",
					Port:       int32(centralServerServicePort),
					TargetPort: intstr.FromString("default"),
				},
			},
			ClusterIP: corev1.ClusterIPNone,
			Selector: map[string]string{
				"cluster": cluster.Name,
				"app":     CentralServerSelectorApp,
			},
		},
	}

	if err := ctrl.SetControllerReference(cluster, service, r.Scheme); err != nil {
		return service, err
	}

	return service, nil
}
