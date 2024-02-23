package main

import (
	"fmt"

	"github.com/kuskoman/chart-proxy/pkg/config"
)

func main() {
	configManager := config.NewConfigManager("config.hcl")

	err := configManager.LoadConfig()

	if err != nil {
		panic(err)
	}

	config := configManager.GetConfig()
	fmt.Printf("%+v\n", config)
}
