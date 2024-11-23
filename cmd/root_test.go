package cmd

import (
	"github.com/spf13/cobra"
	"log/slog"
)

func ExampleExecute() {
	// 模拟设置参数
	Root().SetArgs([]string{
		"--source=none",
		"--bar=ab",
		"foo",
	})

	//  使用flag的方式添加新的参数
	var bar string
	Root().PersistentFlags().StringVar(&bar, "bar", "bar", "this is an usage")

	// 添加子命令
	var foo bool
	Root().AddCommand(&cobra.Command{
		Use:   "foo",
		Short: "a sub command",
		Run: func(cmd *cobra.Command, args []string) {
			foo = true
		},
	})
	slog.Info("ori", slog.Bool("foo", foo))

	// 解析命令
	err := Root().Execute()
	if err != nil {
		panic(err)
	}

	slog.Info("now", slog.Bool("foo", foo))
}
