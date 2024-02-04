package auth

import (
	"context"
	"fmt"
	"math"
	"net/http"
	"strings"

	"github.com/MelnikovNA/noolingo-api-gw/internal/pkg/apierrors"
	"github.com/sirupsen/logrus"
)

const (
	defaultRole = "admin"
	anyRole     = "any"
)

type AuthorizeFunc func(ctx context.Context, token string) (Authorization, error)
type AnnotateContextFunc func(ctx context.Context, request *http.Request, token string, authorization Authorization) context.Context

type Authorization struct {
	UserID string
	Role   string
}

type AuthorizedHandler struct {
	logger              *logrus.Logger
	accessMap           map[string]string
	rolesAccess         map[string]int
	accessPrefix        []string
	handler             http.Handler
	authorizeFunc       AuthorizeFunc
	annotateContextFunc AnnotateContextFunc
}

func NewAuthorizedHandler(accessMap map[string]string, rolesAccess map[string]int, accessPrefix []string, handler http.Handler,
	authorizeFunc AuthorizeFunc, annotateContextFunc AnnotateContextFunc, logger *logrus.Logger) *AuthorizedHandler {
	return &AuthorizedHandler{
		accessMap:           accessMap,
		rolesAccess:         rolesAccess,
		handler:             handler,
		authorizeFunc:       authorizeFunc,
		annotateContextFunc: annotateContextFunc,
		accessPrefix:        accessPrefix,
		logger:              logger,
	}
}

func (h *AuthorizedHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {

	if request.Method == "OPTIONS" {
		h.handler.ServeHTTP(writer, request)
		return
	}

	reqUri := request.URL.Path

	for _, pfx := range h.accessPrefix {
		if strings.HasPrefix(reqUri, pfx) {
			reqUri = pfx
			break
		}
	}

	role, ok := h.accessMap[reqUri]
	if !ok {
		withMethod := fmt.Sprintf("%v@%v", request.Method, reqUri)
		role, ok = h.accessMap[withMethod]
		if !ok {
			role, ok = h.accessMap[defaultRole]
		}
		if !ok {
			role = defaultRole
			ok = true
		}
	}

	if !ok || role == anyRole {
		h.handler.ServeHTTP(writer, request)
		return
	}

	token := request.Header.Get(AuthorizationHeaderKey)
	if token == "" {
		h.logger.Warning("empty token")
		writer.WriteHeader(http.StatusUnauthorized)
		return
	}

	token = strings.TrimPrefix(token, "Bearer ")
	auth, err := h.authorizeFunc(request.Context(), token)
	if err != nil {
		h.logger.WithError(err).Warning("authorization error")
		apierrors.ErrorHandler(request.Context(), nil, nil, writer, request, err)
		return
	}
	uriAccessRequirement, ok := h.rolesAccess[role]
	if !ok {
		uriAccessRequirement = math.MaxInt
	}
	userAccessLevel, ok := h.rolesAccess[auth.Role]
	if !ok {
		userAccessLevel = math.MaxInt
	}

	if userAccessLevel < uriAccessRequirement {
		h.logger.Warning("wrong access level")
		writer.WriteHeader(http.StatusForbidden)
		return
	}

	ctx := h.annotateContextFunc(request.Context(), request, token, auth)

	h.handler.ServeHTTP(writer, request.WithContext(ctx))
}
