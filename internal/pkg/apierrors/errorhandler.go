package apierrors

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/MelnikovNA/noolingoproto/codegen/go/apierrors"
	"github.com/MelnikovNA/noolingoproto/codegen/go/common"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/sirupsen/logrus"
)

//type CustomError struct {
//	ErrName string `json:"err_name"`
//	// ErrSrc  string `json:"err_src"`
//	Error string `json:"error"`
//}

func ErrorHandler(_ context.Context, _ *runtime.ServeMux, _ runtime.Marshaler,
	w http.ResponseWriter, req *http.Request, inErr error) {

	apiError := apierrors.FromError(inErr)
	ce := &common.Response{
		Result: false,
		Error: &common.Error{
			Error:     apiError.Name(),
			ErrorText: apiError.Message(),
		},
	}

	js, err := json.Marshal(ce)
	if err != nil {
		logrus.Error(err)
	}

	logrus.Errorf("grpc error: %v code:%v (%v) src: %v", req.URL, apiError.Code(), apiError.Name(), inErr.Error())

	w.WriteHeader(apiError.Code())
	w.Header().Set("Content-Type", "application/json")

	_, err = w.Write(js)
	if err != nil {
		logrus.Error(err)
	}
}
