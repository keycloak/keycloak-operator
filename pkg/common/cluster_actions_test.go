package common

import (
	"context"
	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/core/v1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"testing"
)

const TestNamespace = "test-namespace"
const TestApplicationName = "test-namespace"

var TestService = v1.Service{
	ObjectMeta: v12.ObjectMeta{
		Name:      TestApplicationName,
		Namespace: TestNamespace,
		Labels: map[string]string{
			"application": TestApplicationName,
		},
	},
	Spec: v1.ServiceSpec{
		Ports: []v1.ServicePort{
			{
				Port:       8443,
				TargetPort: intstr.Parse("8443"),
			},
		},
	},
}

var TestServiceWithModifiedPort = v1.Service{
	ObjectMeta: v12.ObjectMeta{
		Name:      TestApplicationName,
		Namespace: TestNamespace,
		Labels: map[string]string{
			"application": TestApplicationName,
		},
	},
	Spec: v1.ServiceSpec{
		Ports: []v1.ServicePort{
			{
				Port:       8551,
				TargetPort: intstr.Parse("8551"),
			},
		},
	},
}

var TestServiceSelector = client.ObjectKey{
	Name:      TestApplicationName,
	Namespace: TestNamespace,
}

func TestClusterActionRunner_Test_Create_Action(t *testing.T) {
	// given
	mockClient := fake.NewFakeClient()
	actionRunner := NewClusterActionRunner(mockClient)

	testedAction := GenericCreateAction{
		Ref: TestService.DeepCopy(),
		Msg: "non-important-text",
	}

	// when
	runnerError := actionRunner.RunAll(DesiredClusterState{testedAction})
	errorFromGetter := mockClient.Get(context.TODO(), TestServiceSelector, testedAction.Ref)

	// then
	assert.Nil(t, runnerError)
	assert.Nil(t, errorFromGetter)
}

func TestClusterActionRunner_Test_Update_Action(t *testing.T) {
	// given
	mockClient := fake.NewFakeClient(TestServiceWithModifiedPort.DeepCopy())
	actionRunner := NewClusterActionRunner(mockClient)

	testedAction := GenericUpdateAction{
		Ref: TestService.DeepCopy(),
		Msg: "non-important-text",
	}

	// when
	runnerError := actionRunner.RunAll(DesiredClusterState{testedAction})
	errorFromGetter := mockClient.Get(context.TODO(), TestServiceSelector, testedAction.Ref)

	// then
	assert.Nil(t, runnerError)
	assert.Nil(t, errorFromGetter)
	assert.Equal(t, testedAction.Ref, TestService.DeepCopy())
}
