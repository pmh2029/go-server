package logging

import (
	"os"

	"github.com/sirupsen/logrus"
)

// NewLogger creates a new instance of a logger with the given configuration
func NewLogger() *logrus.Logger {
	// create a new instance of a logger
	logger := logrus.New()

	logger.SetOutput(os.Stdout)

	// enable caller reporting
	logger.SetReportCaller(true)

	// set the log formatter to JSON with pretty printing
	logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "15:04:05 02/01/2006",
		PrettyPrint:     true,
	})

	// set the log level to info
	logger.SetLevel(logrus.InfoLevel)

	return logger
}
