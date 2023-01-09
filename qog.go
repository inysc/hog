package qog

import (
	"bytes"
	"context"
	"os"
	"runtime"
	"sync"
	"time"
)

type (
	Level  = uint8
	bufnum struct{ buf [22]byte }
)

var (
	bpl     = sync.Pool{New: func() any { return bytes.NewBuffer(make([]byte, 0, 1024)) }}
	numpl   = sync.Pool{New: func() any { return new(bufnum) }}
	eventpl = sync.Pool{New: func() any { return &event{bytes.NewBuffer(make([]byte, 0, 1024))} }}
)

type Event interface {
	Any(string, any) Event
	Error(string, error) Event
	String(string, string) Event
	Stringp(string, *string) Event
	Duration(string, time.Duration) Event
	Durationp(string, *time.Duration) Event
	Bool(string, bool) Event
	Boolp(string, *bool) Event
	Int(string, int) Event
	Intp(string, *int) Event
	Int8(string, int8) Event
	Int8p(string, *int8) Event
	Int16(string, int16) Event
	Int16p(string, *int16) Event
	Int32(string, int32) Event
	Int32p(string, *int32) Event
	Int64(string, int64) Event
	Int64p(string, *int64) Event
	Uint(string, uint) Event
	Uintp(string, *uint) Event
	Uint8(string, uint8) Event
	Uint8p(string, *uint8) Event
	Uint16(string, uint16) Event
	Uint16p(string, *uint16) Event
	Uint32(string, uint32) Event
	Uint32p(string, *uint32) Event
	Uint64(string, uint64) Event
	Uint64p(string, *uint64) Event
	Float32(string, float32) Event
	Float32p(string, *float32) Event
	Float64(string, float64) Event
	Float64p(string, *float64) Event
	Msg(string)
	Msgf(string, ...any)
}

const (
	TRACE Level = iota
	DEBUG
	INFO
	WARN
	ERROR
	FATAL
)

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

func appendCaller(bf *bytes.Buffer) {
	pc, file, line, ok := runtime.Caller(4)
	if ok {
		var a, b, c = 0, 0, 0
		for i := 0; i < len(file); i++ {
			if file[i] == '/' {
				a = b
				b = i
			}
		}
		bf.WriteString(file[a+1:])
		funcName := runtime.FuncForPC(pc).Name()
		for i := 0; i < len(funcName); i++ {
			if funcName[i] == '.' {
				c = i
			}
		}
		bf.WriteString(funcName[c:])
		bf.WriteByte(':')
		bf.Write(transNum(line))
	} else {
		bf.WriteString("???:??")
	}
}

// [0, 99]
func tinyNum[T ~int](bf *bytes.Buffer, num T) {
	if num < 10 {
		bf.WriteByte('0')
		bf.WriteByte(smallsString[num*2+1])
	} else {
		num *= 2
		bf.WriteByte(smallsString[num+0])
		bf.WriteByte(smallsString[num+1])
	}
}

func transNum[T ~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64](num T) []byte {
	var to = numpl.Get().(*bufnum)
	idx := 22
	for num >= 100 {
		is := num % 100 * 2
		num /= 100
		idx -= 2
		to.buf[idx+1] = smallsString[is+1]
		to.buf[idx+0] = smallsString[is+0]
	}
	// us < 100
	is := num * 2
	idx--
	to.buf[idx] = smallsString[is+1]
	if num >= 10 {
		idx--
		to.buf[idx] = smallsString[is]
	}
	defer numpl.Put(to)
	return to.buf[idx:]
}

func appendTime(b *bytes.Buffer, now time.Time) {

	Y, M, D := now.Date()
	h, m, s := now.Clock()

	b.Write(transNum(Y))
	b.WriteByte('-')
	tinyNum(b, M)
	b.WriteByte('-')
	tinyNum(b, D)
	b.WriteByte(' ')

	tinyNum(b, h)
	b.WriteByte(':')
	tinyNum(b, m)
	b.WriteByte(':')
	tinyNum(b, s)
	b.WriteByte('.')

	nano := now.Nanosecond()
	if nano < 10 {
		b.WriteByte(smallsString[nano*2+1])
		b.WriteByte('0')
		b.WriteByte('0')
	} else if nano < 100 {
		nano *= 2
		b.WriteByte(smallsString[nano+1])
		b.WriteByte(smallsString[nano+1])
		b.WriteByte('0')
	} else {
		b.Write(transNum(nano)[:3])
	}
}

// ************* TRACE ****************
func (l *logger) Trace(msg string, args ...interface{}) {
	if l.lvl > TRACE {
		return
	}
	l.write(context.TODO(), "|TRACE|", msg, args)
}

func (l *logger) TraceT(ctx context.Context, format string, args ...any) {
	if l.lvl > TRACE {
		return
	}
	l.write(ctx, "|TRACE|", format, args)
}

// ************* DEBUG ****************
func (l *logger) Debug(msg string, args ...interface{}) {
	if l.lvl > DEBUG {
		return
	}
	l.write(context.TODO(), "|DEBUG|", msg, args)
}

func (l *logger) DebugT(ctx context.Context, format string, args ...any) {
	if l.lvl > DEBUG {
		return
	}
	l.write(ctx, "|DEBUG|", format, args)
}

// ************* INFO ****************
func (l *logger) Info(msg string, args ...interface{}) {
	if l.lvl > INFO {
		return
	}
	l.write(context.TODO(), "|INFO |", msg, args)
}

func (l *logger) InfoT(ctx context.Context, format string, args ...any) {
	if l.lvl > INFO {
		return
	}
	l.write(ctx, "|INFO |", format, args)
}

// ************* Warn ****************
func (l *logger) Warn(msg string, args ...interface{}) {
	if l.lvl > WARN {
		return
	}
	l.write(context.TODO(), "|WARN |", msg, args)
}

func (l *logger) WarnT(ctx context.Context, format string, args ...any) {
	if l.lvl > WARN {
		return
	}
	l.write(ctx, "|WARN |", format, args)
}

// ************* ERROR ****************
func (l *logger) Error(msg string, args ...interface{}) {
	if l.lvl > ERROR {
		return
	}
	l.write(context.TODO(), "|ERROR|", msg, args)
}

func (l *logger) ErrorT(ctx context.Context, format string, args ...any) {
	if l.lvl > ERROR {
		return
	}
	l.write(ctx, "|ERROR|", format, args)
}

// ************* FATAL ****************
func (l *logger) Fatal(msg string, args ...interface{}) {
	if l.lvl > FATAL {
		return
	}
	l.write(context.TODO(), "|ERROR|", msg, args)
	os.Exit(1)
}

func (l *logger) FatalT(ctx context.Context, format string, args ...any) {
	if l.lvl > FATAL {
		return
	}
	l.write(ctx, "|FATAL|", format, args)
	os.Exit(1)
}
