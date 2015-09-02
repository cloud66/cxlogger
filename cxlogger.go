package cxlogger

// see http://godoc.org/github.com/inconshreveable/log15 for more info
import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	log "gopkg.in/inconshreveable/log15.v2"
)

const tabWidth = 2

var contextIndentations = make(map[string]int)

// func main() {
// 	Initialize("STDOUT", "debug")
// 	Debug("test!")
// }

type Logger struct {
	log.Logger
	Context string
	Level   log.Lvl
}

const (
	LvlCrit  = log.LvlCrit
	LvlError = log.LvlError
	LvlWarn  = log.LvlWarn
	LvlInfo  = log.LvlInfo
	LvlDebug = log.LvlDebug
)

var Log *Logger

func NewWithContext(context string, logOut string, lvl interface{}) (*Logger, error) {
	l := Logger{}
	return &l, l.InitializeWithContext(context, logOut, lvl)
}

func New(logOut string, lvl interface{}) (*Logger, error) {
	l := Logger{}
	return &l, l.Initialize(logOut, lvl)
}

func (l *Logger) InitializeWithContext(context string, logOut string, lvl interface{}) error {
	l.Context = context
	return l.Initialize(logOut, lvl)
}

func (l *Logger) Initialize(logOut string, lvl interface{}) error {
	var (
		level log.Lvl
		err   error
	)

	if l.Context == "" {
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		randInt := r.Int()
		t := time.Now()
		l.Context = fmt.Sprintf("%s#%d", t.Format(time.RFC3339), randInt)
	}

	if str, ok := lvl.(string); ok {
		level, err = log.LvlFromString(str)
		if err != nil {
			return err
		}
	} else {
		level = lvl.(log.Lvl)
	}
	l.Logger = log.New(log.Ctx{"context": l.Context})
	l.Level = level

	if logOut == "STDOUT" {
		normalHandler := log.LvlFilterHandler(level, log.StdoutHandler)
		errorHandler := log.LvlFilterHandler(level, log.CallerStackHandler("%+v", log.StdoutHandler))
		handler := ErrorMultiHandler(normalHandler, errorHandler)
		l.SetHandler(handler)
	} else if logOut == "NONE" {
		l.SetHandler(log.DiscardHandler())
	} else {
		fileHandler := log.Must.FileHandler(logOut, log.LogfmtFormat())
		normalHandler := log.LvlFilterHandler(level, fileHandler)
		errorHandler := log.LvlFilterHandler(level, log.CallerStackHandler("%+v", fileHandler))
		handler := ErrorMultiHandler(normalHandler, errorHandler)
		l.SetHandler(handler)
	}

	return nil
}

func (l *Logger) Debug(v ...interface{}) {
	l.output(func(l *log.Logger, msg string, v ...interface{}) {
		(*l).Debug(msg, v...)
	}, v...)
}

func (l *Logger) Info(v ...interface{}) {
	l.output(func(l *log.Logger, msg string, v ...interface{}) {
		(*l).Info(msg, v...)
	}, v...)
}

func (l *Logger) Warn(v ...interface{}) {
	l.output(func(l *log.Logger, msg string, v ...interface{}) {
		(*l).Warn(msg, v...)
	}, v...)
}

func (l *Logger) Error(v ...interface{}) {
	l.output(func(l *log.Logger, msg string, v ...interface{}) {
		(*l).Error(msg, v...)
	}, v...)
}

func (l *Logger) Crit(v ...interface{}) {
	l.output(func(l *log.Logger, msg string, v ...interface{}) {
		(*l).Crit(msg, v...)
	}, v...)
}

func ErrorMultiHandler(normalHandler, errorHandler log.Handler) log.Handler {
	return log.FuncHandler(func(r *log.Record) error {
		if len(r.Ctx) > 1 {
			_, ok := r.Ctx[1].(error)
			if ok {
				r.Ctx = r.Ctx[2:]
				errorHandler.Log(r)
			} else {
				normalHandler.Log(r)
			}
		} else {
			normalHandler.Log(r)
		}
		return nil
	})
}

