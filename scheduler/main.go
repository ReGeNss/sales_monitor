package main

import (
	"log"
	"os"
	"path/filepath"
	"syscall"

	"github.com/joho/godotenv"

	"sales_monitor/scheduler/app"
	"sales_monitor/scheduler/config"
	"sales_monitor/scheduler/reloader"
	"sales_monitor/scheduler/scheduler"
	"sales_monitor/scheduler/signalctx"
	"sales_monitor/scheduler/watcher"
	"sales_monitor/scheduler/worker"
)

func main() {
	logger := log.New(os.Stdout, "", log.LstdFlags)
	if err := run(logger); err != nil {
		logger.Fatalf("scheduler stopped: %v", err)
	}
}

func run(logger *log.Logger) error {
	if err := godotenv.Load(); err != nil {
		return err
	}

	configPath := config.RequireEnv("SCRAPER_CONFIG_PATH")
	workerCmd := config.RequireEnv("SCRAPER_WORKER_CMD")

	absConfigPath, err := filepath.Abs(configPath)
	if err != nil {
		return err
	}

	ctx, stop := signalctx.CreateContext(syscall.SIGTERM)
	defer stop()

	loader := config.FileConfigLoader{Path: absConfigPath}
	runner := worker.NewWorkerLauncher(ctx, workerCmd, absConfigPath, logger)
	jobScheduler := scheduler.NewCronJobScheduler(runner, logger)
	configReloader := reloader.NewConfigReloader(loader, jobScheduler, logger)
	cfgWatcher := watcher.NewConfigWatcher(absConfigPath, logger)

	schedulerApp := app.NewApp(configReloader, jobScheduler, cfgWatcher, logger)
	return schedulerApp.Run(ctx)
}
