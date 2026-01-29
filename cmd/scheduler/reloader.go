package main

import (
	"log"
)

type ConfigReloader struct {
	loader    ConfigLoader
	scheduler JobScheduler
	logger    *log.Logger
}

func NewConfigReloader(loader ConfigLoader, scheduler JobScheduler, logger *log.Logger) *ConfigReloader {
	return &ConfigReloader{
		loader:    loader,
		scheduler: scheduler,
		logger:    logger,
	}
}

func (r *ConfigReloader) Reload() (int, error) {
	cfg, err := r.loader.Load()
	if err != nil {
		return 0, err
	}

	location, err := cfg.Location()
	if err != nil {
		return 0, err
	}

	resolvedJobs, err := cfg.ResolveJobs()
	if err != nil {
		return 0, err
	}

	if err := r.scheduler.Schedule(location, resolvedJobs); err != nil {
		return 0, err
	}

	r.logger.Printf("applied scheduler config")
	return len(resolvedJobs), nil
}
