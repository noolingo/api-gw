package apierrors

import (
	"errors"
	gRPCErrCodes "google.golang.org/grpc/codes"
	"net/http"
)

const TypeRPCErr = "rpc error"

var GRPCErrorsByCode = map[gRPCErrCodes.Code]string{
	3: "InvalidArgument",
}

var GRPCInvalidArgument = errors.New("invalid request argument provided")

var GRPCErrCodes = map[error]int{
	GRPCInvalidArgument: http.StatusBadRequest,
}
