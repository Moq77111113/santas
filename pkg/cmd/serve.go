package cmd

import (
	"log"

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

	return nil
}
