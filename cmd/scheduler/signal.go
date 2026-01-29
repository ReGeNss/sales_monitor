package main

import (
	"context"
	"os"
	"os/signal"
)

func notifyContext(signals ...os.Signal) (context.Context, context.CancelFunc) {
	allSignals := append([]os.Signal{os.Interrupt}, signals...)
	return signal.NotifyContext(context.Background(), allSignals...)
}
