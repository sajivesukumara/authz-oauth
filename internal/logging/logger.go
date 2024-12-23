package logger

type Logger interface {
	// Used when an error has occurred that is not recoverable, and will most likely
	// involve returning an error to the consumer/user. Implementations must include a stacktrace at this level.
	Error(msg string, err error)

	// Used when a potential issue may exist, but the system can continue to function.
	Warn(msg string)

	// Used when something of interest has occurred that is useful to have logged in a
	// production setting.
	Info(msg string)

	// Used when providing information on specific code paths with the application that are
	// being executed that are not required in a production setting.
	Debug(msg string)

	// WithField returns a new instance of the Logger that has the specified field attached
	// in all subsequent messages.
	WithField(key string, value any) Logger

	// WithError provides a wrapper around WithField to add an error field to the logger,
	// ensuring consistency of error message keys.
	// WithError(err error) Logger

	// WithFields returns a new instance of the Logger that has the specified fields attached
	// in all subsequent messages.
	WithFields(fields Fields) Logger

	// Flush ensures that any pending log messages are written out. For some implementations
	// this function will be a no-op.
	Flush() error
}
