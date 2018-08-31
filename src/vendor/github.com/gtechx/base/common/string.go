package common

import (
	"reflect"
	"strconv"
)

func getStringData(value reflect.Value) string {
	switch f := value; value.Kind() {
	case reflect.Bool:
		rvalue := f.Bool()
		if rvalue {
			return "1"
		} else {
			return ""
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		rvalue := f.Int()
		return strconv.FormatInt(rvalue, 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		rvalue := f.Uint()
		return strconv.FormatUint(rvalue, 10)
	case reflect.Float32, reflect.Float64:
		rvalue := f.Float()
		return strconv.FormatFloat(rvalue, 'f', -1, 64)
	case reflect.String:
		rvalue := f.String()
		return rvalue
	case reflect.Interface:
		value := f.Elem()
		if f.IsNil() {
			return ""
		}
		if !value.IsValid() {
			goto err //panic("can't convert to string of type " + value.Kind().String())
		} else {
			return getStringData(value)
		}
	case reflect.Array, reflect.Slice:
		bbyte := false
		for i := 0; i < f.Len(); i++ {
			value := f.Index(i)
			if value.Kind() == reflect.Uint8 {
				bbyte = true
			}
			break
		}

		if bbyte {
			buffer := make([]byte, f.Len())
			for i := 0; i < f.Len(); i++ {
				buffer[i] = byte(f.Index(i).Uint())
			}
			return string(buffer)
		} else {
			result := ""
			n := f.Len()
			for i := 0; i < n; i++ {
				result += getStringData(f.Index(i))

				if i < n-1 {
					result += " "
				}
			}
			return result
		}
	case reflect.Struct:
		result := ""
		n := f.NumField()
		for i := 0; i < n; i++ {
			val := f.Field(i)
			if val.Kind() == reflect.Interface && !val.IsNil() {
				val = val.Elem()
			}
			result += getStringData(val)

			if i < n-1 {
				result += " "
			}
		}
		return result
	case reflect.Map:
		keys := f.MapKeys()
		result := ""
		n := len(keys)
		for i, key := range keys {
			result += getStringData(key) + ":" + getStringData(f.MapIndex(key))

			if i < n-1 {
				result += " "
			}
		}
		return result
	case reflect.Invalid:
		return ""
	case reflect.Ptr:
		// pointer to array or slice or struct?  ok at top level
		// but not embedded (avoid loops)
		if f.Pointer() != 0 {
			return getStringData(f.Elem())
			// switch a := f.Elem(); a.Kind() {
			// case reflect.Array, reflect.Slice, reflect.Struct, reflect.Map:
			// 	p.buf.WriteByte('&')
			// 	p.printValue(a, verb, depth+1)
			// 	return
			// }
		}
		fallthrough
	//case reflect.Chan, reflect.Func, reflect.UnsafePointer, reflect.Complex64, reflect.Complex128:
	default:
	}
err:
	panic("can't convert to string of type " + value.Kind().String())
}

func String(data interface{}) string {
	if reflect.TypeOf(data) == nil {
		return ""
	}
	return getStringData(reflect.ValueOf(data))
}
