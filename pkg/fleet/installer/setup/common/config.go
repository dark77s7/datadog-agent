// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

package common

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

const (
	disclaimerGenerated = `# This file was generated by the Datadog Installer.
# Other configuration options are available, see https://docs.datadoghq.com/agent/guide/agent-configuration-files/ for more information.`

	configDir              = "/etc/datadog-agent"
	datadogConfFile        = "datadog.yaml"
	injectTracerConfigFile = "inject/tracer.yaml"
)

func writeConfigs(config Config, configDir string) error {
	err := writeConfig(filepath.Join(configDir, datadogConfFile), config.DatadogYAML, 0640, true)
	if err != nil {
		return fmt.Errorf("could not write datadog.yaml: %w", err)
	}
	err = writeConfig(filepath.Join(configDir, injectTracerConfigFile), config.InjectTracerYAML, 0644, false)
	if err != nil {
		return fmt.Errorf("could not write tracer.yaml: %w", err)
	}
	for name, config := range config.IntegrationConfigs {
		err = writeConfig(filepath.Join(configDir, "conf.d", name), config, 0644, false)
		if err != nil {
			return fmt.Errorf("could not write %s.yaml: %w", name, err)
		}
	}
	return nil
}

func writeConfig(path string, config any, perms os.FileMode, merge bool) error {
	serializedNewConfig, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("could not serialize config: %w", err)
	}
	var newConfig map[string]interface{}
	err = yaml.Unmarshal(serializedNewConfig, &newConfig)
	if err != nil {
		return fmt.Errorf("could not unmarshal config: %w", err)
	}
	if len(newConfig) == 0 {
		return nil
	}
	err = os.MkdirAll(filepath.Dir(path), 0755)
	if err != nil {
		return fmt.Errorf("could not create config directory: %w", err)
	}
	var existingConfig map[string]interface{}
	if merge {
		serializedExistingConfig, err := os.ReadFile(path)
		if err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("could not read existing config: %w", err)
		}
		if err == nil {
			err = yaml.Unmarshal(serializedExistingConfig, &existingConfig)
			if err != nil {
				return fmt.Errorf("could not unmarshal existing config: %w", err)
			}
		}
	}
	merged, err := mergeConfig(existingConfig, newConfig)
	if err != nil {
		return fmt.Errorf("could not merge config: %w", err)
	}
	serializedMerged, err := yaml.Marshal(merged)
	if err != nil {
		return fmt.Errorf("could not serialize merged config: %w", err)
	}
	if len(existingConfig) == 0 {
		serializedMerged = []byte(disclaimerGenerated + "\n\n" + string(serializedMerged))
	}
	err = os.WriteFile(path, serializedMerged, perms)
	if err != nil {
		return fmt.Errorf("could not write config: %w", err)
	}
	return nil
}

// Config represents the configuration to write in /etc/datadog-agent
type Config struct {
	// DatadogYAML is the content of the datadog.yaml file
	DatadogYAML DatadogConfig
	// InjectTracerYAML is the content of the inject/tracer.yaml file
	InjectTracerYAML InjectTracerConfig
	// IntegrationConfigs is the content of the integration configuration files under conf.d/
	IntegrationConfigs map[string]IntegrationConfig
}

// DatadogConfig represents the configuration to write in /etc/datadog-agent/datadog.yaml
type DatadogConfig struct {
	APIKey               string                     `yaml:"api_key"`
	Hostname             string                     `yaml:"hostname,omitempty"`
	Site                 string                     `yaml:"site,omitempty"`
	Env                  string                     `yaml:"env,omitempty"`
	Tags                 []string                   `yaml:"tags,omitempty"`
	LogsEnabled          bool                       `yaml:"logs_enabled,omitempty"`
	DJM                  DatadogConfigDJM           `yaml:"djm,omitempty"`
	ProcessConfig        DatadogConfigProcessConfig `yaml:"process_config,omitempty"`
	ExpectedTagsDuration string                     `yaml:"expected_tags_duration,omitempty"`
}

// DatadogConfigDJM represents the configuration for the Data Jobs Monitoring
type DatadogConfigDJM struct {
	Enabled bool `yaml:"enabled,omitempty"`
}

// DatadogConfigProcessConfig represents the configuration for the process agent
type DatadogConfigProcessConfig struct {
	ExpvarPort int `yaml:"expvar_port,omitempty"`
}

// IntegrationConfig represents the configuration for an integration under conf.d/
type IntegrationConfig struct {
	InitConfig []any                   `yaml:"init_config"`
	Instances  []any                   `yaml:"instances,omitempty"`
	Logs       []IntegrationConfigLogs `yaml:"logs,omitempty"`
}

// IntegrationConfigLogs represents the configuration for the logs of an integration
type IntegrationConfigLogs struct {
	Type    string `yaml:"type,omitempty"`
	Path    string `yaml:"path,omitempty"`
	Service string `yaml:"service,omitempty"`
	Source  string `yaml:"source,omitempty"`
}

// IntegrationConfigInstanceSpark represents the configuration for the Spark integration
type IntegrationConfigInstanceSpark struct {
	SparkURL         string `yaml:"spark_url"`
	SparkClusterMode string `yaml:"spark_cluster_mode"`
	ClusterName      string `yaml:"cluster_name"`
	StreamingMetrics bool   `yaml:"streaming_metrics"`
}

// InjectTracerConfig represents the configuration to write in /etc/datadog-agent/inject/tracer.yaml
type InjectTracerConfig struct {
	Version                        int                        `yaml:"version,omitempty"`
	ConfigSources                  string                     `yaml:"config_sources,omitempty"`
	AdditionalEnvironmentVariables []InjectTracerConfigEnvVar `yaml:"additional_environment_variables,omitempty"`
}

// InjectTracerConfigEnvVar represents an environment variable to inject
type InjectTracerConfigEnvVar struct {
	Key   string `yaml:"key"`
	Value string `yaml:"value"`
}

// mergeConfig merges the current config with the setup config.
//
// The values are merged as follows:
// - Scalars: the override value is used
// - Lists: the override list is used
// - Maps: the override map is recursively merged into the base map
func mergeConfig(base interface{}, override interface{}) (interface{}, error) {
	if base == nil {
		return override, nil
	}
	if override == nil {
		// this allows to override a value with nil
		return nil, nil
	}
	if isScalar(base) && isScalar(override) {
		return override, nil
	}
	if isList(base) && isList(override) {
		return override, nil
	}
	if isMap(base) && isMap(override) {
		return mergeMap(base.(map[string]interface{}), override.(map[string]interface{}))
	}
	return nil, fmt.Errorf("could not merge %T with %T", base, override)
}

func mergeMap(base, override map[string]interface{}) (map[string]interface{}, error) {
	merged := make(map[string]interface{})
	for k, v := range base {
		merged[k] = v
	}
	for k := range override {
		v, err := mergeConfig(base[k], override[k])
		if err != nil {
			return nil, fmt.Errorf("could not merge key %v: %w", k, err)
		}
		merged[k] = v
	}
	return merged, nil
}

func isList(i interface{}) bool {
	_, ok := i.([]interface{})
	return ok
}

func isMap(i interface{}) bool {
	_, ok := i.(map[string]interface{})
	return ok
}

func isScalar(i interface{}) bool {
	return !isList(i) && !isMap(i)
}
