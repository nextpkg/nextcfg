package memory

import (
	"bytes"
	"container/list"
	"fmt"
	"log/slog"
	"strings"
	"sync"
	"time"

	"github.com/nextpkg/nextcfg/loader"
	"github.com/nextpkg/nextcfg/reader"
	"github.com/nextpkg/nextcfg/reader/json"
	"github.com/nextpkg/nextcfg/source"
	"github.com/pkg/errors"
)

type memory struct {
	exit chan bool
	opts loader.Options

	sync.RWMutex
	// the current snapshot
	snap *loader.Snapshot
	// the current values
	val reader.Values
	// all the changesets
	sets []*source.ChangeSet
	// all the sources
	sources  []source.Source
	watchers *list.List
}

type updateValue struct {
	version string
	value   reader.Value
}

type watcher struct {
	exit    chan bool
	path    []string
	value   reader.Value
	reader  reader.Reader
	version string
	updates chan updateValue
}

func (m *memory) watch(idx int, s source.Source) {
	// watches a source for changes
	watch := func(idx int, s source.Watcher) error {
		for {
			// get change set
			cs, err := s.Next()
			if err != nil {
				return err
			}

			m.Lock()

			// save
			m.sets[idx] = cs

			// merge sets
			set, err := m.opts.Reader.Merge(m.sets...)
			if err != nil {
				m.Unlock()
				return err
			}

			// set values
			m.val, _ = m.opts.Reader.Values(set)
			m.snap = &loader.Snapshot{
				ChangeSet: set,
				Version:   genVer(),
			}
			m.Unlock()

			// send watch updates
			m.update()
		}
	}

	for {
		// get watcher
		w, err := s.Watch()
		if err != nil {
			slog.Warn("memory.watch() failed.", slog.Any("err", err))
			time.Sleep(time.Second)
			continue
		}

		done := make(chan bool)

		// the stop watch func
		go func() {
			select {
			case <-done:
			case <-m.exit:
			}

			if err = w.Stop(); err != nil {
				slog.Error("Stop failed: ", slog.Any("err", err))
			}
		}()

		// block watch
		if err = watch(idx, w); err != nil {
			// do something better
			time.Sleep(time.Second)
		}

		// close done chan
		close(done)

		// if the config is closed exit
		select {
		case <-m.exit:
			return
		default:
		}
	}
}

func (m *memory) loaded() bool {
	var loaded bool
	m.RLock()
	if m.val != nil {
		loaded = true
	}
	m.RUnlock()
	return loaded
}

// reload reads the sets and creates new values
func (m *memory) reload() error {
	m.Lock()

	// merge sets
	set, err := m.opts.Reader.Merge(m.sets...)
	if err != nil {
		m.Unlock()
		return err
	}

	// set values
	m.val, _ = m.opts.Reader.Values(set)
	m.snap = &loader.Snapshot{
		ChangeSet: set,
		Version:   genVer(),
	}

	m.Unlock()

	// update watchers
	m.update()

	return nil
}

func (m *memory) update() {
	watchers := make([]*watcher, 0, m.watchers.Len())

	m.RLock()
	for e := m.watchers.Front(); e != nil; e = e.Next() {
		watchers = append(watchers, e.Value.(*watcher))
	}

	val := m.val
	snap := m.snap
	m.RUnlock()

	for _, w := range watchers {
		if w.version >= snap.Version {
			continue
		}

		uv := updateValue{
			version: m.snap.Version,
			value:   val.Get(w.path...),
		}

		select {
		case w.updates <- uv:
		default:
		}
	}
}

// Snapshot returns a snapshot of the current loaded config
func (m *memory) Snapshot() (*loader.Snapshot, error) {
	if m.loaded() {
		m.RLock()
		snap := loader.Copy(m.snap)
		m.RUnlock()
		return snap, nil
	}

	// not loaded, sync
	if err := m.Sync(); err != nil {
		return nil, err
	}

	// make copy
	m.RLock()
	snap := loader.Copy(m.snap)
	m.RUnlock()

	return snap, nil
}

