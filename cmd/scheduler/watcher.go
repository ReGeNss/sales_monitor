package main

import (
	"context"
	"log"
	"path/filepath"
	"time"

	"github.com/fsnotify/fsnotify"
)

const reloadDebounce = 500 * time.Millisecond

type ConfigWatcher struct {
	path     string
	logger   *log.Logger
	debounce time.Duration
}

func NewConfigWatcher(path string, logger *log.Logger) *ConfigWatcher {
	return &ConfigWatcher{
		path:     path,
		logger:   logger,
		debounce: reloadDebounce,
	}
}

func (w *ConfigWatcher) Start(ctx context.Context, onChange func()) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	dir := filepath.Dir(w.path)
	if err := watcher.Add(dir); err != nil {
		_ = watcher.Close()
		return err
	}

	go w.watchLoop(ctx, watcher, onChange)
	return nil
}

func (w *ConfigWatcher) watchLoop(ctx context.Context, watcher *fsnotify.Watcher, onChange func()) {
	defer func() {
		if err := watcher.Close(); err != nil {
			w.logger.Printf("config watcher close error: %v", err)
		}
	}()

	var timer *time.Timer
	var timerCh <-chan time.Time

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			if w.shouldReload(event) {
				timer, timerCh = w.resetTimer(timer)
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			w.logger.Printf("config watcher error: %v", err)
		case <-timerCh:
			timer = nil
			timerCh = nil
			onChange()
		case <-ctx.Done():
			if timer != nil {
				timer.Stop()
			}
			return
		}
	}
}

func (w *ConfigWatcher) shouldReload(event fsnotify.Event) bool {
	if event.Op&(fsnotify.Write|fsnotify.Create|fsnotify.Rename) == 0 {
		return false
	}
	return filepath.Clean(event.Name) == w.path
}

func (w *ConfigWatcher) resetTimer(timer *time.Timer) (*time.Timer, <-chan time.Time) {
	if timer == nil {
		timer = time.NewTimer(w.debounce)
		return timer, timer.C
	}
	if !timer.Stop() {
		select {
		case <-timer.C:
		default:
		}
	}
	timer.Reset(w.debounce)
	return timer, timer.C
}
