// Package file is a file source. Expected format is json
package file

import (
	"github.com/nextpkg/nextcfg"
	"github.com/nextpkg/nextcfg/registry"
	"io/ioutil"
	"os"

	"github.com/nextpkg/nextcfg/source"
	"log"
)

type file struct {
	path string
	opts source.Options
}

var (
	// DefaultPath 默认文件名
	DefaultPath = "config.json"
)

const sourceName = "file"

func init() {
	registry.SetCfgSource(sourceName)
	registry.SetCfgLoader(sourceName, func(target string) nextcfg.Loader {
		return GetLoader(target)
	})
}

// Read 读取文件
func (f *file) Read() (*source.ChangeSet, error) {
	fh, err := os.Open(f.path)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := fh.Close(); err != nil {
			log.Println("close file failed:", err)
		}
	}()

	var b []byte
	b, err = ioutil.ReadAll(fh)
	if err != nil {
		return nil, err
	}

	var info os.FileInfo
	info, err = fh.Stat()
	if err != nil {
		return nil, err
	}

	cs := &source.ChangeSet{
		Format:    format(f.path, f.opts.Encoder),
		Source:    f.String(),
		Timestamp: info.ModTime(),
		Data:      b,
	}
	cs.Checksum = cs.Sum()

	return cs, nil
}

// String file
func (f *file) String() string {
	return sourceName
}

// Watch ...
func (f *file) Watch() (source.Watcher, error) {
	if _, err := os.Stat(f.path); err != nil {
		return nil, err
	}
	return newWatcher(f)
}

// Write Empty
func (f *file) Write(*source.ChangeSet) error {
	return nil
}

// NewSource 文件数据源
func NewSource(opts ...source.Option) source.Source {
	options := source.NewOptions(opts...)
	path := DefaultPath
	f, ok := options.Context.Value(filePathKey{}).(string)
	if ok {
		path = f
	}
	return &file{opts: options, path: path}
}

// GetLoader sets file source
func GetLoader(whereIs ...string) nextcfg.Loader {
	return func(l *nextcfg.Loaders) {
		for _, path := range whereIs {
			log.Println("load path:", path)
			err := l.GetCfg().Load(NewSource(WithPath(path)))
			if err != nil {
				log.Println(err)
			} else {
				l.GetCfg().SetState(true)
			}
		}
	}
}

// LoadFile is shorthand for creating a file source and loading it
func LoadFile(path string) error {
	return nextcfg.Load(NewSource(
		WithPath(path),
	))
}
