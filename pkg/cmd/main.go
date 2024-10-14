package cmd

import (
	"github.com/spf13/cobra"
)

// Execute is the entry point for the aura command
func NewChmolyCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "chmoly",
		Short: "Secret santa generator",
		Long:  "Chmolys santas is a secret santa generator",
	}
}
