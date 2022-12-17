package qog

import (
	"bytes"
	"fmt"
	"strconv"
	"time"
)

type event struct{ bf *bytes.Buffer }

var _ Event = (*event)(nil)

func NewEvent() Event { return eventpl.Get().(*event) }
func FreeEvent(e Event) {
	switch e := e.(type) {
	case *event:
		e.bf.Reset()
		eventpl.Put(e)
	}
}

func (e *event) String(key string, val string) Event {
	e.bf.WriteString(key)
	e.bf.WriteString(": ")
	e.bf.WriteString(val)
	e.bf.WriteByte(',')
	return e
}

func (e *event) Stringp(key string, val *string) Event {
	e.bf.WriteString(key)
	e.bf.WriteString(": ")
	if val != nil {
		e.bf.WriteString(*val)
		e.bf.WriteByte(',')
	} else {
		e.bf.WriteString(" <nil>,")
	}
	return e
}

func (e *event) Duration(key string, val time.Duration) Event {
	e.bf.WriteString(key)
	e.bf.WriteString(": ")
	e.bf.WriteString(val.String())
	e.bf.WriteByte(',')
	return e
}

func (e *event) Durationp(key string, val *time.Duration) Event {
	e.bf.WriteString(key)
	e.bf.WriteString(": ")
	if val != nil {
		e.bf.WriteString(val.String())
		e.bf.WriteByte(',')
	} else {
		e.bf.WriteString(" <nil>,")
	}
	return e
}

func (e *event) Bool(key string, val bool) Event {
	e.bf.WriteString(key)
	if val {
		e.bf.WriteString(": true,")
	} else {
		e.bf.WriteString(": false,")
	}
	return e
}

func (e *event) Boolp(key string, val *bool) Event {
	e.bf.WriteString(key)
	if val != nil {
		if *val {
			e.bf.WriteString(": true,")
		} else {
			e.bf.WriteString(": false,")
		}
	} else {
		e.bf.WriteString(": <nil>,")
	}
	return e
}

func (e *event) Int(key string, val int) Event {
	e.bf.WriteString(key)
	e.bf.WriteString(": ")
	e.bf.Write(transNum(val))
	e.bf.WriteByte(',')

	return e
}

func (e *event) Intp(key string, val *int) Event {
	e.bf.WriteString(key)
	e.bf.WriteString(": ")
	if val == nil {
		e.bf.WriteString(": <nil>,")
	} else {
		e.bf.Write(transNum(*val))
		e.bf.WriteByte(',')
	}

	return e
}

func (e *event) Int8(key string, val int8) Event {
	e.bf.WriteString(key)
	e.bf.WriteString(": ")
	e.bf.Write(transNum(val))
	e.bf.WriteByte(',')

	return e
}

func (e *event) Int8p(key string, val *int8) Event {
	e.bf.WriteString(key)
	e.bf.WriteString(": ")
	if val == nil {
		e.bf.WriteString(": <nil>,")
	} else {
		e.bf.Write(transNum(*val))
		e.bf.WriteByte(',')
	}

	return e
}

func (e *event) Int16(key string, val int16) Event {
	e.bf.WriteString(key)
	e.bf.WriteString(": ")
	e.bf.Write(transNum(val))
	e.bf.WriteByte(',')

	return e
}

func (e *event) Int16p(key string, val *int16) Event {
	e.bf.WriteString(key)
	e.bf.WriteString(": ")
	if val == nil {
		e.bf.WriteString(": <nil>,")
	} else {
		e.bf.Write(transNum(*val))
		e.bf.WriteByte(',')
	}

	return e
}

func (e *event) Int32(key string, val int32) Event {
	e.bf.WriteString(key)
	e.bf.WriteString(": ")
	e.bf.Write(transNum(val))
	e.bf.WriteByte(',')

	return e
}

func (e *event) Int32p(key string, val *int32) Event {
	e.bf.WriteString(key)
	e.bf.WriteString(": ")
	if val == nil {
		e.bf.WriteString(": <nil>,")
	} else {
		e.bf.Write(transNum(*val))
		e.bf.WriteByte(',')
	}

	return e
}

func (e *event) Int64(key string, val int64) Event {
	e.bf.WriteString(key)
	e.bf.WriteString(": ")
	e.bf.Write(transNum(val))
	e.bf.WriteByte(',')

	return e
}

func (e *event) Int64p(key string, val *int64) Event {
	e.bf.WriteString(key)
	e.bf.WriteString(": ")
	if val == nil {
		e.bf.WriteString(": <nil>,")
	} else {
		e.bf.Write(transNum(*val))
		e.bf.WriteByte(',')
	}

	return e
}

func (e *event) Uint(key string, val uint) Event {
	e.bf.WriteString(key)
	e.bf.WriteString(": ")
	e.bf.Write(transNum(val))
	e.bf.WriteByte(',')

	return e
}

