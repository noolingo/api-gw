package auth

import (
	"context"
	"encoding/base64"
	"encoding/json"

	"github.com/MelnikovNA/noolingoproto/codegen/go/apierrors"
	"github.com/MelnikovNA/noolingoproto/codegen/go/noolingo"
	"github.com/sirupsen/logrus"
)

type user struct {
	UserID   string `json:"user_id"`
	UserName string `json:"user_name"`
}

const (
	AuthorizationHeaderKey = "Authorization" //временная затычка
)

func DefaultAuthorizeFunc(client noolingo.UserClient, lg *logrus.Logger) AuthorizeFunc {
	return func(ctx context.Context, token string) (Authorization, error) {
		var auth Authorization
		var userData user
		data, err := base64.StdEncoding.DecodeString(token)
		if err != nil {
			return auth, apierrors.ErrInvalidPayload
		}
		if err = json.Unmarshal(data, &userData); err != nil {
			return auth, apierrors.ErrInvalidPayload
		}

		auth.UserID = userData.UserID
		auth.Role = RoleUser

		return auth, nil
	}
}