// Sync loads all the sources, calls the parser and updates the config
func (m *memory) Sync() error {
	//nolint:prealloc
	var sets []*source.ChangeSet

	m.Lock()

	// read the s
	var gErr []string

	for _, s := range m.sources {
		ch, err := s.Read()
		if err != nil {
			gErr = append(gErr, err.Error())
			continue
		}
		sets = append(sets, ch)
	}

	// merge sets
	set, err := m.opts.Reader.Merge(sets...)
	if err != nil {
		m.Unlock()
		return err
	}

	// set values
	var val reader.Values
	val, err = m.opts.Reader.Values(set)
	if err != nil {
		m.Unlock()
		return err
	}
	m.val = val
	m.snap = &loader.Snapshot{
		ChangeSet: set,
		Version:   genVer(),
	}

	m.Unlock()

	// update watchers
	m.update()

	if len(gErr) > 0 {
		return fmt.Errorf("loading errors: %s", strings.Join(gErr, "\n"))
	}

	return nil
}

// Close 终止
func (m *memory) Close() error {
	select {
	case <-m.exit:
		return nil
	default:
		close(m.exit)
	}
	return nil
}

// Get 获取指定路径的配置
func (m *memory) Get(path ...string) (reader.Value, error) {
	if !m.loaded() {
		if err := m.Sync(); err != nil {
			return nil, err
		}
	}

	m.Lock()
	defer m.Unlock()

	if m.val != nil {
		return m.val.Get(path...), nil
	}

	// assuming val is nil, create new val
	ch := m.snap.ChangeSet

	// we are truly screwed, trying to load in a hacked way
	v, err := m.opts.Reader.Values(ch)
	if err != nil {
		return nil, err
	}

	m.val = v

	if m.val != nil {
		return m.val.Get(path...), nil
	}

	// ok we're going hardcore now
	return nil, errors.New("no values")
}

// Load 加载数据源
func (m *memory) Load(sources ...source.Source) error {
	var gErr []string

	for _, s := range sources {
		set, err := s.Read()
		if err != nil {
			gErr = append(gErr, fmt.Sprintf("loading %s error: %v", s, err))
			// continue processing
			continue
		}
		m.Lock()
		m.sources = append(m.sources, s)
		m.sets = append(m.sets, set)
		idx := len(m.sets) - 1
		m.Unlock()
		go m.watch(idx, s)
	}

	if err := m.reload(); err != nil {
		gErr = append(gErr, err.Error())
	}

	// Return errors
	if len(gErr) != 0 {
		return errors.New(strings.Join(gErr, "\n"))
	}
	return nil
}

// Watch 监听路径
func (m *memory) Watch(path ...string) (loader.Watcher, error) {
	value, err := m.Get(path...)
	if err != nil {
		return nil, err
	}

	m.Lock()

	w := &watcher{
		exit:    make(chan bool),
		path:    path,
		value:   value,
		reader:  m.opts.Reader,
		updates: make(chan updateValue, 1),
		version: m.snap.Version,
	}

	e := m.watchers.PushBack(w)

	m.Unlock()

	go func() {
		<-w.exit
		m.Lock()
		m.watchers.Remove(e)
		m.Unlock()
	}()

	return w, nil
}

// String memory
func (m *memory) String() string {
	return "memory"
}

// Next 下一个变更的快照
func (w *watcher) Next() (*loader.Snapshot, error) {
	update := func(v reader.Value) *loader.Snapshot {
		w.value = v

		cs := &source.ChangeSet{
			Data:      v.Bytes(),
			Format:    w.reader.String(),
			Source:    "memory",
			Timestamp: time.Now(),
		}
		cs.Checksum = cs.Sum()

		return &loader.Snapshot{
			ChangeSet: cs,
			Version:   w.version,
		}

	}

	for {
		select {
		case <-w.exit:
			return nil, errors.New("watcher stopped")

		case uv := <-w.updates:
			if uv.version <= w.version {
				continue
			}

			v := uv.value

			w.version = uv.version

			if bytes.Equal(w.value.Bytes(), v.Bytes()) {
				continue
			}

			return update(v), nil
		}
	}
}

// Stop 终止监听器
func (w *watcher) Stop() error {
	select {
	case <-w.exit:
	default:
		close(w.exit)
		close(w.updates)
	}

	return nil
}

func genVer() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

// NewLoader memory配置加载
func NewLoader(opts ...loader.Option) loader.Loader {
	options := loader.Options{
		Reader: json.NewReader(),
	}

	for _, o := range opts {
		o(&options)
	}

	m := &memory{
		exit:     make(chan bool),
		opts:     options,
		watchers: list.New(),
		sources:  options.Source,
	}

	m.sets = make([]*source.ChangeSet, len(options.Source))

	for i, s := range options.Source {
		m.sets[i] = &source.ChangeSet{Source: s.String()}
		go m.watch(i, s)
	}

	return m
}
