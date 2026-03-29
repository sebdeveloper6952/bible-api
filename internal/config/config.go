package config

import (
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Server ServerConfig `yaml:"server"`
	CORS   CORSConfig   `yaml:"cors"`
}

type ServerConfig struct {
	Addr   string `yaml:"addr"`
	DBPath string `yaml:"db_path"`
}

type CORSConfig struct {
	AllowedOrigins []string `yaml:"allowed_origins"`
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// MatchesOrigin checks whether hostname matches any of the allowed origin
// patterns. Patterns may use a leading wildcard (e.g. "*.example.com").
func (c *CORSConfig) MatchesOrigin(hostname string) bool {
	for _, pattern := range c.AllowedOrigins {
		if strings.HasPrefix(pattern, "*.") {
			suffix := pattern[1:] // e.g. ".sebdev.io"
			base := pattern[2:]   // e.g. "sebdev.io"
			if hostname == base || strings.HasSuffix(hostname, suffix) {
				return true
			}
		} else if hostname == pattern {
			return true
		}
	}
	return false
}
