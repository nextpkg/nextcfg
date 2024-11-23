package json

import (
	"testing"

	"github.com/nextpkg/nextcfg/source"
	"github.com/stretchr/testify/require"
)

func TestReader(t *testing.T) {
	at := require.New(t)

	data := []byte(`{"foo": "bar", "baz": {"bar": "cat"}}`)

	testData := []struct {
		path  []string
		value string
	}{
		{
			[]string{"foo"},
			"bar",
		},
		{
			[]string{"baz", "bar"},
			"cat",
		},
	}

	r := NewReader()

	c, err := r.Merge(&source.ChangeSet{Data: data}, &source.ChangeSet{})
	at.Nil(err)

	values, err := r.Values(c)
	at.Nil(err)

	for _, test := range testData {
		v := values.Get(test.path...).String("")
		at.Equal(test.value, v)
	}
}
