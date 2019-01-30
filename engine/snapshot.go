package engine

import (
	"context"
	"encoding/gob"
	"fmt"
	"io"
	"os"
	"strconv"
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

		go func() {
			defer close(done)
			err := encoder.Encode(newDs)
			if err != nil {
				fmt.Fprintf(os.Stderr, "ERROR: encoding snapshot: %v\n", err)
			}
		}()
	}

	return fn
}

func DefaultEncoderFactory(writer io.Writer) Encoder {
	return gob.NewEncoder(writer)
}

func DefaultWriterFactory() io.WriteCloser {
	filePrefix := "/usr/local/var/bifrost/"
	fileName := filePrefix + "data-" + strconv.Itoa(int(time.Now().Unix())) + ".bin"

	// TODO: Fix ignoring error
	file, _ := os.Create(fileName)
	return file
}
