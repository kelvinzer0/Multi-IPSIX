package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// InterfaceConfig defines the structure for a single interface in the YAML file.
// It includes the interface name, the priority IP, and a list of all addresses to manage.
type InterfaceConfig struct {
	Name        string   `yaml:"name"`
	PriorityIP  string   `yaml:"priority_ip"`
	Addresses   []string `yaml:"addresses"`
	MonitorIPv6 bool     `yaml:"monitor_ipv6"`
}

// Config is the top-level structure for the entire YAML configuration.
type Config struct {
	Interfaces []InterfaceConfig `yaml:"interfaces"`
}

// LoadConfig reads a YAML file from the given path and unmarshals it into a Config struct.
// It returns the loaded configuration or an error if the file cannot be read or parsed.
func LoadConfig(filePath string) (*Config, error) {
	// Read the YAML file
	yamlFile, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	// Parse the YAML file
	var config Config
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
