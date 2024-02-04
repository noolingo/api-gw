package auth

import (
	"context"
	"fmt"
	"google.golang.org/grpc/metadata"
	"net/http"
)

func DefaultAnnotateContextFunc() AnnotateContextFunc {
	return func(ctx context.Context, request *http.Request, token string, authorization Authorization) context.Context {
		md := metadata.New(map[string]string{
			"access_token": token,
			"user_id":      authorization.UserID,
			"is_admin":     fmt.Sprintf("%t", authorization.Role == RoleAdmin),
		})
		ctx = metadata.NewIncomingContext(ctx, md)
		return ctx
	}
}
