package ttl_writer

import (
	"runtime/debug"
	"sync"
	"time"
)

var (
	WriteBufferSize = 1024 // write buffer size
	SleepDuration   = 50 * time.Millisecond
	FlushInterval   = 2 * time.Second
)

type TTLWriter struct {
	opts               *Options
	datas              []any
	closeChan          chan bool
	stopChan           chan bool
	flushImmediateChan chan bool
	d                  time.Duration
	ticker             *time.Ticker
	once               sync.Once
	mu                 sync.Mutex
}

func NewTTLWriter(opts ...Option) *TTLWriter {
	w := &TTLWriter{
		opts:               defaultOptions(),
		closeChan:          make(chan bool, 1),
		stopChan:           make(chan bool, 1),
		flushImmediateChan: make(chan bool, 1),
	}

	for _, o := range opts {
		o(w.opts)
	}

	w.datas = make([]any, 0, w.opts.writeBufferSize)
	w.ticker = time.NewTicker(w.opts.flushInterval)

	w.run()
	return w
}

func (w *TTLWriter) ResetFlushInterval(d time.Duration) {
	w.opts.flushInterval = d
	w.ticker.Reset(d)
}

func (w *TTLWriter) Write(data any) {
	w.mu.Lock()
	w.datas = append(w.datas, data)
	len := len(w.datas)
	w.mu.Unlock()

	if len >= WriteBufferSize {
		w.flush()
	}
}

func (w *TTLWriter) Flush() {
	w.flushImmediateChan <- true
}

func (w *TTLWriter) Stop() {
	w.once.Do(func() {
		w.ticker.Stop()
		close(w.closeChan)
		<-w.stopChan
	})
}

func (w *TTLWriter) run() {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				stack := string(debug.Stack())
				w.opts.logger.Printf("catch exception:%v, panic recovered with stack:%s", err, stack)
			}

			w.stopChan <- true
		}()

		for {
			select {
			case <-w.closeChan:
				w.flush()
				return
			case <-w.ticker.C:
				w.flush()
			case <-w.flushImmediateChan:
				w.flush()
			}
		}
	}()
}

func (w *TTLWriter) flush() {
	w.mu.Lock()
	defer w.mu.Unlock()

	if len(w.datas) <= 0 {
		return
	}

	err := w.opts.flushHandler(w.datas)
	if err != nil {
		w.opts.logger.Printf("flush failed due to %v", err)
	}
	w.datas = w.datas[:0]
}
