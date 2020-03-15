package worker

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"path/filepath"

	"github.com/sergivb01/acmecopy/api"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type Worker struct {
	grpcServer *grpc.Server

	buildsDir string
	creds     credentials.TransportCredentials
	cfg       Config

	log *zap.Logger
}

func NewWorker() (*Worker, error) {
	c, err := loadConfig()
	if err != nil {
		return nil, fmt.Errorf("could not load config: %w", err)
	}

	logConfig := zap.NewDevelopmentConfig()
	if c.Production {
		logConfig = zap.NewProductionConfig()
	}
	logConfig.DisableStacktrace = true
	logConfig.DisableCaller = true

	logger, err := logConfig.Build()
	if err != nil {
		return nil, fmt.Errorf("could not create logger: %w", err)
	}

	creds, err := credentials.NewServerTLSFromFile(c.TLSCert, c.TLSKey)
	if err != nil {
		return nil, fmt.Errorf("error creating TLS credentials: %w", err)
	}

	pwd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("could not get current wd: %w", err)
	}

	return &Worker{
		cfg:       *c,
		buildsDir: filepath.Join(pwd, "builds"),
		log:       logger,
		creds:     creds,
	}, nil
}

func (w *Worker) Listen() error {
	w.grpcServer = grpc.NewServer(grpc.Creds(w.creds))
	api.RegisterCompilerServer(w.grpcServer, w)

	lis, err := net.Listen("tcp", w.cfg.Listen)
	if err != nil {
		return fmt.Errorf("failed to net listen: %w", err)
	}

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught
	signal.Notify(c, os.Interrupt)

	go func(c chan<- os.Signal) {
		w.log.Info("started listening worker", zap.String("address", w.cfg.Listen))
		if err := w.grpcServer.Serve(lis); err != nil {
			w.log.Fatal("failed to start grpc server", zap.Error(err))
			c <- nil
		}
	}(c)

	// Block until we receive our signal
	<-c

	w.log.Debug("gracefully closing grpc server")
	w.grpcServer.GracefulStop()
	w.log.Debug("closed grpc server")

	return nil
}
