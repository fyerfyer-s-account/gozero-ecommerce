package zerolog

import "github.com/zeromicro/go-zero/core/logx"

type LogWrapper struct {
	logx.Logger
}

func (l *LogWrapper) Info(msg string, keysAndValues ...interface{}) {
	l.Logger.Infof(msg, keysAndValues...)
}

func (l *LogWrapper) Error(msg string, keysAndValues ...interface{}) {
	l.Logger.Errorf(msg, keysAndValues...)
}