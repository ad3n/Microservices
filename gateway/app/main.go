// Krakend-ce sets up a complete KrakenD API Gateway ready to serve

package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	krakend "github.com/devopsfaith/krakend-ce"
	cmd "github.com/devopsfaith/krakend-cobra"
	viper "github.com/devopsfaith/krakend-viper"
	"github.com/devopsfaith/krakend/config"
)

func main() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		select {
		case sig := <-sigs:
			log.Println("Signal intercepted:", sig)
			cancel()
		case <-ctx.Done():
		}
	}()

	krakend.RegisterEncoders()

	var cfg config.Parser
	cfg = viper.New()
	cmd.Execute(cfg, krakend.NewExecutor(ctx))
}
