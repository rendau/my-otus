package logger

import (
	"go.uber.org/zap"
	"log"
)

// Logger - is type for logger
type Logger struct {
	l  *zap.Logger
	sl *zap.SugaredLogger
}

// NewLogger - creates new Logger instance
func NewLogger(filePath, level string, debug, test bool) (*Logger, error) {
	var err error

	logger := &Logger{}

	if test {
		logger.l = zap.NewExample(zap.AddCallerSkip(1))
	} else if debug {
		logger.l, err = zap.NewDevelopment(zap.AddCallerSkip(1))
		if err != nil {
			return nil, err
		}
	} else {
		cfg := zap.NewProductionConfig()

		if filePath != "" {
			cfg.OutputPaths = []string{filePath}
			cfg.ErrorOutputPaths = []string{filePath}
		}

		if level != "" {
			switch level {
			case "error":
				cfg.Level.SetLevel(zap.ErrorLevel)
				break
			case "warn":
				cfg.Level.SetLevel(zap.WarnLevel)
				break
			case "info":
				cfg.Level.SetLevel(zap.InfoLevel)
				break
			case "debug":
				cfg.Level.SetLevel(zap.DebugLevel)
				break
			}
		}

		logger.l, err = cfg.Build(zap.AddCallerSkip(1))
		if err != nil {
			return nil, err
		}
	}

	logger.sl = logger.l.Sugar()

	return logger, nil
}

// Fatal - is for fatal
func (lg *Logger) Fatal(args ...interface{}) {
	lg.sl.Fatal(args...)
}

// Fatalf - is for fatalf
func (lg *Logger) Fatalf(tmpl string, args ...interface{}) {
	lg.sl.Fatalf(tmpl, args...)
}

// Fatalw - is for fatalw
func (lg *Logger) Fatalw(msg string, args ...interface{}) {
	lg.sl.Fatalw(msg, args...)
}

// Error - is for error
func (lg *Logger) Error(args ...interface{}) {
	lg.sl.Error(args...)
}

// Errorf - is for errorf
func (lg *Logger) Errorf(tmpl string, args ...interface{}) {
	lg.sl.Errorf(tmpl, args...)
}

// Errorw - is for errorw
func (lg *Logger) Errorw(msg string, args ...interface{}) {
	lg.sl.Errorw(msg, args...)
}

// Warn - is for warn
func (lg *Logger) Warn(args ...interface{}) {
	lg.sl.Warn(args...)
}

// Warnf - is for warnf
func (lg *Logger) Warnf(tmpl string, args ...interface{}) {
	lg.sl.Warnf(tmpl, args...)
}

// Warnw - is for warnw
func (lg *Logger) Warnw(msg string, args ...interface{}) {
	lg.sl.Warnw(msg, args...)
}

// Info - is for info
func (lg *Logger) Info(args ...interface{}) {
	lg.sl.Info(args...)
}

// Infof - is for infof
func (lg *Logger) Infof(tmpl string, args ...interface{}) {
	lg.sl.Infof(tmpl, args...)
}

// Infow - is for infow
func (lg *Logger) Infow(msg string, args ...interface{}) {
	lg.sl.Infow(msg, args...)
}

// Debug - is for debug
func (lg *Logger) Debug(args ...interface{}) {
	lg.sl.Debug(args...)
}

// Debugf - is for debugf
func (lg *Logger) Debugf(tmpl string, args ...interface{}) {
	lg.sl.Debugf(tmpl, args...)
}

// Debugw - is for debugw
func (lg *Logger) Debugw(msg string, args ...interface{}) {
	lg.sl.Debugw(msg, args...)
}

// Sync - is flush last data
func (lg *Logger) Sync() {
	err := lg.sl.Sync()
	if err != nil {
		log.Fatalln("Fail to sync zap-logger", err)
	}
}
