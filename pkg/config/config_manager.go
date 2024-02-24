package config

import (
	"sync"

	"github.com/hashicorp/hcl/v2/hclsimple"
)

type ConfigManager struct {
	location        string
	config          *Config
	mutex           sync.RWMutex
	reloadHooks     []ReloadHook
	reloadErrorChan chan error
}

type ReloadHook func(config *Config, errorChan chan error)

func NewConfigManager(location string, watch bool, reloadErrorChan chan error) *ConfigManager {
	return &ConfigManager{
		location:        location,
		mutex:           sync.RWMutex{},
		reloadHooks:     []ReloadHook{},
		reloadErrorChan: reloadErrorChan,
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
	defer cm.mutex.Unlock()

	cm.config = mergedConfig
	cm.handleReloadHooks()

	return nil
}

// RegisterReloadHook registers a function that will be called after the configuration is reloaded.
func (cm *ConfigManager) RegisterReloadHook(hook ReloadHook) {
	cm.reloadHooks = append(cm.reloadHooks, hook)
}

func (cm *ConfigManager) handleReloadHooks() {
	for _, hook := range cm.reloadHooks {
		hook(cm.config, cm.reloadErrorChan)
	}
}

func (cm *ConfigManager) readConfig() (*Config, error) {
	config := &Config{}
	err := hclsimple.DecodeFile(cm.location, nil, config)

	if err != nil {
		return nil, err
	}

	return config, nil
}
