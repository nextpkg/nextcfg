package nextcfg

import (
	"time"

	"github.com/nextpkg/nextcfg/reader"
)

type value struct{}

func newValue() reader.Value {
	return new(value)
}

// Bool Empty
func (v *value) Bool(bool) bool {
	return false
}

// Int Empty
func (v *value) Int(int) int {
	return 0
}

// String Empty
func (v *value) String(string) string {
	return ""
}

// Float64 Empty
func (v *value) Float64(float64) float64 {
	return 0.0
}

// Duration Empty
func (v *value) Duration(time.Duration) time.Duration {
	return time.Duration(0)
}

// StringSlice Empty
func (v *value) StringSlice([]string) []string {
	return nil
}

// StringMap Empty
func (v *value) StringMap(map[string]string) map[string]string {
	return map[string]string{}
}

// Scan Empty
func (v *value) Scan(interface{}) error {
	return nil
}

// Bytes Empty
func (v *value) Bytes() []byte {
	return nil
}
