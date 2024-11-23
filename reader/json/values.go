package json

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	simple "github.com/bitly/go-simplejson"
	"github.com/bytedance/sonic"
	"github.com/nextpkg/nextcfg/reader"
	"github.com/nextpkg/nextcfg/source"
)

type jsonValues struct {
	ch *source.ChangeSet
	sj *simple.Json
}

type jsonValue struct {
	*simple.Json
}

func newValues(ch *source.ChangeSet) (reader.Values, error) {
	sj := simple.New()
	data, _ := reader.ReplaceEnvVars(ch.Data)
	if err := sj.UnmarshalJSON(data); err != nil {
		sj.SetPath(nil, string(ch.Data))
	}
	return &jsonValues{ch, sj}, nil
}

// Get 获取Json节点
func (j *jsonValues) Get(path ...string) reader.Value {
	return &jsonValue{j.sj.GetPath(path...)}
}

// Del 删除Json节点
func (j *jsonValues) Del(path ...string) {
	// delete the tree?
	if len(path) == 0 {
		j.sj = simple.New()
		return
	}

	if len(path) == 1 {
		j.sj.Del(path[0])
		return
	}

	val := j.sj.GetPath(path[:len(path)-1]...)
	val.Del(path[len(path)-1])
	j.sj.SetPath(path[:len(path)-1], val.Interface())
	return
}

// Set 设置Json节点
func (j *jsonValues) Set(val interface{}, path ...string) {
	j.sj.SetPath(path, val)
}

// Bytes To Bytes
func (j *jsonValues) Bytes() []byte {
	b, _ := j.sj.MarshalJSON()
	return b
}

// Map To Map
func (j *jsonValues) Map() map[string]interface{} {
	m, _ := j.sj.Map()
	return m
}

// Scan To Anything
func (j *jsonValues) Scan(v interface{}) error {
	b, err := j.sj.MarshalJSON()
	if err != nil {
		return err
	}
	return sonic.Unmarshal(b, v)
}

// String JSON
func (j *jsonValues) String() string {
	return "json"
}

// Bool To Bool
func (j *jsonValue) Bool(def bool) bool {
	b, err := j.Json.Bool()
	if err == nil {
		return b
	}

	str, ok := j.Interface().(string)
	if !ok {
		return def
	}

	b, err = strconv.ParseBool(str)
	if err != nil {
		return def
	}

	return b
}

// Int To Int
func (j *jsonValue) Int(def int) int {
	i, err := j.Json.Int()
	if err == nil {
		return i
	}

	str, ok := j.Interface().(string)
	if !ok {
		return def
	}

	i, err = strconv.Atoi(str)
	if err != nil {
		return def
	}

	return i
}

// String To String
func (j *jsonValue) String(def string) string {
	return j.Json.MustString(def)
}

// Float64 To Float64
func (j *jsonValue) Float64(def float64) float64 {
	f, err := j.Json.Float64()
	if err == nil {
		return f
	}

	str, ok := j.Interface().(string)
	if !ok {
		return def
	}

	f, err = strconv.ParseFloat(str, 64)
	if err != nil {
		return def
	}

	return f
}

// Duration To Duration
func (j *jsonValue) Duration(def time.Duration) time.Duration {
	v, err := j.Json.String()
	if err != nil {
		return def
	}

	var pd time.Duration
	pd, err = time.ParseDuration(v)
	if err != nil {
		return def
	}

	return pd
}

// StringSlice To Slice
func (j *jsonValue) StringSlice(def []string) []string {
	v, err := j.Json.String()
	if err == nil {
		sl := strings.Split(v, ",")
		if len(sl) > 1 {
			return sl
		}
	}
	return j.Json.MustStringArray(def)
}

// StringMap To Map
func (j *jsonValue) StringMap(def map[string]string) map[string]string {
	m, err := j.Json.Map()
	if err != nil {
		return def
	}

	res := map[string]string{}

	for k, v := range m {
		res[k] = fmt.Sprintf("%v", v)
	}

	return res
}

// Scan To Anything
func (j *jsonValue) Scan(v interface{}) error {
	b, err := j.Json.MarshalJSON()
	if err != nil {
		return err
	}
	return sonic.Unmarshal(b, v)
}

// Bytes To Bytes
func (j *jsonValue) Bytes() []byte {
	b, err := j.Json.Bytes()
	if err != nil {
		// try return marshalled
		b, err = j.Json.MarshalJSON()
		if err != nil {
			return []byte{}
		}
		return b
	}
	return b
}
