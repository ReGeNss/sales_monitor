package main

import (
	"log"
	"os"
	"path/filepath"
	"syscall"

	"github.com/joho/godotenv"
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

	configPath, err := requireEnv("SCRAPER_CONFIG_PATH")
	if err != nil {
		return err
	}
	workerCmd, err := requireEnv("SCRAPER_WORKER_CMD")
	if err != nil {
		return err
	}

	absConfigPath, err := filepath.Abs(configPath)
	if err != nil {
		return err
	}

	loader := FileConfigLoader{Path: absConfigPath}
	runner := NewWorkerLauncher(workerCmd, absConfigPath, logger)
	scheduler := NewCronJobScheduler(runner, logger)
	reloader := NewConfigReloader(loader, scheduler, logger)
	watcher := NewConfigWatcher(absConfigPath, logger)

	ctx, stop := notifyContext(syscall.SIGTERM)
	defer stop()

	app := NewApp(reloader, scheduler, watcher, logger)
	return app.Run(ctx)
}
