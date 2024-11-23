package cmd

import (
	"github.com/spf13/pflag"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/spf13/cobra"
)

func TestAppendCommand(t *testing.T) {
	count := 0
	fn := func(cli *cobra.Command, args []string) error {
		count++
		return nil
	}

	AppendCommand(nil, fn)
	AppendCommand(&(Root().PreRunE), fn)
	AppendCommand(&(Root().PreRunE), fn)

	Convey("case #1", t, func() {
		So(Execute(), ShouldBeNil)
		So(count, ShouldEqual, 2)
	})
}

func TestAddSubFlags(t *testing.T) {
	Convey("case #1", t, func() {
		// 父命令
		var bar string
		Root().PersistentFlags().StringVar(&bar, "key1", "val1", "this is an usage")

		// 子命令
		val := "test1"
		AddSubFlags("key1", "val1", func() *FlagSet {
			fs := NewFlagSet("--key=val", pflag.ContinueOnError)
			fs.StringVar(&val, "addr", val, "target address")
			return fs
		})

		Root().SetArgs([]string{
			"--key=val",
			"--addr=test2",
		})
		So(Execute(), ShouldBeNil)
		So(val, ShouldEqual, "test2")
	})
}
