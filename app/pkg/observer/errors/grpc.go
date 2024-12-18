package errs

import (
	"errors"
	"example/config"
	"google.golang.org/grpc/status"
)

func (e *Errs) ToGrpcError() error {
	cfg := config.GetConfig()
	grpcCode := e.code.ToGrpcCode()
	if cfg.Environment != "dev" {
		e.err = errors.New(e.msg)
	}
	grpcErr := status.New(grpcCode, e.err.Error())
	return grpcErr.Err()
}
