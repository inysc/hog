package qog

import (
	"runtime"
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

type integer interface {
	~int | ~int8 | ~int16 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint64
}

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
		appendNum(bf, line)
	}
}

func appendNum[T integer](b *buf, num T) {
	var to = toPl.Get().(*bufTo) // +1 for sign of 64bit value in base 2
	to.buf[21] = ']'
	i := 21
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
	b.buf = append(b.buf, to.buf[i:]...)
	toPl.Put(to)
}
