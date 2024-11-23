package json

import (
	"github.com/bytedance/sonic"
	"github.com/nextpkg/nextcfg/encoder"
)

type jsonEncoder struct{}

// Encode Json编码...
func (j jsonEncoder) Encode(v interface{}) ([]byte, error) {
	return sonic.Marshal(v)
}

// Decode Json解码...
func (j jsonEncoder) Decode(d []byte, v interface{}) error {
	return sonic.Unmarshal(d, v)
}

// String JSON
func (j jsonEncoder) String() string {
	return "json"
}

// NewEncoder Json编解码器
func NewEncoder() encoder.Encoder {
	return jsonEncoder{}
}
