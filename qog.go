package qog

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"sync"
	"time"
)

type (
	Level uint8
	buf   struct{ buf []byte }
	bufTo struct{ buf [22]byte }
)

var (
	bpl  = &sync.Pool{New: func() any { return &buf{buf: make([]byte, maxLen)} }}
	toPl = &sync.Pool{New: func() any { return &bufTo{} }}
)

const (
	TRACE Level = iota
	DEBUG
	INFO
	WARN
	ERROR
)

type Logger struct {
	writer io.Writer
	lvl    Level
	name   string
	trace  string
	debug  string
	info   string
	warn   string
	err    string
}

func (l *Logger) write(info string, msg string) (n int, err error) {
	bs := bpl.Get().(*buf)
	bs.buf = bs.buf[:0]

	// 写入时间戳
	appendFormat(bs)

	// 写入等级及服务名
	bs.buf = append(bs.buf, info...)

	// 写入调用信息
	appendCaller(bs)

	// 写入 goid
	bs.buf = append(bs.buf, "goid:"...)
	bs.buf = append(bs.buf, transNum(runtime.Goid(), true)...)

	// 写入日志信息
	bs.buf = append(bs.buf, msg...)
	bs.buf = append(bs.buf, '\n')

	// 写入到文件
	n, err = l.writer.Write(bs.buf)

	// 太大就抛弃了
	if len(bs.buf) < maxLen<<2 {
		bpl.Put(bs)
	}
	return
}

func (l *Logger) SetLevel(lvl Level) {
	l.lvl = lvl
}

func New(name string, lvl Level, ws ...io.Writer) (l *Logger) {
	name = name + "|"
	var w io.Writer
	switch len(ws) {
	case 0:
		w = os.Stdout
	default:
		w = io.MultiWriter(ws...)
	}
	return &Logger{
		writer: w,
		lvl:    lvl,
		name:   name,
		trace:  name + "TRACE|",
		debug:  name + "DEBUG|",
		info:   name + "INFO |",
		warn:   name + "WARN |",
		err:    name + "ERROR|",
	}
}

const maxLen = 1024

const smallsString = "00010203040506070809" +
	"10111213141516171819" +
	"20212223242526272829" +
	"30313233343536373839" +
	"40414243444546474849" +
	"50515253545556575859" +
	"60616263646566676869" +
	"70717273747576777879" +
	"80818283848586878889" +
	"90919293949596979899"

func appendCaller(bf *buf) {
	pc, file, line, ok := runtime.Caller(3)
	if ok {
		var a, b, c = 0, 0, 0
		for i := 0; i < len(file); i++ {
			if file[i] == '/' {
				a = b
				b = i
			}
		}
		bf.buf = append(bf.buf, file[a+1:]...)
		funcName := runtime.FuncForPC(pc).Name()
		for i := 0; i < len(funcName); i++ {
			if funcName[i] == '.' {
				c = i
			}
		}
		bf.buf = append(bf.buf, funcName[c:]...)
		bf.buf = append(bf.buf, ':')
		bf.buf = append(bf.buf, transNum(line, true)...)
	}
}

func transNum[T int | int64](num T, needSep bool) []byte {
	var to = toPl.Get().(*bufTo) // +1 for sign of 64bit value in base 2
	i := 22
	if needSep {
		i--
		to.buf[21] = '|'
	}
	for num >= 100 {
		is := num % 100 * 2
		num /= 100
		i -= 2
		to.buf[i+1] = smallsString[is+1]
		to.buf[i+0] = smallsString[is+0]
	}
	// us < 100
	is := num * 2
	i--
	to.buf[i] = smallsString[is+1]
	if num >= 10 {
		i--
		to.buf[i] = smallsString[is]
	}
	defer toPl.Put(to)
	return to.buf[i:]
}

func appendFormat(b *buf) {
	now := time.Now()

	b.buf = append(b.buf, transNum(now.Year(), false)...)
	b.buf = append(b.buf, '-')
	tinyNum(b, now.Month())
	b.buf = append(b.buf, '-')
	tinyNum(b, now.Day())
	b.buf = append(b.buf, ' ')
	tinyNum(b, now.Hour())
	b.buf = append(b.buf, ':')
	tinyNum(b, now.Minute())
	b.buf = append(b.buf, ':')
	tinyNum(b, now.Second())
	b.buf = append(b.buf, '.')
	nano := now.Nanosecond()
	if nano < 10 {
		b.buf = append(b.buf, smallsString[nano*2+1])
		b.buf = append(b.buf, '0')
		b.buf = append(b.buf, '0')
	} else if nano < 100 {
		nano *= 2
		b.buf = append(b.buf, smallsString[nano+1])
		b.buf = append(b.buf, smallsString[nano+1])
		b.buf = append(b.buf, '0')
	} else {
		b.buf = append(b.buf, transNum(nano, false)[:3]...)
	}
	b.buf = append(b.buf, '|')
}

// [0, 99]
func tinyNum[T ~int](b *buf, num T) {
	if num < 10 {
		b.buf = append(b.buf, '0')
		b.buf = append(b.buf, smallsString[num*2+1])
	} else {
		num *= 2
		b.buf = append(b.buf, smallsString[num+0])
		b.buf = append(b.buf, smallsString[num+1])
	}
}

// ************* TRACE ****************
func (l *Logger) Trace(msg string) {
	if l.lvl > TRACE {
		return
	}
	l.write(l.trace, msg)
}

func (l *Logger) Tracef(format string, args ...interface{}) {
	if l.lvl > TRACE {
		return
	}
	msg := fmt.Sprintf(format, args...)
	l.write(l.trace, msg)
}

// ************* DEBUG ****************
func (l *Logger) Debug(msg string) {
	if l.lvl > DEBUG {
		return
	}
	l.write(l.debug, msg)
}

func (l *Logger) Debugf(format string, args ...interface{}) {
	if l.lvl > DEBUG {
		return
	}
	msg := fmt.Sprintf(format, args...)
	l.write(l.debug, msg)
}

// ************* INFO ****************
func (l *Logger) Info(msg string) {
	if l.lvl > INFO {
		return
	}
	l.write(l.info, msg)
}

func (l *Logger) Infof(format string, args ...interface{}) {
	if l.lvl > INFO {
		return
	}
	msg := fmt.Sprintf(format, args...)
	l.write(l.info, msg)
}

// ************* Warn ****************
func (l *Logger) Warn(msg string) {
	if l.lvl > WARN {
		return
	}
	l.write(l.warn, msg)
}

func (l *Logger) Warnf(format string, args ...interface{}) {
	if l.lvl > WARN {
		return
	}
	l.write(l.warn, fmt.Sprintf(format, args...))
}

// ************* ERROR ****************
func (l *Logger) Error(msg string) {
	if l.lvl > ERROR {
		return
	}
	l.write(l.err, msg)
}

func (l *Logger) Errorf(format string, args ...interface{}) {
	if l.lvl > ERROR {
		return
	}
	l.write(l.err, fmt.Sprintf(format, args...))
}
