## 能力

| encoder | loader | reader | secrets   | source  |
|---------|--------|--------|-----------|---------|
| hcl     | memory | json   | box       | consul  |
| json    |        |        | secretbox | env     |
| toml    |        |        |           | file    |
| xml     |        |        |           | flag    |
| yaml    |        |        |           | memory  |
|         |        |        |           | rainbow |
|         |        |        |           | url     |

## Import

```go
import "github.com/nextpkg/nextcfg"
```

## Initialize

### File source

```go
// 配置中心根据后缀名来判断用的解析器是哪个，yaml后缀名将使用yaml格式解析，json将以json格式解析
config.Init(NewConfig(), nextcfg.WithFileSource("my.yaml"))
```

### Multi config source

```go
// 多数据源的配置会被合并，按从前到后的顺序合并
config.Init(NewConfig(), nextcfg.WithFileSource("my.yaml"))
```

### Load new source

```go
s := toml.NewSource()

err := nextcfg.Load(s)
if err != nil {
    panic(err)
}
```

## Example

```go
// config is framework configuration
type config struct {
    // 示例：需用json作为字段的tag，但配置文件可以不用json格式
    Console  string  `json:"console"`
}

// Get 获取配置（原子操作）
func Get() *config {
    return config.GetCopy().(*config)
}

// NewConfig get framework configuration
func NewConfig() *config {
    return &config{
        // 示例：允许使用默认值，在配置中心不可用时将使用以下默认值
        Console: "abc",
    }
}

// Validate 验证配置是否有效，在配置被合理检查后才会被更新(未导出变量无法被复制)
func (c *config) Validate() error {
    // 如果返回值不是nil，则不会触发配置更新
    return nil
}

// Revoke 在新配置生效后，旧配置会调用Revoke取消
func (c *config) Revoke() {
		log.Info("revoke changed? ", c.Console)
}
```