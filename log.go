package qog

import (
	"io"
	"os"
	"runtime"
	"time"
)

type Level uint8

const (
	TRACE Level = iota
	DEBUG
	INFO
	WARN
	ERROR
)

type Logger struct {
	io.Writer
	lvl   Level
	name  string
	trace string
	debug string
	info  string
	warn  string
	err   string
}

func (l *Logger) write(info string, msg string) (n int, err error) {
	bs := bpl.Get().(*buf)
	bs.buf = bs.buf[:0]

	// 写入时间戳
	bs.buf = time.Now().AppendFormat(bs.buf, "[2006-01-02 15:04:05.000]") // 写入时间

	// 写入等级及服务名
	bs.buf = append(bs.buf, info...)

	// 写入调用信息
	appendCaller(bs)

	// 写入 goid
	bs.buf = append(bs.buf, "[goid:"...)
	appendNum(bs, runtime.Goid())

	// 写入日志信息
	bs.buf = append(bs.buf, msg...)
	bs.buf = append(bs.buf, '\n')

	// 写入
	n, err = l.Write(bs.buf)
	if len(bs.buf) < maxLen<<2 { // 太大就抛弃了
		bpl.Put(bs)
	}
	return
}

func (l *Logger) SetLevel(lvl Level) {
	l.lvl = lvl
}

func New(name string, lvl Level, w io.Writer) (l *Logger) {
	name = "[" + name + "]"
	if w == nil {
		w = os.Stdout
	}
	return &Logger{
		Writer: w,
		lvl:    lvl,
		name:   name,
		trace:  name + "[TRACE][",
		debug:  name + "[DEBUG][",
		info:   name + "[INFO][",
		warn:   name + "[WARN][",
		err:    name + "[ERROR][",
	}
}
