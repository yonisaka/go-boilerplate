package cmd

import "github.com/urfave/cli/v2"

// Build is a method
func (cmd *Command) Build() []*cli.Command {
	cmd.registerCLI(cmd.httpStart())

	return cmd.CLI
}
