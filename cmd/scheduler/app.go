package main

import (
	"context"
	"log"
)

type App struct {
	reloader  *ConfigReloader
	scheduler JobScheduler
	watcher   *ConfigWatcher
	logger    *log.Logger
}

func NewApp(reloader *ConfigReloader, scheduler JobScheduler, watcher *ConfigWatcher, logger *log.Logger) *App {
	return &App{
		reloader:  reloader,
		scheduler: scheduler,
		watcher:   watcher,
		logger:    logger,
	}
}

func (a *App) Run(ctx context.Context) error {
	count, err := a.reloader.Reload()
	if err != nil {
		return err
	}
	a.logger.Printf("scheduler loaded %d jobs", count)

	if err := a.watcher.Start(ctx, func() {
		count, err := a.reloader.Reload()
		if err != nil {
			a.logger.Printf("reload failed: %v", err)
			return
		}
		a.logger.Printf("scheduler reloaded %d jobs", count)
	}); err != nil {
		return err
	}

	<-ctx.Done()
	a.scheduler.Stop()
	return nil
}
