package cli

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
)

// setupSignalHandling configures handling for SIGINT and SIGTERM signals.
func setupSignalHandling(mainCancel context.CancelFunc, logger *logrus.Entry) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigChan
		logger.Infof("Received signal %s, initiating shutdown...", sig)
		mainCancel()
	}()
}
