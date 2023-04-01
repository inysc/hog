package qog

import (
	"io"
	"os"
	"runtime"
	"time"
)

type logger struct {
	io.Writer
	lvl Level
}

func (l *logger) newEvent(lvl Level) Event {
	if lvl >= l.lvl {
		e := getevent()
		e.Writer = l.Writer
		appendTime(e.Buffer, time.Now())

		switch lvl {
		case TRACE:
			e.Buffer.WriteString("|TRACE|")
		case DEBUG:
			e.Buffer.WriteString("|DEBUG|")
		case INFO:
			e.Buffer.WriteString("|INFO|")
		case WARN:
			e.Buffer.WriteString("|WARN|")
		case ERROR:
			e.Buffer.WriteString("|ERROR|")
		case FATAL:
			e.Buffer.WriteString("|FATAL|")
		}

		// 写入调用信息
		appendCaller(e.Buffer)

		// 写入 goid
		e.WriteString("|goid:")
		e.Buffer.Write(transNum(runtime.Goid()))
		e.WriteByte('|')
		e.WriteByte(',')

		return e
	}
	return nilevent{}
}

func (l *logger) Trace() Event { return l.newEvent(TRACE) }
func (l *logger) Debug() Event { return l.newEvent(DEBUG) }
func (l *logger) Info() Event  { return l.newEvent(INFO) }
func (l *logger) Warn() Event  { return l.newEvent(WARN) }
func (l *logger) Error() Event { return l.newEvent(ERROR) }
func (l *logger) Fatal() Event { return l.newEvent(FATAL) }

func (l *logger) SetLevel(lvl Level) { l.lvl = lvl }

// needStd 是否需要同时输出到标准输出（仅在 ppid != 1 时生效 ）
func New(needStd bool, lvl Level, ws ...io.Writer) *logger {
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
	return &logger{w, lvl}
}

func Simple(filename string) *logger {
	return New(true, DEBUG, &LoggerFile{
		Filename:   filename,
		MaxSize:    100,
		MaxAge:     30,
		MaxBackups: 11,
		LocalTime:  false,
		Compress:   false,
	})
}
