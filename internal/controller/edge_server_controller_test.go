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

func edgeServerCreateOrUpdateTest(t *testing.T) {
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
			EdgeServer: v1alpha1.Device{
				Replica: 1,
				Resources: v1alpha1.Resources{
					Cpu:    "1000m",
					Memory: "5Gi",
				},
			},
		},
	}
	r := &FLClusterReconciler{
		Client: fake.NewClientBuilder().WithScheme(scheme).WithObjects(flCluster).Build(),
		Scheme: scheme,
	}
	err := r.createOrUpdateEdgeServer(context.Background(), flCluster)
	if err != nil {
		t.Fatalf("failed to create or update edge server")
	}
	deployment := &appsv1.Deployment{}
	err1 := r.Get(context.Background(), client.ObjectKey{
		Name:      flCluster.Name + "-" + edgeServer,
		Namespace: flCluster.Namespace,
	}, deployment)
	if err1 != nil {
		t.Errorf("failed to get deployment")
	}
	service := &corev1.Service{}
	err2 := r.Get(context.Background(), client.ObjectKey{
		Name:      flCluster.Name + "-" + edgeServer,
		Namespace: flCluster.Namespace,
	}, service)
	if err2 != nil {
		t.Errorf("failed to get service")
	}
}

func DesiredEdgeServerDeploymentTest(t *testing.T) {
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
			EdgeServer: v1alpha1.Device{
				Replica: 1,
				Resources: v1alpha1.Resources{
					Cpu:    "1000m",
					Memory: "5Gi",
				},
			},
		},
	}
	r := &FLClusterReconciler{
		Client: fake.NewClientBuilder().WithScheme(scheme).WithObjects(flCluster).Build(),
		Scheme: scheme,
	}
	expectedName := flCluster.Name + "-" + edgeServer
	expectedNamespace := flCluster.Namespace
	expectedResource := v1alpha1.Resources{
		Cpu:    "1000m",
		Memory: "5Gi",
	}
	expectedLabels := map[string]string{
		"cluster": flCluster.Name,
		"app":     edgeServerSelectorApp,
	}
	expectedContainerName := flCluster.Name + "-" + edgeServer
	deployment, err := r.desiredEdgeServerDeployment(flCluster)
	assert.Nil(t, err)
	assert.Equal(t, expectedName, deployment.Name)
	assert.Equal(t, expectedNamespace, deployment.Namespace)
	assert.Equal(t, expectedResource, deployment.Spec.Template.Spec.Containers[0].Resources)
	assert.Equal(t, expectedLabels, deployment.Labels)
	assert.Equal(t, expectedLabels, deployment.Spec.Selector)
	assert.Equal(t, expectedLabels, deployment.Spec.Template.Labels)
	assert.Equal(t, expectedContainerName, deployment.Spec.Template.Spec.Containers[0].Name)

	if *deployment.Spec.Replicas != 1 {
		t.Errorf("expected 1 edge server deployment replica go %d", *deployment.Spec.Replicas)
	}
}

func DesiredEdgeServerServiceTest(t *testing.T) {
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
			EdgeServer: v1alpha1.Device{
				Replica: 1,
				Resources: v1alpha1.Resources{
					Cpu:    "1000m",
					Memory: "5Gi",
				},
			},
		},
	}
	r := &FLClusterReconciler{
		Client: fake.NewClientBuilder().WithScheme(scheme).WithObjects(flCluster).Build(),
		Scheme: scheme,
	}

	expectedName := flCluster.Name + "-" + edgeServer
	expectedNamespace := flCluster.Namespace
	expectedLabels := map[string]string{
		"cluster": flCluster.Name,
		"app":     edgeServerSelectorApp,
	}

	service, err := r.desiredEdgeServerService(flCluster)

	assert.Nil(t, err)
	assert.Equal(t, expectedName, service.Name)
	assert.Equal(t, expectedNamespace, service.Namespace)
	assert.Equal(t, expectedLabels, service.Labels)
	assert.Equal(t, expectedLabels, service.Spec.Selector)
}
