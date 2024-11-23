// Package url loads change sets from a url
package url

import (
	"github.com/nextpkg/nextcfg"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/nextpkg/nextcfg/cmd"
	"github.com/nextpkg/nextcfg/source"
	"github.com/pkg/errors"
	"github.com/spf13/pflag"
	"log"
)

// sourceName 数据源名称
const sourceName = "url"

// DefaultURL 默认目标
var DefaultURL = "http://config-center/render/zhiwei/" + filepath.Base(os.Args[0])

func init() {
	registry.SetDefaultSource(sourceName)
	cmd.AddSubFlags(registry.From, sourceName, func() *cmd.FlagSet {
		fs := cmd.NewFlagSet("--source=url", pflag.ContinueOnError)
		fs.StringVar(&DefaultURL, "config_address", DefaultURL, "url system target address")
		return fs
	})
	registry.SetToHub(sourceName, func(target string) nextcfg.Loader {
		if target != "" {
			target = DefaultURL + "/" + target
		} else {
			target = DefaultURL
		}

		return GetLoader(WithURL(target))
	})
}

type urlSource struct {
	url  string
	opts source.Options
}

// Read 使用GET方法获取配置（通过后缀判断格式）
func (u *urlSource) Read() (*source.ChangeSet, error) {
	rsp, err := http.Get(u.url)
	if err != nil {
		return nil, errors.Wrapf(err, "get url %s failed", u.url)
	}
	defer func() {
		err = rsp.Body.Close()
		if err != nil {
			log.Fatalf("close URL:%s config source failed: %s", u.url, err)
		}
	}()

	if rsp.StatusCode != http.StatusOK {
		return nil, errors.New(rsp.Status)
	}

	b, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return nil, err
	}

	cs := &source.ChangeSet{
		Data:      b,
		Format:    formatUrl(u.url, u.opts),
		Timestamp: time.Now(),
		Source:    u.String(),
	}
	cs.Checksum = cs.Sum()

	return cs, nil
}

// Watch 定时检查最新版本
func (u *urlSource) Watch() (source.Watcher, error) {
	return newWatcher(u)
}

// Write is unsupported
func (u *urlSource) Write(*source.ChangeSet) error {
	return nil
}

// String URL
func (u *urlSource) String() string {
	return sourceName
}

// NewSource URL数据源
func NewSource(opts ...source.Option) source.Source {
	options := source.NewOptions(opts...)

	url, ok := options.Context.Value(urlKey{}).(string)
	if !ok || url == "" {
		url = DefaultURL
	}

	return &urlSource{
		url:  url,
		opts: options,
	}
}

// GetLoader sets url source
func GetLoader(opts ...source.Option) nextcfg.Loader {
	return func(l *nextcfg.Loaders) {
		err := l.GetCfg().Load(NewSource(opts...))
		if err != nil {
			log.Fatalln(err)
		} else {
			l.GetCfg().SetState(true)
		}
	}
}
