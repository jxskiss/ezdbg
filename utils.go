package ezdbg

import (
	"reflect"
	"runtime"
	"unicode"
	"unsafe"
)

// Caller returns function name, filename, and the line number of the caller.
// The argument skip is the number of stack frames to ascend, with 0
// identifying the caller of Caller.
func Caller(skip int) (name, file string, line int) {
	pc, file, line, _ := runtime.Caller(skip + 1)
	name = runtime.FuncForPC(pc).Name()
	for i := len(name) - 1; i >= 0; i-- {
		if name[i] == '/' {
			name = name[i+1:]
			break
		}
	}
	pathSepCnt := 0
	for i := len(file) - 1; i >= 0; i-- {
		if file[i] == '/' {
			pathSepCnt++
			if pathSepCnt == 2 {
				file = file[i+1:]
				break
			}
		}
	}
	return
}

// CallerName returns the function name of the direct caller.
// This is a convenient wrapper around Caller.
func CallerName() string {
	name, _, _ := Caller(1)
	return name
}

func isBasicType(typ reflect.Type) bool {
	switch typ.Kind() {
	case reflect.Bool, reflect.String,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
		reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128:
		return true
	}
	return false
}

// isNil tells whether v is nil or the underlying data is nil.
func isNil(v interface{}) bool {
	if v == nil {
		return true
	}
	// eface { _rtype, word }
	ef := *(*[2]unsafe.Pointer)(unsafe.Pointer(&v))
	typ := reflect.TypeOf(v)
	if typ.Kind() == reflect.Slice {
		return *(*unsafe.Pointer)(ef[1]) == nil
	}
	return ef[1] == nil
}

func stringToBytes(s string) []byte {
	type Slice struct {
		Data unsafe.Pointer
		Len  int
		Cap  int
	}
	type String struct {
		Data unsafe.Pointer
		Len  int
	}
	sh := (*String)(unsafe.Pointer(&s))
	bh := &Slice{
		Data: sh.Data,
		Len:  sh.Len,
		Cap:  sh.Len,
	}
	return *(*[]byte)(unsafe.Pointer(bh))
}

func bytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func toSnakeCase(in string) string {
	runes := []rune(in)
	length := len(runes)

	var out []rune
	for i := 0; i < length; i++ {
		if i > 0 && unicode.IsUpper(runes[i]) &&
			((i+1 < length && unicode.IsLower(runes[i+1])) || unicode.IsLower(runes[i-1])) {
			out = append(out, '_')
		}
		out = append(out, unicode.ToLower(runes[i]))
	}

	return string(out)
}
