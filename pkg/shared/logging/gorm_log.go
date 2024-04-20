package logging

import (
	"context"
	"errors"
	"go-server/pkg/shared/logging/hooks"
	"path/filepath"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	gormLog "gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
)

// Logger based on logrus, but compatible with gorm
type GormLogger struct {
	Logger                    *logrus.Entry
	LogLevel                  gormLog.LogLevel
	IgnoreRecordNotFoundError bool
	SlowThreshold             time.Duration
	FileWithLineNumField      string
}

// NewGormLogger creates a new GormLogger instance with the provided options.
// It initializes the logger if not provided and sets the log level to Silent if not specified.
func NewGormLogger(
	opts GormLogger,
) *GormLogger {
	if opts.Logger == nil {
		opts.Logger = logrus.NewEntry(logrus.New())
	}

	if opts.LogLevel == 0 {
		opts.LogLevel = gormLog.Silent
	}
	return &opts
}

// LogMode sets the log level for the GormLogger and returns the updated logger instance.
func (l *GormLogger) LogMode(level gormLog.LogLevel) gormLog.Interface {
	newlogger := *l
	newlogger.LogLevel = level
	return &newlogger
}

// Info logs an informational message
func (l *GormLogger) Info(ctx context.Context, s string, args ...interface{}) {
	if l.LogLevel >= gormLog.Info {
		l.Logger.WithContext(ctx).Infof(s, args...)
	}
}

// Warn logs a warning message
func (l *GormLogger) Warn(ctx context.Context, s string, args ...interface{}) {
	if l.LogLevel >= gormLog.Warn {
		l.Logger.WithContext(ctx).Warnf(s, args...)
	}
}

// Error logs an error message
func (l *GormLogger) Error(ctx context.Context, s string, args ...interface{}) {
	if l.LogLevel >= gormLog.Error {
		l.Logger.WithContext(ctx).Errorf(s, args...)
	}
}

// Trace logs detailed information about a database query execution, including SQL, duration, and potential errors.
// We want the SQL logs with the info level, while it's defined as trace by gorm
// It checks the log level to determine whether to log the information and formats the log message accordingly.
// Additionally, it handles cases such as slow queries and error messages related to database operations.
func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.LogLevel <= gormLog.Silent {
		return
	}

	fields := logrus.Fields{}
	fields[l.FileWithLineNumField] = filepath.Base(utils.FileWithLineNum())

	sql, rows := fc()
	if rows == -1 {
		fields["rows"] = "-"
	} else {
		fields["rows"] = rows
	}
	fields["sql"] = sql

	elapsed := time.Since(begin)
	fields["durations"] = elapsed.String()

	if requestID, ok := ctx.Value(hooks.RequestIDField).(string); ok {
		fields[hooks.RequestIDField] = requestID
	}

	switch {
	case err != nil && (!errors.Is(err, gorm.ErrRecordNotFound) || !l.IgnoreRecordNotFoundError) && l.LogLevel >= gormLog.Error:
		fields["error"] = err.Error()
		l.Logger.WithContext(ctx).WithFields(fields).Error("SQL Query failed")
	case l.SlowThreshold != 0 && elapsed > l.SlowThreshold && l.LogLevel >= gormLog.Warn:
		l.Logger.WithContext(ctx).WithFields(fields).Warn("Performed SLOW SQL Query")
	case l.LogLevel == gormLog.Info:
		l.Logger.WithContext(ctx).WithFields(fields).Info("Performed SQL Query")
	}
}
