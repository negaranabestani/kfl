package controller

import (
	"context"
	"github.com/negaranabestani/kfl/api/v1alpha1"
	"github.com/stretchr/testify/assert"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"testing"
)

func CreateOrUpdateEdgeClientTest(t *testing.T) {
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
			EdgeClient: v1alpha1.Device{
				Replica: 1,
				Resources: v1alpha1.Resources{
					Cpu:    "1",
					Memory: "5Gi",
				},
			},
		},
	}
	r := &FLClusterReconciler{
		Client: fake.NewClientBuilder().WithScheme(scheme).WithObjects(flCluster).Build(),
		Scheme: scheme,
	}
	err := r.createOrUpdateEdgeClient(context.Background(), flCluster)
	if err != nil {
		t.Fatalf("failed to create or update edge client")
	}
	deployment := &appsv1.Deployment{}
	err1 := r.Get(context.Background(), client.ObjectKey{
		Name:      flCluster.Name + "-" + EdgeClient,
		Namespace: flCluster.Namespace,
	}, deployment)
	if err1 != nil {
		t.Errorf("failed to get deployment")
	}
	service := &corev1.Service{}
	err2 := r.Get(context.Background(), client.ObjectKey{
		Name:      flCluster.Name + "-" + EdgeClient,
		Namespace: flCluster.Namespace,
	}, service)
	if err2 != nil {
		t.Errorf("failed to get service")
	}
}
func DesiredEdgeClientDeploymentTest(t *testing.T) {
	scheme := runtime.NewScheme()
	_ = corev1.AddToScheme(scheme)
	_ = appsv1.AddToScheme(scheme)
	_ = v1alpha1.AddToScheme(scheme)

	flCluster := &v1alpha1.FLCluster{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-fl",
			Namespace: "sample",
		},
		Spec: v1alpha1.FLClusterSpec{
			EdgeClient: v1alpha1.Device{
				Replica: 1,
				Resources: v1alpha1.Resources{
					Cpu:    "1",
					Memory: "5Gi",
				},
			},
		},
	}
	r := &FLClusterReconciler{
		Client: fake.NewClientBuilder().WithScheme(scheme).WithObjects(flCluster).Build(),
		Scheme: scheme,
	}
	expectedName := flCluster.Name + "-" + EdgeClient
	expectedNamespace := flCluster.Namespace
	expectedResource := v1alpha1.Resources{
		Cpu:    "1",
		Memory: "5Gi",
	}
	expectedLabels := map[string]string{
		"cluster": flCluster.Name,
		"app":     EdgeClient,
	}
	expectedContainerName := flCluster.Name + "-central-server-container"
	deployment, err := r.desiredEdgeClientDeployment(flCluster)
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

func DesiredEdgeClientServiceTest(t *testing.T) {
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
			EdgeClient: v1alpha1.Device{
				Replica: 1,
				Resources: v1alpha1.Resources{
					Cpu:    "1",
					Memory: "5Gi",
				},
			},
		},
	}
	r := &FLClusterReconciler{
		Client: fake.NewClientBuilder().WithScheme(scheme).WithObjects(flCluster).Build(),
		Scheme: scheme,
	}

	expectedName := flCluster.Name + "-" + EdgeClient
	expectedNamespace := flCluster.Namespace
	expectedLabels := map[string]string{
		"cluster": flCluster.Name,
		"app":     EdgeClientSelectorApp,
	}

	service, err := r.desiredEdgeClientService(flCluster)

	assert.Nil(t, err)
	assert.Equal(t, expectedName, service.Name)
	assert.Equal(t, expectedNamespace, service.Namespace)
	assert.Equal(t, expectedLabels, service.Labels)
	assert.Equal(t, expectedLabels, service.Spec.Selector)
}
