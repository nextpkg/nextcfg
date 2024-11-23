package xml

import (
	"encoding/xml"

	"github.com/nextpkg/nextcfg/encoder"
)

type xmlEncoder struct{}

// Encode Xml编码器...
func (x xmlEncoder) Encode(v interface{}) ([]byte, error) {
	return xml.Marshal(v)
}

// Decode Xml解码器...
func (x xmlEncoder) Decode(d []byte, v interface{}) error {
	return xml.Unmarshal(d, v)
}

// String XML
func (x xmlEncoder) String() string {
	return "xml"
}

// NewEncoder Xml编解码器...
func NewEncoder() encoder.Encoder {
	return xmlEncoder{}
}
