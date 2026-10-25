package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

var log = logrus.New()

// Init configures the global logger.
func Init() {
	// Log as JSON for structured logging.
	log.SetFormatter(&logrus.JSONFormatter{})
	log.SetOutput(os.Stdout)
	// Set default level; could be extended to read from env.
	log.SetLevel(logrus.InfoLevel)
}

// L returns the configured logger instance.
func L() *logrus.Logger {
	return log
}