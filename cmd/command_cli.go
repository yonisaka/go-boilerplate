package cmd

import (
	"github.com/urfave/cli/v2"
	"github.com/yonisaka/go-boilerplate/internal/consts"
	"github.com/yonisaka/go-boilerplate/internal/server"
	"github.com/yonisaka/go-boilerplate/pkg/logger"
)

// httpStart is a method to start http server
func (cmd *Command) httpStart() *cli.Command {
	return &cli.Command{
		Name:  "http:start",
		Usage: "A command to start http server",
		Action: func(c *cli.Context) error {
			httpServer := server.NewHTTPServer()

			logger.Info(logger.MessageFormat("starting document-service services... %d", cmd.App.Port), logger.EventName(consts.LogEventNameServiceStarting))
			if err := httpServer.Run(); err != nil {
				return err
			}

			return nil
		},
	}
}
