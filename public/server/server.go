package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	v1 "github.com/fujin-io/fujin-grpc-gateway/internal/v1"
	"github.com/fujin-io/fujin-grpc-gateway/public/config"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/tmc/grpc-websocket-proxy/wsproxy"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Server struct {
	conf config.Config
	l    *slog.Logger
}

func NewServer(conf config.Config, l *slog.Logger) *Server {
	return &Server{conf: conf, l: l}
}

func (s *Server) ListenAndServe(ctx context.Context) error {
	mux := runtime.NewServeMux()
	err := v1.RegisterFujinServiceHandlerFromEndpoint(ctx, mux, s.conf.GRPC.Addr,
		[]grpc.DialOption{
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		},
	)
	if err != nil {
		return fmt.Errorf("register fujin service handler from endpoint: %w", err)
	}

	srv := http.Server{
		Addr:    s.conf.Addr,
		Handler: wsproxy.WebsocketProxy(mux),
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			if err != http.ErrServerClosed {
				s.l.Error("listen and serve", "error", err)
			}
		}
	}()
	s.l.Info("fujin grpc gateway server started")
	<-ctx.Done()
	srv.Shutdown(ctx)
	s.l.Info("fujin grpc gateway server stopped")
	return nil
}
