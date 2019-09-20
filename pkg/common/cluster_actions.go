package common

import (
	"context"
	"fmt"
	"github.com/go-logr/logr"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type ActionRunner interface {
	RunAll(desiredState DesiredClusterState) error
	Create(obj runtime.Object) error
	Update(obj runtime.Object) error
}

type ClusterAction interface {
	Run(runner ActionRunner) (string, error)
}

type ClusterActionRunner struct {
	client client.Client
	logger logr.Logger
}

func NewClusterActionRunner(client client.Client, logger logr.Logger) ActionRunner {
	return &ClusterActionRunner{
		client: client,
		logger: logger,
	}
}

func (i *ClusterActionRunner) RunAll(desiredState DesiredClusterState) error {
	for index, action := range desiredState {
		msg, err := action.Run(i)
		if err != nil {
			i.logger.Info(fmt.Sprintf("(%5d) %10s %s", index, "FAILED", msg))
			return err
		}
		i.logger.Info(fmt.Sprintf("(%5d) %10s %s", index, "SUCCESS", msg))
	}

	return nil
}

func (i *ClusterActionRunner) Create(obj runtime.Object) error {
	return i.client.Create(context.TODO(), obj)
}

func (i *ClusterActionRunner) Update(obj runtime.Object) error {
	return i.client.Update(context.TODO(), obj)
}

// An action to create generic kubernetes resources
// (resources that don't require special treatment)
type GenericCreateAction struct {
	Ref runtime.Object
	Msg string
}

// An action to update generic kubernetes resources
// (resources that don't require special treatment)
type GenericUpdateAction struct {
	Ref runtime.Object
	Msg string
}

// Services need to have a valid ClusterIP and ResourceVersion in
// order to be updated
type ServiceUpdateAction struct {
	Ref             *v1.Service
	Msg             string
	ClusterIP       string
	ResourceVersion string
}

func (i GenericCreateAction) Run(runner ActionRunner) (string, error) {
	return i.Msg, runner.Create(i.Ref)
}

func (i GenericUpdateAction) Run(runner ActionRunner) (string, error) {
	return i.Msg, runner.Update(i.Ref)
}

func (i ServiceUpdateAction) Run(runner ActionRunner) (string, error) {
	i.Ref.ResourceVersion = i.ResourceVersion
	i.Ref.Spec.ClusterIP = i.ClusterIP
	return i.Msg, runner.Update(i.Ref)
}
