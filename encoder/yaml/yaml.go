package yaml

import (
	"github.com/nextpkg/nextcfg/encoder"
	"gopkg.in/yaml.v3"
)

type yamlEncoder struct{}

// Encode Yaml编码器
func (y yamlEncoder) Encode(v interface{}) ([]byte, error) {
	return yaml.Marshal(v)
}

// Decode Yaml解码器
func (y yamlEncoder) Decode(d []byte, v interface{}) error {
	return yaml.Unmarshal(d, v)
}

// String YAML
func (y yamlEncoder) String() string {
	return "yaml"
}

// NewEncoder Yaml编解码器
func NewEncoder() encoder.Encoder {
	return yamlEncoder{}
}
