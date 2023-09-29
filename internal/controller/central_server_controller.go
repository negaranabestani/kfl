package controller

import (
	"context"
	"github.com/negaranabestani/kfl/api/v1alpha1"
	utils "github.com/negaranabestani/kfl/internal/controller/util"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"reflect"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

const (
	CentralServer                  = "central-server"
	CentralServerSelectorApp       = "central-server"
	CentralServerImage             = "kennethreitz/httpbin"
	CentralServerContainerPort     = 9000
	CentralServerContainerPortName = "httpbin"
	centralServerServicePort       = 9000
	centralServerMountPath         = "/results"
)

func (r *FLClusterReconciler) createOrUpdateCentralServer(ctx context.Context, cluster *v1alpha1.FLCluster) error {
	desiredPVC, er0 := r.desiredCentralServerPVC(cluster)
	desiredDep, er2 := r.desiredCentralServerDeployment(cluster)
	desiredService, er1 := r.desiredCentralServerService(cluster)
	logger := log.FromContext(ctx)
	if er0 != nil {
		return er0
	}
	if er1 != nil {
		return er1
	}
	if er2 != nil {
		logger.Info(er2.Error())
		return er2
	}

	existingPVC := &corev1.PersistentVolumeClaim{}
	err3 := r.Get(ctx, client.ObjectKeyFromObject(desiredPVC), existingPVC)
	if err3 != nil && !errors.IsNotFound(err3) {
		logger.Info(err3.Error())
		return err3
	}
	if errors.IsNotFound(err3) {
		if err := r.Create(ctx, desiredPVC); err != nil {
			logger.Info(err.Error())
			return err
		}
	}
	if !reflect.DeepEqual(existingPVC, desiredPVC) {
		existingPVC = desiredPVC
		if err4 := r.Update(ctx, existingPVC); err4 != nil {
			logger.Info(err4.Error())
			return err4
		}
	}

	existingDep := &appsv1.Deployment{}
	err7 := r.Get(ctx, client.ObjectKeyFromObject(desiredDep), existingDep)
	if err7 != nil && !errors.IsNotFound(err7) {
		logger.Info(err7.Error())
		return err7
	}
	if errors.IsNotFound(err7) {
		if err := r.Create(ctx, desiredDep); err != nil {
			logger.Info(err.Error())
			return err
		}
	}
	if !reflect.DeepEqual(existingDep, desiredDep) {
		existingDep = desiredDep
		if err4 := r.Update(ctx, existingDep); err4 != nil {
			logger.Info(err4.Error())
			return err4
		}
	}

	existingSer := &corev1.Service{}
	err5 := r.Get(ctx, client.ObjectKeyFromObject(desiredService), existingSer)
	if err5 != nil && !errors.IsNotFound(err5) {
		logger.Info(err5.Error())
		return err5
	}
	if errors.IsNotFound(err5) {
		if err := r.Create(ctx, desiredService); err != nil {
			logger.Info(err.Error())
			return err
		}
	}
	if !reflect.DeepEqual(existingSer, desiredService) {
		existingSer = desiredService
		if err4 := r.Update(ctx, existingSer); err4 != nil {
			logger.Info(err4.Error())
			return err4
		}
	}

	err6 := r.Status().Update(ctx, cluster)
	if err6 != nil {
		logger.Info(err6.Error())
		return err6
	}
	return nil
}

func (r *FLClusterReconciler) deleteCentralServer(ctx context.Context, cluster v1alpha1.FLCluster) error {
	//TODO implement
	return nil
}

func (r *FLClusterReconciler) desiredCentralServerDeployment(cluster *v1alpha1.FLCluster) (*appsv1.Deployment, error) {
	resources, _ := utils.ResourceRequirements(cluster.Spec.CentralServer.Resources)
	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cluster.Name + "-" + CentralServer,
			Namespace: cluster.Namespace,
			Labels: map[string]string{
				"cluster": cluster.Name,
				"app":     CentralServer,
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: utils.Int32ptr(cluster.Spec.CentralServer.Replica),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"cluster": cluster.Name,
					"app":     CentralServer,
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"cluster": cluster.Name,
						"app":     CentralServer,
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  cluster.Name + "-" + CentralServer,
							Image: CentralServerImage,
							Ports: []corev1.ContainerPort{
								{
									Name:          CentralServerContainerPortName,
									ContainerPort: CentralServerContainerPort,
								},
							},
							Resources: *resources,
							VolumeMounts: []corev1.VolumeMount{
								{
									MountPath: centralServerMountPath,
									Name:      cluster.Name + "-data",
								},
							},
						},
					},
					Volumes: []corev1.Volume{
						{
							Name: cluster.Name + "-data",
							VolumeSource: corev1.VolumeSource{
								PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
									ClaimName: cluster.Name + "-" + CentralServer,
								},
							},
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

func (r *FLClusterReconciler) desiredCentralServerService(cluster *v1alpha1.FLCluster) (*corev1.Service, error) {
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: cluster.Namespace,
			Name:      cluster.Name + "-" + CentralServer,
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

func (r *FLClusterReconciler) desiredCentralServerPVC(cluster *v1alpha1.FLCluster) (*corev1.PersistentVolumeClaim, error) {
	storage, err := resource.ParseQuantity("1M")
	if err != nil {
		return nil, err
	}

	sName := "fast"
	pvc := &corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cluster.Name + "-" + CentralServer,
			Namespace: cluster.Namespace,
			Labels: map[string]string{
				"cluster": cluster.Name,
				"app":     CentralServerSelectorApp,
			},
		},
		Spec: corev1.PersistentVolumeClaimSpec{
			StorageClassName: &sName,
			AccessModes: []corev1.PersistentVolumeAccessMode{
				"ReadWriteOnce",
			},
			Resources: corev1.ResourceRequirements{
				Limits: nil,
				Requests: corev1.ResourceList{
					corev1.ResourceStorage: storage,
				},
			},
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"cluster": cluster.Name,
					"app":     CentralServerSelectorApp,
				},
			},
		},
	}

	if err := ctrl.SetControllerReference(cluster, pvc, r.Scheme); err != nil {
		return pvc, err
	}
	return pvc, nil
}
