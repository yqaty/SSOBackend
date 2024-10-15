package server

import (
	"github.com/UniqueStudio/UniqueSSOBackend/internal/handler"
	"github.com/UniqueStudio/UniqueSSOBackend/internal/proto/sso"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
)

func InitGRPCServer() *grpc.Server {
	s := grpc.NewServer(
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			grpc_recovery.StreamServerInterceptor(),
			otelgrpc.StreamServerInterceptor(),
		)),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_recovery.UnaryServerInterceptor(),
			otelgrpc.UnaryServerInterceptor(),
		)),
	)

	sso.RegisterSSOServiceServer(s, handler.NewSSOServer())

	return s
}
