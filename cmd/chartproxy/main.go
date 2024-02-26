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

	configManager, err := config.NewConfigManager(*configFile, fatalErrorChan)
	if err != nil {
		slog.Error("error creating configuration manager", "error", err)
		panic(err)
	}

	configManager.RegisterReloadHook(logging.SetupLogging)

	serverManager := server.NewServerManager()

	configManager.RegisterReloadHook(serverManager.StartOrRestartServer)

	if err := configManager.LoadConfig(); err != nil {
		slog.Error("error loading configuration", "error", err)
		panic(err)
	}

	if *watch {
		go func() {
			err = configManager.WatchConfig()
			if err != nil {
				slog.Error("error watching configuration", "error", err)
				panic(err)
			}
		}()
	}

	for range fatalErrorChan {
		err := <-fatalErrorChan
		slog.Error("fatal error", "error", err)
		panic(err)
	}
}