func (l *Logger) Debugf(format string, v ...interface{}) { l.Debug(fmt.Sprintf(format, v...)) }
func (l *Logger) Infof(format string, v ...interface{})  { l.Info(fmt.Sprintf(format, v...)) }
func (l *Logger) Warnf(format string, v ...interface{})  { l.Warn(fmt.Sprintf(format, v...)) }
func (l *Logger) Errorf(format string, v ...interface{}) { l.Error(fmt.Sprintf(format, v...)) }
func (l *Logger) Critf(format string, v ...interface{})  { l.Crit(fmt.Sprintf(format, v...)) }

func (l *Logger) DebugIndent(indentation int, v ...interface{}) {
	l.outputWithIndentation(func(l *Logger, v ...interface{}) {
		l.Debug(v...)
	}, indentation, v...)
}

func (l *Logger) InfoIndent(indentation int, v ...interface{}) {
	l.outputWithIndentation(func(l *Logger, v ...interface{}) {
		l.Info(v...)
	}, indentation, v...)
}

func (l *Logger) WarnIndent(indentation int, v ...interface{}) {
	l.outputWithIndentation(func(l *Logger, v ...interface{}) {
		l.Warn(v...)
	}, indentation, v...)
}

func (l *Logger) ErrorIndent(indentation int, v ...interface{}) {
	l.outputWithIndentation(func(l *Logger, v ...interface{}) {
		l.Error(v...)
	}, indentation, v...)
}

func (l *Logger) CritIndent(indentation int, v ...interface{}) {
	l.outputWithIndentation(func(l *Logger, v ...interface{}) {
		l.Crit(v...)
	}, indentation, v...)
}

func (l *Logger) IncreaseIndentation() {
	contextIndentations[l.Context] += 1
}

func (l *Logger) DecreaseIndentation() {
	contextIndentations[l.Context] -= 1
}

func (l *Logger) output(logFunc func(l *log.Logger, msg string, v ...interface{}), v ...interface{}) {
	err, ok := v[0].(error)
	if ok {
		logFunc(&l.Logger, l.currentIndentation()+err.Error(), "err", err)
	} else {
		msg := v[0].(string)
		if len(v) > 1 {
			logFunc(&l.Logger, l.currentIndentation()+msg, v[1:]...)
		} else {
			logFunc(&l.Logger, l.currentIndentation()+msg)
		}
	}
}

func (l *Logger) outputWithIndentation(logFunc func(l *Logger, v ...interface{}), indentation int, v ...interface{}) {
	oldIndentation := contextIndentations[l.Context]
	contextIndentations[l.Context] = indentation
	logFunc(l, v...)
	contextIndentations[l.Context] = oldIndentation
}
func (l *Logger) currentIndentation() string {
	return Indentation(contextIndentations[l.Context])
}

func Indentation(tabs int) string {
	return strings.Repeat(" ", tabWidth*tabs)
}

func Initialize(logOut string, lvl interface{}) error {
	if Log == nil {
		Log = &Logger{
			Logger: log.New(),
		}
	}
	return Log.InitializeWithContext("main", logOut, lvl)
}

func Debug(v ...interface{}) { Log.Debug(v...) }
func Info(v ...interface{})  { Log.Info(v...) }
func Warn(v ...interface{})  { Log.Warn(v...) }
func Error(v ...interface{}) { Log.Error(v...) }
func Crit(v ...interface{})  { Log.Crit(v...) }

func Debugf(format string, v ...interface{}) { Log.Debugf(format, v...) }
func Infof(format string, v ...interface{})  { Log.Infof(format, v...) }
func Warnf(format string, v ...interface{})  { Log.Warnf(format, v...) }
func Errorf(format string, v ...interface{}) { Log.Errorf(format, v...) }
func Critf(format string, v ...interface{})  { Log.Critf(format, v...) }

func IncreaseIndentation() { Log.IncreaseIndentation() }
func DecreaseIndentation() { Log.DecreaseIndentation() }
