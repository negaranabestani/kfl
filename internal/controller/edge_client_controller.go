package controller

import (
	"context"
	"github.com/negaranabestani/kfl/api/v1alpha1"
	utils "github.com/negaranabestani/kfl/internal/util"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
)

const (
	EdgeClient      = "edge-client"
	EdgeClientImage = "something"
)

func (r *FLClusterReconciler) createOrUpdateEdgeClient(ctx context.Context, cluster v1alpha1.FLCluster) error {
	//TODO implement and add required input params
	return nil
}

func (r *FLClusterReconciler) deleteEdgeClient(ctx context.Context, cluster v1alpha1.FLCluster) error {
	//TODO implement and add required input params
	return nil
}

func (r *FLClusterReconciler) desiredEdgeClientDeployment(cluster *v1alpha1.FLCluster) (*appsv1.Deployment, error) {

	resources, _ := utils.ResourceRequirements(cluster.Spec.EdgeClient.Resources)
	deploymentTemplate := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cluster.Name + "-" + EdgeClient,
			Namespace: cluster.Namespace,
			Labels: map[string]string{
				"cluster": cluster.Name,
				"app":     EdgeClient,
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: utils.Int32ptr(cluster.Spec.EdgeClient.Replica),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"cluster": cluster.Name,
					"app":     EdgeClient,
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"cluster": cluster.Name,
						"app":     EdgeClient,
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  cluster.Name + "-" + EdgeClient,
							Image: EdgeClientImage,
							Ports: []corev1.ContainerPort{
								{
									Name:          "http",
									ContainerPort: 8080,
								},
							},
							Resources: *resources,
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "",
									MountPath: "",
								},
							},
						},
					},
					Volumes: []corev1.Volume{
						{
							Name:         "",
							VolumeSource: corev1.VolumeSource{},
						},
					},
				},
			},
		},
	}
	if err := ctrl.SetControllerReference(cluster, deploymentTemplate, r.Scheme); err != nil {
		return deploymentTemplate, err
	}
	return deploymentTemplate, nil

}
