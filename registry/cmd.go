package registry

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/nextpkg/nextcfg/cmd"
	"github.com/spf13/cobra"
)

var cfgSource string

// cfgFlag 数据源的命令行标志
const cfgFlag = "cfg"

func init() {
	// --cfg
	cmd.Root().PersistentFlags().StringVar(&cfgSource, cfgFlag, "file", "")

	// enrich usage
	next := cmd.Root().HelpFunc()
	cmd.Root().SetHelpFunc(func(cli *cobra.Command, args []string) {
		flag := cli.Flag(cfgFlag)
		flag.Usage = fmt.Sprintf("support %s", append(GetRegistryList(), "none"))
		next(cli, args)
	})
}

func GetCfgSource() string {
	return cfgSource
}

func SetCfgSource(sourceName string) {
	flag := cmd.Root().Flag(cfgFlag)
	if flag == nil {
		slog.Error("--cfg is undefined")
		os.Exit(1)
	}

	flag.DefValue = sourceName
	err := flag.Value.Set(cfgFlag)
	if err != nil {
		slog.Error("reset --cfg value failed. ", slog.Any("err", err))
		os.Exit(1)
	}
}