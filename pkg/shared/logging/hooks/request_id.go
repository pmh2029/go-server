package hooks

import (
	"github.com/sirupsen/logrus"
)

const (
	RequestIDField = "request-id"
)

// RequestIDHook is a Logrus hook for including request ID in log entries
type RequestIDHook struct{}

// Levels returns the logging levels for which this hook should be called
func (hook *RequestIDHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

// Fire is called when a log entry is made
func (hook *RequestIDHook) Fire(entry *logrus.Entry) error {
	// If request ID is present in the log entry context, add it to the log entry fields
	ctx := entry.Context
	if ctx == nil {
		return nil
	}

	if requestID, ok := ctx.Value(RequestIDField).(string); ok {
		entry.Data[RequestIDField] = requestID
	}
	return nil
}
