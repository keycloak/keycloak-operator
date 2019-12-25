package e2e

import (
	"context"
	"fmt"
	"testing"
	"time"

	framework "github.com/operator-framework/operator-sdk/pkg/test"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	dynclient "sigs.k8s.io/controller-runtime/pkg/client"
)

// Stolen from https://github.com/kubernetes/kubernetes/blob/master/test/e2e/framework/util.go
func WaitForStatefulSetReplicasReady(t *testing.T, c kubernetes.Interface, statefulSetName, ns string) error {
	t.Logf("waiting up to %v for StatefulSet %s to have all replicas ready", pollTimeout, statefulSetName)
	for start := time.Now(); time.Since(start) < pollTimeout; time.Sleep(pollRetryInterval) {
		sts, err := c.AppsV1().StatefulSets(ns).Get(statefulSetName, metav1.GetOptions{})
		if err != nil {
			t.Logf("get StatefulSet %s failed, ignoring for %v: %v", statefulSetName, pollRetryInterval, err)
			continue
		}
		if sts.Status.ReadyReplicas == *sts.Spec.Replicas {
			t.Logf("all %d replicas of StatefulSet %s are ready. (%v)", sts.Status.ReadyReplicas, statefulSetName, time.Since(start))
			return nil
		}
		t.Logf("statefulSet %s found but there are %d ready replicas and %d total replicas.", statefulSetName, sts.Status.ReadyReplicas, *sts.Spec.Replicas)
	}
	return fmt.Errorf("statefulSet %s still has unready pods within %v", statefulSetName, pollTimeout)
}

func WaitForPersistentVolumeClaimCreated(t *testing.T, c kubernetes.Interface, persistentVolumeClaimName, ns string) error {
	t.Logf("waiting up to %v for PersistentVolumeClaim %s to be created", pollTimeout, persistentVolumeClaimName)
	for start := time.Now(); time.Since(start) < pollTimeout; time.Sleep(pollRetryInterval) {
		pvc, err := c.CoreV1().PersistentVolumeClaims(ns).Get(persistentVolumeClaimName, metav1.GetOptions{})
		if err != nil {
			t.Logf("get PersistentVolumeClaim %s failed, ignoring for %v: %v", persistentVolumeClaimName, pollRetryInterval, err)
			continue
		}
		if pvc.Status.Phase == "Bound" {
			t.Logf("PersistentVolumeClaim is bound")
			return nil
		}
		t.Logf("persistentVolumeClaim %s found but is not bound", persistentVolumeClaimName)
	}
	return fmt.Errorf("persistentVolumeClaim %s still not created or not bound within %v", persistentVolumeClaimName, pollTimeout)
}

func Create(f *framework.Framework, obj runtime.Object, ctx *framework.TestCtx) error {
	return f.Client.Create(context.TODO(), obj, &framework.CleanupOptions{TestContext: ctx, Timeout: cleanupTimeout, RetryInterval: cleanupRetryInterval})
}

func Get(f *framework.Framework, key dynclient.ObjectKey, obj runtime.Object, ctx *framework.TestCtx) error {
	return f.Client.Get(context.TODO(), key, obj)
}
