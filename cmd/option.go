package cmd

import (
	"github.com/yonisaka/go-boilerplate/config"
)

// Option is an option type
type Option func(c *Command)

// WithConfig is a function option
func WithConfig(cfg *config.Config) Option {
	return func(c *Command) {
		c.Config = cfg
	}
}
