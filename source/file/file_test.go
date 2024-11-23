package file_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/nextpkg/nextcfg"
	"github.com/nextpkg/nextcfg/source/file"
	"github.com/stretchr/testify/require"
)

func TestConfig(t *testing.T) {
	at := require.New(t)

	path := filepath.Join(os.TempDir(), fmt.Sprintf("file1.%d", time.Now().UnixNano()))
	fh, err := os.Create(path)
	at.Nil(err)

	defer func() {
		at.Nil(fh.Close())
		at.Nil(os.Remove(path))
	}()

	_, err = fh.Write([]byte(`{"foo1": "bar1"}`))
	at.Nil(err)

	conf, err := nextcfg.NewConfig()
	at.Nil(err)
	at.Nil(conf.Load(file.NewSource(file.WithPath(path))))

	// simulate multiple close
	go func() {
		at.Nil(conf.Close())
	}()
	go func() {
		at.Nil(conf.Close())
	}()
}

func TestFile(t *testing.T) {
	at := require.New(t)

	path := filepath.Join(os.TempDir(), fmt.Sprintf("file2.%d", time.Now().UnixNano()))
	fh, err := os.Create(path)
	at.Nil(err)

	defer func() {
		at.Nil(fh.Close())
		at.Nil(os.Remove(path))
	}()

	data := []byte(`{"foo2": "bar2"}`)
	_, err = fh.Write(data)
	at.Nil(err)

	f := file.NewSource(file.WithPath(path))
	c, err := f.Read()
	at.Nil(err)

	at.Equal(data, c.Data)
}
