package zaplogger

import (
	"context"
	"log"

	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"github.com/myesui/uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func LoggingUnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		id := uuid.NewV4()

		ctxzap.AddFields(
			ctx,
			zap.String("requestID", id.String()),
		)

		log.Printf("New context now: %v", ctx)

		return handler(ctx, req)
	}
}
