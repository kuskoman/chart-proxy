package main

import (
	"flag"
	"fmt"

	"github.com/kuskoman/chart-proxy/internal/logging"
	"github.com/kuskoman/chart-proxy/pkg/config"
)

func main() {
	watch := flag.Bool("watch", false, "Watch for changes in the configuration file")
	configFile := flag.String("config", "config.hcl", "Path to the configuration file")
	flag.Parse()

	configManager := config.NewConfigManager(*configFile, *watch)
	configManager.RegisterReloadHook(logging.SetupSlog)

	err := configManager.LoadConfig()

	if err != nil {
		panic(err)
	}

	config := configManager.GetConfig()
	fmt.Printf("%+v\n", config)
}
