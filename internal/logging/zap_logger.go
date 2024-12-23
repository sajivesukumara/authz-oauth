package logger

import (
    "go.uber.org/zap"
    "go.uber.org/zap/zapcore"
    "os"
)

// LogLevel defines the level of logging
type LogLevel int8

type Fields map[string]any

const (
    DEBUG LogLevel = iota - 1
    INFO
    WARNING
    ERROR
)

// ZapJSONLogger is an implementation of the logging repository that
// uses zap's sugared logger.
type ZapJSONLogger struct {
	logger *zap.Logger
}

func getLoggerLevel(value string) LogLevel {
    switch value {
    case "DEBUG":
     return DEBUG
    case "ERROR":
     return ERROR
    case "INFO":
        return INFO  
    default:
     return INFO
    }
}

// GetLogger will return the SingletonLogger (after InitializeConfigurations is called.).
// otherwise, it will create a default logger with level "defaultLoggingLevel"
func GetLogger() Logger {
	// see comment above for usage of SingletonLogger
    logLevel := getLoggerLevel(os.Getenv("LOG_LEVEL"))
	logger, _ := NewZapJSONLogger(logLevel)
	return logger
}

func GetZapLogger() *zap.Logger {
	// see comment above for usage of SingletonLogger
    logLevel := getLoggerLevel(os.Getenv("LOG_LEVEL"))
	logger, _ := NewZapJSONLogger(logLevel)
	return logger.logger
}

func NewZapJSONLogger(logLevel LogLevel) (*ZapJSONLogger, error) {
	zapConfig, err := stdJSONLoggerConfig(logLevel)
	if err != nil {
		return nil, err
	}

    logger, err := zapConfig.Build(
        zap.AddCallerSkip(1),
        zap.AddStacktrace(zap.ErrorLevel),
    )
    if err != nil {
        return nil, err
    }

	return &ZapJSONLogger{
		logger: logger,
	}, nil

}

// SetLogLevel sets the log level for the logger
func (l *ZapJSONLogger) SetLogLevel(level LogLevel) {
    l.logger = l.logger.WithOptions(zap.IncreaseLevel(zapcore.Level(level)))
}

// Debug logs a debug message
func (l *ZapJSONLogger) Debug(msg string) {
	l.logger.Debug(msg)
}

// Info logs an info message
func (l *ZapJSONLogger) Info(msg string) {
        l.logger.Info(msg)
}

// Warning logs a warning message
func (l *ZapJSONLogger) Warn(msg string) {
        l.logger.Warn(msg)
}

// Error logs an error message
func (l *ZapJSONLogger) Error(msg string, err error) {
	l.logger.Error(msg)
}

// WithField returns a new logger with the specified key-value pair attached for
// subsequent logging operations.
//
// This function returns a repositories logger interface rather than the explicit
// ZapJSONLogger to allow it to satisfy the Logger interface.
func (l *ZapJSONLogger) WithField(key string, value any) Logger {
	return &ZapJSONLogger{
		l.logger.With(zap.Any(key, value)),
	}
}

func (z *ZapJSONLogger) WithFields(fields Fields) Logger {
	fieldList := make([]zap.Field, 0, len(fields))
	for key, value := range fields {
		fieldList = append(fieldList, zap.Any(key, value))
	}

	return &ZapJSONLogger{
		z.logger.With(fieldList...),
	}
}

// Flush flushes any pending log statements. This is a no-op as logs are written to STDOUT and
// synchonization is not supported on STDOUT/STDERR.
func (z *ZapJSONLogger) Flush() error {
	return nil
}

func stdJSONLoggerConfig(level LogLevel) (zap.Config, error) {
	
	return zap.Config{
		Level:             zap.NewAtomicLevelAt(zapcore.Level(level)),
		Development:       false,
		DisableStacktrace: false,
		Encoding:          "json",
		OutputPaths:       []string{"stdout"},
		ErrorOutputPaths:  []string{"stderr"},
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "timestamp",
			LevelKey:       "level",
			MessageKey:     "message",
			NameKey:        "name",
			StacktraceKey:  "",
			//CallerKey:      "caller",
			LineEnding:     "\n",
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
			EncodeCaller:   zapcore.FullCallerEncoder,
		},
	}, nil
}
