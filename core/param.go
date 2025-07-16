package core

import (
	"errors"
	"html/template"
	"strconv"
	"strings"
)

// Params Params
type Params map[string]string

// Set sets the key to value. It replaces any existing
// values.
func (p Params) Set(key string, value string) {
	p[key] = value
}

// Get gets the first value associated with the given key.
// If there are no values associated with the key, Get returns
// the empty string. To access multiple values, use the map
// directly.
func (p Params) Get(key string) string {
	if p == nil {
		return ""
	}
	return p[key]
}

// String returns request form as string
func (p Params) String(key string) string {
	if p == nil {
		return ""
	}
	return p[key]
}

// String returns request form as string
func (p Params) Strings(key string) []string {
	return strings.Split(p.Get(key), ",")
}

// Escape returns request form as escaped string
func (p Params) Escape(key string) (string, error) {
	vv := p.Get(key)
	if len(vv) > 0 {
		return template.HTMLEscapeString(vv), nil
	}

	return "", errors.New("not exist")
}

// Int returns request form as int
func (p Params) Int(key string) (int, error) {
	return strconv.Atoi(p.Get(key))
}

// Int32 returns request form as int32
func (p Params) Int32(key string) (int32, error) {
	vv, err := strconv.ParseInt(p.Get(key), 10, 32)
	return int32(vv), err
}

// Int64 returns request form as int64
func (p Params) Int64(key string) (int64, error) {
	return strconv.ParseInt(p.Get(key), 10, 64)
}

// Uint returns request form as uint
func (p Params) Uint(key string) (uint, error) {
	vv, err := strconv.ParseUint(p.Get(key), 10, 64)
	return uint(vv), err
}

// Uint32 returns request form as uint32
func (p Params) Uint32(key string) (uint32, error) {
	vv, err := strconv.ParseUint(p.Get(key), 10, 32)
	return uint32(vv), err
}

// Uint64 returns request form as uint64
func (p Params) Uint64(key string) (uint64, error) {
	return strconv.ParseUint(p.Get(key), 10, 64)
}

// Float32 returns request form as float32
func (p Params) Float32(key string) (float32, error) {
	vv, err := strconv.ParseFloat(p.Get(key), 32)
	return float32(vv), err
}

// Float64 returns request form as float64
func (p Params) Float64(key string) (float64, error) {
	return strconv.ParseFloat(p.Get(key), 64)
}

// Bool returns request form as bool
func (p Params) Bool(key string) (bool, error) {
	return strconv.ParseBool(p.Get(key))
}

// MustString returns request form as slice of string with default
func (p Params) MustString(key string, defaults ...string) string {
	if v := p.Get(key); len(v) > 0 {
		return v
	}

	if len(defaults) > 0 {
		return defaults[0]
	}
	return ""
}

// MustStrings returns request form as slice of string with default
func (p Params) MustStrings(key string, defaults ...[]string) []string {
	v := p.Strings(key)
	if len(v) > 0 {
		return v
	}

	if len(defaults) > 0 {
		return defaults[0]
	}
	return []string{}
}

// MustEscape returns request form as escaped string with default
func (p Params) MustEscape(key string, defaults ...string) string {

	if v, err := p.Escape(key); err == nil {
		return v
	}

	if len(defaults) > 0 {
		return defaults[0]
	}
	return ""
}

// MustInt returns request form as int with default
func (p Params) MustInt(key string, defaults ...int) int {
	if v, err := p.Int(key); err == nil {
		return v
	}
	if len(defaults) > 0 {
		return defaults[0]
	}
	return 0
}

// MustInt32 returns request form as int32 with default
func (p Params) MustInt32(key string, defaults ...int32) int32 {
	if v, err := p.Int32(key); err == nil {
		return v
	}
	if len(defaults) > 0 {
		return defaults[0]
	}
	return 0
}

// MustInt64 returns request form as int64 with default
func (p Params) MustInt64(key string, defaults ...int64) int64 {
	if v, err := p.Int64(key); err == nil {
		return v
	}
	if len(defaults) > 0 {
		return defaults[0]
	}
	return 0
}

// MustUint returns request form as uint with default
func (p Params) MustUint(key string, defaults ...uint) uint {
	if v, err := p.Uint(key); err == nil {
		return v
	}
	if len(defaults) > 0 {
		return defaults[0]
	}
	return 0
}

// MustUint32 returns request form as uint32 with default
func (p Params) MustUint32(key string, defaults ...uint32) uint32 {
	if v, err := p.Uint32(key); err == nil {
		return v
	}
	if len(defaults) > 0 {
		return defaults[0]
	}
	return 0
}

// MustUint64 returns request form as uint64 with default
func (p Params) MustUint64(key string, defaults ...uint64) uint64 {
	if v, err := p.Uint64(key); err == nil {
		return v
	}
	if len(defaults) > 0 {
		return defaults[0]
	}
	return 0
}

// MustFloat32 returns request form as float32 with default
func (p Params) MustFloat32(key string, defaults ...float32) float32 {
	if v, err := p.Float32(key); err == nil {
		return v
	}
	if len(defaults) > 0 {
		return defaults[0]
	}
	return 0
}

// MustFloat64 returns request form as float64 with default
func (p Params) MustFloat64(key string, defaults ...float64) float64 {
	if v, err := p.Float64(key); err == nil {
		return v
	}
	if len(defaults) > 0 {
		return defaults[0]
	}
	return 0
}

// MustBool returns request form as bool with default
func (p Params) MustBool(key string, defaults ...bool) bool {
	if v, err := p.Bool(key); err == nil {
		return v
	}
	if len(defaults) > 0 {
		return defaults[0]
	}
	return false
}
