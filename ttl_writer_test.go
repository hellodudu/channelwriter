package ttl_writer

import (
	"fmt"
	"testing"
	"time"
)

func TestChannelWriter(t *testing.T) {
	flushFn := func(datas []any) error {
		fmt.Println("begin flush")
		for _, data := range datas {
			fmt.Println(data)
		}
		return nil
	}

	cw := NewTTLWriter(
		WithFlushInterval(time.Second),
		WithFlushHandler(flushFn),
	)

	cw.Write(1)
	cw.Write(2)

	time.Sleep(2 * time.Second)

	cw.Write(3)
	cw.Stop()
}
