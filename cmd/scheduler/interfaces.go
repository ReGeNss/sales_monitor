package main

import (
	"time"

	"sales_monitor/internal/schedulerconfig"
)

type ConfigLoader interface {
	Load() (*schedulerconfig.Config, error)
}

type JobRunner interface {
	Run(jobID string)
}

type JobScheduler interface {
	Schedule(location *time.Location, jobs []schedulerconfig.ResolvedJob) error
	Stop()
}
