package middleware

import (
	"context"
	"errors"
	errs "example/pkg/observer/errors"
	"example/pkg/observer/logger"
	trace "example/pkg/observer/tracing"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const (
	traceIdHeader = "X-Trace-Id"
)

func Start(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	ctx, span := trace.StartName(ctx, info.FullMethod)
	defer span.End()

	err := grpc.SendHeader(ctx, metadata.Pairs(traceIdHeader, span.GetTraceId()))
	if err != nil {
		logger.Log.Infop("failed to send trace id to header")
	}

	res, err := handler(ctx, req)
	if err != nil {
		e, ok := errs.AsErrs(err)
		if !ok {
			errors.As(errs.New().Msg("unwrapped error").Log().Span(span).Wrap(err), &e)
		}
		return res, e.ToGrpcError()
	}

	return res, nil
}
