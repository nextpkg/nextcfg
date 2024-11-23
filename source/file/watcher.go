//go:build !linux
// +build !linux

package file

import (
	"os"

	"github.com/fsnotify/fsnotify"
	"github.com/nextpkg/nextcfg/source"
	"log"
)

type watcher struct {
	f    *file
	fw   *fsnotify.Watcher
	exit chan bool
}

func newWatcher(f *file) (source.Watcher, error) {
	fw, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	if err := fw.Add(f.path); err != nil {
		log.Println("add notify path failed:", err)
	}

	return &watcher{
		f:    f,
		fw:   fw,
		exit: make(chan bool),
	}, nil
}

// Next ...
func (w *watcher) Next() (*source.ChangeSet, error) {
	// is it closed?
	select {
	case <-w.exit:
		return nil, source.ErrWatcherStopped
	default:
	}

	// try get the event
	select {
	case event, _ := <-w.fw.Events:
		if event.Op == fsnotify.Rename {
			// check existence of file, and add watch again
			_, err := os.Stat(event.Name)
			if err == nil || os.IsExist(err) {
				if err := w.fw.Add(event.Name); err != nil {
					log.Println("add notify path failed:", err)
				}
			}
		}

		c, err := w.f.Read()
		if err != nil {
			return nil, err
		}
		return c, nil
	case err := <-w.fw.Errors:
		return nil, err
	case <-w.exit:
		return nil, source.ErrWatcherStopped
	}
}

// Stop ...
func (w *watcher) Stop() error {
	return w.fw.Close()
}
