package hog

import (
	"bytes"
	"path/filepath"
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
	eventpl = sync.Pool{New: func() any { return &event{bytes.NewBuffer(make([]byte, 0, 1024)), nil, 0, nil} }}
)

type Event interface {
	Any(string, any) Event
	Error(string, error) Event
	IgError(string, error) Event
	String(string, string) Event
	Strings(string, []string) Event
	Stringp(string, *string) Event
	Duration(string, time.Duration) Event
	Durations(string, []time.Duration) Event
	Durationp(string, *time.Duration) Event
	Byte(string, byte) Event
	Bytes(string, []byte) Event
	Bytep(string, *byte) Event
	Bool(string, bool) Event
	Bools(string, []bool) Event
	Boolp(string, *bool) Event
	Int(string, int) Event
	Ints(string, []int) Event
	Intp(string, *int) Event
	Int8(string, int8) Event
	Int8s(string, []int8) Event
	Int8p(string, *int8) Event
	Int16(string, int16) Event
	Int16s(string, []int16) Event
	Int16p(string, *int16) Event
	Int32(string, int32) Event
	Int32s(string, []int32) Event
	Int32p(string, *int32) Event
	Int64(string, int64) Event
	Int64s(string, []int64) Event
	Int64p(string, *int64) Event
	Uint(string, uint) Event
	Uints(string, []uint) Event
	Uintp(string, *uint) Event
	Uint8(string, uint8) Event
	Uint8s(string, []uint8) Event
	Uint8p(string, *uint8) Event
	Uint16(string, uint16) Event
	Uint16s(string, []uint16) Event
	Uint16p(string, *uint16) Event
	Uint32(string, uint32) Event
	Uint32s(string, []uint32) Event
	Uint32p(string, *uint32) Event
	Uint64(string, uint64) Event
	Uint64s(string, []uint64) Event
	Uint64p(string, *uint64) Event
	Float32(string, float32) Event
	Float32s(string, []float32) Event
	Float32p(string, *float32) Event
	Float64(string, float64) Event
	Float64s(string, []float64) Event
	Float64p(string, *float64) Event
	Done()
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
	PANIC

	OP
)

const maxLen = 1 << 13

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

func appendCaller(bf *bytes.Buffer, skip int) {
	_, file, line, ok := runtime.Caller(skip)
	if ok {
		bf.WriteString(filepath.Base(file))
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

type number interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}

// BUG: 无法处理负数
func transNum[T number](num T) []byte {
	var to = numpl.Get().(*bufnum)
	defer numpl.Put(to)

	neg := num < 0
	if neg {
		num = -num
	}

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
	if neg {
		idx--
		to.buf[idx] = '-'
	}
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
