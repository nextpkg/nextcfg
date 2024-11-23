// Package encoder handles source encoding formats
package encoder

// Encoder 配置编解码接口
type Encoder interface {
	Encode(interface{}) ([]byte, error)
	Decode([]byte, interface{}) error
	String() string
}
