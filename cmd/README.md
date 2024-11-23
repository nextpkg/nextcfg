# 1. 这是什么

cmd.Root()是对`cobra`库的封装，服务于本项目，用法可以直接参考`cobra`文档

# 2. 怎么用

## 2.1 引入cmd包

```
import "github.com/nextpkg/nextcfg/cmd"
```

## 2.2 添加全局Flag参数

```
var bar string
cmd.Root().PersistentFlags().StringVar(&bar, "foo", "bar", "this is a foo usage")
```

## 2.3 添加扩展Flag参数

```
val := "test1"
cmd.AddSubFlags("foo", "bar", func() *FlagSet {
    fs := NewFlagSet("--foo=bar", pflag.ContinueOnError)
    fs.StringVar(&val, "addr", val, "target address")
    return fs
})
```

## 2.3 添加命令

```
var foo bool
cmd.Root().AddCommand(&cobra.Command{
	Use:   "foo",
	Short: "a sub command",
	Run: func(cmd *cobra.Command, args []string) {
		foo = true
	},
})
```

## 2.4 传参

> 注意：仅作测试，用来模拟命令行

```
cmd.Root().SetArgs([]string{
    "--foo=ab",
    "foo",
})
```

## 2.5 设置默认值

```
cmd.Root().Flag("foo").DefValue = "cd"
```

## 2.6 解析命令行

```
cmd.Root().Execute()
```
