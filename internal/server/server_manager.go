package server

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/kuskoman/chart-proxy/pkg/config"
)

type ServerManager struct {
	server *http.Server
	mutex  *sync.RWMutex
}

func NewServerManager() *ServerManager {
	return &ServerManager{
		mutex: &sync.RWMutex{},
	}
}

func (manager *ServerManager) StartOrRestartServer(cfg *config.Config) error {
	manager.mutex.Lock()
	defer manager.mutex.Unlock()

	if manager.server != nil {
		slog.Info("shutting down server...")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := manager.server.Shutdown(ctx); err != nil {
			return err
		}
	}

	manager.server = newAppServer(cfg)

	go func() {
		slog.Info("starting server on", "address", manager.server.Addr)
		if err := manager.server.ListenAndServe(); err != http.ErrServerClosed {
			slog.Error("error starting server", "error", err)
			os.Exit(1)
		}
	}()

	return nil
}
