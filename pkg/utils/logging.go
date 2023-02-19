package utils

import (
	"io"
	"log"
	"time"
)

type prefixWriter struct {
	f func() string
	w io.Writer
}

func (p prefixWriter) Write(b []byte) (n int, err error) {
	if n, err = p.w.Write([]byte(p.f())); err != nil {
		return
	}
	nn, err := p.w.Write(b)
	return n + nn, err
}

func InitLogging() {
	log.SetFlags(0)
	log.SetOutput(prefixWriter{
		f: func() string { return time.Now().Format(time.RFC3339Nano) + " " },
		w: log.Writer(),
	})
}
