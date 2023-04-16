package qog

import (
	"io"
	"os"
	"runtime"
	"time"
)

type Logger interface {
	Trace(...bool) Event
	Debug(...bool) Event
	Info(...bool) Event
	Warn(...bool) Event
	Error(...bool) Event
	Fatal(...bool) Event

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

func (l *logger) newEvent(lvl Level, flag []bool) Event {
	if lvl >= l.lvl {
		e := getevent()
		e.Writer = l.Writer
		if len(flag) == 0 || !flag[0] {
			appendTime(e.Buffer, time.Now())

			switch lvl {
			case TRACE:
				e.Buffer.WriteString("|TRACE|")
			case DEBUG:
				e.Buffer.WriteString("|DEBUG|")
			case INFO:
				e.Buffer.WriteString("|INFO |")
			case WARN:
				e.Buffer.WriteString("|WARN |")
			case ERROR:
				e.Buffer.WriteString("|ERROR|")
			case FATAL:
				e.Buffer.WriteString("|FATAL|")
			}

			// 写入调用信息
			appendCaller(e.Buffer, l.skip)

			// 写入 goid
			e.WriteString("|goid:")
			e.Buffer.Write(transNum(runtime.Goid()))
			e.WriteByte('|')
			e.WriteByte(',')
		}
		return e
	}
	return nilevent{}
}

func (l *logger) Trace(trimperfix ...bool) Event { return l.newEvent(TRACE, trimperfix) }
func (l *logger) Debug(trimperfix ...bool) Event { return l.newEvent(DEBUG, trimperfix) }
func (l *logger) Info(trimperfix ...bool) Event  { return l.newEvent(INFO, trimperfix) }
func (l *logger) Warn(trimperfix ...bool) Event  { return l.newEvent(WARN, trimperfix) }
func (l *logger) Error(trimperfix ...bool) Event { return l.newEvent(ERROR, trimperfix) }
func (l *logger) Fatal(trimperfix ...bool) Event { return l.newEvent(FATAL, trimperfix) }

func (l *logger) SetLevel(lvl Level) { l.lvl = lvl }
func (l *logger) AddSkip(skip int)   { l.skip += skip }

// needStd 是否需要同时输出到标准输出（仅在 ppid != 1 时生效 ）
func New(needStd bool, lvl Level, ws ...io.Writer) Logger {
	var w io.Writer
	switch len(ws) {
	case 0:
		w = os.Stdout
	default:
		if needStd && os.Getppid() != 1 {
			ws = append(ws, os.Stdout)
		}
		if len(ws) == 1 {
			w = ws[0]
		} else {
			w = io.MultiWriter(ws...)
		}
	}
	return &logger{w, lvl, 3}
}

func Simple(filename string) Logger {
	return New(true, DEBUG, &LoggerFile{
		Filename:   filename,
		MaxSize:    100,
		MaxAge:     30,
		MaxBackups: 11,
		LocalTime:  false,
		Compress:   false,
	})
}

func (nillogger) Trace(...bool) Event { return nilevent{} }
func (nillogger) Debug(...bool) Event { return nilevent{} }
func (nillogger) Info(...bool) Event  { return nilevent{} }
func (nillogger) Warn(...bool) Event  { return nilevent{} }
func (nillogger) Error(...bool) Event { return nilevent{} }
func (nillogger) Fatal(...bool) Event { return nilevent{} }
func (nillogger) SetLevel(Level)      {}
func (nillogger) AddSkip(int)         {}
