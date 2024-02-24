package config

import (
	"sync"

	"github.com/hashicorp/hcl/v2/hclsimple"
)

type ConfigManager struct {
	location string
	config   *Config
	mutex    sync.RWMutex
}

func NewConfigManager(location string) *ConfigManager {
	return &ConfigManager{
		location: location,
		mutex:    sync.RWMutex{},
	}
}

func (cm *ConfigManager) GetConfig() *Config {
	return cm.config
}

// LoadConfig reads the HCL configuration file and unmarshals it into a Config struct.
// Can be called multiple times to reload the configuration.
func (cm *ConfigManager) LoadConfig() error {
	userDefinedConfig, err := cm.readConfig()
	if err != nil {
		return err
	}

	defaultConfig := getDefaultConfig()
	mergedConfig, err := mergeConfigs(defaultConfig, userDefinedConfig)
	if err != nil {
		return err
	}

	cm.mutex.Lock()

	cm.config = mergedConfig
	cm.handleReloadHooks()

	cm.mutex.Unlock()

	return nil
}

func (cm *ConfigManager) handleReloadHooks() {
	config := cm.GetConfig()

	setupSlog(config.Logging)
}

func (cm *ConfigManager) readConfig() (*Config, error) {
	config := &Config{}
	err := hclsimple.DecodeFile(cm.location, nil, config)

	if err != nil {
		return nil, err
	}

	return config, nil
}
