package engine

import (
	"context"
	"io"
	"time"
)

type WriterFactory func() io.WriteCloser
type EncoderFactory func(io.Writer) Encoder

type Encoder interface {
	Encode(interface{}) error
}

func SnapshotTimer(
	ctx context.Context,
	interval time.Duration,
	newEncoder EncoderFactory,
	newWriteCloser WriterFactory,
) <-chan Event {
	ch := make(chan Event)

	go func() {
		defer close(ch)
		for {
			select {
			case <-time.After(interval):
			case <-ctx.Done():
				return
			}

			done := make(chan struct{})

			writer := newWriteCloser()
			defer writer.Close()
			encoder := newEncoder(writer)

			ch <- Snapshot(encoder, done)
			<-done
		}
	}()

	return ch
}

func Snapshot(encoder Encoder, done chan struct{}) Event {
	var fn EventFn = func(ds *DataStore) {
		newDs := ds.Copy()

		// TODO: Pseudo code...
		go func() {
			defer close(done)
			encoder.Encode(newDs)
		}()
	}

	return fn
}
