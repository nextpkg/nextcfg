package nextcfg

import (
	"context"
	"reflect"
	"sync"

	"github.com/nextpkg/nextcfg/reader"
	"go.uber.org/atomic"
)

// Loader 配置加载器
type Loader func(o *Loaders)

// Loaders 配置加载器
type Loaders struct {
	ctx  context.Context
	once sync.Once
	data atomic.Value
	scan Scanner
	cfg  Config
}

// Scanner scans the config into struct
type Scanner func(reader.Value, interface{}) error

// InitLoader 初始化配置加载器
func InitLoader(lds ...Loader) *Loaders {
	cfg, err := NewConfig()
	if err != nil {
		panic(err)
	}

	l := &Loaders{
		ctx: context.Background(),
		cfg: cfg,
	}

	for _, o := range lds {
		o(l)
	}
	if nil == ld {
		ldLock.Lock()
		defer ldLock.Unlock()
		if nil == ld {
			ld = l
			DefaultConfig = ld.cfg
		}
	}
	return l
}

// Init inits configurator
func Init(t interface{}, lds ...Loader) *Loaders {
	if t == nil {
		panic("you should not use nil template")
	}
	if ld == nil && len(lds) == 0 {
		lds = append(lds, WithContext(context.Background()))
	}

	l := InitLoader(lds...)

	watchDataAddr(t, l)

	return l
}

// watch data addr
func watchDataAddr(t interface{}, l *Loaders) {
	data := t

	variable := reflect.ValueOf(data)

	if variable.Kind() != reflect.Ptr {
		if variable.Kind() != reflect.Struct {
			panic("non-struct")
		}

		pointer := reflect.New(variable.Type())
		pointer.Elem().Set(variable)
		data = pointer.Interface()
	}

	if variable.Kind() == reflect.Ptr {
		if variable.Elem().Kind() != reflect.Struct {
			panic("non-valid-struct")
		}
	}
	if nil == ld.data.Load() {
		ldLock.Lock()
		defer ldLock.Unlock()
		if nil == ld.data.Load() {
			ld.data.Store(data)
		}
	}
	l.data.Store(data)
}

// Reload reload config
func Reload(t interface{}, lds ...Loader) {
	Init(t, lds...)
	GetCopy()
}

// WithContext attach context
func WithContext(ctx context.Context) Loader {
	return func(l *Loaders) {
		l.ctx = ctx
	}
}

// WithScanner sets custom scanner
func WithScanner(scanner Scanner) Loader {
	return func(l *Loaders) {
		l.scan = scanner
	}
}
