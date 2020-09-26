package grpc

import (
	"fmt"

	"github.com/layer5io/gokit/errors"
)

var (
	ErrRequestInvalid = errors.New("603", "Apply Request invalid")
)

func ErrPanic(r interface{}) error {
	return errors.New(errors.ErrPanic, fmt.Sprintf("%v", r))
}

func ErrGrpcListener(err error) error {
	return errors.New(errors.ErrGrpcListener, fmt.Sprintf("Error during grpc listener initialization : %v", err))
}

func ErrGrpcServer(err error) error {
	return errors.New(errors.ErrGrpcServer, fmt.Sprintf("Error during grpc server initialization : %v", err))
}
