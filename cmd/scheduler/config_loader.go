package main

import "sales_monitor/internal/schedulerconfig"

type FileConfigLoader struct {
	Path string
}

func (l FileConfigLoader) Load() (*schedulerconfig.Config, error) {
	return schedulerconfig.LoadConfig(l.Path)
}
