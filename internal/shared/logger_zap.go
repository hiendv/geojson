package shared

import (
	"go.uber.org/zap"
)

// LoggerZap wraps Zap's SugaredLogger.
type LoggerZap struct {
	*zap.SugaredLogger
}

// Clone copies itself.
func (logger *LoggerZap) Clone() Logger {
	return &LoggerZap{logger.With()}
}

// NewLoggerZap constructs a LoggerZap with verbosity specification.
func NewLoggerZap(verbose bool) (Logger, error) {
	config := zap.Config{
		Level:            zap.NewAtomicLevelAt(zap.InfoLevel),
		Development:      false,
		DisableCaller:    true,
		Encoding:         "console",
		EncoderConfig:    zap.NewDevelopmentEncoderConfig(),
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}

	if verbose {
		config.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	}

	logCore, err := config.Build()
	if err != nil {
		return nil, err
	}

	return &LoggerZap{logCore.Sugar()}, nil
}
