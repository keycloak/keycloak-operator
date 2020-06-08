package common

import (
	v1 "github.com/openshift/api/route/v1"
	"github.com/pkg/errors"
	v12 "k8s.io/api/apps/v1"
	v13 "k8s.io/api/batch/v1"
)

const (
	ConditionStatusSuccess = "True"
)

func IsRouteReady(route *v1.Route) bool {
	if route == nil {
		return false
	}
	// A route has a an array of Ingress where each have an array of conditions
	for _, ingress := range route.Status.Ingress {
		for _, condition := range ingress.Conditions {
			// A successful route will have the admitted condition type as true
			if condition.Type == v1.RouteAdmitted && condition.Status != ConditionStatusSuccess {
				return false
			}
		}
	}
	return true
}

func IsStatefulSetReady(statefulSet *v12.StatefulSet) (bool, error) {
	if statefulSet == nil {
		return false, nil
	}
	// Check the correct number of replicas match and are ready
	numOfReplicasMatch := *statefulSet.Spec.Replicas == statefulSet.Status.Replicas
	allReplicasReady := statefulSet.Status.Replicas == statefulSet.Status.ReadyReplicas
	revisionsMatch := statefulSet.Status.CurrentRevision == statefulSet.Status.UpdateRevision

	return numOfReplicasMatch && allReplicasReady && revisionsMatch, nil
}

func IsDeploymentReady(deployment *v12.Deployment) (bool, error) {
	if deployment == nil {
		return false, nil
	}
	// A deployment has an array of conditions
	for _, condition := range deployment.Status.Conditions {
		// One failure condition exists, if this exists, return the Reason
		if condition.Type == v12.DeploymentReplicaFailure {
			return false, errors.Errorf(condition.Reason)
			// A successful deployment will have the progressing condition type as true
		} else if condition.Type == v12.DeploymentProgressing && condition.Status != ConditionStatusSuccess {
			return false, nil
		}
	}
	return true, nil
}

func IsJobReady(job *v13.Job) (bool, error) {
	if job == nil {
		return false, nil
	}

	return job.Status.Succeeded == 1, nil
}
