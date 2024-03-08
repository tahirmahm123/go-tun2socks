package log

var logger Logger

func RegisterLogger(l Logger) {
	logger = l
}

func Debugf(args ...interface{}) {
	if logger != nil {
		logger.Debug(args...)
	}
}

func Infof(args ...interface{}) {
	if logger != nil {
		logger.Info(args...)
	}
}

func Warnf(args ...interface{}) {
	if logger != nil {
		logger.Warning(args...)
	}
}

func Errorf(args ...interface{}) {
	if logger != nil {
		logger.Error(args...)
	}
}

func Fatalf(args ...interface{}) {
	if logger != nil {
		logger.Panic(args...)
	}
}
