package errs

import (
	"errors"
	"example/config"
	"example/pkg/observer/logger"
	"google.golang.org/grpc/status"
)

// ToGrpcError returns error state depending on app enironment
func (e *Errs) ToGrpcError() error {
	cfg := config.GetConfig()
	grpcCode := e.code.ToGrpcCode()

	var res string
	if cfg.Environment == "dev" {
		info, err := e.ToJson()
		if err != nil {
			logger.Log.Errorf("failed to marshal error to json", "error", err)
			res = e.err.Error()
		}
		res = string(info)
		grpcErr := status.New(grpcCode, res)
		return grpcErr.Err()
	}

	e.err = errors.New(e.msg)
	res = e.err.Error()
	grpcErr := status.New(grpcCode, res)
	return grpcErr.Err()
}
