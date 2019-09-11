package core

import (
	"errors"
	"html/template"
	"strconv"
	"strings"
)

// Values maps a string key to a list of values.
// It is typically used for query parameters and form values.
// Unlike in the http.Header map, the keys in a Values map
// are case-sensitive.
type Values string

func (v Values) get() string {
	return string(v)
}

// String returns request form as string
func (v Values) String() string {
	return string(v)
}

// Strings returns request form as []string
func (v Values) Strings() []string {
	return strings.Split(string(v), ",")
}

// Escape returns request form as escaped string
func (v Values) Escape() (string, error) {
	if len(v) > 0 {
		return template.HTMLEscapeString(v.get()), nil
	}

	return "", errors.New("not exist")
}

// Int returns request form as int
func (v Values) Int() (int, error) {
	return strconv.Atoi(v.get())
}

// Int32 returns request form as int32
func (v Values) Int32() (int32, error) {
	vv, err := strconv.ParseInt(v.get(), 10, 32)
	return int32(vv), err
}

// Int64 returns request form as int64
func (v Values) Int64() (int64, error) {
	return strconv.ParseInt(v.get(), 10, 64)
}

// Ints returns request form as []int
func (v Values) Ints(sep string) (data []int, err error) {
	vv := strings.Split(v.get(), sep)
	data = make([]int, 0)
	var i int
	for _, d := range vv {
		i, err = strconv.Atoi(d)
		if err != nil {
			continue
		}
		data = append(data, i)
	}

	return
}

// Int32s returns request form as []int32
func (v Values) Int32s(sep string) (data []int32, err error) {
	vv := strings.Split(v.get(), sep)
	data = make([]int32, 0)
	var i int64
	for _, d := range vv {
		i, err = strconv.ParseInt(d, 10, 32)
		if err != nil {
			continue
		}
		data = append(data, int32(i))
	}

	return
}

// Int64s returns request form as []int64
func (v Values) Int64s(sep string) (data []int64, err error) {
	vv := strings.Split(v.get(), sep)
	data = make([]int64, 0)
	var i int64
	for _, d := range vv {
		i, err = strconv.ParseInt(d, 10, 64)
		if err != nil {
			continue
		}
		data = append(data, i)
	}

	return
}

// UInt returns request form as uint
func (v Values) UInt() (uint, error) {
	vv, err := strconv.ParseUint(v.get(), 10, 64)
	return uint(vv), err
}

// UInt32 returns request form as uint32
func (v Values) UInt32() (uint32, error) {
	vv, err := strconv.ParseUint(v.get(), 10, 32)
	return uint32(vv), err
}

// UInt64 returns request form as uint64
func (v Values) UInt64() (uint64, error) {
	return strconv.ParseUint(v.get(), 10, 64)
}

// Bool returns request form as bool
func (v Values) Bool() (bool, error) {
	return strconv.ParseBool(v.get())
}

// Float32 returns request form as float32
func (v Values) Float32() (float32, error) {
	vv, err := strconv.ParseFloat(v.get(), 32)
	return float32(vv), err
}

// Float64 returns request form as float64
func (v Values) Float64() (float64, error) {
	return strconv.ParseFloat(v.get(), 64)
}

// MustString returns request form as slice of string with default
func (v Values) MustString(defaults ...string) string {
	if v != "" {
		return string(v)
	}

	if len(defaults) > 0 {
		return defaults[0]
	}
	return ""
}

// MustEscape returns request form as escaped string with default
func (v Values) MustEscape(defaults ...string) string {
	if v, err := v.Escape(); err == nil {
		return v
	}

	if len(defaults) > 0 {
		return defaults[0]
	}
	return ""
}

// MustInt returns request form as int with default
func (v Values) MustInt(defaults ...int) int {
	if v, err := v.Int(); err == nil {
		return v
	}
	if len(defaults) > 0 {
		return defaults[0]
	}
	return 0
}

// MustInt32 returns request form as int32 with default
func (v Values) MustInt32(defaults ...int32) int32 {
	if v, err := v.Int32(); err == nil {
		return v
	}
	if len(defaults) > 0 {
		return defaults[0]
	}
	return 0
}

// MustInt64 returns request form as int64 with default
func (v Values) MustInt64(defaults ...int64) int64 {
	if v, err := v.Int64(); err == nil {
		return v
	}
	if len(defaults) > 0 {
		return defaults[0]
	}
	return 0
}

// MustUInt returns request form as uint with default
func (v Values) MustUInt(defaults ...uint) uint {
	if v, err := v.UInt(); err == nil {
		return v
	}
	if len(defaults) > 0 {
		return defaults[0]
	}
	return 0

}

// MustUInt32 returns request form as uint32 with default
func (v Values) MustUInt32(defaults ...uint32) uint32 {
	if v, err := v.UInt32(); err == nil {
		return v
	}
	if len(defaults) > 0 {
		return defaults[0]
	}
	return 0
}

// MustUInt64 returns request form as uint64 with default
func (v Values) MustUInt64(defaults ...uint64) uint64 {
	if v, err := v.UInt64(); err == nil {
		return v
	}
	if len(defaults) > 0 {
		return defaults[0]
	}
	return 0
}

// MustFloat32 returns request form as float32 with default
func (v Values) MustFloat32(defaults ...float32) float32 {
	if v, err := v.Float32(); err == nil {
		return v
	}
	if len(defaults) > 0 {
		return defaults[0]
	}
	return 0
}

// MustFloat64 returns request form as float64 with default
func (v Values) MustFloat64(defaults ...float64) float64 {
	if v, err := v.Float64(); err == nil {
		return v
	}
	if len(defaults) > 0 {
		return defaults[0]
	}
	return 0
}

// MustInts returns request form as int with default
func (v Values) MustInts(defaults ...[]int) []int {
	if v, err := v.Ints(","); err == nil {
		return v
	}
	if len(defaults) > 0 {
		return defaults[0]
	}
	return nil
}

// // MustInt32s returns request form as int32 with default
// func (v Values) MustInt32s(defaults ...[]int32) []int32 {
// 	if v, err := v.Int32s(); err == nil {
// 		return v
// 	}
// 	if len(defaults) > 0 {
// 		return defaults[0]
// 	}
// 	return nil
// }

// // MustInt64s returns request form as int64 with default
// func (v Values) MustInt64s(defaults ...[]int64) []int64 {
// 	if v, err := v.Int64s(); err == nil {
// 		return v
// 	}
// 	if len(defaults) > 0 {
// 		return defaults[0]
// 	}
// 	return nil
// }

// MustBool returns request form as bool with default
func (v Values) MustBool(defaults ...bool) bool {
	if v, err := v.Bool(); err == nil {
		return v
	}
	if len(defaults) > 0 {
		return defaults[0]
	}
	return false
}
