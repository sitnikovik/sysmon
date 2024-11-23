package main

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"

	"github.com/sitnikovik/sysmon/internal/metrics"
)

// config - struct to hold the configuration of the sysmon.
type config struct {
	Interval int `yaml:"interval"`
	Margin   int `yaml:"margin"`
	GRPCPort int `yaml:"grpcPort"`
	Exclude  struct {
		Metrics []string `yaml:"metrics"`
	} `yaml:"exclude"`
}

func loadConfig(path string) (*config, error) {
	bb, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read the configuration file: %w", err)
	}

	c := &config{}
	if err = yaml.Unmarshal(bb, c); err != nil {
		return nil, fmt.Errorf("failed to load the configuration: %w", err)
	}

	if err = c.validate(); err != nil {
		return nil, fmt.Errorf("failed to validate the configuration: %w", err)
	}

	return c, nil
}

// Validate validates the configuration.
func (c *config) validate() error {
	if c.Interval <= 0 {
		return fmt.Errorf("invalid interval: %d", c.Interval)
	}

	if c.Margin <= 0 {
		return fmt.Errorf("invalid margin: %d", c.Margin)
	}

	if c.GRPCPort <= 0 {
		return fmt.Errorf("invalid gRPC port: %d", c.GRPCPort)
	}

	for _, metric := range c.Exclude.Metrics {
		if metrics.NameToType(metric) == metrics.Undefined {
			return fmt.Errorf("invalid metric name: %s", metric)
		}
	}

	return nil
}

// GetMetricsToParse returns the metrics to parse.
func (c *config) GetMetricsToParse(allMetrics []string) []string {
	excludedMetrics := make(map[string]struct{})
	for _, metric := range c.Exclude.Metrics {
		excludedMetrics[metric] = struct{}{}
	}

	var metricsToParse []string
	for _, metric := range allMetrics {
		if _, excluded := excludedMetrics[metric]; !excluded {
			metricsToParse = append(metricsToParse, metric)
		}
	}

	return metricsToParse
}
