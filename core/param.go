package core

//Params Params
type Params interface {
	Get(key string) string
	String(key string) (string, error)
	Strings(key string) ([]string, error)
	Escape(key string) (string, error)
	Int(key string) (int, error)
	Int32(key string) (int32, error)
	Int64(key string) (int64, error)
	Uint(key string) (uint, error)
	Uint32(key string) (uint32, error)
	Uint64(key string) (uint64, error)
	Bool(key string) (bool, error)
	Float32(key string) (float32, error)
	Float64(key string) (float64, error)
	MustString(key string, defaults ...string) string
	MustStrings(key string, defaults ...[]string) []string
	MustEscape(key string, defaults ...string) string
	MustInt(key string, defaults ...int) int
	MustInt32(key string, defaults ...int32) int32
	MustInt64(key string, defaults ...int64) int64
	MustUint(key string, defaults ...uint) uint
	MustUint32(key string, defaults ...uint32) uint32
	MustUint64(key string, defaults ...uint64) uint64
	MustFloat32(key string, defaults ...float32) float32
	MustFloat64(key string, defaults ...float64) float64
	MustBool(key string, defaults ...bool) bool
	// Param(key string, defaults ...string) string
	// ParamStrings(key string, defaults ...[]string) []string
	// ParamEscape(key string, defaults ...string) string
	// ParamInt(key string, defaults ...int) int
	// ParamInt32(key string, defaults ...int32) int32
	// ParamInt64(key string, defaults ...int32) int64
	// ParamUint(key string, defaults ...int) uint
	// ParamUint32(key string, defaults ...int32) uint32
	// ParamUint64(key string, defaults ...int32) uint64
	// ParamFloat32(key string, defaults ...float32) float32
	// ParamFloat64(key string, defaults ...float64) float64
	// ParamBool(key string, defaults ...bool) bool

	//Set(key string, value interface{})
	//SetParams(params []Param)
}
