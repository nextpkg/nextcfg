package nextcfg

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/nextpkg/nextcfg/source"
	"github.com/nextpkg/nextcfg/source/env"
	"github.com/nextpkg/nextcfg/source/file"
	"github.com/nextpkg/nextcfg/source/memory"
	"github.com/stretchr/testify/assert"
)

func createFileForIssue18(t *testing.T, content string) *os.File {

	at := assert.New(t)

	data := []byte(content)
	path := filepath.Join(os.TempDir(), fmt.Sprintf("file.%d", time.Now().UnixNano()))
	fh, err := os.Create(path)
	at.Nil(err)

	_, err = fh.Write(data)
	at.Nil(err)

	return fh
}

func createFileForTest(t *testing.T) *os.File {

	at := assert.New(t)

	data := []byte(`{"foo": "bar"}`)
	path := filepath.Join(os.TempDir(), fmt.Sprintf("file.%d", time.Now().UnixNano()))
	fh, err := os.Create(path)
	at.Nil(err)

	_, err = fh.Write(data)
	at.Nil(err)

	return fh
}

func TestConfigLoadWithGoodFile(t *testing.T) {

	at := assert.New(t)

	fh := createFileForTest(t)
	path := fh.Name()
	defer func() {
		at.Nil(fh.Close())
		at.Nil(os.Remove(path))
	}()

	// Create new config
	conf, err := NewConfig()
	at.Nil(err)

	// Load file source
	err = conf.Load(file.NewSource(
		file.WithPath(path),
	))
	at.Nil(err)
}

func TestConfigLoadWithInvalidFile(t *testing.T) {

	at := assert.New(t)

	fh := createFileForTest(t)
	path := fh.Name()
	defer func() {
		at.Nil(fh.Close())
		at.Nil(os.Remove(path))
	}()

	// Create new config
	conf, err := NewConfig()
	if err != nil {
		t.Fatalf("Expected no error but got %v", err)
	}
	// Load file source
	err = conf.Load(file.NewSource(
		file.WithPath(path),
		file.WithPath("/i/do/not/exists.json"),
	))
	at.NotNil(err)

	if !strings.Contains(fmt.Sprintf("%v", err), "/i/do/not/exists.json") {
		t.Fatalf("Expected error to contain the unexisting file but got %v", err)
	}
}

func TestConfigMerge(t *testing.T) {

	at := assert.New(t)

	fh := createFileForIssue18(t, `{
  "amqp": {
    "host": "rabbit.platform",
    "port": 80
  },
  "handler": {
    "exchange": "springCloudBus"
  }
}`)
	path := fh.Name()
	defer func() {
		at.Nil(fh.Close())
		at.Nil(os.Remove(path))
	}()
	at.Nil(os.Setenv("AMQP_HOST", "rabbit.testing.com"))

	conf, err := NewConfig()
	at.Nil(err)

	err = conf.Load(
		file.NewSource(
			file.WithPath(path),
		),
		env.NewSource(),
	)
	at.Nil(err)

	actualHost := conf.Get("amqp", "host").String("backup")
	at.Equal("rabbit.testing.com", actualHost)
}

func TestConfigWatcherDirtyOverride(t *testing.T) {

	at := assert.New(t)

	n := runtime.GOMAXPROCS(0)
	defer runtime.GOMAXPROCS(n)

	runtime.GOMAXPROCS(1)

	l := 100

	ss := make([]source.Source, l, l)

	for i := 0; i < l; i++ {
		ss[i] = memory.NewSource(memory.WithJSON([]byte(fmt.Sprintf(`{"key%d": "val%d"}`, i, i))))
	}

	conf, _ := NewConfig()

	for _, s := range ss {
		_ = conf.Load(s)
	}
	runtime.Gosched()

	for i := range ss {
		k := fmt.Sprintf("key%d", i)
		v := fmt.Sprintf("val%d", i)
		at.Equal(v, conf.Get(k).String(""))
	}
}
