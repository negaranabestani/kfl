package controller

import (
	"context"
	"github.com/negaranabestani/kfl/api/v1alpha1"
	"github.com/stretchr/testify/require"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"testing"
)

func flClusterCreateOrUpdateComponentsTest(t *testing.T) {
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
			EdgeServer: []*v1alpha1.Device{
				{
					Replica: 1,
					Resources: v1alpha1.Resources{
						Cpu:    "1000m",
						Memory: "128Mi",
					},
				},
			},
			EdgeClient: []v1alpha1.Device{
				{
					Replica: 1,
					Resources: v1alpha1.Resources{
						Cpu:    "1000m",
						Memory: "128Mi",
					},
				},
			},
		},
	}
	r := &FLClusterReconciler{
		Client: fake.NewClientBuilder().WithScheme(scheme).WithObjects(flCluster).Build(),
		Scheme: scheme,
	}

	ctx := context.Background()
	logger := log.FromContext(ctx)

	err := r.createOrUpdateComponents(ctx, flCluster, logger)
	require.NoError(t, err)
}
