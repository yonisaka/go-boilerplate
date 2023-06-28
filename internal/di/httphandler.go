package di

import "github.com/yonisaka/go-boilerplate/internal/adapters/httphandler"

// GetHealthHandler is a function to get health handler
func GetHealthHandler() httphandler.Handler {
	return httphandler.NewHealthHandler()
}
