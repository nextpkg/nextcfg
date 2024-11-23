package url

import (
	"context"

	"github.com/nextpkg/nextcfg/source"
)

type urlKey struct{}

// WithURL ...
func WithURL(u string) source.Option {
	return func(o *source.Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, urlKey{}, u)
	}
}
