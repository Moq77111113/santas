package cmd

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/moq77111113/chmoly-santas/pkg/handlers"
	"github.com/moq77111113/chmoly-santas/pkg/services"
	"github.com/spf13/cobra"
)

// NewServeCmd starts the server
func NewServeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Start the  server",
		RunE:  serve,
	}

	return cmd
}

func serve(cmd *cobra.Command, args []string) error {

	c := services.NewContainer()

	defer func() {
		if err := c.Shutdown(); err != nil {
			log.Fatalf("shutdown error: %v", err)
		}
	}()

	if err := handlers.Bootstrap(c); err != nil {
		log.Fatalf("failed to bootstrap the server: %v", err)
	}

	srv := http.Server{
		Addr:         fmt.Sprintf("%s:%d", c.Config.Http.Hostname, c.Config.Http.Port),
		Handler:      c.Web,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, syscall.SIGTERM)

	go func() {

		if err := c.Web.StartServer(&srv); !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("shutting down the server: %v", err)
		}
	}()

	<-exit

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()
		if err := c.Web.Shutdown(ctx); err != nil {
			log.Fatalf("shutdown error: %v", err)
		}
	}()

	wg.Wait()
	return nil
}
