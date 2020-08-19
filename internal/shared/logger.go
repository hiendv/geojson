package shared

// Logger is the contract of an application logger.
type Logger interface {
	Infow(string, ...interface{})
	Debugw(string, ...interface{})
	Warnw(string, ...interface{})
	Error(...interface{})
	Errorw(string, ...interface{})
	Clone() Logger
	Sync() error
}

type loggerNoop struct{}

func (logger *loggerNoop) Infow(string, ...interface{})  {}
func (logger *loggerNoop) Debugw(string, ...interface{}) {}
func (logger *loggerNoop) Warnw(string, ...interface{})  {}
func (logger *loggerNoop) Error(...interface{})          {}
func (logger *loggerNoop) Errorw(string, ...interface{}) {}
func (logger *loggerNoop) Sync() error {
	return nil
}

func (logger *loggerNoop) Clone() Logger {
	return logger
}

// LoggerNoop is a logger which does nothing.
var LoggerNoop = &loggerNoop{}
