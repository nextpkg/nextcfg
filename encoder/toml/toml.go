package toml

import (
	"bytes"

	"github.com/BurntSushi/toml"
	"github.com/nextpkg/nextcfg/encoder"
)

type tomlEncoder struct{}

// Encode Toml编码...
func (t tomlEncoder) Encode(v interface{}) ([]byte, error) {
	b := bytes.NewBuffer(nil)
	defer b.Reset()
	err := toml.NewEncoder(b).Encode(v)
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

// Decode Toml解码...
func (t tomlEncoder) Decode(d []byte, v interface{}) error {
	return toml.Unmarshal(d, v)
}

// String TOML
func (t tomlEncoder) String() string {
	return "toml"
}

// NewEncoder Toml编解码器...
func NewEncoder() encoder.Encoder {
	return tomlEncoder{}
}
