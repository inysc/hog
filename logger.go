package hog

import (
	"io"
	"os"
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

func (l *logger) NewEvent(skip int, lvl Level, flag bool) Event {
	if lvl >= l.lvl {
		e := getevent()
		e.Init(flag, l.skip, lvl, l.Writer)

		return e
	}
	return nilevent{}
}

func (l *logger) Trace() Event { return l.NewEvent(l.skip, TRACE, false) }
func (l *logger) Debug() Event { return l.NewEvent(l.skip, DEBUG, false) }
func (l *logger) Info() Event  { return l.NewEvent(l.skip, INFO, false) }
func (l *logger) Warn() Event  { return l.NewEvent(l.skip, WARN, false) }
func (l *logger) Error() Event { return l.NewEvent(l.skip, ERROR, false) }
func (l *logger) Fatal() Event { return l.NewEvent(l.skip, FATAL, false) }
func (l *logger) Panic() Event { return l.NewEvent(l.skip, PANIC, false) }
func (l *logger) Op() Event    { return l.NewEvent(l.skip, OP, true) }

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
