package logz

import (
	"encoding/json"
	"fmt"
	"io"
	_log "log"
	"os"
	"path"
	"runtime"
	"strings"
	"sync"
)

type LogzLevel int

const (
	Debug LogzLevel = iota
	Info
	Warn
	Error
	Fatal
)

type Log struct {
	dest     io.Writer
	_package string
	Level    LogzLevel
	_info    *_log.Logger
	_warn    *_log.Logger
	_erro    *_log.Logger
	_fatl    *_log.Logger
	_debg    *_log.Logger
	m        sync.Mutex
}

func DefaultLog(pkg string) *Log {
	return NewLog(os.Stdout, pkg, Warn)
}

func NewLog(d io.Writer, pkg string, level LogzLevel) *Log {
	prefix := fmt.Sprintf("[%s] ", pkg)
	flags := _log.LstdFlags | _log.Lmicroseconds | _log.Lmsgprefix
	l := Log{
		dest:     d,
		_package: prefix,
		Level:    level,
	}
	l._debg = _log.New(d, "[DEBUG] "+prefix, flags)
	l._info = _log.New(d, "[INFO ] "+prefix, flags)
	l._warn = _log.New(d, "[WARN ] "+prefix, flags)
	l._erro = _log.New(d, "[ERROR] "+prefix, flags)
	l._fatl = _log.New(d, "[FATAL] "+prefix, flags)

	return &l
}

func fmtCaller() string {
	_, fn, line, _ := runtime.Caller(3)
	caller := fmt.Sprintf("%s:%d:", path.Base(fn), line)
	return caller
}

func jsonify(i interface{}) string {
	b, err := json.Marshal(i)
	if err != nil {
		return " <msg>{\"error\": 1}</msg>"
	}
	return fmt.Sprintf(" | <msg>%s</msg>", string(b))
}

func combine(args ...interface{}) string {
	if len(args) == 1 {
		return fmt.Sprint(args[0])
	} else {
		var el []string
		for _, v := range args {
			el = append(el, fmt.Sprint(v))
		}
		x := strings.Join(el, " ")
		fmt.Println(x)
		return strings.Join(el, " ")
	}
}

func (l *Log) call(lvl LogzLevel, msg string) {
	l.m.Lock()
	defer l.m.Unlock()
	caller := fmtCaller()
	if lvl == Debug && l.Level <= Debug {
		l._debg.Println(caller, msg)
	} else if lvl == Info && l.Level <= Info {
		l._info.Println(caller, msg)
	} else if lvl == Warn && l.Level <= Warn {
		l._warn.Println(caller, msg)
	} else if lvl == Error && l.Level <= Error {
		l._erro.Panic(caller, " ", msg)
	} else if lvl == Fatal && l.Level <= Fatal {
		l._fatl.Fatal(caller, " ", msg)
	}
}

func (l *Log) Debug(args ...interface{}) {
	l.call(Debug, combine(args...))
}

func (l *Log) DebugMsg(st interface{}, args ...interface{}) {
	l.call(Debug, combine(args...)+jsonify(st))
}

func (l *Log) Info(args ...interface{}) {
	l.call(Info, combine(args...))
}

func (l *Log) InfoMsg(st interface{}, args ...interface{}) {
	l.call(Info, combine(args...)+jsonify(st))
}

func (l *Log) Warn(args ...interface{}) {
	l.call(Warn, combine(args...))
}

func (l *Log) WarnMsg(st interface{}, args ...interface{}) {
	l.call(Warn, combine(args...)+jsonify(st))
}

func (l *Log) Error(args ...interface{}) {
	l.call(Error, combine(args...))
}

func (l *Log) ErrorMsg(st interface{}, args ...interface{}) {
	l.call(Error, combine(args...)+jsonify(st))
}

func (l *Log) Fatal(args ...interface{}) {
	l.call(Fatal, combine(args...))
}

func (l *Log) FatalMsg(st interface{}, args ...interface{}) {
	l.call(Fatal, combine(args...)+jsonify(st))
}
