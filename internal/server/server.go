package server

import (
	"fmt"
	"net/http"

	"github.com/kuskoman/chart-proxy/pkg/config"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func newAppServer(config *config.Config) *http.Server {
	mux := http.NewServeMux()

	mux.Handle("/metrics", promhttp.Handler())

	listenUrl := fmt.Sprintf("%s:%d", config.Server.Host, config.Server.Port)
	server := &http.Server{
		Addr:    listenUrl,
		Handler: mux,
	}

	return server
}
