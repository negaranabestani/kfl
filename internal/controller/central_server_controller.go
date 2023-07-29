package controller

import (
	"context"
	"github.com/negaranabestani/kfl/api/v1alpha1"
	"github.com/negaranabestani/kfl/values"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: "central-server-deployment",
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &cluster.Spec.CentralServer.Replica,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": values.CentralServerSelectorApp,
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": values.CentralServerSelectorApp,
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "central-server-container",
							Image: values.CentralServerImage,
							Ports: []corev1.ContainerPort{
								{
									Name:          values.CentralServerContainerPortName,
									ContainerPort: values.CentralServerContainerPort,
								},
							},
							Resources: corev1.ResourceRequirements{
								Limits: corev1.ResourceList{
									corev1.ResourceMemory: {
										//	TODO set memory
									},
									corev1.ResourceCPU: {
										//	TODO set cpu
									},
								},
								Requests: corev1.ResourceList{
									corev1.ResourceMemory: {
										//	TODO set memory
									},
									corev1.ResourceCPU: {
										//	TODO set cpu
									},
								},
							},
						},
					},
				},
			},
		},
	}
	return dep, nil
}
