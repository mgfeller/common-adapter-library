package grpc

import (
	"net"
	"time"

	"google.golang.org/grpc/reflection"

	"github.com/mgfeller/common-adapter-library/adapter"
	"github.com/mgfeller/common-adapter-library/api/tracing"
	"github.com/mgfeller/common-adapter-library/meshes"

	"fmt"

	middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	otelgrpc "go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc"

	apitrace "go.opentelemetry.io/otel/api/trace"
	"google.golang.org/grpc"
)

// Service object holds all the information about the server parameters.
type Service struct {
	Name      string    `json:"name"`
	Port      string    `json:"port"`
	Version   string    `json:"version"`
	StartedAt time.Time `json:"startedat"`
	TraceURL  string    `json:"traceurl"`
	Handler   adapter.Handler
	Channel   chan interface{}
}

// panicHandler is the handler function to handle panic errors
func panicHandler(r interface{}) error {
	fmt.Println("600 Error")
	return ErrPanic(r)
}

// Start starts grpc server
func Start(s *Service, tr tracing.Handler) error {
	address := fmt.Sprintf(":%s", s.Port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return ErrGrpcListener(err)
	}

	middlewares := middleware.ChainUnaryServer(
		grpc_recovery.UnaryServerInterceptor(
			grpc_recovery.WithRecoveryHandler(panicHandler),
		),
	)
	if tr != nil {
		middlewares = middleware.ChainUnaryServer(
			otelgrpc.UnaryServerInterceptor(tr.Tracer(s.Name).(apitrace.Tracer)),
		)
	}

	server := grpc.NewServer(
		grpc.UnaryInterceptor(middlewares),
	)
	// Reflection is enabled to simplify accessing the gRPC service using gRPCurl, e.g.
	//    grpcurl --plaintext localhost:10002 meshes.MeshService.SupportedOperations
	// If the use of reflection is not desirable, the parameters '-import-path ./meshes/ -proto meshops.proto' have
	//    to be added to each grpcurl request, with the appropriate import path.
	reflection.Register(server)

	//Register Proto
	meshes.RegisterMeshServiceServer(server, s)

	// Start serving requests
	if err = server.Serve(listener); err != nil {
		return ErrGrpcServer(err)
	}
	return nil
}
