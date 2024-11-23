package memory

import (
	"github.com/nextpkg/nextcfg/source"
)

type watcher struct {
	Id      string
	Updates chan *source.ChangeSet
	Source  *memory
}

// Next ...
func (w *watcher) Next() (*source.ChangeSet, error) {
	cs := <-w.Updates
	return cs, nil
}

// Stop ...
func (w *watcher) Stop() error {
	w.Source.Lock()
	delete(w.Source.Watchers, w.Id)
	w.Source.Unlock()
	return nil
}