func (e *event) Uintp(key string, val *uint) Event {
	e.bf.WriteString(key)
	e.bf.WriteString(": ")
	if val == nil {
		e.bf.WriteString(": <nil>,")
	} else {
		e.bf.Write(transNum(*val))
		e.bf.WriteByte(',')
	}

	return e
}

func (e *event) Uint8(key string, val uint8) Event {
	e.bf.WriteString(key)
	e.bf.WriteString(": ")
	e.bf.Write(transNum(val))
	e.bf.WriteByte(',')

	return e
}

func (e *event) Uint8p(key string, val *uint8) Event {
	e.bf.WriteString(key)
	e.bf.WriteString(": ")
	if val == nil {
		e.bf.WriteString(": <nil>,")
	} else {
		e.bf.Write(transNum(*val))
		e.bf.WriteByte(',')
	}

	return e
}

func (e *event) Uint16(key string, val uint16) Event {
	e.bf.WriteString(key)
	e.bf.WriteString(": ")
	e.bf.Write(transNum(val))
	e.bf.WriteByte(',')

	return e
}

func (e *event) Uint16p(key string, val *uint16) Event {
	e.bf.WriteString(key)
	e.bf.WriteString(": ")
	if val == nil {
		e.bf.WriteString(": <nil>,")
	} else {
		e.bf.Write(transNum(*val))
		e.bf.WriteByte(',')
	}

	return e
}

func (e *event) Uint32(key string, val uint32) Event {
	e.bf.WriteString(key)
	e.bf.WriteString(": ")
	e.bf.Write(transNum(val))
	e.bf.WriteByte(',')

	return e
}

func (e *event) Uint32p(key string, val *uint32) Event {
	e.bf.WriteString(key)
	e.bf.WriteString(": ")
	if val == nil {
		e.bf.WriteString(": <nil>,")
	} else {
		e.bf.Write(transNum(*val))
		e.bf.WriteByte(',')
	}

	return e
}

func (e *event) Uint64(key string, val uint64) Event {
	e.bf.WriteString(key)
	e.bf.WriteString(": ")
	e.bf.Write(transNum(val))
	e.bf.WriteByte(',')

	return e
}

func (e *event) Uint64p(key string, val *uint64) Event {
	e.bf.WriteString(key)
	e.bf.WriteString(": ")
	if val == nil {
		e.bf.WriteString(": <nil>,")
	} else {
		e.bf.Write(transNum(*val))
		e.bf.WriteByte(',')
	}

	return e
}

func (e *event) Any(key string, val any) Event {
	e.bf.WriteString(key)
	e.bf.WriteString(": ")
	if val == nil {
		e.bf.WriteString(": <nil>,")
	} else {
		if v, ok := val.(interface{ String() string }); ok {
			e.bf.WriteString(v.String())
		} else {
			e.bf.WriteString(fmt.Sprintf("%+v", val))
		}
		e.bf.WriteByte(',')
	}
	return e
}

func (e *event) Float32(key string, val float32) Event {
	e.bf.WriteString(key)
	e.bf.WriteString(": ")
	e.bf.Write(strconv.AppendFloat([]byte{}, float64(val), 'f', -1, 32))
	e.bf.WriteByte(',')

	return e
}

func (e *event) Float32p(key string, val *float32) Event {
	e.bf.WriteString(key)
	e.bf.WriteString(": ")
	if val == nil {
		e.bf.WriteString(": <nil>,")
	} else {
		e.bf.Write(strconv.AppendFloat([]byte{}, float64(*val), 'f', -1, 32))
		e.bf.WriteByte(',')
	}

	return e
}

func (e *event) Float64(key string, val float64) Event {
	e.bf.WriteString(key)
	e.bf.WriteString(": ")
	e.bf.Write(strconv.AppendFloat([]byte{}, val, 'f', -1, 64))
	e.bf.WriteByte(',')
	return e
}

func (e *event) Float64p(key string, val *float64) Event {
	e.bf.WriteString(key)
	e.bf.WriteString(": ")
	if val == nil {
		e.bf.WriteString(": <nil>,")
	} else {
		e.bf.Write(strconv.AppendFloat([]byte{}, *val, 'f', -1, 64))
		e.bf.WriteByte(',')
	}

	return e
}

func (e *event) Error(key string, val error) Event {
	e.bf.WriteString(key)
	e.bf.WriteString(": ")
	e.bf.WriteString(val.Error())
	e.bf.WriteByte(',')
	return e
}

func (e *event) Msg(msg string) {
	e.bf.WriteString(msg)
}

func (e *event) Msgf(format string, args ...any) {
	e.bf.WriteString(fmt.Sprintf(format, args...))
}
