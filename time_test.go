package qog

import (
	"bytes"
	"fmt"
	"testing"
	"time"
)

func assert(b bool, msg string) {
	if !b {
		panic(msg)
	}
}

func TestTimeAppend(t *testing.T) {
	for i := 0; i < 1000; i++ {
		if i < 10 {
			assert(fmt.Sprintf("%03d", i) == fmt.Sprintf("00%c", smallsString[i*2+1]), fmt.Sprintf("i=%d", i))
		} else if i < 100 {
			nano := i * 2
			assert(
				fmt.Sprintf("%03d", i) ==
					fmt.Sprintf("0%c%c", smallsString[nano+0], smallsString[nano+1]),
				fmt.Sprintf("i=%d, 0%c%c", i, smallsString[nano+0], smallsString[nano+1]))
		} else {
			assert(
				fmt.Sprintf("%03d", i) ==
					string(transNum(i)),
				fmt.Sprintf("i=%d,  %s", i, transNum(i)))
		}
	}
}

func BenchmarkIf1(b *testing.B) {
	a := 0
	for i := 0; i < b.N; i++ {
		if a == 0 {
			a += 0
		} else {
			a = 1
		}
	}
}

func BenchmarkIf2(b *testing.B) {
	a := 0
	for i := 0; i < b.N; i++ {
		if a != 1 {
			a += 0
		} else {
			a = 1
		}
	}
}

func BenchmarkAppendByteStr(b *testing.B) {
	var be = make([]byte, 0, 1024)
	for i := 0; i < b.N; i++ {
		be = be[:0]
		be = append(be, "<nil>"...)
	}
}

func BenchmarkAppendBytes(b *testing.B) {
	var be = make([]byte, 0, 1024)
	for i := 0; i < b.N; i++ {
		be = be[:0]
		be = append(be, '<', 'n', 'i', 'l', '>')
	}
}

func BenchmarkAppendStrCopy(b *testing.B) {
	var be = make([]byte, 0, 1024)
	for i := 0; i < b.N; i++ {
		be = be[:0]
		copy(be, "<nil>")
	}
}

func TestLogger(t *testing.T) {
}

// goos: linux
// goarch: amd64
// pkg: github.com/inysc/qog
// cpu: AMD Ryzen 5 1600 Six-Core Processor
// BenchmarkTimeAppendFormat1-12            2391530               511.6 ns/op            56 B/op          3 allocs/op
// BenchmarkTimeAppendFormat2-12            3178784               382.5 ns/op             0 B/op          0 allocs/op
func BenchmarkTimeAppendFormat1(b *testing.B) {
	bf := bpl.Get().(*bytes.Buffer)
	for i := 0; i < b.N; i++ {
		bf.Reset()
		bf.Write(time.Now().AppendFormat([]byte{}, "2006-01-02 15:04:05.000"))
	}
	bpl.Put(bf)
}

func BenchmarkTimeAppendFormat2(b *testing.B) {
	bf := bpl.Get().(*bytes.Buffer)
	for i := 0; i < b.N; i++ {
		bf.Reset()
		appendTime(bf, time.Now())
	}
	bpl.Put(bf)
}

func TestTransNum(t *testing.T) {
	for i := 0; i < 100; i++ {
		t.Logf("%s", transNum(i))
	}
}

func TestDebug(t *testing.T) {
	lg := Simple("")
	lg.lvl = TRACE
	lg.Trace().Any("any=", nil).Bool("a=", true).Bools("||c=", []bool{true, false, true}).Msg("")
	lg.Debug().Bool("b=", true).Bools("||c=", []bool{true, false, true}).Msg("")
	lg.Info().Bool("c=", true).Bools("||c=", []bool{true, false, true}).Msg("")

	t.Log(1 << 13)
}
