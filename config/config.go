package config

import (
	"io/ioutil"

	"github.com/prometheus/common/log"
	yaml "gopkg.in/yaml.v2"
)

// Config represents the exporter configuration
type Config struct {
	Applications []application `yaml:"applications,omitempty"`
}

type application struct {
	ID   int64  `yaml:"id,omitempty"`
	Name string `yaml:"name,omitempty"`
}

// Parse reads and parse a given configuration file to a new Config
func Parse(path string) Config {
	var config Config

	bts, err := ioutil.ReadFile(path)
	if err != nil {
		log.With("path", path).Fatalf("Failed to read configuration file: %v", err)
	}
	if err := yaml.Unmarshal(bts, &config); err != nil {
		log.With("path", path).Fatalf("Failed to unmarshall configuration file: %v", err)
	}

	return config
}
