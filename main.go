package main

import (
	"github.com/yonisaka/go-boilerplate/cmd"
	"github.com/yonisaka/go-boilerplate/internal/di"
	"log"
	"os"
)

func main() {
	cfg := di.GetConfig()
	command := cmd.NewCommand(
		cmd.WithConfig(cfg),
	)

	app := cmd.NewCLI()
	app.Commands = command.Build()

	err := app.Run(os.Args)
	if err != nil {
		log.Fatalf("Unable to run CLI command, err: %v", err)
	}
}
