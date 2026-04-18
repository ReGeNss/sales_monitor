package config

import config "sales_monitor/internal/scheduler_config"

type ConfigLoader interface {
	Load() (*config.Config)
}

type FileConfigLoader struct {
	Path string
}

func (l FileConfigLoader) Load() (*config.Config) {
	return config.LoadConfig(l.Path)
}
