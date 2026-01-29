package main

import (
	"log"
	"os"
	"os/exec"
	"strings"
)

type WorkerLauncher struct {
	cmd        string
	configPath string
	env        []string
	logger     *log.Logger
}

func NewWorkerLauncher(cmd, configPath string, logger *log.Logger) *WorkerLauncher {
	return &WorkerLauncher{
		cmd:        cmd,
		configPath: configPath,
		env:        os.Environ(),
		logger:     logger,
	}
}

func (w *WorkerLauncher) Run(jobID string) {
	cmdParts := strings.Fields(w.cmd)
	if len(cmdParts) == 0 {
		w.logger.Printf("worker command is empty; job %q skipped", jobID)
		return
	}

	args := append(cmdParts[1:], "--config", w.configPath, "--job-id", jobID)
	cmd := exec.Command(cmdParts[0], args...)
	cmd.Env = w.env
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		w.logger.Printf("failed to start worker for job %q: %v", jobID, err)
		return
	}

	go func() {
		if err := cmd.Wait(); err != nil {
			w.logger.Printf("worker for job %q exited with error: %v", jobID, err)
		}
	}()
}
