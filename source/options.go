package source

import (
	"context"

	"github.com/nextpkg/nextcfg/encoder"
	"github.com/nextpkg/nextcfg/encoder/json"
)

// Options 数据源配置
type Options struct {
	// Encoder
	Encoder encoder.Encoder

	// for alternative data
	Context context.Context
}

// Option 数据源配置
type Option func(o *Options)

// NewOptions 新的数据源配置
func NewOptions(opts ...Option) Options {
	options := Options{
		Encoder: json.NewEncoder(),
		Context: context.Background(),
	}

	for _, o := range opts {
		o(&options)
	}

	return options
}

// WithEncoder sets the source encoder
func WithEncoder(e encoder.Encoder) Option {
	return func(o *Options) {
		o.Encoder = e
	}
}
