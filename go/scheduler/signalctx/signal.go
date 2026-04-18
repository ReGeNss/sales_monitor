package signalctx

import (
	"context"
	"os"
	"os/signal"
)

func CreateContext(signals ...os.Signal) (context.Context, context.CancelFunc) {
	allSignals := append([]os.Signal{os.Interrupt}, signals...)
	return signal.NotifyContext(context.Background(), allSignals...)
}
