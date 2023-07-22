package hog

import (
	"io"
	"os"
	"time"
)

type Logger interface {
	Trace() Event
	Debug() Event
	Info() Event
	Warn() Event
	Error() Event
	Fatal() Event
	Panic() Event

	Op() Event

	SetLevel(Level)
	AddSkip(int)
}

type nillogger struct{}

var DefaultLogger Logger = nillogger{}

type logger struct {
	io.Writer
	lvl  Level
	skip int
}

func (l *logger) newEvent(lvl Level, flag bool) Event {
	if lvl >= l.lvl {
		e := getevent()
		e.lvl = lvl
		e.Writer = l.Writer
		if !flag {
			appendTime(e.Buffer, time.Now())

			switch lvl {
			case TRACE:
				e.Buffer.WriteString(" TRACE ")
			case DEBUG:
				e.Buffer.WriteString(" DEBUG ")
			case INFO:
				e.Buffer.WriteString(" INFO ")
			case WARN:
				e.Buffer.WriteString(" WARN ")
			case ERROR:
				e.Buffer.WriteString(" ERROR ")
			case FATAL:
				e.Buffer.WriteString(" FATAL ")
			}

			// 写入调用信息
			appendCaller(e.Buffer, l.skip)

			// 写入 goid
			// e.WriteString("|goid:")
			// e.Buffer.Write(transNum(runtime.Goid()))
			e.WriteByte(' ')
		}
		return e
	}
	return nilevent{}
}

func (l *logger) Trace() Event { return l.newEvent(TRACE, false) }
func (l *logger) Debug() Event { return l.newEvent(DEBUG, false) }
func (l *logger) Info() Event  { return l.newEvent(INFO, false) }
func (l *logger) Warn() Event  { return l.newEvent(WARN, false) }
func (l *logger) Error() Event { return l.newEvent(ERROR, false) }
func (l *logger) Fatal() Event { return l.newEvent(FATAL, false) }
func (l *logger) Panic() Event { return l.newEvent(PANIC, false) }
func (l *logger) Op() Event    { return l.newEvent(OP, true) }

func (l *logger) SetLevel(lvl Level) { l.lvl = lvl }
func (l *logger) AddSkip(skip int)   { l.skip += skip }

func New(lvl Level, ws ...io.Writer) Logger {
	var w io.Writer
	switch len(ws) {
	case 0:
		w = os.Stdout
	case 1:
		w = ws[0]
	default:
		w = io.MultiWriter(ws...)
	}
	return &logger{w, lvl, 3}
}

func Simple(filename string) Logger {
	return New(DEBUG, &LoggerFile{
		Filename:   filename,
		MaxSize:    100,
		MaxAge:     30,
		MaxBackups: 11,
		LocalTime:  false,
		Compress:   false,
	})
}

func (nillogger) Trace() Event   { return nilevent{} }
func (nillogger) Debug() Event   { return nilevent{} }
func (nillogger) Info() Event    { return nilevent{} }
func (nillogger) Warn() Event    { return nilevent{} }
func (nillogger) Error() Event   { return nilevent{} }
func (nillogger) Fatal() Event   { return fpevent{} }
func (nillogger) Panic() Event   { return fpevent{} }
func (nillogger) Op() Event      { return nilevent{} }
func (nillogger) SetLevel(Level) {}
func (nillogger) AddSkip(int)    {}
