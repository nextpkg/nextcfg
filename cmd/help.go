package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"log/slog"
)

func init() {
	var help = rootCmd.HelpFunc()
	rootCmd.SetHelpFunc(func(cli *cobra.Command, args []string) {
		help(cli, args)
		os.Exit(0)
	})
}

// CobraFunc 命令行函数
type CobraFunc = func(cli *cobra.Command, args []string) error

// CobraPtr 命令行函数指针
type CobraPtr = *CobraFunc

// AddSubFlags 新增子flags
func AddSubFlags(matchKey, matchValue string, flagSetFunc func() *FlagSet) {
	current := func(cli *cobra.Command, args []string) error {
		parent := cli.Flag(matchKey)
		if parent == nil || parent.Value.String() != matchValue {
			return nil
		}

		cli.Flags().AddFlagSet(flagSetFunc().FlagSet)
		return cli.Flags().Parse(args)
	}
	AppendCommand(&(rootCmd.RunE), current)
}

// AppendCommand 在cobra的运行链中加入新命令
func AppendCommand(origin CobraPtr, current CobraFunc) {
	if origin == nil {
		return
	}

	var help = rootCmd.HelpFunc()
	rootCmd.SetHelpFunc(func(cli *cobra.Command, args []string) {
		if err := current(cli, args); err != nil {
			slog.Error("command error", slog.Any("err", err))
			os.Exit(1)
		}

		help(cli, args)
	})

	runner := *origin
	*origin = func(cli *cobra.Command, args []string) error {
		if err := current(cli, args); err != nil {
			return err
		}

		if runner != nil {
			return runner(cli, args)
		}

		return nil
	}
}
