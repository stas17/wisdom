package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"wisdom/internal/server/config"
	"wisdom/internal/server/logger"
	"wisdom/internal/server/services/wisdom"

	"github.com/pkg/errors"
)

func main() {
	ctx, cancel := context.WithCancel(context.TODO())

	// Setting a hardcoded log level here because we haven't loaded the config yet.
	// The logger gets replaced with a new instance using the log level from the config after loading.
	mainLogger := logger.New(logger.InfoLevel)
	configLoader := config.NewConfigLoader(os.Getenv(config.ENVConfigFile))
	if err := configLoader.Load(); err != nil {
		mainLogger.Err(err).Msg("got an error while loading config")
	}

	appConfig := configLoader.Get()

	mainLogger = logger.New(logger.LevelFromString(appConfig.Logging.Level))
	mainLogger.Info().Msg("Start")
	wsds := wisdom.NewWisdomService(configLoader, mainLogger)
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		s := <-ch
		mainLogger.Printf("got a termination signal: %v", s.String())
		_ = wsds.Stop(ctx)
		cancel()
	}()

	// Run wisdom service
	if err := wsds.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		mainLogger.Err(err).Msg("can't start the wisdom service")
	}
}
