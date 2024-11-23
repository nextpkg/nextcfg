package cmd

import (
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type FlagSet struct {
	*pflag.FlagSet
}

func NewFlagSet(name string, handling pflag.ErrorHandling) *FlagSet {
	return &FlagSet{
		FlagSet: pflag.NewFlagSet(name, handling),
	}
}

// StringEnvVar 同时设置字符串命令行与环境变量
// name：环境变量名称（自动转大写）
// usage：参数用途
// def: 参数值
func (f *FlagSet) StringEnvVar(def *string, name string, usage string) {
	if def == nil {
		panic("[StringSliceEnvVar]paramVal is nil")
	}

	envName := strings.ToUpper(name)

	err := viper.BindEnv(envName)
	if err != nil {
		panic(err)
	}

	value := viper.GetString(envName)
	if value == "" {
		value = *def
	}

	f.StringVar(def, name, value, usage)
}

// StringSliceEnvVar 同时设置数组命令行与环境变量
// name：环境变量名称（自动转大写）
// usage：参数用途
// def: 参数值
func (f *FlagSet) StringSliceEnvVar(def *[]string, name, usage string) {
	if def == nil {
		panic("[StringSliceEnvVar]param p is nil")
	}

	envName := strings.ToUpper(name)

	err := viper.BindEnv(envName)
	if err != nil {
		panic(err)
	}

	value := viper.GetStringSlice(envName)
	if len(value) == 0 {
		value = *def
	}

	f.StringSliceVar(def, name, value, usage)
}
