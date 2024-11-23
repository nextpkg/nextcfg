package nextcfg

import (
	"bytes"
	"log/slog"
	"sync"
	"time"

	"github.com/nextpkg/nextcfg/loader"
	"github.com/nextpkg/nextcfg/loader/memory"
	"github.com/nextpkg/nextcfg/reader"
	"github.com/nextpkg/nextcfg/reader/json"
	"github.com/nextpkg/nextcfg/source"
	"github.com/pkg/errors"
)

type config struct {
	exit  chan bool
	opts  Options
	state bool

	sync.RWMutex
	// the current snapshot
	snap *loader.Snapshot
	// the current values
	val reader.Values
}

type watcher struct {
	lw    loader.Watcher
	rd    reader.Reader
	path  []string
	value reader.Value
}

func newConfig(opts ...Option) (Config, error) {
	var c config

	if err := c.Init(opts...); err != nil {
		return nil, err
	}

	go c.run()

	return &c, nil
}

// Init 初始化配置
func (c *config) Init(opts ...Option) error {
	c.opts = Options{
		Reader: json.NewReader(),
	}
	c.exit = make(chan bool)
	for _, o := range opts {
		o(&c.opts)
	}

	// default loader uses the configured reader
	if c.opts.Loader == nil {
		c.opts.Loader = memory.NewLoader(memory.WithReader(c.opts.Reader))
	}

	err := c.opts.Loader.Load(c.opts.Source...)
	if err != nil {
		return err
	}

	c.snap, err = c.opts.Loader.Snapshot()
	if err != nil {
		return err
	}

	c.val, err = c.opts.Reader.Values(c.snap.ChangeSet)
	if err != nil {
		return err
	}

	return nil
}

// Options 配置选项
func (c *config) Options() Options {
	return c.opts
}

func (c *config) run() {
	watch := func(w loader.Watcher) error {
		for {
			// get change set
			snap, err := w.Next()
			if err != nil {
				return err
			}

			c.Lock()

			if c.snap.Version >= snap.Version {
				c.Unlock()
				continue
			}

			// save
			c.snap = snap

			// set values
			c.val, _ = c.opts.Reader.Values(snap.ChangeSet)

			c.Unlock()
		}
	}

	for {
		w, err := c.opts.Loader.Watch()
		if err != nil {
			time.Sleep(time.Second)
			continue
		}

		done := make(chan bool)

		// the stop watch func
		go func() {
			select {
			case <-done:
			case <-c.exit:
			}

			if err = w.Stop(); err != nil {
				slog.Error("Stop failed.", slog.String("err", err.Error()))
			}
		}()

		// block watch
		if err = watch(w); err != nil {
			// do something better
			time.Sleep(time.Second)
		}

		// close done chan
		close(done)

		// if the config is closed exit
		select {
		case <-c.exit:
			return
		default:
		}
	}
}

// Map To map
func (c *config) Map() map[string]interface{} {
	c.RLock()
	defer c.RUnlock()
	return c.val.Map()
}

// Scan Scan anything
func (c *config) Scan(v interface{}) error {
	c.RLock()
	defer c.RUnlock()
	return c.val.Scan(v)
}

// Sync loads all the sources, calls the parser and updates the config
func (c *config) Sync() error {
	if err := c.opts.Loader.Sync(); err != nil {
		return errors.Wrap(err, "sync() failed")
	}

	snap, err := c.opts.Loader.Snapshot()
	if err != nil {
		return err
	}

	c.Lock()
	defer c.Unlock()

	c.snap = snap

	var val reader.Values
	val, err = c.opts.Reader.Values(snap.ChangeSet)
	if err != nil {
		return err
	}
	c.val = val

	return nil
}

// Close 关闭配置
func (c *config) Close() error {
	select {
	case <-c.exit:
		return nil
	default:
		close(c.exit)
	}
	return nil
}

// Get 获取配置项对应的内容
func (c *config) Get(path ...string) reader.Value {
	c.RLock()
	defer c.RUnlock()

	// did sync actually work?
	if c.val != nil {
		return c.val.Get(path...)
	}

	// no value
	return newValue()
}

// Set 设置配置项及内容
func (c *config) Set(val interface{}, path ...string) {
	c.Lock()
	defer c.Unlock()

	if c.val != nil {
		c.val.Set(val, path...)
	}

	return
}

// Del 删除配置项
func (c *config) Del(path ...string) {
	c.Lock()
	defer c.Unlock()

	if c.val != nil {
		c.val.Del(path...)
	}

	return
}

// Bytes 以[]bytes格式返回配置内容
func (c *config) Bytes() []byte {
	c.RLock()
	defer c.RUnlock()

	if c.val == nil {
		return []byte{}
	}

	return c.val.Bytes()
}

// Load 加载配置
func (c *config) Load(sources ...source.Source) error {
	if err := c.opts.Loader.Load(sources...); err != nil {
		return errors.Wrap(err, "load() failed")
	}

	snap, err := c.opts.Loader.Snapshot()
	if err != nil {
		return err
	}

	c.Lock()
	defer c.Unlock()

	c.snap = snap

	var val reader.Values
	val, err = c.opts.Reader.Values(snap.ChangeSet)
	if err != nil {
		return err
	}
	c.val = val

	return nil
}

// SetState 设置服务状态
func (c *config) SetState(state bool) {
	c.state = state
}

// GetState 获取服务状态
func (c *config) GetState() bool {
	return c.state
}

// Watch 监听器
func (c *config) Watch(path ...string) (Watcher, error) {
	w, err := c.opts.Loader.Watch(path...)
	if err != nil {
		return nil, errors.Wrap(err, "watch() failed")
	}

	return &watcher{
		lw:    w,
		rd:    c.opts.Reader,
		path:  path,
		value: c.Get(path...),
	}, nil
}

// String config
func (c *config) String() string {
	return "config"
}

// Next 监听下一个配置变更
func (w *watcher) Next() (reader.Value, error) {
	for {
		s, err := w.lw.Next()
		if err != nil {
			return nil, err
		}

		// only process changes
		if bytes.Equal(w.value.Bytes(), s.ChangeSet.Data) {
			continue
		}

		v, err := w.rd.Values(s.ChangeSet)
		if err != nil {
			return nil, err
		}

		w.value = v.Get()
		return w.value, nil
	}
}

// Stop 停止监听配置变更
func (w *watcher) Stop() error {
	return w.lw.Stop()
}
