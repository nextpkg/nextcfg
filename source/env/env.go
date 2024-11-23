package env

import (
	"github.com/nextpkg/nextcfg"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"dario.cat/mergo"
	"github.com/nextpkg/nextcfg/source"
)

var (
	// DefaultPrefixes 默认前缀
	DefaultPrefixes []string
)

type env struct {
	prefixes         []string
	strippedPrefixes []string
	opts             source.Options
}

const sourceName = "env"

func init() {
	registry.SetDefaultSource(sourceName)
	registry.SetToHub(sourceName, func(target string) nextcfg.Loader {
		return GetLoader()
	})
}

// Read ...
func (e *env) Read() (*source.ChangeSet, error) {
	var changes map[string]interface{}

	for _, env := range os.Environ() {

		if len(e.prefixes) > 0 || len(e.strippedPrefixes) > 0 {
			notFound := true

			if _, ok := matchPrefix(e.prefixes, env); ok {
				notFound = false
			}

			if match, ok := matchPrefix(e.strippedPrefixes, env); ok {
				env = strings.TrimPrefix(env, match)
				notFound = false
			}

			if notFound {
				continue
			}
		}

		pair := strings.SplitN(env, "=", 2)
		value := pair[1]
		keys := strings.Split(strings.ToLower(pair[0]), "_")
		reverse(keys)

		tmp := make(map[string]interface{})
		for i, k := range keys {
			if i == 0 {
				if intValue, err := strconv.Atoi(value); err == nil {
					tmp[k] = intValue
				} else if boolValue, err := strconv.ParseBool(value); err == nil {
					tmp[k] = boolValue
				} else {
					tmp[k] = value
				}
				continue
			}

			tmp = map[string]interface{}{k: tmp}
		}

		if err := mergo.Map(&changes, tmp); err != nil {
			return nil, err
		}
	}

	b, err := e.opts.Encoder.Encode(changes)
	if err != nil {
		return nil, err
	}

	cs := &source.ChangeSet{
		Format:    e.opts.Encoder.String(),
		Data:      b,
		Timestamp: time.Now(),
		Source:    e.String(),
	}
	cs.Checksum = cs.Sum()

	return cs, nil
}

func matchPrefix(pre []string, s string) (string, bool) {
	for _, p := range pre {
		if strings.HasPrefix(s, p) {
			return p, true
		}
	}

	return "", false
}

func reverse(ss []string) {
	for i := len(ss)/2 - 1; i >= 0; i-- {
		opp := len(ss) - 1 - i
		ss[i], ss[opp] = ss[opp], ss[i]
	}
}

// Watch ...
func (e *env) Watch() (source.Watcher, error) {
	return newWatcher()
}

// Write Empty
func (e *env) Write(*source.ChangeSet) error {
	return nil
}

// String env
func (e *env) String() string {
	return sourceName
}

// NewSource returns a config source for parsing ENV variables.
// Underscores are delimiters for nesting, and all keys are lower cased.
//
// Example:
//
//	"DATABASE_SERVER_HOST=localhost" will convert to
//
//	{
//	    "database": {
//	        "server": {
//	            "host": "localhost"
//	        }
//	    }
//	}
func NewSource(opts ...source.Option) source.Source {
	options := source.NewOptions(opts...)

	var sp []string
	var pre []string
	if p, ok := options.Context.Value(strippedPrefixKey{}).([]string); ok {
		sp = p
	}

	if p, ok := options.Context.Value(prefixKey{}).([]string); ok {
		pre = p
	}

	if len(sp) > 0 || len(pre) > 0 {
		pre = append(pre, DefaultPrefixes...)
	}
	return &env{prefixes: pre, strippedPrefixes: sp, opts: options}
}

// GetLoader sets env source
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
