package main

import (
	"flag"
	"log/slog"

	"github.com/kuskoman/chart-proxy/internal/logging"
	"github.com/kuskoman/chart-proxy/internal/server"
	"github.com/kuskoman/chart-proxy/pkg/config"
)

func main() {
	watch := flag.Bool("watch", false, "Watch for changes in the configuration file")
	configFile := flag.String("config", "config.hcl", "Path to the configuration file")
	flag.Parse()

	fatalErrorChan := make(chan error)

	configManager := config.NewConfigManager(*configFile, *watch, fatalErrorChan)
	configManager.RegisterReloadHook(logging.SetupLogging)

	serverManager := server.NewServerManager()

	configManager.RegisterReloadHook(serverManager.StartOrRestartServer)

	err := configManager.LoadConfig()
	if err != nil {
		slog.Error("error loading configuration", "error", err)
		panic(err)
	}

	for range fatalErrorChan {
		err := <-fatalErrorChan
		slog.Error("fatal error", "error", err)
		panic(err)
	}
}
