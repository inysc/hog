package qog

import "sync"

type buf struct {
	buf []byte
}

type bufTo struct {
	buf [22]byte
}

var bpl = &sync.Pool{
	New: func() any {
		return &buf{buf: make([]byte, maxLen)}
	},
}

var toPl = &sync.Pool{
	New: func() any {
		return &bufTo{}
	},
}
