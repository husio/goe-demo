package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	mrand "math/rand"
	"os"
	"os/signal"
	"time"

	"github.com/husio/goe-demo/pkg/goe"
	"github.com/husio/goe-demo/pkg/x"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type configuration struct {
	ConsumerService string
}

func main() {
	conf := configuration{
		ConsumerService: env("CONSUMER_SERVICE", "consumer:12345"),
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	if err := run(ctx, conf); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(ctx context.Context, conf configuration) error {
	conn, err := grpc.Dial(
		conf.ConsumerService,
		// consumer server does not support authentication or SSL
		grpc.WithInsecure(),
	)
	if err != nil {
		return fmt.Errorf("connect to randmer service: %w", err)
	}
	defer conn.Close()

	ip, err := x.ExternalIP()
	if err != nil {
		return fmt.Errorf("acquire external IP: %w", err)
	}

	logger := log.New(os.Stdout, fmt.Sprintf("Service %s: ", ip), 0)

	cli := goe.NewRandomerClient(conn)

	for {
		sleep := time.Duration(mrand.Intn(5)) * time.Second

		select {
		case <-time.After(sleep):
			req := &goe.RandomRequest{
				CreatedAt: timestamppb.New(now()),
				Id:        generateID(),
				Data:      randomData(),
			}
			if _, err := cli.GenerateRandom(ctx, req); err != nil {
				log.Printf("generate random call: %s", err)
			} else {
				logger.Printf(
					"sending timestamp: %d, id: %x, data: %x",
					req.CreatedAt.AsTime().UnixNano(), req.Id, req.Data)
			}
		case <-ctx.Done():
			return nil
		}
	}
}

func env(name, fallback string) string {
	if v, ok := os.LookupEnv(name); ok {
		return v
	}
	return fallback
}

func now() time.Time {
	return time.Now().UTC()
}

func randomData() []byte {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	return b
}

func generateID() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	return hex.EncodeToString(b)
}
