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

func NewWriter(ws ...io.Writer) io.Writer {
	switch len(ws) {
	case 0:
		return os.Stdout
	case 1:
		return ws[0]
	default:
		mw := MultiWriter{}
		mw = append(mw, ws...)
		return mw
	}
}

type MultiWriter []io.Writer

type errs []error

func (err *errs) Error() string {
	switch len(*err) {
	case 0:
		return "<nil>"
	case 1:
		return (*err)[0].Error()
	default:
		bs := bpl.Get().(*buf)
		bs.buf = bs.buf[:0]
		bs.buf = append(bs.buf, []byte("multi error\n")...)
		for _, v := range *err {
			bs.buf = append(bs.buf, []byte(" - ")...)
			bs.buf = append(bs.buf, []byte(v.Error())...)
			bs.buf = append(bs.buf, '\n')
		}
		ret := string(bs.buf)
		bpl.Put(bs)
		return ret
	}
}

func (err *errs) Append(e error) {
	if err == nil {
		err = &errs{}
	}
	if e != nil {
		*err = append(*err, e)
	}
}

func (mw MultiWriter) Write(p []byte) (n int, err error) {
	var me *errs
	for i := range mw {
		n, err = mw[i].Write(p)
		me.Append(err)
	}
	return n, me
}
