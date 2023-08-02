package controller

import (
	"github.com/negaranabestani/kfl/api/v1alpha1"
	"github.com/stretchr/testify/assert"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"testing"
)

func DesiredDeploymentTest(t *testing.T) {
	scheme := runtime.NewScheme()
	_ = corev1.AddToScheme(scheme)
	_ = appsv1.AddToScheme(scheme)
	_ = v1alpha1.AddToScheme(scheme)

	flCluster := &v1alpha1.FLCluster{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-fl",
			Namespace: "default",
		},
		Spec: v1alpha1.FLClusterSpec{
			CentralServer: v1alpha1.Device{
				Replica: 1,
				Resources: v1alpha1.Resources{
					Cpu:    "1000m",
					Memory: "128Mi",
				},
			},
		},
	}
	r := &FLClusterReconciler{
		Client: fake.NewClientBuilder().WithScheme(scheme).WithObjects(flCluster).Build(),
		Scheme: scheme,
	}
	expectedName := flCluster.Name + "-central-server-deployment"
	expectedNamespace := flCluster.Namespace
	expectedResource := v1alpha1.Resources{
		Cpu:    "1000m",
		Memory: "128Mi",
	}
	expectedLabels := map[string]string{
		"cluster": flCluster.Name,
		"app":     CentralServerSelectorApp,
	}
	expectedContainerName := flCluster.Name + "-central-server-container"
	deployment, err := r.centralServerDesiredDeployment(flCluster)
	assert.Nil(t, err)
	assert.Equal(t, expectedName, deployment.Name)
	assert.Equal(t, expectedNamespace, deployment.Namespace)
	assert.Equal(t, expectedResource, deployment.Spec.Template.Spec.Containers[0].Resources)
	assert.Equal(t, expectedLabels, deployment.Labels)
	assert.Equal(t, expectedLabels, deployment.Spec.Selector)
	assert.Equal(t, expectedLabels, deployment.Spec.Template.Labels)
	assert.Equal(t, expectedContainerName, deployment.Spec.Template.Spec.Containers[0].Name)

	if *deployment.Spec.Replicas != 1 {
		t.Errorf("expected 1 central server deployment replica go %d", *deployment.Spec.Replicas)
	}
}

func DesiredServiceTest(t *testing.T) {
	scheme := runtime.NewScheme()
	_ = corev1.AddToScheme(scheme)
	_ = appsv1.AddToScheme(scheme)
	_ = v1alpha1.AddToScheme(scheme)

	flCluster := &v1alpha1.FLCluster{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-fl",
			Namespace: "default",
		},
		Spec: v1alpha1.FLClusterSpec{
			CentralServer: v1alpha1.Device{
				Replica: 1,
				Resources: v1alpha1.Resources{
					Cpu:    "1000m",
					Memory: "128Mi",
				},
			},
		},
	}
	r := &FLClusterReconciler{
		Client: fake.NewClientBuilder().WithScheme(scheme).WithObjects(flCluster).Build(),
		Scheme: scheme,
	}

	expectedName := flCluster.Name + "-central-server-service"
	expectedNamespace := flCluster.Namespace
	expectedLabels := map[string]string{
		"cluster": flCluster.Name,
		"app":     CentralServerSelectorApp,
	}

	service, err := r.centralServerDesiredService(flCluster)

	assert.Nil(t, err)
	assert.Equal(t, expectedName, service.Name)
	assert.Equal(t, expectedNamespace, service.Namespace)
	assert.Equal(t, expectedLabels, service.Labels)
	assert.Equal(t, expectedLabels, service.Spec.Selector)
}

func DesiredPVC(t *testing.T) {
	scheme := runtime.NewScheme()
	_ = corev1.AddToScheme(scheme)
	_ = appsv1.AddToScheme(scheme)
	_ = v1alpha1.AddToScheme(scheme)

	flCluster := &v1alpha1.FLCluster{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-fl",
			Namespace: "default",
		},
		Spec: v1alpha1.FLClusterSpec{
			CentralServer: v1alpha1.Device{
				Replica: 1,
				Resources: v1alpha1.Resources{
					Cpu:    "1000m",
					Memory: "128Mi",
				},
			},
		},
	}
	r := &FLClusterReconciler{
		Client: fake.NewClientBuilder().WithScheme(scheme).WithObjects(flCluster).Build(),
		Scheme: scheme,
	}

	expectedName := flCluster.Name + "-pvc"
	pvc, err := r.centralServerDesiredPVC(flCluster)
	assert.Nil(t, err)
	assert.Equal(t, expectedName, pvc.Name)

}
