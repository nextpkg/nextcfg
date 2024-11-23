package registry

import (
	"github.com/nextpkg/nextcfg"
	"sync"
)

type Registry func(target string) nextcfg.Loader

var (
	cfgRegistry    = map[string]Registry{}
	cfgRegistryMtx = sync.RWMutex{}
)

// SetCfgLoader 设置配置加载器
// sourceName 数据源名称，例如 file/url...
// sourceLoad 数据源加载器，例如 file/url 的具体加载的实例
func SetCfgLoader(sourceName string, sourceLoad Registry) {
	cfgRegistryMtx.Lock()
	defer cfgRegistryMtx.Unlock()

	cfgRegistry[sourceName] = sourceLoad
}

// GetCfgLoader 获取配置加载器
// sourceName 数据源名称，例如 file/url...
// loaderPara 加载器所使用的参数
func GetCfgLoader(sourceName, loaderPara string) *nextcfg.Loaders {
	cfgRegistryMtx.RLock()
	defer cfgRegistryMtx.RUnlock()

	cc, ok := cfgRegistry[sourceName]
	if !ok {
		return nil
	}

	return nextcfg.InitLoader(cc(loaderPara))
}

func GetRegistryList() []string {
	cfgRegistryMtx.RLock()
	defer cfgRegistryMtx.RUnlock()

	var l []string
	for k := range cfgRegistry {
		l = append(l, k)
	}

	return l
}
