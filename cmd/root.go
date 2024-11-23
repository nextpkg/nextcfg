package cmd

import (
	"github.com/mitchellh/go-homedir"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log/slog"
)

// Root represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   filepath.Base(os.Args[0]),
	Short: "nextpkg framework",
	// 运行顺序
	PersistentPreRun:   func(cmd *cobra.Command, args []string) {},
	PersistentPreRunE:  func(cmd *cobra.Command, args []string) error { return nil },
	PreRun:             func(cmd *cobra.Command, args []string) {},
	PreRunE:            func(cmd *cobra.Command, args []string) error { return nil },
	Run:                func(cmd *cobra.Command, args []string) {},
	RunE:               func(cmd *cobra.Command, args []string) error { return nil },
	PostRun:            func(cmd *cobra.Command, args []string) {},
	PostRunE:           func(cmd *cobra.Command, args []string) error { return nil },
	PersistentPostRun:  func(cmd *cobra.Command, args []string) {},
	PersistentPostRunE: func(cmd *cobra.Command, args []string) error { return nil },
	// 忽略错误
	FParseErrWhitelist: cobra.FParseErrWhitelist{UnknownFlags: true},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() error {
	cmd, err := rootCmd.ExecuteC()
	if err != nil {
		return err
	}

	if cmd != rootCmd {
		os.Exit(0)
	}

	return nil
}

// Root 根命令
func Root() *cobra.Command {
	return rootCmd
}

func init() {
	cobra.OnInitialize(initConfig)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	// Find home directory.
	home, err := homedir.Dir()
	if err != nil {
		slog.Error("init config failed.", slog.Any("err", err))
		os.Exit(1)
	}

	// Search config in home directory with name ".cmd" (without extension).
	viper.AddConfigPath(home)
	viper.SetConfigName(".cmd")

	// read in environment variables that match
	viper.AutomaticEnv()

	// If a config file is found, read it in.
	err = viper.ReadInConfig()
	if err == nil {
		slog.Info("Using config file.", slog.String("file", viper.ConfigFileUsed()))
	}

	rootCmd.DisableFlagParsing = true
}
