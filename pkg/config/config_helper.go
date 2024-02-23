package config

import (
	"log/slog"

	"dario.cat/mergo"
)

func mergeConfigs(configs ...*Config) (*Config, error) {
	serverConfigs := make([]ServerConfig, 0, len(configs))
	loggingConfigs := make([]LoggingConfig, 0, len(configs))
	aliaseConfigs := make([][]RepositoryAlias, 0, len(configs))
	mappingConfigs := make([][]Mapping, 0, len(configs))

	for _, config := range configs {
		serverConfigs = append(serverConfigs, *config.Server)
		loggingConfigs = append(loggingConfigs, *config.Logging)
		aliaseConfigs = append(aliaseConfigs, *config.Aliases)
		mappingConfigs = append(mappingConfigs, *config.Mappings)
	}

	mergedServerConfig, err := mergeServerConfigs(serverConfigs...)
	if err != nil {
		return nil, err
	}
	mergedLoggingConfig, err := mergeLoggingConfigs(loggingConfigs...)
	if err != nil {
		return nil, err
	}
	mergedAliases := mergeRepositoryAliases(aliaseConfigs...)
	mergedMappings := mergeMappings(mappingConfigs...)

	mergedConfig := &Config{
		Server:   mergedServerConfig,
		Logging:  mergedLoggingConfig,
		Aliases:  mergedAliases,
		Mappings: mergedMappings,
	}

	return mergedConfig, nil
}

func mergeServerConfigs(configs ...ServerConfig) (*ServerConfig, error) {
	mergedConfig := ServerConfig{}

	for _, config := range configs {
		err := mergo.Merge(&mergedConfig, config, mergo.WithOverride)
		if err != nil {
			return nil, err
		}
	}

	return &mergedConfig, nil
}

func mergeLoggingConfigs(configs ...LoggingConfig) (*LoggingConfig, error) {
	mergedConfig := LoggingConfig{}

	for _, config := range configs {
		err := mergo.Merge(&mergedConfig, config, mergo.WithOverride)
		if err != nil {
			return nil, err
		}
	}

	return &mergedConfig, nil
}

func mergeRepositoryAliases(aliases ...[]RepositoryAlias) *[]RepositoryAlias {
	repositoryAliasesMap := make(map[string]RepositoryAlias)

	for _, aliasList := range aliases {
		for _, alias := range aliasList {
			if _, exists := repositoryAliasesMap[alias.Name]; exists {
				slog.Warn("repository alias already exists, overriding", "alias", alias.Name, "old_url", repositoryAliasesMap[alias.Name].URL, "new_url", alias.URL)
			}
			repositoryAliasesMap[alias.Name] = alias
		}
	}

	mergedAliases := make([]RepositoryAlias, 0, len(repositoryAliasesMap))
	for _, alias := range repositoryAliasesMap {
		mergedAliases = append(mergedAliases, alias)
	}

	return &mergedAliases
}

func mergeMappings(mappings ...[]Mapping) *[]Mapping {
	mappingsMap := make(map[string]Mapping)

	for _, mappingList := range mappings {
		for _, mapping := range mappingList {
			if _, exists := mappingsMap[mapping.Name]; exists {
				slog.Warn("mapping already exists, overriding", "mapping", mapping.Name, "old_upstream", mappingsMap[mapping.Name].Upstream, "new_upstream", mapping.Upstream, "old_downstream", mappingsMap[mapping.Name].Downstream, "new_downstream", mapping.Downstream)
			}
			mappingsMap[mapping.Name] = mapping
		}
	}

	mergedMappings := make([]Mapping, 0, len(mappingsMap))
	for _, mapping := range mappingsMap {
		mergedMappings = append(mergedMappings, mapping)
	}

	return &mergedMappings
}
