package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/noolingo/api-gw/internal/domain"
	"github.com/noolingo/api-gw/internal/pkg/gateway"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/sirupsen/logrus"
)

func Run(configPath string) error {
	cfg := new(domain.Config)

	err := cleanenv.ReadConfig(configPath, cfg)
	if err != nil {
		return fmt.Errorf("error creating configs: %w", err)
	}

	parseFlags(cfg)

	log := logrus.New()
	lvl, err := logrus.ParseLevel(cfg.Log.Level["any"])
	if err != nil {
		return fmt.Errorf("error parsing log level: %w", err)
	}
	log.SetLevel(lvl)
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	log.SetReportCaller(true)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	g, err := gateway.NewGateway(ctx, &gateway.Configs{
		Host:         cfg.Listen.Host,
		HttpPort:     cfg.Listen.Ports.Http,
		AccessMap:    cfg.App.AccessMap,
		RolesAccess:  cfg.App.RolesAccess,
		AccessPrefix: cfg.App.AccessPrefix,
		GrpcClients: gateway.GrpcClients{
			UserService:      cfg.Grpc.Clients.UserService,
			CardService:      cfg.Grpc.Clients.CardService,
			DeckService:      cfg.Grpc.Clients.DeckService,
			StatisticService: cfg.Grpc.Clients.StatisticService,
		},
		Secret: cfg.Auth.AccessSecretKey,
		Cors:   cfg.App.Cors,
	}, log)

	if err != nil {
		return fmt.Errorf("error intializing gateway: %w", err)
	}

	errCh := make(chan error)

	go func() {
		errCh <- g.Serve()
	}()

	log.Infof("Http listener has started on port %s", cfg.Listen.Ports.Http)

	select {
	case err = <-errCh:
		log.Error(err)
		log.Info("stopping gracefully...")
	case q := <-quit:
		log.Infof("%s signal received, stopping gracefully...", q.String())
	}

	log.Info("grpc listener has closed")

	return nil
}
