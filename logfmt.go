package ezdbg

import (
	"bytes"
	"fmt"
	"reflect"
	"sort"
	"strings"
	"unicode/utf8"
)

// Logfmt converts given object to a string in logfmt format, it never
// returns error. Note that only struct and map of basic types are
// supported, non-basic types are simply ignored.
func Logfmt(v interface{}) string {
	if isNil(v) {
		return "null"
	}
	var src []byte
	switch v := v.(type) {
	case []byte:
		src = v
	case string:
		src = stringToBytes(v)
	}
	if src != nil && utf8.Valid(src) {
		srcstr := string(src)
		if bytes.IndexFunc(src, needsQuoteValueRune) != -1 {
			return JSON(srcstr)
		}
		return srcstr
	}

	// simple values
	val := reflect.Indirect(reflect.ValueOf(v))
	if !val.IsValid() {
		return "null"
	}
	if isBasicType(val.Type()) {
		return fmt.Sprint(val)
	}
	if val.Kind() != reflect.Struct && val.Kind() != reflect.Map {
		return "<error: unsupported logfmt type>"
	}

	keyValues := make([]interface{}, 0)
	if val.Kind() == reflect.Map {
		keys := make([]string, 0, val.Len())
		values := make(map[string]interface{}, val.Len())
		for iter := val.MapRange(); iter.Next(); {
			k, v := iter.Key(), reflect.Indirect(iter.Value())
			if !isBasicType(k.Type()) || !v.IsValid() {
				continue
			}
			v = reflect.ValueOf(v.Interface())
			if !v.IsValid() {
				continue
			}
			kstr := fmt.Sprint(k.Interface())
			if isBasicType(v.Type()) {
				keys = append(keys, kstr)
				values[kstr] = v.Interface()
				continue
			}
			if bv, ok := v.Interface().([]byte); ok {
				if len(bv) > 0 && utf8.Valid(bv) {
					keys = append(keys, kstr)
					values[kstr] = string(bv)
				}
				continue
			}
			if v.Kind() == reflect.Slice && isBasicType(v.Elem().Type()) {
				keys = append(keys, kstr)
				values[kstr] = JSON(v.Interface())
				continue
			}
		}
		sort.Strings(keys)
		for _, k := range keys {
			v := values[k]
			keyValues = append(keyValues, k, v)
		}
	} else { // reflect.Struct
		typ := val.Type()
		fieldNum := val.NumField()
		for i := 0; i < fieldNum; i++ {
			field := typ.Field(i)
			// ignore unexported fields which we can't take interface
			if len(field.PkgPath) != 0 {
				continue
			}
			fk := toSnakeCase(field.Name)
			fv := reflect.Indirect(val.Field(i))
			if !(fv.IsValid() && fv.CanInterface()) {
				continue
			}
			if isBasicType(fv.Type()) {
				keyValues = append(keyValues, fk, fv.Interface())
				continue
			}
			if bv, ok := fv.Interface().([]byte); ok {
				if len(bv) > 0 && utf8.Valid(bv) {
					keyValues = append(keyValues, fk, string(bv))
				}
				continue
			}
			if fv.Kind() == reflect.Slice && isBasicType(fv.Elem().Type()) {
				keyValues = append(keyValues, fk, JSON(fv.Interface()))
				continue
			}
		}
	}
	if len(keyValues) == 0 {
		return ""
	}

	buf := &strings.Builder{}
	needSpace := false
	for i := 0; i < len(keyValues); i += 2 {
		k, v := keyValues[i], keyValues[i+1]
		if needSpace {
			buf.WriteByte(' ')
		}
		addLogfmtString(buf, k)
		buf.WriteByte('=')
		addLogfmtString(buf, v)
		needSpace = true
	}
	return buf.String()
}

func addLogfmtString(buf *strings.Builder, val interface{}) {
	str, ok := val.(string)
	if !ok {
		str = fmt.Sprint(val)
	}
	if strings.IndexFunc(str, needsQuoteValueRune) != -1 {
		str = JSON(str)
	}
	buf.WriteString(str)
}

func needsQuoteValueRune(r rune) bool {
	switch r {
	case '\\', '"', '=', '\n', '\r', '\t':
		return true
	default:
		return r <= ' '
	}
}
