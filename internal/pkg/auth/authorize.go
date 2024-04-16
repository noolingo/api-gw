package auth

import (
	"context"

	"github.com/MelnikovNA/noolingo-api-gw/internal/pkg/parsetoken"
)

const (
	AuthorizationHeaderKey = "Authorization" //временная затычка
)

func DefaultAuthorizeFunc(secret string) AuthorizeFunc {
	return func(ctx context.Context, tokenString string) (Authorization, error) {
		userID, err := parsetoken.ParseToken(tokenString, secret)
		return Authorization{UserID: userID}, err
	}
}
