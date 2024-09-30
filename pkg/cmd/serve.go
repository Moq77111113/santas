package cmd

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/moq77111113/chmoly-santas/internal/apis"
	"github.com/moq77111113/chmoly-santas/internal/core"
	httpserver "github.com/moq77111113/chmoly-santas/pkg/http"

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

	app := core.NewApp()

	go app.Notifier.Listen()

	router, err := apis.Init(app)

	if err != nil {
		return err
	}

	server := httpserver.New(router)

	server.Run()

	go func() {
		err := <-server.Notify()
		if err != nil {
			log.Printf("Error: %v", err)
		}
	}()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	<-signals

	if err := server.Shutdown(); err != nil {
		log.Printf("Error: %v", err)
	}
	return nil

}
