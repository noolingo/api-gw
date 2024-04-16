package parsetoken

import (
	"errors"

	"github.com/golang-jwt/jwt"
)

func ParseToken(tokenString, secret string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Check the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secret), nil
	})
	if err != nil {
		return "", err
	}
	if !token.Valid {
		return "", errors.New("token is invalid") //переписать
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", err
	}
	userID, ok := claims["userID"].(string)
	if !ok {
		return "", errors.New("invalid user id")
	}
	return userID, err
}
