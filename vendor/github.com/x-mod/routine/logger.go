package routine

import (
	"context"
)

type _logger struct{}

//Logger declare
type Logger interface {
	Debug(args ...interface{})
	Trace(args ...interface{})
	Info(args ...interface{})
	Error(args ...interface{})
	Warn(args ...interface{})
}

//WithLogger context with logger
func WithLogger(ctx context.Context, logger Logger) context.Context {
	if ctx != nil {
		return context.WithValue(ctx, _logger{}, logger)
	}
	return context.WithValue(context.TODO(), _logger{}, logger)
}

//Debug ctx
func Debug(ctx context.Context, args ...interface{}) {
	if ctx != nil {
		logger := ctx.Value(_logger{})
		if logger != nil {
			logger.(Logger).Debug(args...)
		}
	}
}

//Trace ctx
func Trace(ctx context.Context, args ...interface{}) {
	if ctx != nil {
		logger := ctx.Value(_logger{})
		if logger != nil {
			logger.(Logger).Trace(args...)
		}
	}
}

//Info ctx
func Info(ctx context.Context, args ...interface{}) {
	if ctx != nil {
		logger := ctx.Value(_logger{})
		if logger != nil {
			logger.(Logger).Info(args...)
		}
	}
}

//Error ctx
func Error(ctx context.Context, args ...interface{}) {
	if ctx != nil {
		logger := ctx.Value(_logger{})
		if logger != nil {
			logger.(Logger).Error(args...)
		}
	}
}

//Warn ctx
func Warn(ctx context.Context, args ...interface{}) {
	if ctx != nil {
		logger := ctx.Value(_logger{})
		if logger != nil {
			logger.(Logger).Warn(args...)
		}
	}
}
