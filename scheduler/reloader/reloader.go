package reloader

import (
	"log"

	"sales_monitor/scheduler/config"
	"sales_monitor/scheduler/scheduler"
)

type ConfigReloader struct {
	loader    config.ConfigLoader
	scheduler scheduler.JobScheduler
	logger    *log.Logger
}

func NewConfigReloader(loader config.ConfigLoader, scheduler scheduler.JobScheduler, logger *log.Logger) *ConfigReloader {
	return &ConfigReloader{
		loader:    loader,
		scheduler: scheduler,
		logger:    logger,
	}
}

func (r *ConfigReloader) Reload() (int, error) {
	cfg := r.loader.Load()

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
