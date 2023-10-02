package controller

import (
	"context"
	"github.com/negaranabestani/kfl/api/v1alpha1"
	utils "github.com/negaranabestani/kfl/internal/controller/util"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"reflect"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"strconv"
)

const (
	EdgeClient                  = "edge-client"
	EdgeClientSelectorApp       = "edge-client"
	EdgeClientImage             = "curlimages/curl"
	EdgeClientContainerPort     = 9001
	EdgeClientContainerPortName = "httpbin"
	EdgeClientServicePort       = 9001
)

func (r *FLClusterReconciler) createOrUpdateEdgeClient(ctx context.Context, cluster *v1alpha1.FLCluster, i int) error {
	desiredDep, er2 := r.desiredEdgeClientDeployment(cluster, i)
	desiredService, er1 := r.desiredEdgeClientService(cluster, i)
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

func (r *FLClusterReconciler) deleteEdgeClient(ctx context.Context, cluster v1alpha1.FLCluster) error {
	//TODO implement and add required input params
	return nil
}

func (r *FLClusterReconciler) desiredEdgeClientDeployment(cluster *v1alpha1.FLCluster, i int) (*appsv1.Deployment, error) {

	resources, _ := utils.ResourceRequirements(cluster.Spec.EdgeClient[i].Resources)
	deploymentTemplate := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cluster.Name + "-" + EdgeClient + string(strconv.Itoa(i)),
			Namespace: cluster.Namespace,
			Labels: map[string]string{
				"cluster": cluster.Name,
				"app":     EdgeClient,
				"device":  EdgeClient + string(strconv.Itoa(i)),
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: utils.Int32ptr(cluster.Spec.EdgeClient[i].Replica),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"cluster": cluster.Name,
					"app":     EdgeClient,
					"device":  EdgeClient + string(strconv.Itoa(i)),
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"cluster": cluster.Name,
						"app":     EdgeClient,
						"device":  EdgeClient + string(strconv.Itoa(i)),
					},
				},
				Spec: corev1.PodSpec{
					Affinity: &corev1.Affinity{
						PodAntiAffinity: &corev1.PodAntiAffinity{
							PreferredDuringSchedulingIgnoredDuringExecution: []corev1.WeightedPodAffinityTerm{
								{
									Weight: 100,
									PodAffinityTerm: corev1.PodAffinityTerm{
										LabelSelector: &metav1.LabelSelector{
											MatchLabels: map[string]string{
												"cluster": cluster.Name,
												"app":     CentralServerSelectorApp,
											},
										},
										TopologyKey: "kubernetes.io/hostname",
									},
								},
								{
									Weight: 100,
									PodAffinityTerm: corev1.PodAffinityTerm{
										LabelSelector: &metav1.LabelSelector{
											MatchLabels: map[string]string{
												"cluster": cluster.Name,
												"app":     edgeServerSelectorApp,
											},
										},
										TopologyKey: "kubernetes.io/hostname",
									},
								},
							},
						},
					},
					Containers: []corev1.Container{
						{
							Name:  cluster.Name + "-" + EdgeClient + string(strconv.Itoa(i)),
							Image: EdgeClientImage,
							Ports: []corev1.ContainerPort{
								{
									Name:          "http",
									ContainerPort: EdgeClientContainerPort,
								},
							},
							Resources: *resources,
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

func (r *FLClusterReconciler) desiredEdgeClientService(cluster *v1alpha1.FLCluster, i int) (*corev1.Service, error) {
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: cluster.Namespace,
			Name:      cluster.Name + "-" + EdgeClient + string(strconv.Itoa(i)),
			Labels: map[string]string{
				"cluster": cluster.Name,
				"app":     EdgeClientSelectorApp,
				"device":  EdgeClient + string(strconv.Itoa(i)),
			},
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Name:       "default",
					Port:       EdgeClientServicePort,
					TargetPort: intstr.FromString("default"),
				},
			},
			Selector: map[string]string{
				"cluster": cluster.Name,
				"app":     EdgeClientSelectorApp,
				"device":  EdgeClient + string(strconv.Itoa(i)),
			},
		},
	}

	if err := ctrl.SetControllerReference(cluster, service, r.Scheme); err != nil {
		return service, err
	}

	return service, nil
}
