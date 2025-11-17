package main

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/fujin-io/fujin-grpc-gateway/public/service"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer cancel()

	service.RunCLI(ctx)
}
