package ttl_writer

import (
	"io"
	"log"
	"time"
)

type FlushHandler func([]any) error

type Option func(*Options)
type Options struct {
	writeBufferSize int
	flushInterval   time.Duration
	sleepDuration   time.Duration
	logger          *log.Logger
	flushHandler    FlushHandler
}

func defaultOptions() *Options {
	return &Options{
		writeBufferSize: 1024,
		flushInterval:   2 * time.Second,
		sleepDuration:   50 * time.Millisecond,
		logger:          log.Default(),
		flushHandler:    func([]any) error { return nil },
	}
}

func WithWriteBufferSize(sz int) Option {
	return func(o *Options) {
		o.writeBufferSize = sz
	}
}

func WithFlushInterval(t time.Duration) Option {
	return func(o *Options) {
		o.flushInterval = t
	}
}

func WithSleepDuration(d time.Duration) Option {
	return func(o *Options) {
		o.sleepDuration = d
	}
}

func WithLogger(w io.Writer) Option {
	return func(o *Options) {
		o.logger = log.New(w, "ttl_writer: ", log.Lmsgprefix|log.LstdFlags)
		log.SetOutput(w)
	}
}

func WithFlushHandler(h FlushHandler) Option {
	return func(o *Options) {
		o.flushHandler = h
	}
}
