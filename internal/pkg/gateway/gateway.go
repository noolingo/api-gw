package gateway

import (
	"context"
	"net/http"
	"strings"

	"github.com/MelnikovNA/noolingo-api-gw/internal/pkg/apierrors"
	"github.com/MelnikovNA/noolingo-api-gw/internal/pkg/auth"
	"github.com/MelnikovNA/noolingoproto/codegen/go/noolingo"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

type Configs struct {
	Host         string
	HttpPort     string
	GrpcClients  GrpcClients
	Cors         bool
	AccessMap    map[string]string
	RolesAccess  map[string]int
	AccessPrefix []string
	Secret       string
}

type GrpcClients struct {
	UserService      string
	CardService      string
	DeckService      string
	StatisticService string
}

type Gateway struct {
	handler http.Handler
	host    string
	port    string
}

func NewGateway(ctx context.Context, cfg *Configs, lg *logrus.Logger) (*Gateway,
	error) {

	mux := runtime.NewServeMux(
		runtime.WithErrorHandler(apierrors.ErrorHandler),
		runtime.WithUnescapingMode(runtime.UnescapingModeAllExceptReserved),
	)

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
			md, ok := metadata.FromIncomingContext(ctx)
			if !ok {
				return invoker(ctx, method, req, reply, cc, opts...)
			}
			for key, values := range md {
				for _, value := range values {
					ctx = metadata.AppendToOutgoingContext(ctx, key, value)
				}
			}
			return invoker(ctx, method, req, reply, cc, opts...)
		}),
	}

	err := noolingo.RegisterCardsHandlerFromEndpoint(
		ctx,
		mux,
		cfg.GrpcClients.CardService,
		opts,
	)
	if err != nil {
		return nil, err
	}

	err = noolingo.RegisterUserHandlerFromEndpoint(
		ctx,
		mux,
		cfg.GrpcClients.UserService,
		opts,
	)

	if err != nil {
		return nil, err
	}
	err = noolingo.RegisterDeckHandlerFromEndpoint(
		ctx,
		mux,
		cfg.GrpcClients.DeckService,
		opts,
	)

	if err != nil {
		return nil, err
	}

	err = noolingo.RegisterStatisticHandlerFromEndpoint(
		ctx,
		mux,
		cfg.GrpcClients.StatisticService,
		opts,
	)

	if err != nil {
		return nil, err
	}

	var handler http.Handler
	handler = mux

	if cfg.Cors {
		handler = cors.AllowAll().Handler(mux)
	}

	handler = auth.NewAuthorizedHandler(
		cfg.AccessMap,
		cfg.RolesAccess,
		cfg.AccessPrefix,
		handler,
		auth.DefaultAuthorizeFunc(cfg.Secret),
		auth.DefaultAnnotateContextFunc(),
		lg,
	)

	return &Gateway{handler: handler, host: cfg.Host, port: cfg.HttpPort}, nil
}

func (g *Gateway) Serve() error {
	err := http.ListenAndServe(
		strings.Join([]string{g.host, g.port}, ":"),
		g.handler,
	)

	return err
}
