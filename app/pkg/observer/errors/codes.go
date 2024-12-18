package errs

import (
	grpcCodes "google.golang.org/grpc/codes"
	"net/http"
)

type Code uint8

const (
	OK Code = iota
	Canceled
	//Unknown
	InvalidArgument
	DeadlineExceeded
	NotFound
	AlreadyExists
	PermissionDenied
	ResourceExhausted
	FailedPrecondition
	Aborted
	OutOfRange
	Unimplemented
	Internal
	Unavailable
	DataLoss
	Unauthenticated
)

func (c *Code) ToGrpcCode() grpcCodes.Code {
	var res grpcCodes.Code

	switch *c {
	case OK:
		res = grpcCodes.OK
	case Canceled:
		res = grpcCodes.Canceled
	case InvalidArgument:
		res = grpcCodes.InvalidArgument
	case DeadlineExceeded:
		res = grpcCodes.DeadlineExceeded
	case NotFound:
		res = grpcCodes.NotFound
	case AlreadyExists:
		res = grpcCodes.AlreadyExists
	case PermissionDenied:
		res = grpcCodes.PermissionDenied
	case ResourceExhausted:
		res = grpcCodes.ResourceExhausted
	case FailedPrecondition:
		res = grpcCodes.FailedPrecondition
	case Aborted:
		res = grpcCodes.Aborted
	case OutOfRange:
		res = grpcCodes.OutOfRange
	case Unimplemented:
		res = grpcCodes.Unimplemented
	case Internal:
		res = grpcCodes.Internal
	case Unavailable:
		res = grpcCodes.Unavailable
	case DataLoss:
		res = grpcCodes.DataLoss
	case Unauthenticated:
		res = grpcCodes.Unauthenticated
	default:
		res = grpcCodes.Unknown
	}

	return res
}

func (c *Code) ToHttpCode() int {
	var res int

	switch *c {
	case OK:
		res = http.StatusOK
	case Canceled:
		res = 499
	case InvalidArgument:
		res = http.StatusBadRequest
	case DeadlineExceeded:
		res = http.StatusGatewayTimeout
	case NotFound:
		res = http.StatusNotFound
	case AlreadyExists:
		res = http.StatusConflict
	case PermissionDenied:
		res = http.StatusForbidden
	case ResourceExhausted:
		res = http.StatusTooManyRequests
	case FailedPrecondition:
		res = http.StatusBadRequest
	case Aborted:
		res = http.StatusConflict
	case OutOfRange:
		res = http.StatusBadRequest
	case Unimplemented:
		res = http.StatusNotImplemented
	case Internal:
		res = http.StatusInternalServerError
	case Unavailable:
		res = http.StatusServiceUnavailable
	case DataLoss:
		res = http.StatusInternalServerError
	case Unauthenticated:
		res = http.StatusUnauthorized
	default:
		res = http.StatusInternalServerError
	}

	return res
}
