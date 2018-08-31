package common

import (
	"reflect"
)

func getBoolData(value reflect.Value) interface{} {
	switch f := value; value.Kind() {
	case reflect.Bool:
		return f.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return f.Int() != 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return f.Uint() != uint64(0)
	case reflect.Float32, reflect.Float64:
		return f.Float() != float64(0)
	case reflect.String:
		rvalue := f.String()
		if rvalue == "0" || rvalue == "" {
			return false
		} else {
			return true
		}
	case reflect.Interface:
		value := f.Elem()
		if !value.IsValid() {
			goto err //panic("can't convert to bool of type " + value.Kind().String())
		} else {
			return getBoolData(value)
		}
	case reflect.Array, reflect.Slice:
		for i := 0; i < f.Len(); i++ {
			value := f.Index(i)
			switch value.Kind() {
			case reflect.Bool:
				return value.Bool()
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				return value.Int() != 0
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
				return value.Uint() != uint64(0)
			case reflect.Float32, reflect.Float64:
				return value.Float() != Float64(0)
			case reflect.String:
				rval := value.String()
				return rval != "" && rval != "0"
			default:
			}
			break
		}
	case reflect.Ptr:
		// pointer to array or slice or struct?  ok at top level
		// but not embedded (avoid loops)
		if f.Pointer() != 0 {
			return getBoolData(f.Elem())
			// switch a := f.Elem(); a.Kind() {
			// case reflect.Array, reflect.Slice, reflect.Struct, reflect.Map:
			// 	p.buf.WriteByte('&')
			// 	p.printValue(a, verb, depth+1)
			// 	return
			// }
		}
		fallthrough
	//case reflect.Chan, reflect.Func, reflect.UnsafePointer, reflect.Ptr, reflect.Array, reflect.Slice, reflect.Map, reflect.Struct, reflect.Complex64, reflect.Complex128, reflect.Invalid:
	default:
	}
err:
	panic("can't convert to bool of type " + value.Kind().String())
}

func Bool(data interface{}) bool {
	if reflect.TypeOf(data) == nil {
		return false
	}
	num, ok := getBoolData(reflect.ValueOf(data)).(bool)
	if ok {
		return num
	} else {
		panic("can't convert to bool of type " + reflect.TypeOf(data).Kind().String())
	}
}
