package reader

import (
	"github.com/nextpkg/nextcfg/encoder"
	"github.com/nextpkg/nextcfg/encoder/hcl"
	"github.com/nextpkg/nextcfg/encoder/json"
	"github.com/nextpkg/nextcfg/encoder/toml"
	"github.com/nextpkg/nextcfg/encoder/xml"
	"github.com/nextpkg/nextcfg/encoder/yaml"
)

// Options 选项
type Options struct {
	Encoding map[string]encoder.Encoder
}

// Option 选项
type Option func(o *Options)

// NewOptions 编解码器选项
func NewOptions(opts ...Option) Options {
	options := Options{
		Encoding: map[string]encoder.Encoder{
			"json": json.NewEncoder(),
			"yaml": yaml.NewEncoder(),
			"toml": toml.NewEncoder(),
			"xml":  xml.NewEncoder(),
			"hcl":  hcl.NewEncoder(),
			"yml":  yaml.NewEncoder(),
		},
	}
	for _, o := range opts {
		o(&options)
	}
	return options
}

// WithEncoder 新的编解码器
func WithEncoder(e encoder.Encoder) Option {
	return func(o *Options) {
		if o.Encoding == nil {
			o.Encoding = make(map[string]encoder.Encoder)
		}
		o.Encoding[e.String()] = e
	}
}
