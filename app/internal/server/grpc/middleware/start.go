package middleware

import (
	"context"
	errs "example/pkg/observer/errors"
	"example/pkg/observer/logger"
	trace "example/pkg/observer/tracing"
	"google.golang.org/grpc"
)

func Start(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	ctx, span := trace.Start(ctx, info.FullMethod)
	defer span.End()

	res, err := handler(ctx, req)
	if err == nil {
		return res, nil
	}
	traceId := trace.GetTraceIdFromCtx(ctx)
	e, ok := errs.ToErrs(err)
	if ok {
		logger.Build.Err().Pairs("traceId", traceId, "info", e.ToMap()).Err(err)
	}

	return res, e.ToGrpcError()
}
