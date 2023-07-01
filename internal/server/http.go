package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/yonisaka/go-boilerplate/internal/consts"
	"github.com/yonisaka/go-boilerplate/internal/di"
	"github.com/yonisaka/go-boilerplate/pkg/logger"
)

type httpServer struct {
	srv *http.Server
}

func NewHTTPServer() Server {
	return &httpServer{}
}

func (h *httpServer) Run() error {
	var err error

	ctx := context.Background()

	cfg := di.GetConfig()

	router := di.NewRouter()

	server := http.Server{
		Addr:         fmt.Sprintf("0.0.0.0:%d", cfg.App.Port),
		Handler:      router.Route(),
		ReadTimeout:  time.Duration(cfg.App.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.App.WriteTimeout) * time.Second,
	}

	go func() {
		err = server.ListenAndServe()
		if err != http.ErrServerClosed {
			logger.Error(logger.MessageFormat("http server got %v", err), logger.EventName(consts.LogEventNameServiceStarting))
		}
	}()

	<-ctx.Done()

	ctxShutDown, cancel := context.WithTimeout(context.Background(), 6*time.Second)
	defer func() {
		cancel()
	}()

	if err = server.Shutdown(ctxShutDown); err != nil {
		logger.Fatal(logger.MessageFormat("server Shutdown Failed:%v", err), logger.EventName(consts.LogEventNameServiceTerminated))
	}

	logger.Info("server exited properly", logger.EventName(consts.LogEventNameServiceTerminated))

	if err == http.ErrServerClosed {
		err = nil
	}

	return err
}

func (h *httpServer) GracefulStop() {
	_ = h.srv.Shutdown(context.Background())
}
