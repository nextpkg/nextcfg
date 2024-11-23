package nextcfg

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testCfg struct {
	Test1 string `json:"test1"`
	is    bool
}

// Validate ...
func (c *testCfg) Validate() error {
	c.is = true
	return nil
}

func TestCfg_Get(t *testing.T) {
	at := assert.New(t)

	f, err := ioutil.TempFile("", "*.yaml")
	at.Nil(err)
	_, err = f.WriteString("test1: abcd")
	at.Nil(err)
	at.Nil(f.Close())

	Init(&testCfg{Test1: "dd"}, WithFileSource(f.Name()))

	data := GetCopy().(*testCfg)
	at.Equal("abcd", data.Test1)
	at.True(data.is)
}

func TestCfg_GetMulti(t *testing.T) {
	at := assert.New(t)

	f1, err := ioutil.TempFile("", "*.yaml")
	f2, err := ioutil.TempFile("", "*.yaml")
	at.Nil(err)
	_, err = f1.WriteString("test1: abcd")
	_, err = f2.WriteString("test1: edfg")
	at.Nil(err)
	at.Nil(f1.Close())
	at.Nil(f2.Close())

	loader1 := Init(&testCfg{Test1: "dd"}, WithFileSource(f1.Name()))
	loader2 := Init(&testCfg{Test1: "ee"}, WithFileSource(f2.Name()))

	data1 := loader1.GetCopy().(*testCfg)
	data2 := loader2.GetCopy().(*testCfg)
	at.Equal("abcd", data1.Test1)
	at.Equal("edfg", data2.Test1)
	at.True(data1.is)
	at.True(data2.is)
}
