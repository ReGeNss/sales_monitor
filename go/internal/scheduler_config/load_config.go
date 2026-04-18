package scheduler_config

import (
	"log"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

func LoadConfig(path string) (*Config) {
	if path == "" {
		log.Fatalln("config path is empty")
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		log.Fatalf("resolve config path: %v", err)
	}

	data, err := os.ReadFile(absPath)
	if err != nil {
		log.Fatalf("read config: %v", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		log.Fatalf("parse yaml: %v", err)
	}

	if err := cfg.Validate(); err != nil {
		log.Fatalf("validate config: %v", err)
	}

	return &cfg
}
