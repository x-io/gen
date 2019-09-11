package core

import (
	"errors"
)

//ViewData ViewData
type ViewData map[string]interface{}

// Set sets the key to value. It replaces any existing
// values.
func (v ViewData) Set(key string, value interface{}) {
	v[key] = value
}

// Get gets the first value associated with the given key.
// If there are no values associated with the key, Get returns
// the empty string. To access multiple values, use the map
// directly.
func (v ViewData) Get(key string) interface{} {
	if v == nil {
		return nil
	}
	return v[key]
}

// String returns request form as string
func (v ViewData) String(key string) (string, error) {
	if vv, ok := v.Get(key).(string); ok {
		return vv, nil
	}

	return "", errors.New("not exist")
}
