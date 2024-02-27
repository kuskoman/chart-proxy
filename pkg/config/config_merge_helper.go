package config

import (
	"dario.cat/mergo"
)

func mergeConfigs(configs ...*Config) (*Config, error) {
	mergedAliases, mergedMappings := []RepositoryAlias{}, []Mapping{}

	extractServerConfig := func(c *Config) ServerConfig { return c.Server }
	extractLoggingConfig := func(c *Config) LoggingConfig { return c.Logging }

	extractedServerConfig := extractConfigs(configs, extractServerConfig)
	extractedLoggingConfig := extractConfigs(configs, extractLoggingConfig)

	mergedServerConfig, err := merge(extractedServerConfig...)
	if err != nil {
		return nil, err
	}

	mergedLoggingConfig, err := merge(extractedLoggingConfig...)
	if err != nil {
		return nil, err
	}

	mergeSlices(&mergedAliases, extractSliceConfigs(configs, func(c *Config) []RepositoryAlias { return c.Aliases })...)
	mergeSlices(&mergedMappings, extractSliceConfigs(configs, func(c *Config) []Mapping { return c.Mappings })...)

	return &Config{
		Server:   *mergedServerConfig,
		Logging:  *mergedLoggingConfig,
		Aliases:  mergedAliases,
		Mappings: mergedMappings,
	}, nil
}

func merge[T any](configs ...T) (*T, error) {
	var mergedConfig T

	for _, config := range configs {
		if err := mergo.Merge(mergedConfig, config, mergo.WithOverride); err != nil {
			return nil, err
		}
	}
	return &mergedConfig, nil
}

func mergeSlices[T comparable](mergedSlice *[]T, slices ...[]T) {
	uniqueMap := make(map[T]struct{})

	for _, slice := range slices {
		for _, item := range slice {
			uniqueMap[item] = struct{}{}
		}
	}

	for item := range uniqueMap {
		*mergedSlice = append(*mergedSlice, item)
	}
}

func extractConfigs[T any](configs []*Config, extractor func(*Config) T) []T {
	result := make([]T, 0, len(configs))
	for _, config := range configs {
		result = append(result, extractor(config))
	}
	return result
}

func extractSliceConfigs[T any](configs []*Config, extractor func(*Config) []T) [][]T {
	result := make([][]T, 0, len(configs))
	for _, config := range configs {
		result = append(result, extractor(config))
	}
	return result
}
