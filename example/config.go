// Package example 示例的配置结构体
package example

import (
	"github.com/nextpkg/nextcfg"
	"log"
)

// Config is framework configuration
type Config struct {
	MysqlHost string `json:"mysql_host"` // 示例：需用json作为字段的tag，但配置文件可以不用json格式
	test      string
}

// NewConfig get framework configuration
func NewConfig() *Config {
	return &Config{
		MysqlHost: "old.db", // 示例：允许使用默认值，在配置中心不可用时将使用以下默认值
		test:      "test1",  // 示例：未导出变量
	}
}

// Validate 验证配置是否有效，配置接受检查后才会被更新
func (c *Config) Validate() error {
	// 如果配置没问题，返回nil；否则，返回error；
	// 如果返回值是error则不会触发更新
	log.Println("load mysql host: ", c.MysqlHost)

	return nil
}

// Revoke 在新配置生效后，旧配置会调用Revoke取消
func (c *Config) Revoke() {
	log.Println("revoke mysql host: ", c.MysqlHost)
}

// Get 获取配置（原子操作）
func Get() *Config {
	return nextcfg.GetCopy().(*Config)
}
