package main

import (

	"os"

	"github.com/moq77111113/chmoly-santas/pkg/cmd"

)

func main() {
	root := cmd.NewChmolyCmd()

	root.AddCommand(cmd.NewServeCmd())

	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}