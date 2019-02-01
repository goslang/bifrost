package engine

import (
	"context"
	"time"
)

// Config captures configurable options for the Engine.
type Config struct {
	ctx context.Context

	snapshotInterval time.Duration
	snapshotFilename string
}

type Opt func(Config) Config

// SnapshotInterval specifies how often to write the current engine state to
// the file system.
func SnapshotInterval(i time.Duration) Opt {
	return func(conf Config) Config {
		conf.snapshotInterval = i
		return conf
	}
}

// SnapshotFilename specifies where to save snapshots.
func SnapshotFilename(name string) Opt {
	return func(conf Config) Config {
		conf.snapshotFilename = name
		return conf
	}
}

// Context sets the Engine's context, the Engine will exit when it is
// `Done()`.
func Context(ctx context.Context) Opt {
	return func(conf Config) Config {
		conf.ctx = ctx
		return conf
	}
}
