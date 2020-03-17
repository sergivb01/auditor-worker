package worker

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"

	"github.com/thoas/bokchoy"
	"github.com/thoas/bokchoy/middleware"
	"go.uber.org/zap"
)

type Worker struct {
	buildsDir string
	cfg       Config

	engine *bokchoy.Bokchoy

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

	pwd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("could not get current wd: %w", err)
	}

	engine, err := bokchoy.New(context.TODO(), bokchoy.Config{
		Queues: []bokchoy.QueueConfig{{Name: "compiler"}},
		Broker: bokchoy.BrokerConfig{
			Type: "redis",
			Redis: bokchoy.RedisConfig{
				Type: "client",
				Client: bokchoy.RedisClientConfig{
					Addr: c.RedisAddr,
				},
			},
		}}, bokchoy.WithConcurrency(5),
		bokchoy.WithMaxRetries(1))

	if err != nil {
		return nil, fmt.Errorf("error creating bokchoy worker engine: %w", err)
	}

	engine.Use(middleware.Recoverer)
	engine.Use(middleware.RequestID)
	engine.Use(middleware.DefaultLogger)

	engine.Queue("compiler.send").HandleFunc(handleQueueJob)

	return &Worker{
		buildsDir: filepath.Join(pwd, "builds"),
		cfg:       *c,
		engine:    engine,
		log:       logger,
	}, nil
}

func (w *Worker) Start() error {
	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught
	signal.Notify(c, os.Interrupt)

	go func(c chan<- os.Signal) {
		w.log.Info("worker started", zap.Strings("queueNames", w.engine.QueueNames()))
		if err := w.engine.Run(context.TODO()); err != nil {
			w.log.Fatal("failed to start grpc server", zap.Error(err))
			c <- nil
		}
	}(c)

	// Block until we receive our signal
	<-c

	w.log.Debug("gracefully closing grpc server")
	w.engine.Stop(context.TODO())
	w.log.Debug("closed grpc server")

	return nil
}
