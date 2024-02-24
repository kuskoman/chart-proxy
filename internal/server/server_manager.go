package server

import (
	"context"
	"log/slog"
	"net/http"
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

func (manager *ServerManager) StartOrRestartServer(cfg *config.Config, errorChan chan error) {
	manager.mutex.Lock()
	defer manager.mutex.Unlock()

	if err := manager.shutdownIfRunning(); err != nil {
		errorChan <- err
		return
	}

	manager.startServer(cfg, errorChan)
}

func (manager *ServerManager) startServer(cfg *config.Config, errorChan chan error) {
	manager.server = newAppServer(cfg)

	go func() {
		slog.Info("starting server on", "address", manager.server.Addr)
		if err := manager.server.ListenAndServe(); err != http.ErrServerClosed {
			errorChan <- err
		}
		slog.Debug("server finished work")
	}()
}

func (manager *ServerManager) shutdownIfRunning() error {
	if manager.server != nil {
		slog.Info("shutting down server...")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := manager.server.Shutdown(ctx); err != nil {
			return err
		}
		slog.Info("server shut down")
	}

	return nil
}
