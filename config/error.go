package config

import (
	"fmt"

	"github.com/layer5io/gokit/errors"
)

var (
	ErrEmptyConfig = errors.New(errors.ErrEmptyConfig, "Config not initialized")
)

func ErrViper(err error) error {
	return errors.New(errors.ErrViper, fmt.Sprintf("Viper initialization failed with error: %s", err.Error()))
}
