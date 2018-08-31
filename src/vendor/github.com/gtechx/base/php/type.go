//这个包是关于数据类型转换的
package php

import (
	"fmt"
	"strconv"
	"strings"
)

// func Intval(val interface{}) (int, error) {
// 	switch v := val.(type) {
// 	case int:
// 		return v, nil
// 	case int8:
// 	case int16:
// 	case int32:
// 	case int64:
// 		return int(v), nil
// 	default:
// 		return strconv.Atoi(fmt.Sprintf("%v", v))
// 	}
// 	return 0, nil
// }

func Type(args ...interface{}) string {
	return strings.Replace(fmt.Sprintf("%T", args[0]), " ", "", -1)
}

func Intval(val interface{}) int {
	switch v := val.(type) {
	case int:
		return v
	case int8:
	case int16:
	case int32:
	case int64:
		return int(v)
	default:
		ret, _ := strconv.Atoi(fmt.Sprintf("%v", v))
		return ret //strconv.Atoi(fmt.Sprintf("%v", v))
	}
	return 0
}

func Strval(val interface{}) string {
	return fmt.Sprintf("%v", val)
}

func Boolval(val interface{}) bool {

	switch v := val.(type) {
	case int, int8, int16, int32, int64:
		if v != 0 {
			return true
		}
		return false
	case uint, uint8, uint16, uint32, uint64:
		if v != 0 {
			return true
		}
		return false
	case bool:
		return v
	case complex64, complex128:
		if v != complex128(0) {
			return true
		}
		return false
	case float32, float64:
		if v != float64(0) {
			return true
		}
	default:
		return v == nil
	}
	return false
}

// func Intval(x float64) int {
// 	return int(x)
// }
