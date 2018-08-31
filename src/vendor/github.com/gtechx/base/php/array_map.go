/*
 *此包内主要是php中关联数组相关的操作
 *本包内的数组均指map[string]interface{}
 */
package php

import (
	"reflect"
	"strings"
)

type ArrayMap map[string]interface{}

type Case int

const (
	CASE_UPPER Case = 1
	CASE_LOWER Case = 0
)

/**
计算数组的差集
在第一个ArrayMap中但不在剩余ArrayMap中的元素
*/
func Array_diff(m ...interface{}) interface{} {
	//todo
	if len(m) == 0 {
		return nil
	}

	if len(m) == 1 {
		return m[0]
	}

	kind := reflect.TypeOf(m[0]).Kind()

	if kind != reflect.Slice && kind != reflect.Array && kind != reflect.Map {
		return nil
	}

	diffkind := make([]bool, len(m))
	diffcount := 0
	for i := 1; i < len(m); i += 1 {
		newkind := reflect.TypeOf(m[i]).Kind()
		if kind != newkind || (newkind != reflect.Slice && newkind != reflect.Map && newkind != reflect.Array) {
			diffkind[i] = true
			diffcount += 1
		}
	}

	if len(m)-diffcount <= 1 {
		return m[0]
	}

	m0val := reflect.ValueOf(m[0])

	if kind == reflect.Map {
		keys := m0val.MapKeys()
		existarr := make([]bool, len(keys))
		for i, key := range keys {
			value := m0val.MapIndex(key)
			for j := 1; j < len(m); j += 1 {
				if !diffkind[j] {
					if In_array(value.Interface(), m[j]) {
						existarr[i] = true
					}
				}
			}
		}

		ret := reflect.MakeMap(m0val.Type())
		for i, key := range keys {
			if !existarr[i] {
				value := m0val.MapIndex(key)
				ret.SetMapIndex(key, value)
			}
		}
		return ret.Interface()
	} else if kind == reflect.Slice || kind == reflect.Array {
		existarr := make([]bool, m0val.Len())
		for i := 0; i < m0val.Len(); i++ {
			//m0val.Index(i)
			value := m0val.Index(i)
			for j := 1; j < len(m); j += 1 {
				if !diffkind[j] {
					if In_array(value.Interface(), m[j]) {
						existarr[i] = true
					}
				}
			}
		}

		ret := reflect.MakeSlice(m0val.Type(), 0, 0)

		for i := 0; i < m0val.Len(); i++ {
			if !existarr[i] {
				value := m0val.Index(i)
				ret = reflect.Append(ret, value)
			}
		}

		return ret.Interface()
	}

	return nil
}

func Array_diff_assoc(m ...interface{}) interface{} {
	//todo
	if len(m) == 0 {
		return nil
	}

	if len(m) == 1 {
		return m[0]
	}

	kind := reflect.TypeOf(m[0]).Kind()

	if kind != reflect.Slice && kind != reflect.Array && kind != reflect.Map {
		return nil
	}

	diffkind := make([]bool, len(m))
	diffcount := 0
	for i := 1; i < len(m); i += 1 {
		newkind := reflect.TypeOf(m[i]).Kind()
		if kind != newkind || (newkind != reflect.Slice && newkind != reflect.Map && newkind != reflect.Array) {
			diffkind[i] = true
			diffcount += 1
		}
	}

	if len(m)-diffcount <= 1 {
		return m[0]
	}

	m0val := reflect.ValueOf(m[0])

	if kind == reflect.Map {
		keys := m0val.MapKeys()
		existarr := make([]bool, len(keys))
		for i, key := range keys {
			value := m0val.MapIndex(key)
			for j := 1; j < len(m); j += 1 {
				comkeys := reflect.ValueOf(m[j]).MapKeys()
				if !diffkind[j] {
					if In_array(value.Interface(), m[j]) && is_in_array(key, comkeys) {
						existarr[i] = true
					}
				}
			}
		}

		ret := reflect.MakeMap(m0val.Type())
		for i, key := range keys {
			if !existarr[i] {
				value := m0val.MapIndex(key)
				ret.SetMapIndex(key, value)
			}
		}
		return ret.Interface()
	} else if kind == reflect.Slice || kind == reflect.Array {
		existarr := make([]bool, m0val.Len())
		for i := 0; i < m0val.Len(); i++ {
			//m0val.Index(i)
			value := m0val.Index(i)
			for j := 1; j < len(m); j += 1 {
				if !diffkind[j] {
					if pos_in_array(value.Interface(), m[j]) == i {
						existarr[i] = true
					}
				}
			}
		}

		ret := reflect.MakeSlice(m0val.Type(), 0, 0)

		for i := 0; i < m0val.Len(); i++ {
			if !existarr[i] {
				value := m0val.Index(i)
				ret = reflect.Append(ret, value)
			}
		}

		return ret.Interface()
	}

	return nil
}

