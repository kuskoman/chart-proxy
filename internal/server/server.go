package server

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/kuskoman/chart-proxy/pkg/config"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func newAppServer(config *config.Config) *http.Server {
	slog.Info("creating new server", "host", config.Server.Host, "port", config.Server.Port)

	mux := http.NewServeMux()

	mux.Handle("/metrics", promhttp.Handler())

	listenUrl := fmt.Sprintf("%s:%d", config.Server.Host, config.Server.Port)
	server := &http.Server{
		Addr:    listenUrl,
		Handler: mux,
	}

	return server
}
