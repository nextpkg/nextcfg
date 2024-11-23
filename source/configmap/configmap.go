package configmap

import (
	"fmt"
	"strings"

	"github.com/nextpkg/nextcfg"
	"github.com/nextpkg/nextcfg/cmd"
	"github.com/nextpkg/nextcfg/hub"
	"github.com/nextpkg/nextcfg/source"
	"github.com/spf13/pflag"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"log"
)

type configmap struct {
	opts       source.Options
	client     *kubernetes.Clientset
	err        error
	group      string
	name       string
	namespace  string
	configPath string
	format     string
}

// Predefined variables
var (
	DefaultConfigPath = ""
	DefaultNamespace  = "default"
	DefaultGroup      = "config"
)

const sourceName = "configmap"

func init() {
	hub.SetDefaultSource(sourceName)

	// 此处依赖于hub的初始化参数--source
	cmd.AddSubFlags(hub.SourceFlag, sourceName, func() *pflag.FlagSet {
		fs := pflag.NewFlagSet("--source=configmap", pflag.ContinueOnError)
		fs.StringVar(&DefaultConfigPath, "config_path", DefaultConfigPath, "configmap system path")
		fs.StringVar(&DefaultNamespace, "config_namespace", DefaultNamespace, "configmap system namespace")
		fs.StringVar(&DefaultGroup, "config_group", DefaultGroup, "configmap system group")
		return fs
	})

	// 注册
	hub.SetToHub(sourceName, func(target string) nextcfg.Loader {
		return GetLoader(DefaultGroup, target,
			WithConfigPath(DefaultConfigPath),
			WithNamespace(DefaultNamespace),
			WithGroup(DefaultGroup),
		)
	})
}

// Read 读取配置
func (k *configmap) Read() (*source.ChangeSet, error) {
	if k.err != nil {
		return nil, k.err
	}

	cmp, err := k.client.CoreV1().ConfigMaps(k.namespace).Get(k.opts.Context, k.group, v1.GetOptions{})
	if err != nil {
		return nil, err
	}

	data, ok := cmp.Data[k.name]
	if !ok {
		return nil, fmt.Errorf("group:'%s' -> no such key: '%s'", k.group, k.name)
	}

	cs := &source.ChangeSet{
		Format:    k.format,
		Source:    k.String(),
		Data:      []byte(data),
		Timestamp: cmp.CreationTimestamp.Time,
	}

	cs.Checksum = cs.Sum()

	return cs, nil
}

// Write is unsupported
func (k *configmap) Write(*source.ChangeSet) error {
	return nil
}

// String ...
func (k *configmap) String() string {
	return sourceName
}

// Watch 监控配置变化
func (k *configmap) Watch() (source.Watcher, error) {
	if k.err != nil {
		return nil, k.err
	}

	return newWatcher(k.group, k.name, k.namespace, k.format, k.client)
}

// NewSource is a factory function
func NewSource(opts ...source.Option) source.Source {
	var (
		options    = source.NewOptions(opts...)
		name       = ""
		group      = DefaultGroup
		configPath = DefaultConfigPath
		namespace  = DefaultNamespace
	)

	cfg, ok := options.Context.Value(groupKey{}).(string)
	if ok {
		group = cfg
	}
	if cfg == "" {
		group = DefaultGroup
	}

	cfg, ok = options.Context.Value(nameKey{}).(string)
	if ok {
		name = cfg
	}

	cfg, ok = options.Context.Value(configPathKey{}).(string)
	if ok {
		configPath = cfg
	}

	cfg, ok = options.Context.Value(namespaceKey{}).(string)
	if ok {
		namespace = cfg
	}

	// TODO handle if the client fails what to do current return does not support error
	client, err := getClient(configPath)

	// 提取格式
	format := options.Encoder.String()
	parts := strings.Split(name, ".")
	if len(parts) > 1 {
		format = parts[len(parts)-1]
	}

	return &configmap{
		err:        err,
		client:     client,
		opts:       options,
		name:       name,
		group:      group,
		configPath: configPath,
		namespace:  namespace,
		format:     format,
	}
}

// GetLoader sets configmap source
func GetLoader(group, name string, opts ...source.Option) nextcfg.Loader {
	return func(l *nextcfg.Loaders) {
		options := []source.Option{
			WithGroup(group),
			WithName(name),
		}

		if len(opts) > 0 {
			options = append(options, opts...)
		}

		err := l.GetCfg().Load(NewSource(options...))
		if err != nil {
			log.Println(err)
		} else {
			l.GetCfg().SetState(true)
		}
	}
}
