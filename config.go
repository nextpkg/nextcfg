// Package nextcfg is an interface for dynamic configuration.
package nextcfg

import (
	"context"

	"github.com/nextpkg/nextcfg/loader"
	"github.com/nextpkg/nextcfg/reader"
	"github.com/nextpkg/nextcfg/source"
)

// Config is an interface abstraction for dynamic configuration
type Config interface {
	// Values provide the reader.Values interface
	reader.Values
	// Init the config
	Init(opts ...Option) error
	// Options in the config
	Options() Options
	// Close Stop the config loader/watcher
	Close() error
	// Load config sources
	Load(source ...source.Source) error
	// Sync Force a source change set sync
	Sync() error
	// Watch a value for changes
	Watch(path ...string) (Watcher, error)
	// SetState 设置服务状态
	SetState(state bool)
	// GetState 获取服务状态
	GetState() bool
}

// Watcher is the config watcher
type Watcher interface {
	Next() (reader.Value, error)
	Stop() error
}

// Options ...
type Options struct {
	Loader loader.Loader
	Reader reader.Reader
	Source []source.Source

	// for alternative data
	Context context.Context
}

// Option ...
type Option func(o *Options)

// Default Config Manager
var DefaultConfig, _ = NewConfig()

// NewConfig returns new config
func NewConfig(opts ...Option) (Config, error) {
	return newConfig(opts...)
}

// Bytes Return config as raw json
func Bytes() []byte {
	return DefaultConfig.Bytes()
}

// Map Return config as a map
func Map() map[string]interface{} {
	return DefaultConfig.Map()
}

// Scan values to a go type
func Scan(v interface{}) error {
	return DefaultConfig.Scan(v)
}

// Sync Force a source ChangeSet sync
func Sync() error {
	return DefaultConfig.Sync()
}

// Get a value from the config
func Get(path ...string) reader.Value {
	return DefaultConfig.Get(path...)
}

// Load config sources
func Load(source ...source.Source) error {
	return DefaultConfig.Load(source...)
}

// Watch a value for changes
func Watch(path ...string) (Watcher, error) {
	return DefaultConfig.Watch(path...)
}
