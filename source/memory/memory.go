// Package memory is a memory source
package memory

import (
	"errors"
	"github.com/nextpkg/nextcfg"
	"github.com/nextpkg/nextcfg/registry"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/nextpkg/nextcfg/source"
	"log"
)

type memory struct {
	sync.RWMutex
	ChangeSet *source.ChangeSet
	Watchers  map[string]*watcher
}

const sourceName = "memory"

func init() {
	registry.SetCfgSource(sourceName)
	registry.SetCfgLoader(sourceName, func(target string) nextcfg.Loader {
		return GetLoader()
	})
}

// Read 读取配置
func (s *memory) Read() (*source.ChangeSet, error) {
	s.RLock()
	defer s.RUnlock()

	if s.ChangeSet == nil {
		return nil, errors.New("call memory.WithJSON() at first")
	}

	cs := &source.ChangeSet{
		Format:    s.ChangeSet.Format,
		Timestamp: s.ChangeSet.Timestamp,
		Data:      s.ChangeSet.Data,
		Checksum:  s.ChangeSet.Checksum,
		Source:    s.ChangeSet.Source,
	}
	return cs, nil
}

// Watch ...
func (s *memory) Watch() (source.Watcher, error) {
	w := &watcher{
		Id:      uuid.New().String(),
		Updates: make(chan *source.ChangeSet, 100),
		Source:  s,
	}

	s.Lock()
	s.Watchers[w.Id] = w
	s.Unlock()
	return w, nil
}

// Write ...
func (s *memory) Write(cs *source.ChangeSet) error {
	s.Update(cs)
	return nil
}

// Update allows manual updates of the config data.
func (s *memory) Update(c *source.ChangeSet) {
	// don't process nil
	if c == nil {
		return
	}

	// hash the file
	s.Lock()
	// update change set
	s.ChangeSet = &source.ChangeSet{
		Data:      c.Data,
		Format:    c.Format,
		Source:    "memory",
		Timestamp: time.Now(),
	}
	s.ChangeSet.Checksum = s.ChangeSet.Sum()

	// update watchers
	for _, w := range s.Watchers {
		select {
		case w.Updates <- s.ChangeSet:
		default:
		}
	}
	s.Unlock()
}

// String memory
func (s *memory) String() string {
	return sourceName
}

// NewSource ...
func NewSource(opts ...source.Option) source.Source {
	var options source.Options
	for _, o := range opts {
		o(&options)
	}

	s := &memory{
		Watchers: make(map[string]*watcher),
	}

	if options.Context != nil {
		c, ok := options.Context.Value(changeSetKey{}).(*source.ChangeSet)
		if ok {
			s.Update(c)
		}
	}

	return s
}

// GetLoader sets memory source
func GetLoader(opts ...source.Option) nextcfg.Loader {
	return func(l *nextcfg.Loaders) {
		err := l.GetCfg().Load(NewSource(opts...))
		if err != nil {
			log.Println(err)
		} else {
			l.GetCfg().SetState(true)
		}
	}
}
