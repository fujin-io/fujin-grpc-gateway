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
	"google.golang.org/grpc/credentials"
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

	grpcOptions := []grpc.DialOption{}
	if s.conf.TLS.Enabled {
		grpcOptions = append(grpcOptions, grpc.WithTransportCredentials(credentials.NewTLS(s.conf.TLS.Config)))
	} else {
		grpcOptions = append(grpcOptions, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	err := v1.RegisterFujinServiceHandlerFromEndpoint(ctx, mux, s.conf.GRPC.Addr,
		grpcOptions,
	)
	if err != nil {
		return fmt.Errorf("register fujin service handler from endpoint: %w", err)
	}

	handler := wsproxy.WebsocketProxy(mux, wsproxy.WithRequestMutator(func(incoming *http.Request, outgoing *http.Request) *http.Request {
		if outgoing.Method == http.MethodGet {
			outgoing.Method = http.MethodPost
		}
		return outgoing
	}))

	srv := http.Server{
		Addr:    s.conf.Addr,
		Handler: handler,
	}

	go func() {
		var err error
		if s.conf.TLS.Enabled {
			if s.conf.TLS.Config != nil {
				srv.TLSConfig = s.conf.TLS.Config
			}
			certFile := s.conf.TLS.ServerCertPEMPath
			keyFile := s.conf.TLS.ServerKeyPEMPath
			err = srv.ListenAndServeTLS(certFile, keyFile)
		} else {
			err = srv.ListenAndServe()
		}
		if err != nil && err != http.ErrServerClosed {
			s.l.Error("listen and serve", "error", err)
		}
	}()
	s.l.Info("fujin grpc gateway server started", "addr", s.conf.Addr)
	<-ctx.Done()
	srv.Shutdown(ctx)
	s.l.Info("fujin grpc gateway server stopped")
	return nil
}
