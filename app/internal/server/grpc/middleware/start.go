package middleware

import (
	"context"
	errs "example/pkg/observer/errors"
	"example/pkg/observer/logger"
	trace "example/pkg/observer/tracing"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func Start(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	ctx, span := trace.Start(ctx, info.FullMethod)
	defer span.End()

	err := grpc.SendHeader(ctx, metadata.Pairs("traceId", span.GetTraceId()))
	if err != nil {
		logger.Log.Info("failed to send trace id to header")
	}

	res, err := handler(ctx, req)
	if err == nil {
		return res, nil
	}

	e, ok := errs.ToErrs(err)
	if ok {
		logger.Build.Err().Pairs("info", e.ToMap()).Err(err)
	}

	return res, e.ToGrpcError()
}
