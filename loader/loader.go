// Package loader manages loading from multiple sources
package loader

import (
	"context"

	"github.com/nextpkg/nextcfg/reader"
	"github.com/nextpkg/nextcfg/source"
)

// Loader manages loading sources
type Loader interface {
	// Close Stop the loader
	Close() error
	// Load the sources
	Load(...source.Source) error
	// A Snapshot of loaded config
	Snapshot() (*Snapshot, error)
	// Sync Force sync of sources
	Sync() error
	// Watch for changes
	Watch(...string) (Watcher, error)
	// String Name of loader
	String() string
}

// Watcher lets you watch sources and returns a merged ChangeSet
type Watcher interface {
	// Next First call to next may return the current Snapshot
	// If you are watching a path then only the data from
	// that path is returned.
	Next() (*Snapshot, error)
	// Stop watching for changes
	Stop() error
}

// Snapshot is a merged ChangeSet
type Snapshot struct {
	// The merged ChangeSet
	ChangeSet *source.ChangeSet
	// Deterministic and comparable version of the snapshot
	Version string
}

// Options loader选项
type Options struct {
	Reader reader.Reader
	Source []source.Source

	// for alternative data
	Context context.Context
}

// Option loader选项
type Option func(o *Options)

// Copy snapshot
func Copy(s *Snapshot) *Snapshot {
	snapshot := *(s.ChangeSet)
	return &Snapshot{
		ChangeSet: &snapshot,
		Version:   s.Version,
	}
}
