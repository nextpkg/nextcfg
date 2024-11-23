package nextcfg

import (
	"fmt"
	"log/slog"
	"reflect"
	"sync"

	"github.com/mohae/deepcopy"
	"github.com/nextpkg/nextcfg/reader"
	"github.com/pkg/errors"
)

// ld default loaders
var ld *Loaders
var ldLock = &sync.Mutex{}

// Validate 验证配置
type Validate interface {
	Validate() error
}

// Revoke 配置撤销
type Revoke interface {
	Revoke()
}

// GetCopy returns config with watching
func GetCopy() interface{} {
	ld.once.Do(func() {
		err := ld.load(ld.cfg.Get())
		if nil != err {
			panic(fmt.Sprintf("GetCopy Failed: %+v", err))
		}

		go ld.watch()
	})

	return ld.data.Load()
}

// GetOnce returns config once
func GetOnce() interface{} {
	ld.once.Do(func() {
		if err := ld.load(ld.cfg.Get()); nil != err {
			panic(err)
		}
	})

	// reset
	defer func() {
		ld.once = sync.Once{}
	}()

	return ld.data.Load()
}

// GetCopy by loaders returns config with watching
func (l *Loaders) GetCopy() interface{} {
	l.once.Do(func() {
		err := l.load(l.cfg.Get())
		if err != nil {
			panic(err)
		}

		go l.watch()
	})
	return l.data.Load()
}

// GetOnce by loaders returns config once
func (l *Loaders) GetOnce() interface{} {
	l.once.Do(func() {
		if err := l.load(l.cfg.Get()); nil != err {
			panic(err)
		}
	})

	// reset
	defer func() {
		l.once = sync.Once{}
	}()

	return l.data.Load()
}

// GetCfg get loader config
func (l *Loaders) GetCfg() Config {
	return l.cfg
}

func (l *Loaders) load(r reader.Value) error {
	replica := l.data.Load()

	shadow := deepcopy.Copy(replica)

	// l.scan是自定义的配置扫描函数，r.scan是默认的配置扫描函数
	if l.scan == nil {
		err := r.Scan(shadow)
		if err != nil {
			return err
		}
	} else {
		err := l.scan(r, shadow)
		if err != nil {
			return err
		}
	}

	hi := reflect.TypeOf(shadow)
	ht := reflect.TypeOf((*Validate)(nil)).Elem()
	if hi.Implements(ht) {
		err := shadow.(Validate).Validate()
		if err != nil {
			return errors.Wrap(err, "validate() failed")
		}
	}

	l.data.Store(shadow)

	hi = reflect.TypeOf(replica)
	ht = reflect.TypeOf((*Revoke)(nil)).Elem()
	if hi.Implements(ht) {
		replica.(Revoke).Revoke()
	}

	return nil
}

func (l *Loaders) watch() {
	w, err := l.cfg.Watch()
	if err != nil {
		slog.Error("watch failed.", slog.String("err", err.Error()))
	}

	go func() {
		select {
		case <-l.ctx.Done():
			err = w.Stop()
			if err != nil {
				slog.Error("stop watcher failed.", slog.String("err", err.Error()))
			}
		}
	}()

	for {
		r, err := w.Next()
		if err != nil {
			slog.Error("watch next failed.", slog.String("err", err.Error()))
			ld.once = sync.Once{}
			break
		}

		slog.Info("Configuration reloading...")

		err = l.load(r)
		if err != nil {
			slog.Error("load failed.", slog.String("err", err.Error()))
			continue
		}

		slog.Info("Configuration Reloaded!")
	}
}
