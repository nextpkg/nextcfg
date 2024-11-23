package registry

import (
	"github.com/nextpkg/nextcfg"
	"github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestGetSourceType(t *testing.T) {
	convey.Convey("TestGetSourceType", t, func() {
		convey.So(GetCfgSource(), convey.ShouldEqual, "file")
		SetCfgSource("test source")
		convey.So(GetCfgSource(), convey.ShouldEqual, "test source")
	})
}

func TestSetCfgLoader(t *testing.T) {
	convey.Convey("TestSetCfgLoader", t, func() {
		param := ""
		SetCfgLoader("test", func(target string) nextcfg.Loader {
			return func(o *nextcfg.Loaders) {
				param = target
			}
		})

		GetCfgLoader("test", "pa")
		convey.So(param, convey.ShouldEqual, "pa")
		convey.So(GetRegistryList(), convey.ShouldEqual, []string{"test"})
	})
}
