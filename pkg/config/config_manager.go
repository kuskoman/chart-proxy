package config

import (
	"log/slog"
	"path/filepath"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
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

func NewConfigManager(location string, reloadErrorChan chan error) (*ConfigManager, error) {
	absoluteLocation, err := filepath.Abs(location)
	if err != nil {
		return nil, err
	}

	configManager := &ConfigManager{
		location:        absoluteLocation,
		mutex:           sync.RWMutex{},
		reloadHooks:     []ReloadHook{},
		reloadErrorChan: reloadErrorChan,
	}

	return configManager, nil
}

func (cm *ConfigManager) GetConfig() *Config {
	return cm.config
}

func (cm *ConfigManager) WatchConfig() error {
	slog.Debug("watching configuration file", "location", cm.location)

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	defer func() {
		slog.Debug("closing watcher", "location", cm.location)
		err := watcher.Close()
		slog.Error("error closing watcher", "error", err)
	}()

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					slog.Debug("watcher events channel closed")
					return
				}

				if event.Op&fsnotify.Write == fsnotify.Write {
					slog.Debug("configuration file changed, reloading", "location", cm.location)
					err := cm.LoadConfig()
					if err != nil {
						cm.reloadErrorChan <- err
					}
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}

				cm.reloadErrorChan <- err
			}
		}
	}()

	err = watcher.Add(cm.location)
	if err != nil {
		return err
	}

	<-make(chan struct{})
	return nil
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
	reloadStartTime := time.Now()

	for _, hook := range cm.reloadHooks {
		hook(cm.config, cm.reloadErrorChan)
	}

	reloadDuration := time.Since(reloadStartTime)

	slog.Debug("executed post-load hooks", "duration", reloadDuration)
}

func (cm *ConfigManager) readConfig() (*Config, error) {
	config := &Config{}
	err := hclsimple.DecodeFile(cm.location, nil, config)

	if err != nil {
		return nil, err
	}

	return config, nil
}
