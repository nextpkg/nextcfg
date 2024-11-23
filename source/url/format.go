package url

import (
	"path/filepath"

	"github.com/nextpkg/nextcfg/source"
)

func formatUrl(p string, opts source.Options) string {
	parts := filepath.Ext(p)
	if parts == "" {
		return opts.Encoder.String()
	}

	return parts[1:]
}
