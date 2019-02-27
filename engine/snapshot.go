package engine

import (
	"context"
	"encoding/gob"
	"fmt"
	"io"
	"os"
	"time"
)

const snapshotFile = "/usr/local/var/bifrost/snapshot.data"

type WriteCloserFactory func() (io.WriteCloser, error)
type ReadCloserFactory func() (io.ReadCloser, error)
type EncoderFactory func(io.Writer) (Encoder, error)

//type DecoderFactory func(io.Reader) (Decoder, error)

// Encoder is a generalized interface implemented by Go's various encoding
// libraries.
type Encoder interface {
	Encode(interface{}) error
}

type Decoder interface {
	Decode(interface{}) error
}

// SnapshotTimer produces Snapshot Events on the returned channel until ctx
// is Done.
// `interval specifies how often to produce events.
// `newEncoder` is a factory that builds encoders to serialize the Engine State
// `newWriteCloser` is a factory that builds io.WriteClosers to save the
// serialized state to.
func SnapshotTimer(
	ctx context.Context,
	interval time.Duration,
	newEncoder EncoderFactory,
	newWriteCloser WriteCloserFactory,
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

			writer, err := newWriteCloser()
			if err != nil {
				fmt.Fprintf(os.Stderr, "ERROR: making writer to save snapshot: %v\n", err)
				return
			}
			defer writer.Close()

			encoder, err := newEncoder(writer)
			if err != nil {
				fmt.Fprintf(os.Stderr, "ERROR: Making encoder to save snapshot %v\n", err)
				return
			}

			ch <- Snapshot(encoder, done)
			<-done
		}
	}()

	return ch
}

// Snapshot returns a new Snapshot event that will use the provided encoder
// to serialize the Engine state. It will close() the `done` channel when the
// snapshot is complete.
func Snapshot(encoder Encoder, done chan struct{}) Event {
	var fn EventFn = func(ds *DataStore) ChangeSet {
		newDs := ds.Copy()

		go func() {
			defer close(done)
			err := encoder.Encode(newDs)
			if err != nil {
				fmt.Fprintf(os.Stderr, "ERROR: encoding snapshot: %v\n", err)
			}
		}()

		return ChangeSet{}
	}

	return fn
}

// DefaultEncoder returns a new gob.Encoder for the provided writer.
func DefaultEncoder(writer io.Writer) (Encoder, error) {
	return gob.NewEncoder(writer), nil
}

func DefaultDecoder(reader io.Reader) (Decoder, error) {
	return gob.NewDecoder(reader), nil
}

// DefaultWriteCloserFactory returns a new WriteCloserFactory that will write
// to the file specified by `snapshotFile`. If the file already exists, it
// will be copied to a file ending in `.bkp` before writing the new file.
func DefaultWriteCloserFactory(snapshotFile string) WriteCloserFactory {
	return func() (io.WriteCloser, error) {
		return DefaultWriteCloser(snapshotFile)
	}
}

func DefaultWriteCloser(snapshotFile string) (io.WriteCloser, error) {
	snapshotBackupFile := snapshotFile + ".bkp"
	err := os.Rename(snapshotFile, snapshotBackupFile)
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	err = os.Remove(snapshotFile)
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	file, err := os.Create(snapshotFile)
	if err != nil {
		return nil, err
	}

	return file, nil
}

func DefaultReadCloserFactory(snapshotFile string) ReadCloserFactory {
	return func() (io.ReadCloser, error) {
		return DefaultReadCloser(snapshotFile)
	}
}

func DefaultReadCloser(snapshotFile string) (io.ReadCloser, error) {
	return os.Open(snapshotFile)
}
