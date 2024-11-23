package configmap

import (
	"context"

	"github.com/nextpkg/nextcfg/source"
)

type configPathKey struct{}
type groupKey struct{}
type nameKey struct{}
type namespaceKey struct{}

// WithNamespace is an option to add namespace of configmap
func WithNamespace(s string) source.Option {
	return func(o *source.Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, namespaceKey{}, s)
	}
}

// WithGroup configmap分组
func WithGroup(s string) source.Option {
	return func(o *source.Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, groupKey{}, s)
	}
}

// WithName is an option to add name of configmap
func WithName(s string) source.Option {
	return func(o *source.Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, nameKey{}, s)
	}
}

// WithConfigPath option for setting a custom path to kubeconfig
func WithConfigPath(s string) source.Option {
	return func(o *source.Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, configPathKey{}, s)
	}
}
