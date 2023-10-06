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
	CentralServerImage             = "negaranabestani/fake-server:v1"
	CentralServerContainerPort     = 9000
	CentralServerContainerPortName = "serverport"
	centralServerMountPath         = "/results"
	CentralServerBaseCommand       = "python3 central-server.py"
)

func (r *FLClusterReconciler) createOrUpdateCentralServer(ctx context.Context, cluster *v1alpha1.FLCluster) error {
	desiredPV, er8 := r.desiredCentralServerPV(cluster)
	desiredPVC, er0 := r.desiredCentralServerPVC(cluster)
	desiredDep, er2 := r.desiredCentralServerDeployment(cluster)
	desiredService, er1 := r.desiredCentralServerService(cluster)
	logger := log.FromContext(ctx)
	if er8 != nil {
		return er8
	}
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
	existingPV := &corev1.PersistentVolume{}
	err9 := r.Get(ctx, client.ObjectKeyFromObject(desiredPV), existingPV)
	if err9 != nil && !errors.IsNotFound(err9) {
		logger.Info(err9.Error())
		return err9
	}
	if errors.IsNotFound(err9) {
		if err := r.Create(ctx, desiredPV); err != nil {
			logger.Info(err.Error())
			return err
		}
	}
	if !reflect.DeepEqual(existingPV, desiredPV) {
		existingPV = desiredPV
		if err4 := r.Update(ctx, existingPV); err4 != nil {
			logger.Info(err4.Error())
			return err4
		}
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
					Affinity: &corev1.Affinity{
						NodeAffinity: &corev1.NodeAffinity{
							PreferredDuringSchedulingIgnoredDuringExecution: []corev1.PreferredSchedulingTerm{
								{
									Weight: 100,
									Preference: corev1.NodeSelectorTerm{
										MatchExpressions: []corev1.NodeSelectorRequirement{
											{
												Key:      "fl-role",
												Operator: corev1.NodeSelectorOperator("In"),
												Values: []string{
													"central-server",
												},
											},
										},
									},
								},
							},
						},
					},
					Containers: []corev1.Container{
						{
							Name:  cluster.Name + "-" + CentralServer,
							Image: CentralServerImage,
							Command: []string{
								CentralServerBaseCommand,
								"-ns",
								cluster.Namespace,
								"-cn",
								cluster.Name,
							},
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
					Port:       int32(CentralServerContainerPort),
					TargetPort: intstr.FromString(CentralServerContainerPortName),
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
	storage, err := resource.ParseQuantity("100M")
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
			VolumeName:       cluster.Name + "-" + CentralServer,
			StorageClassName: &sName,
			AccessModes: []corev1.PersistentVolumeAccessMode{
				"ReadWriteMany",
			},
			Resources: corev1.ResourceRequirements{
				Limits: nil,
				Requests: corev1.ResourceList{
					corev1.ResourceStorage: storage,
				},
			},
		},
	}

	if err := ctrl.SetControllerReference(cluster, pvc, r.Scheme); err != nil {
		return pvc, err
	}
	return pvc, nil
}

func (r *FLClusterReconciler) desiredCentralServerPV(cluster *v1alpha1.FLCluster) (*corev1.PersistentVolume, error) {
	storage, err := resource.ParseQuantity("100M")
	if err != nil {
		return nil, err
	}
	sName := "fast"
	fsType := "ext4"
	vMode := corev1.PersistentVolumeMode("Filesystem")
	pv := &corev1.PersistentVolume{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cluster.Name + "-" + CentralServer,
			Namespace: cluster.Namespace,
			Labels: map[string]string{
				"cluster": cluster.Name,
				"app":     CentralServerSelectorApp,
			},
		},
		Spec: corev1.PersistentVolumeSpec{
			VolumeMode:       &vMode,
			StorageClassName: sName,
			Capacity: corev1.ResourceList{
				corev1.ResourceStorage: storage,
			},
			AccessModes: []corev1.PersistentVolumeAccessMode{
				"ReadWriteMany",
			},
			NodeAffinity: &corev1.VolumeNodeAffinity{
				Required: &corev1.NodeSelector{
					NodeSelectorTerms: []corev1.NodeSelectorTerm{
						{
							MatchExpressions: []corev1.NodeSelectorRequirement{
								{
									Key:      "fl-role",
									Operator: corev1.NodeSelectorOperator("In"),
									Values: []string{
										"central-server",
									},
								},
							},
						},
					},
				},
			},
			PersistentVolumeReclaimPolicy: corev1.PersistentVolumeReclaimPolicy("Delete"),
			PersistentVolumeSource: corev1.PersistentVolumeSource{
				Local: &corev1.LocalVolumeSource{
					Path:   centralServerMountPath,
					FSType: &fsType,
				},
			},
		},
	}
	if err := ctrl.SetControllerReference(cluster, pv, r.Scheme); err != nil {
		return pv, err
	}
	return pv, nil
}