func Array_merge(m ...interface{}) interface{} {
	if len(m) == 0 {
		return nil
	}

	if len(m) == 1 {
		return m[0]
	}

	kind := reflect.TypeOf(m[0]).Kind()

	if kind != reflect.Slice && kind != reflect.Array && kind != reflect.Map {
		return nil
	}

	diffkind := make([]bool, len(m))
	diffcount := 0
	for i := 1; i < len(m); i += 1 {
		newkind := reflect.TypeOf(m[i]).Kind()
		if kind != newkind || (newkind != reflect.Slice && newkind != reflect.Map && newkind != reflect.Array) {
			diffkind[i] = true
			diffcount += 1
		}
	}

	if len(m)-diffcount <= 1 {
		return m[0]
	}

	if kind == reflect.Map {
		ret := reflect.MakeMap(reflect.ValueOf(m[0]).Type())

		for j := 0; j < len(m); j += 1 {
			if !diffkind[j] {
				m0val := reflect.ValueOf(m[j])
				keys := m0val.MapKeys()
				for _, key := range keys {
					value := m0val.MapIndex(key)
					ret.SetMapIndex(key, value)

				}
			}
		}

		return ret.Interface()
	} else if kind == reflect.Slice || kind == reflect.Array {
		ret := reflect.MakeSlice(reflect.ValueOf(m[0]).Type(), 0, 0)
		for j := 0; j < len(m); j += 1 {
			if !diffkind[j] {
				m0val := reflect.ValueOf(m[j])

				for i := 0; i < m0val.Len(); i++ {
					value := m0val.Index(i)
					ret = reflect.Append(ret, value)
				}
			}
		}

		return ret.Interface()
	}

	return nil
}

func Array_change_key_case(arr ArrayMap, cs Case) ArrayMap {
	var tmp ArrayMap = ArrayMap{}
	if cs == CASE_LOWER {
		for k, v := range arr {
			tmp[strings.ToLower(k)] = v
		}
	} else if cs == CASE_UPPER {
		for k, v := range arr {
			tmp[strings.ToUpper(k)] = v
		}
	}
	return tmp
}

func Array_keys(array interface{}) []interface{} {
	// keys := make([]interface{}, len(maparr))
	// for k := range maparr {
	// 	keys = append(keys, k)
	// }
	// return keys

	f := reflect.ValueOf(array)
	switch f.Kind() {
	// case reflect.Slice, reflect.Array:
	// 	for i := 0; i < f.Len(); i++ {
	// 		p.printValue(f.Index(i), verb, depth+1)
	// 	}
	case reflect.Map:
		keys := f.MapKeys()
		ret := make([]interface{}, len(keys))
		for i, key := range keys {
			ret[i] = key.Interface()
		}
		return ret
	default:
		return nil
	}
}

