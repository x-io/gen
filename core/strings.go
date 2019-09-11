package core

// import (
// 	"errors"
// 	"html/template"
// 	"strconv"
// )

// // Values maps a string key to a list of values.
// // It is typically used for query parameters and form values.
// // Unlike in the http.Header map, the keys in a Values map
// // are case-sensitive.
// type Values map[string][]string

// // Get gets the first value associated with the given key.
// // If there are no values associated with the key, Get returns
// // the empty string. To access multiple values, use the map
// // directly.
// func (v Values) Get(key string) string {
// 	if v == nil {
// 		return ""
// 	}
// 	vs := v[key]
// 	if len(vs) == 0 {
// 		return ""
// 	}
// 	return vs[0]
// }

// // Set sets the key to value. It replaces any existing
// // values.
// func (v Values) Set(key, value string) {
// 	v[key] = []string{value}
// }

// // Add adds the value to key. It appends to any existing
// // values associated with key.
// func (v Values) Add(key, value string) {
// 	v[key] = append(v[key], value)
// }

// // Del deletes the values associated with key.
// func (v Values) Del(key string) {
// 	delete(v, key)
// }

// // String returns request form as string
// func (v Values) String(key string) (string, error) {
// 	if v := v.Get(key); len(v) > 0 {
// 		return v, nil
// 	}

// 	return "", errors.New("not exist")
// }

// // Escape returns request form as escaped string
// func (v Values) Escape(key string) (string, error) {
// 	if v := v.Get(key); len(v) > 0 {
// 		return template.HTMLEscapeString(v), nil
// 	}

// 	return "", errors.New("not exist")
// }

// // Int returns request form as int
// func (v Values) Int(key string) (int, error) {
// 	return strconv.Atoi(v.Get(key))
// }

// // Int32 returns request form as int32
// func (v Values) Int32(key string) (int32, error) {
// 	vv, err := strconv.ParseInt(v.Get(key), 10, 32)
// 	return int32(vv), err
// }

// // Int64 returns request form as int64
// func (v Values) Int64(key string) (int64, error) {
// 	return strconv.ParseInt(v.Get(key), 10, 64)
// }

// // UInt returns request form as uint
// func (v Values) UInt(key string) (uint, error) {
// 	vv, err := strconv.ParseUint(v.Get(key), 10, 64)
// 	return uint(vv), err
// }

// // UInt32 returns request form as uint32
// func (v Values) UInt32(key string) (uint32, error) {
// 	vv, err := strconv.ParseUint(v.Get(key), 10, 32)
// 	return uint32(vv), err
// }

// // UInt64 returns request form as uint64
// func (v Values) UInt64(key string) (uint64, error) {
// 	return strconv.ParseUint(v.Get(key), 10, 64)
// }

// // Bool returns request form as bool
// func (v Values) Bool(key string) (bool, error) {
// 	return strconv.ParseBool(v.Get(key))
// }

// // Float32 returns request form as float32
// func (v Values) Float32(key string) (float32, error) {
// 	vv, err := strconv.ParseFloat(v.Get(key), 32)
// 	return float32(vv), err
// }

// // Float64 returns request form as float64
// func (v Values) Float64(key string) (float64, error) {
// 	return strconv.ParseFloat(v.Get(key), 64)
// }

// // MustString returns request form as slice of string with default
// func (v Values) MustString(key string, defaults ...string) string {
// 	if v, err := v.String(key); err == nil {
// 		return v
// 	}

// 	if len(defaults) > 0 {
// 		return defaults[0]
// 	}
// 	return ""
// }

// // MustEscape returns request form as escaped string with default
// func (v Values) MustEscape(key string, defaults ...string) string {
// 	if v, err := v.Escape(key); err == nil {
// 		return v
// 	}

// 	if len(defaults) > 0 {
// 		return defaults[0]
// 	}
// 	return ""
// }

// // MustInt returns request form as int with default
// func (v Values) MustInt(key string, defaults ...int) int {
// 	if v, err := v.Int(key); err == nil {
// 		return v
// 	}
// 	if len(defaults) > 0 {
// 		return defaults[0]
// 	}
// 	return 0
// }

// // MustInt32 returns request form as int32 with default
// func (v Values) MustInt32(key string, defaults ...int32) int32 {
// 	if v, err := v.Int32(key); err == nil {
// 		return v
// 	}
// 	if len(defaults) > 0 {
// 		return defaults[0]
// 	}
// 	return 0
// }

// // MustInt64 returns request form as int64 with default
// func (v Values) MustInt64(key string, defaults ...int64) int64 {
// 	if v, err := v.Int64(key); err == nil {
// 		return v
// 	}
// 	if len(defaults) > 0 {
// 		return defaults[0]
// 	}
// 	return 0
// }

// // MustUInt returns request form as uint with default
// func (v Values) MustUInt(key string, defaults ...uint) uint {
// 	if v, err := v.UInt(key); err == nil {
// 		return v
// 	}
// 	if len(defaults) > 0 {
// 		return defaults[0]
// 	}
// 	return 0

// }

// // MustUInt32 returns request form as uint32 with default
// func (v Values) MustUInt32(key string, defaults ...uint32) uint32 {
// 	if v, err := v.UInt32(key); err == nil {
// 		return v
// 	}
// 	if len(defaults) > 0 {
// 		return defaults[0]
// 	}
// 	return 0
// }

// // MustUInt64 returns request form as uint64 with default
// func (v Values) MustUInt64(key string, defaults ...uint64) uint64 {
// 	if v, err := v.UInt64(key); err == nil {
// 		return v
// 	}
// 	if len(defaults) > 0 {
// 		return defaults[0]
// 	}
// 	return 0
// }

// // MustFloat32 returns request form as float32 with default
// func (v Values) MustFloat32(key string, defaults ...float32) float32 {
// 	if v, err := v.Float32(key); err == nil {
// 		return v
// 	}
// 	if len(defaults) > 0 {
// 		return defaults[0]
// 	}
// 	return 0
// }

// // MustFloat64 returns request form as float64 with default
// func (v Values) MustFloat64(key string, defaults ...float64) float64 {
// 	if v, err := v.Float64(key); err == nil {
// 		return v
// 	}
// 	if len(defaults) > 0 {
// 		return defaults[0]
// 	}
// 	return 0
// }

// // MustBool returns request form as bool with default
// func (v Values) MustBool(key string, defaults ...bool) bool {
// 	if v, err := v.Bool(key); err == nil {
// 		return v
// 	}
// 	if len(defaults) > 0 {
// 		return defaults[0]
// 	}
// 	return false
// }
