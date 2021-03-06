package ezdbg

import (
	"context"
	"fmt"
	"log"
	_ "unsafe"
)

func Config(cfg Cfg) {
	_logcfg = cfg
}

var _logcfg Cfg

type Cfg struct {
	EnableDebug func(context.Context) bool
	LoggerFunc  func(context.Context) DebugLogger
}

func (p Cfg) getLogger(ctxp *context.Context) DebugLogger {
	ctx := context.Background()
	if ctxp != nil && *ctxp != nil {
		ctx = *ctxp
	}
	if p.LoggerFunc != nil {
		if lg := p.LoggerFunc(ctx); lg != nil {
			return lg
		}
	}
	return stdLogger{}
}

//go:linkname log_std log.std
var log_std *log.Logger

type stdLogger struct{}

const _stdLogDepth = 2

func (_ stdLogger) Debugf(format string, args ...interface{}) {
	log_std.Output(_stdLogDepth, fmt.Sprintf("[DEBUG]: "+format, args...))
}

// DebugLogger is an interface which log an message at DEBUG level.
// It's implemented by *logrus.Logger, *logrus.Entry, *zap.SugaredLogger,
// and many other logging packages.
type DebugLogger interface {
	Debugf(format string, args ...interface{})
}

// PrintFunc is a function to print the given arguments in format to somewhere.
// It implements the interface `ErrDebugLogger`.
type PrintFunc func(format string, args ...interface{})

func (f PrintFunc) Debugf(format string, args ...interface{}) { f(format, args...) }
