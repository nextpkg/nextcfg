package memory

import (
	"container/list"
	"errors"
	"github.com/nextpkg/nextcfg/loader"
	"github.com/nextpkg/nextcfg/reader/json"
	"github.com/nextpkg/nextcfg/source"
	"github.com/smartystreets/goconvey/convey"
	"sync"
	"testing"
	"time"
)

type mockWatcher struct {
	Id      string
	Updates chan *source.ChangeSet
	Source  *mockSource
}

func (w *mockWatcher) Next() (*source.ChangeSet, error) {
	cs := <-w.Updates
	return cs, nil
}
func (w *mockWatcher) Stop() error {
	w.Source.Lock()
	delete(w.Source.Watchers, w.Id)
	w.Source.Unlock()
	return nil
}

type mockSource struct {
	sync.RWMutex
	ChangeSet *source.ChangeSet
	Watchers  map[string]*mockWatcher
}

func (s *mockSource) Read() (*source.ChangeSet, error) {
	s.RLock()
	defer s.RUnlock()

	if s.ChangeSet == nil {
		return nil, errors.New("empty changeset")
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
func (s *mockSource) Watch() (source.Watcher, error) {
	w := &mockWatcher{
		Id:      "mock id",
		Updates: make(chan *source.ChangeSet, 100),
		Source:  s,
	}

	s.Lock()
	s.Watchers[w.Id] = w
	s.Unlock()
	return w, nil
}
func (s *mockSource) Write(cs *source.ChangeSet) error {
	s.Update(cs)
	return nil
}
func (s *mockSource) Update(c *source.ChangeSet) {
	if c == nil {
		return
	}

	s.Lock()
	s.ChangeSet = &source.ChangeSet{
		Data:      c.Data,
		Format:    c.Format,
		Source:    "mock",
		Timestamp: time.Now(),
	}
	s.ChangeSet.Checksum = s.ChangeSet.Sum()

	for _, w := range s.Watchers {
		select {
		case w.Updates <- s.ChangeSet:
		default:
		}
	}
	s.Unlock()
}
func (s *mockSource) String() string {
	return "mock"
}

func TestDoWatch(t *testing.T) {
	convey.Convey("case #1", t, func() {
		cs := &source.ChangeSet{
			Data:      []byte("{}"),
			Checksum:  "sum",
			Format:    "text",
			Source:    "mock",
			Timestamp: time.Now(),
		}
		src := &mockSource{
			Watchers:  make(map[string]*mockWatcher),
			ChangeSet: cs,
		}
		m := &memory{
			exit: make(chan bool),
			opts: loader.Options{
				Reader: json.NewReader(),
				Source: []source.Source{src},
			},
			watchers: list.New(),
			sources:  []source.Source{src},
			sets:     []*source.ChangeSet{cs},
		}

		convey.So(m.loaded(), convey.ShouldBeFalse)
		convey.So(m.reload(), convey.ShouldBeNil)
		convey.So(m.loaded(), convey.ShouldBeTrue)

		snap, err := m.Snapshot()
		convey.So(err, convey.ShouldBeNil)

		convey.So(snap, convey.ShouldNotBeNil)
		convey.So(snap.ChangeSet.Data, convey.ShouldEqual, []byte("{}"))
		convey.So(m.loaded(), convey.ShouldBeTrue)

		get, err := m.Get()
		convey.So(err, convey.ShouldBeNil)
		convey.So(get.Bytes(), convey.ShouldEqual, []byte("{}"))
	})

}
