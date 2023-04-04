package helper

// primitiveType holds all primitive types in Go.
type primitiveType interface {
	string | bool | float32 | float64 |
		int | int8 | int16 | int32 | int64 |
		uint | uint8 | uint16 | uint32 | uint64
}

// Def return the value of given pointer if not nil otherwise the default value
// will be returned.
func Def[T primitiveType](t *T) (n T) {
	if t == nil {
		return n
	}
	return *t
}

// Ptr return pointer of given t.
func Ptr[T primitiveType](t T) *T {
	return &t
}
