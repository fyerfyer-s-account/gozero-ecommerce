package zerolog

import (
    "context"
    "time"

    "github.com/zeromicro/go-zero/core/logx"
)

type Logger struct{}

var defaultLogger = NewLogger()

// GetLogger returns the default logger instance
func GetLogger() *Logger {
    return defaultLogger
}

func NewLogger(opts ...Option) *Logger {
    config := &options{
        level:      "info",
        timeFormat: time.RFC3339,
    }

    for _, opt := range opts {
        opt(config)
    }

    return &Logger{}
}

type options struct {
    level      string
    timeFormat string
}

type Option func(*options)

func WithLevel(level string) Option {
    return func(o *options) {
        o.level = level
    }
}

func (l *Logger) Debug(ctx context.Context, msg string, fields map[string]interface{}) {
    logx.WithContext(ctx).Debugw(msg, fields2KV(fields)...)
}

func (l *Logger) Info(ctx context.Context, msg string, fields map[string]interface{}) {
    logx.WithContext(ctx).Infow(msg, fields2KV(fields)...)
}

func (l *Logger) Warn(ctx context.Context, msg string, fields map[string]interface{}) {
    logx.WithContext(ctx).Sloww(msg, fields2KV(fields)...)
}

func (l *Logger) Error(ctx context.Context, msg string, err error, fields map[string]interface{}) {
    if err != nil {
        fields["error"] = err.Error()
    }
    logx.WithContext(ctx).Errorw(msg, fields2KV(fields)...)
}

func fields2KV(fields map[string]interface{}) []logx.LogField {
    kvs := make([]logx.LogField, 0, len(fields))
    for k, v := range fields {
        kvs = append(kvs, logx.Field(k, v))
    }
    return kvs
}

func (l *Logger) WithError(ctx context.Context, err error, msg string, fields map[string]interface{}) {
    if fields == nil {
        fields = make(map[string]interface{})
    }
    if err != nil {
        fields["error"] = err.Error()
    }
    l.Error(ctx, msg, err, fields)
}