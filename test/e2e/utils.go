package e2e

import (
	"fmt"
	"testing"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

const (
	PollDuration    = time.Second * 15
	TimeoutDuration = time.Second * 480
)

// Stolen from https://github.com/kubernetes/kubernetes/blob/master/test/e2e/framework/util.go
func WaitForStatefulSetReplicasReady(t *testing.T, c kubernetes.Interface, statefulSetName, ns string) error {
	t.Logf("waiting up to %v for StatefulSet %s to have all replicas ready", TimeoutDuration, statefulSetName)
	for start := time.Now(); time.Since(start) < TimeoutDuration; time.Sleep(PollDuration) {
		sts, err := c.AppsV1().StatefulSets(ns).Get(statefulSetName, metav1.GetOptions{})
		if err != nil {
			t.Logf("get StatefulSet %s failed, ignoring for %v: %v", statefulSetName, PollDuration, err)
			continue
		}
		if sts.Status.ReadyReplicas == *sts.Spec.Replicas {
			t.Logf("all %d replicas of StatefulSet %s are ready. (%v)", sts.Status.ReadyReplicas, statefulSetName, time.Since(start))
			return nil
		}
		t.Logf("statefulSet %s found but there are %d ready replicas and %d total replicas.", statefulSetName, sts.Status.ReadyReplicas, *sts.Spec.Replicas)
	}
	return fmt.Errorf("statefulSet %s still has unready pods within %v", statefulSetName, TimeoutDuration)
}
