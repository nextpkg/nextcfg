package env

import (
	"github.com/nextpkg/nextcfg/source"
)

type watcher struct {
	exit chan struct{}
}

// Next ...
func (w *watcher) Next() (*source.ChangeSet, error) {
	<-w.exit

	return nil, source.ErrWatcherStopped
}

// Stop ...
func (w *watcher) Stop() error {
	close(w.exit)
	return nil
}

func newWatcher() (source.Watcher, error) {
	return &watcher{exit: make(chan struct{})}, nil
}
