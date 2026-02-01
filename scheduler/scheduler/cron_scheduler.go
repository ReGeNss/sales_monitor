package scheduler

import (
	"fmt"
	"log"
	"sync"
	"time"

	config "sales_monitor/internal/scheduler_config"

	"github.com/robfig/cron/v3"
)

type JobRunner interface {
	Run(jobID string)
}

type JobScheduler interface {
	Schedule(location *time.Location, jobs []config.ResolvedJob) error
	Stop()
}

type CronJobScheduler struct {
	mu     sync.Mutex
	cron   *cron.Cron
	runner JobRunner
	logger *log.Logger
}

func NewCronJobScheduler(runner JobRunner, logger *log.Logger) *CronJobScheduler {
	return &CronJobScheduler{
		runner: runner,
		logger: logger,
	}
}

func (s *CronJobScheduler) Schedule(location *time.Location, jobs []config.ResolvedJob) error {
	newCron := cron.New(cron.WithLocation(location))
	for _, job := range jobs {
		jobCopy := job
		if _, err := newCron.AddFunc(job.Cron, func() {
			s.runner.Run(jobCopy.ID)
		}); err != nil {
			return fmt.Errorf("schedule job %q: %w", job.ID, err)
		}
	}

	newCron.Start()

	s.mu.Lock()
	oldCron := s.cron
	s.cron = newCron
	s.mu.Unlock()

	if oldCron != nil {
		ctx := oldCron.Stop()
		<-ctx.Done()
	}

	return nil
}

func (s *CronJobScheduler) Stop() {
	s.mu.Lock()
	current := s.cron
	s.cron = nil
	s.mu.Unlock()

	if current == nil {
		return
	}

	ctx := current.Stop()
	<-ctx.Done()
	s.logger.Printf("scheduler stopped")
}
