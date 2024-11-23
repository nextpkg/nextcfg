package flag

import (
	"errors"
	"flag"
	"github.com/nextpkg/nextcfg"
	"log"
	"strings"
	"time"

	"dario.cat/mergo"
	"github.com/nextpkg/nextcfg/source"
)

type flagSrc struct {
	opts source.Options
}

const sourceName = "flag"

func init() {
	registry.SetDefaultSource(sourceName)
	registry.SetToHub(sourceName, func(target string) nextcfg.Loader {
		return GetLoader()
	})
}

// Read 读取配置
func (fs *flagSrc) Read() (*source.ChangeSet, error) {
	if !flag.Parsed() {
		return nil, errors.New("flags not parsed")
	}

	var changes map[string]interface{}

	visitFn := func(f *flag.Flag) {
		n := strings.ToLower(f.Name)
		keys := strings.FieldsFunc(n, split)
		reverse(keys)

		tmp := make(map[string]interface{})
		for i, k := range keys {
			if i == 0 {
				tmp[k] = f.Value
				continue
			}

			tmp = map[string]interface{}{k: tmp}
		}

		// need to sort error handling
		if err := mergo.Map(&changes, tmp); err != nil {
			log.Println("merge map failed: ", err)
		}

		return
	}

	unset, ok := fs.opts.Context.Value(includeUnsetKey{}).(bool)
	if ok && unset {
		flag.VisitAll(visitFn)
	} else {
		flag.Visit(visitFn)
	}

	b, err := fs.opts.Encoder.Encode(changes)
	if err != nil {
		return nil, err
	}

	cs := &source.ChangeSet{
		Format:    fs.opts.Encoder.String(),
		Data:      b,
		Timestamp: time.Now(),
		Source:    fs.String(),
	}
	cs.Checksum = cs.Sum()

	return cs, nil
}

func split(r rune) bool {
	return r == '-' || r == '_'
}

func reverse(ss []string) {
	for i := len(ss)/2 - 1; i >= 0; i-- {
		opp := len(ss) - 1 - i
		ss[i], ss[opp] = ss[opp], ss[i]
	}
}

// Watch ...
func (fs *flagSrc) Watch() (source.Watcher, error) {
	return source.NewNoopWatcher()
}

// Write Empty
func (fs *flagSrc) Write(*source.ChangeSet) error {
	return nil
}

// String flag
func (fs *flagSrc) String() string {
	return sourceName
}

// NewSource returns a config source for integrating parsed flags.
// Hyphens are delimiters for nesting, and all keys are lower cased.
//
// Example:
//
//	dbHost := flag.String("database-host", "localhost", "the db host name")
//
//	{
//	    "database": {
//	        "host": "localhost"
//	    }
//	}
func NewSource(opts ...source.Option) source.Source {
	return &flagSrc{opts: source.NewOptions(opts...)}
}

// GetLoader sets flag source
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