func Array_values(array interface{}) []interface{} {
	//array_type := array.([]interface{})
	f := reflect.ValueOf(array)
	switch f.Kind() {
	// case reflect.Slice, reflect.Array:
	// 	for i := 0; i < f.Len(); i++ {
	// 		p.printValue(f.Index(i), verb, depth+1)
	// 	}
	case reflect.Map:
		keys := f.MapKeys()
		ret := make([]interface{}, len(keys))
		for i, key := range keys {
			ret[i] = f.MapIndex(key).Interface()
		}
		return ret
	default:
		return nil
	}
}

func isEqual(a, b reflect.Value) bool {
	switch a.Kind() {
	case reflect.Invalid:
		return false
	case reflect.Bool:
		return a.Bool() == b.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return a.Int() == b.Int()
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return a.Uint() == b.Uint()
	case reflect.Float32:
		return a.Float() == b.Float()
	case reflect.Float64:
		return a.Float() == b.Float()
	case reflect.Complex64:
		return a.Complex() == b.Complex()
	case reflect.Complex128:
		return a.Complex() == b.Complex()
	case reflect.String:
		return a.String() == b.String()
	case reflect.Struct:
	case reflect.Interface:
	case reflect.Ptr:
	case reflect.Chan, reflect.Func, reflect.UnsafePointer:
	case reflect.Map:
	case reflect.Array, reflect.Slice:
	default:
	}

	return false
}

//support:
//string, (u)int8 (u)int16 (u)int (u)int32 (u)int64
func In_array(single interface{}, array interface{}) bool {
	kind := reflect.TypeOf(array).Kind()
	if kind != reflect.Map && kind != reflect.Array && kind != reflect.Slice {
		return false
	}

	f := reflect.ValueOf(array)
	switch f.Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < f.Len(); i++ {
			if reflect.TypeOf(single).Kind() != f.Index(i).Kind() {
				return false
			} else {
				if isEqual(reflect.ValueOf(single), f.Index(i)) {
					return true
				}
			}
		}
	case reflect.Map:
		keys := f.MapKeys()
		//ret := make([]interface{}, len(keys))
		for _, key := range keys {
			//ret[i] = key.Interface()
			if reflect.TypeOf(single).Kind() != f.MapIndex(key).Kind() {
				return false
			} else {
				if isEqual(reflect.ValueOf(single), f.MapIndex(key)) {
					return true
				}
			}
		}
	default:
	}

	return false
}

func pos_in_array(single interface{}, array interface{}) int {
	kind := reflect.TypeOf(array).Kind()
	if kind != reflect.Array && kind != reflect.Slice {
		return -1
	}

	f := reflect.ValueOf(array)
	switch f.Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < f.Len(); i++ {
			if reflect.TypeOf(single).Kind() != f.Index(i).Kind() {
				return -1
			} else {
				if isEqual(reflect.ValueOf(single), f.Index(i)) {
					return i
				}
			}
		}
	default:
	}

	return -1
}

func is_in_array(single reflect.Value, arrray []reflect.Value) bool {
	for _, v := range arrray {
		if isEqual(single, v) {
			return true
		}
	}

	return false
}

func Array_key_exists(key interface{}, array interface{}) bool {
	kind := reflect.TypeOf(array).Kind()
	if kind != reflect.Map {
		return false
	}

	f := reflect.ValueOf(array)
	switch f.Kind() {
	case reflect.Map:
		keys := f.MapKeys()
		//ret := make([]interface{}, len(keys))

		return is_in_array(reflect.ValueOf(key), keys)
		// for _, key := range keys {
		// 	//ret[i] = key.Interface()
		// 	if reflect.TypeOf(single).Kind() != f.MapIndex(key).Kind() {
		// 		return false
		// 	} else {
		// 		if isEqual(reflect.ValueOf(single), f.MapIndex(key)) {
		// 			return true
		// 		}
		// 	}
		// }
	default:
	}

	return false
}
