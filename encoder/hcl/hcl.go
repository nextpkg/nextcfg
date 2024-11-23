package hcl

import (
	"encoding/json"

	"github.com/hashicorp/hcl"
	"github.com/nextpkg/nextcfg/encoder"
)

type hclEncoder struct{}

// Encode HCL编码
func (h hclEncoder) Encode(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

// Decode HCL解码
func (h hclEncoder) Decode(d []byte, v interface{}) error {
	return hcl.Unmarshal(d, v)
}

// String HCL
func (h hclEncoder) String() string {
	return "hcl"
}

// NewEncoder HCL编解码器
func NewEncoder() encoder.Encoder {
	return hclEncoder{}
}
