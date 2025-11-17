package service

import (
	"context"
	"fmt"
	"io"
	"log"
	"log/slog"
	"os"
	"strings"

	"github.com/fujin-io/fujin-grpc-gateway/public/config"
	"github.com/fujin-io/fujin-grpc-gateway/public/server"
	"gopkg.in/yaml.v3"
)

var (
	Version string
	conf    config.Config
)

func RunCLI(ctx context.Context) {
	log.Printf("version: %s", Version)

	if len(os.Args) > 2 {
		log.Fatal("invalid args")
	}
	confPath := ""
	if len(os.Args) == 2 {
		confPath = os.Args[1]
	}

	if err := loadConfig(confPath, &conf); err != nil {
		log.Fatal(err)
	}

	logLevel := os.Getenv("FUJIN_GRPC_GATEWAY_LOG_LEVEL")
	logType := os.Getenv("FUJIN_GRPC_GATEWAY_LOG_TYPE")
	logger := configureLogger(logLevel, logType)

	s := server.NewServer(conf, logger)

	if err := s.ListenAndServe(ctx); err != nil {
		logger.Error("listen and serve", "err", err)
	}

}

func loadConfig(filePath string, cfg *config.Config) error {
	paths := []string{}

	if filePath == "" {
		paths = append(paths, "./config.yaml", "conf/config.yaml", "config/config.yaml")
	} else {
		paths = append(paths, filePath)
	}

	for _, p := range paths {
		f, err := os.Open(p)
		if err == nil {
			log.Printf("reading config from: %s\n", p)
			data, err := io.ReadAll(f)
			f.Close()
			if err != nil {
				return fmt.Errorf("read config: %w", err)
			}

			if err := yaml.Unmarshal(data, &cfg); err != nil {
				return fmt.Errorf("unmarshal config: %w", err)
			}

			if err := cfg.GRPC.TLS.Parse(); err != nil {
				return fmt.Errorf("parse tls config: %w", err)
			}

			return nil
		}
	}

	return fmt.Errorf("failed to find config in: %v", paths)
}

func configureLogger(logLevel, logType string) *slog.Logger {
	var parsedLogLevel slog.Level
	switch strings.ToUpper(logLevel) {
	case "DEBUG":
		parsedLogLevel = slog.LevelDebug
	case "WARN":
		parsedLogLevel = slog.LevelWarn
	case "ERROR":
		parsedLogLevel = slog.LevelError
	default:
		parsedLogLevel = slog.LevelInfo
	}

	var handler slog.Handler
	switch strings.ToLower(logType) {
	case "json":
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: parsedLogLevel,
		})
	default:
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: parsedLogLevel,
		})
	}

	return slog.New(handler)
}
