package consul

import (
	"fmt"
	"github.com/nextpkg/nextcfg"
	"net"
	"time"

	"github.com/hashicorp/consul/api"
	"github.com/nextpkg/nextcfg/cmd"
	"github.com/nextpkg/nextcfg/hub"
	"github.com/nextpkg/nextcfg/source"
	"github.com/spf13/pflag"
	"log"
)

// Currently a single consul reader
type consul struct {
	prefix      string
	stripPrefix string
	addr        string
	opts        source.Options
	client      *api.Client
}

var (
	// DefaultPrefix is the prefix that consul keys will be assumed to have if you
	// haven't specified one
	DefaultPrefix     = "/mRPC/config/"
	DefaultAddress    = ""
	DefaultDatacenter = ""
	DefaultToken      = ""
)

const sourceName = "consul"

func init() {
	hub.SetDefaultSource(sourceName)

	// 此处依赖于hub的初始化参数--source
	cmd.AddSubFlags(hub.SourceFlag, sourceName, func() *pflag.FlagSet {
		fs := pflag.NewFlagSet("--source=consul", pflag.ContinueOnError)
		fs.StringVar(&DefaultPrefix, "config_prefix", DefaultPrefix, "consul system prefix")
		fs.StringVar(&DefaultAddress, "config_address", DefaultAddress, "consul system address")
		fs.StringVar(&DefaultDatacenter, "config_datacenter", DefaultDatacenter, "consul system datacenter")
		fs.StringVar(&DefaultToken, "config_token", DefaultToken, "consul system token")
		return fs
	})

	hub.SetToHub(sourceName, func(target string) nextcfg.Loader {
		return func(l *nextcfg.Loaders) {
			var opts = []source.Option{
				WithPrefix(DefaultPrefix),
				WithAddress(DefaultAddress),
				WithDatacenter(DefaultDatacenter),
				WithToken(DefaultToken),
			}

			err := l.GetCfg().Load(NewSource(opts...))
			if err != nil {
				log.Println(err)
			} else {
				l.GetCfg().SetState(true)
			}
		}
	})
}

// Read latest
func (c *consul) Read() (*source.ChangeSet, error) {
	kv, _, err := c.client.KV().List(c.prefix, nil)
	if err != nil {
		return nil, err
	}

	if kv == nil || len(kv) == 0 {
		return nil, fmt.Errorf("source not found: %s", c.prefix)
	}

	data, err := makeMap(c.opts.Encoder, kv, c.stripPrefix)
	if err != nil {
		return nil, fmt.Errorf("error reading data: %v", err)
	}

	b, err := c.opts.Encoder.Encode(data)
	if err != nil {
		return nil, fmt.Errorf("error reading source: %v", err)
	}

	cs := &source.ChangeSet{
		Timestamp: time.Now(),
		Format:    c.opts.Encoder.String(),
		Source:    c.String(),
		Data:      b,
	}
	cs.Checksum = cs.Sum()

	return cs, nil
}

// Write is unsupported
func (c *consul) Write(*source.ChangeSet) error {
	return nil
}

// String consul
func (c *consul) String() string {
	return sourceName
}

// Watch change
func (c *consul) Watch() (source.Watcher, error) {
	w, err := newWatcher(c.prefix, c.addr, c.String(), c.stripPrefix, c.opts.Encoder)
	if err != nil {
		return nil, err
	}
	return w, nil
}

// NewSource creates a new consul source
func NewSource(opts ...source.Option) source.Source {
	options := source.NewOptions(opts...)

	// use default config
	cfg := api.DefaultConfig()

	// use the consul config passed in the options if any
	if co, ok := options.Context.Value(configKey{}).(*api.Config); ok {
		cfg = co
	}

	// check if there are any address
	a, ok := options.Context.Value(addressKey{}).(string)
	if ok {
		addr, port, err := net.SplitHostPort(a)
		if ae, ok := err.(*net.AddrError); ok && ae.Err == "missing port in address" {
			port = "8500"
			addr = a
			cfg.Address = fmt.Sprintf("%s:%s", addr, port)
		} else if err == nil {
			cfg.Address = fmt.Sprintf("%s:%s", addr, port)
		}
	}

	dc, ok := options.Context.Value(dcKey{}).(string)
	if ok {
		cfg.Datacenter = dc
	}

	token, ok := options.Context.Value(tokenKey{}).(string)
	if ok {
		cfg.Token = token
	}

	// create the client
	client, err := api.NewClient(cfg)
	if err != nil {
		panic(err)
	}

	prefix := DefaultPrefix
	sp := ""
	f, ok := options.Context.Value(prefixKey{}).(string)
	if ok {
		prefix = f
	}

	if b, ok := options.Context.Value(stripPrefixKey{}).(bool); ok && b {
		sp = prefix
	}

	return &consul{
		prefix:      prefix,
		stripPrefix: sp,
		addr:        cfg.Address,
		opts:        options,
		client:      client,
	}
}
