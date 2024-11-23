package url

import (
	"testing"

	"github.com/nextpkg/nextcfg/encoder/json"
	"github.com/nextpkg/nextcfg/source"
	"github.com/stretchr/testify/assert"
)

func TestFormat(t *testing.T) {
	at := assert.New(t)

	testCases := []struct {
		url    string
		format string
	}{
		{"configmap/t1", "json"},
		{"configmap/t1.json", "json"},
		{"configmap/t1.yaml", "yaml"},
		{"configmap/t1.unknown", "unknown"},
	}

	defaultEncoder := source.WithEncoder(json.NewEncoder())
	defaultOptions := source.NewOptions(defaultEncoder)
	for _, c := range testCases {
		at.Equal(c.format, formatUrl(c.url, defaultOptions))
	}
}
