package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	"github.com/husio/goe-demo/pkg/goe"
	"github.com/husio/goe-demo/pkg/x"
	"google.golang.org/grpc"
)

type configuration struct {
	ListenAddr string
	Redis      string
}

func main() {
	conf := configuration{
		ListenAddr: env("LISTEN_ADDR", "0:12345"),
		Redis:      env("REDIS", "redis:6379"),
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	if err := run(ctx, conf); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(ctx context.Context, conf configuration) error {
	ls, err := net.Listen("tcp", conf.ListenAddr)
	if err != nil {
		return fmt.Errorf("listen tcp: %w", err)
	}
	defer ls.Close()

	ip, err := x.ExternalIP()
	if err != nil {
		return fmt.Errorf("acquire external IP: %w", err)
	}

	logger := log.New(os.Stdout, fmt.Sprintf("Service %s: ", ip), 0)

	store := goe.NewRedisStore(conf.Redis)

	// For simplicity, no TLS support.

	grpcServer := grpc.NewServer()
	goe.RegisterRandomerServer(grpcServer, goe.NewRandomerServer(logger, store))
	grpcServer.Serve(ls)

	return nil
}

func env(name, fallback string) string {
	if v, ok := os.LookupEnv(name); ok {
		return v
	}
	return fallback
}
