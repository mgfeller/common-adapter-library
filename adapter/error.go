package adapter

import (
	"fmt"

	"github.com/layer5io/gokit/errors"
)

var (
	ErrOpInvalid = errors.New(errors.ErrOpInvalid, "Invalid operation")
)

func ErrInstallMesh(err error) error {
	return errors.New(errors.ErrInstallMesh, fmt.Sprintf("Error installing mesh: %s", err.Error()))
}

func ErrMeshConfig(err error) error {
	return errors.New(errors.ErrMeshConfig, fmt.Sprintf("Error configuration mesh: %s", err.Error()))
}

func ErrPortForward(err error) error {
	return errors.New(errors.ErrPortForward, fmt.Sprintf("Error portforwarding mesh gui: %s", err.Error()))
}

func ErrClientConfig(err error) error {
	return errors.New(errors.ErrClientConfig, fmt.Sprintf("Error setting client Config: %s", err.Error()))
}

func ErrClientSet(err error) error {
	return errors.New(errors.ErrClientSet, fmt.Sprintf("Error setting clientset: %s", err.Error()))
}

func ErrStreamEvent(err error) error {
	return errors.New(errors.ErrStreamEvent, fmt.Sprintf("Error streaming event: %s", err.Error()))
}
