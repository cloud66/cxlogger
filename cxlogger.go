package cxlogger

// see http://godoc.org/github.com/inconshreveable/log15 for more info
import (
	"fmt"

	log "gopkg.in/inconshreveable/log15.v2"
)

// func main() {
// 	Initialize("STDOUT", "debug")
// 	Debug("test!")
// }

type Logger struct {
	log.Logger
	Level log.Lvl
}

const (
	LvlCrit  = log.LvlCrit
	LvlError = log.LvlError
	LvlWarn  = log.LvlWarn
	LvlInfo  = log.LvlInfo
	LvlDebug = log.LvlDebug
)

var Log *Logger

func Initialize(logOut string, lvl interface{}) error {
	var (
		level log.Lvl
		err   error
	)

	if str, ok := lvl.(string); ok {
		level, err = log.LvlFromString(str)
		if err != nil {
			return err
		}
	} else {
		level = lvl.(log.Lvl)
	}
	Log = &Logger{log.New(), level}

	if logOut == "STDOUT" {
		normalHandler := log.LvlFilterHandler(level, log.StdoutHandler)
		errorHandler := log.LvlFilterHandler(level, log.CallerStackHandler("%+v", log.StdoutHandler))
		handler := ErrorMultiHandler(normalHandler, errorHandler)
		Log.SetHandler(handler)
	} else if logOut == "NONE" {
		Log.SetHandler(log.DiscardHandler())
	} else {
		fileHandler := log.Must.FileHandler(logOut, log.LogfmtFormat())
		normalHandler := log.LvlFilterHandler(level, fileHandler)
		errorHandler := log.LvlFilterHandler(level, log.CallerStackHandler("%+v", fileHandler))
		handler := ErrorMultiHandler(normalHandler, errorHandler)
		Log.SetHandler(handler)
	}

	return nil
}

func (l *Logger) Debug(v ...interface{}) {
	if interfaces, ok := v[0].([]interface{}); ok {
		v = interfaces
	}
	err, ok := v[0].(error)
	if ok {
		l.Logger.Debug(err.Error(), "err", err)
	} else {
		msg := v[0].(string)
		if len(v) > 1 {
			l.Logger.Debug(msg, v[1:])
		} else {
			l.Logger.Debug(msg)
		}
	}
}

func (l *Logger) Info(v ...interface{}) {
	if interfaces, ok := v[0].([]interface{}); ok {
		v = interfaces
	}
	err, ok := v[0].(error)
	if ok {
		l.Logger.Info(err.Error(), "err", err)
	} else {
		msg := v[0].(string)
		if len(v) > 1 {
			l.Logger.Info(msg, v[1:])
		} else {
			l.Logger.Info(msg)
		}
	}
}

func (l *Logger) Warn(v ...interface{}) {
	if interfaces, ok := v[0].([]interface{}); ok {
		v = interfaces
	}
	err, ok := v[0].(error)
	if ok {
		l.Logger.Warn(err.Error(), "err", err)
	} else {
		msg := v[0].(string)
		if len(v) > 1 {
			l.Logger.Warn(msg, v[1:])
		} else {
			l.Logger.Warn(msg)
		}
	}
}

func (l *Logger) Error(v ...interface{}) {
	if interfaces, ok := v[0].([]interface{}); ok {
		v = interfaces
	}
	err, ok := v[0].(error)
	if ok {
		l.Logger.Error(err.Error(), "err", err)
	} else {
		msg := v[0].(string)
		if len(v) > 1 {
			l.Logger.Error(msg, v[1:])
		} else {
			l.Logger.Error(msg)
		}
	}
}

func (l *Logger) Crit(v ...interface{}) {
	if interfaces, ok := v[0].([]interface{}); ok {
		v = interfaces
	}
	err, ok := v[0].(error)
	if ok {
		l.Logger.Crit(err.Error(), "err", err)
	} else {
		msg := v[0].(string)
		if len(v) > 1 {
			l.Logger.Crit(msg, v[1:])
		} else {
			l.Logger.Crit(msg)
		}
	}
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

func Debug(v ...interface{}) { Log.Debug(v) }
func Info(v ...interface{})  { Log.Info(v) }
func Warn(v ...interface{})  { Log.Warn(v) }
func Error(v ...interface{}) { Log.Error(v) }
func Crit(v ...interface{})  { Log.Crit(v) }

func Debugf(format string, v ...interface{}) { Log.Debugf(format, v) }
func Infof(format string, v ...interface{})  { Log.Infof(format, v) }
func Warnf(format string, v ...interface{})  { Log.Warnf(format, v) }
func Errorf(format string, v ...interface{}) { Log.Errorf(format, v) }
func Critf(format string, v ...interface{})  { Log.Crit(format, v) }
