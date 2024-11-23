package json

import (
	"testing"

	"github.com/nextpkg/nextcfg/source"
	"github.com/stretchr/testify/require"
)

func TestValues(t *testing.T) {
	at := require.New(t)

	emptyStr := ""
	testData := []struct {
		data   []byte
		path   []string
		accept interface{}
		value  interface{}
	}{
		{
			[]byte(`{"foo": "bar", "baz": {"bar": "cat"}}`),
			[]string{"foo"},
			emptyStr,
			"bar",
		},
		{
			[]byte(`{"foo": "bar", "baz": {"bar": "cat"}}`),
			[]string{"baz", "bar"},
			emptyStr,
			"cat",
		},
	}

	for _, test := range testData {
		values, err := newValues(&source.ChangeSet{
			Data: test.data,
		})
		at.Nil(err)
		at.NotNil(values)

		at.Nil(values.Get(test.path...).Scan(&test.accept))
		at.Equal(test.value, test.accept)
	}
}
func TestStructArray(t *testing.T) {

	at := require.New(t)

	type tt struct {
		Foo string
	}

	var emptyTSlice []tt

	testData := []struct {
		data   []byte
		accept []tt
		value  []tt
	}{
		{
			[]byte(`[{"foo": "bar"}]`),
			emptyTSlice,
			[]tt{{Foo: "bar"}},
		},
	}

	for _, test := range testData {
		values, err := newValues(&source.ChangeSet{
			Data: test.data,
		})
		at.Nil(err)
		at.Nil(values.Get().Scan(&test.accept))

		at.EqualValues(test.value, test.accept)
	}
}
