package main

import (
	"context"
	"fmt"
	"log"
	"net"

	api "github.com/sadief/grpc-log-middleware-test/api"
	zapper "github.com/sadief/grpc-log-middleware-test/zapper"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

type pingServer struct{}

var (
	zapLogger *zap.Logger
)

const (
	port = "localhost:50051"
)

func (s *pingServer) Ping(ctx context.Context, in *api.PingRequest) (*api.PingResponse, error) {
	l := grpc_zap.Extract(ctx)
	l.Debug(fmt.Sprintf("Ping Received: %v", in.Message))
	return &api.PingResponse{Message: "pong"}, nil
}

func main() {

	zapLogger, _ = zap.NewDevelopment()
	defer zapLogger.Sync()

	zapLogger.Named("grpc zap")

	// Make sure that log statements internal to gRPC library are logged using the zapLogger as well.
	grpc_zap.ReplaceGrpcLogger(zapLogger)

	opts := []grpc_zap.Option{
		grpc_zap.WithLevels(codeToLevel),
	}

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer(grpc_middleware.WithUnaryServerChain(
		grpc_ctxtags.UnaryServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
		grpc_zap.UnaryServerInterceptor(zapLogger, opts...),
		zapper.LoggingUnaryServerInterceptor(),
	))

	api.RegisterPingServiceServer(s, &pingServer{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func codeToLevel(code codes.Code) zapcore.Level {
	switch code {
	default:
		fallthrough
	case codes.OK:
		return zap.DebugLevel
	}
}
