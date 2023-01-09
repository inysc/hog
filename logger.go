package qog

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"runtime"
	"time"
)

type logger struct {
	writer io.Writer
	trace  string
	lvl    Level
}

func (l *logger) write(ctx context.Context, info string, msg string, args []any) (n int, err error) {
	bs := bpl.Get().(*bytes.Buffer)
	bs.Reset()

	// 写入时间戳
	// time.Now().AppendFormat(bs, "2006-01-02 15:04:05.000")
	// bs.Write(time.Now().AppendFormat([]byte{}, "2006-01-02 15:04:05.000"))
	appendTime(bs, time.Now())

	// 写入等级及服务名
	bs.WriteString(info)
	// bs.buf = append(bs.buf, info...)

	// 写入调用信息
	appendCaller(bs)

	// 写入 goid
	bs.WriteString("|goid:")
	bs.Write(transNum(runtime.Goid()))
	bs.WriteByte('|')

	if ctx != nil {
		if trace, _ := ctx.Value(l.trace).(string); trace != "" {
			bs.WriteString("trace:")
			bs.WriteString(trace)
		}
	}

	bs.WriteByte(',')

	// 写入日志信息
	if len(args) == 0 {
		bs.WriteString(msg)
	} else {
		bs.WriteString(fmt.Sprintf(msg, args...))
	}
	bs.WriteByte('\n')

	// 写入到文件
	n, err = l.writer.Write(bs.Bytes())

	// 太大就抛弃了
	if bs.Len() < maxLen<<3 {
		bpl.Put(bs)
	}
	return
}

func (l *logger) SetLevel(lvl Level)    { l.lvl = lvl }
func (l *logger) SetTrace(trace string) { l.trace = trace }

// needStd 是否需要同时输出到标准输出（仅在 ppid != 1 时生效 ）
func New(needStd bool, lvl Level, trace_field string, ws ...io.Writer) *logger {
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
	return &logger{w, trace_field, lvl}
}

func Simple(filename string) *logger {
	return New(true, DEBUG, "trace", &LoggerFile{
		Filename:   filename,
		MaxSize:    100,
		MaxAge:     30,
		MaxBackups: 11,
		LocalTime:  false,
		Compress:   false,
	})
}
