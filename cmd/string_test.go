package cmd

import (
	. "github.com/smartystreets/goconvey/convey"
	"github.com/spf13/pflag"
	"testing"
)

func TestStringEnvVar(t *testing.T) {
	Convey("case #1", t, func() {
		set := NewFlagSet("test", pflag.ExitOnError)
		val := "val test"
		set.StringEnvVar(&val, "env", "env test")
		str, err := set.GetString("env")
		So(err, ShouldBeNil)
		So(str, ShouldEqual, "val test")

		val1 := []string{"test1", "test2"}
		set.StringSliceEnvVar(&val1, "slice1", "usage")
		sli, err := set.GetStringSlice("slice1")
		So(err, ShouldBeNil)
		So(sli, ShouldResemble, val1)
	})
}
