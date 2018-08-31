package common

import (
//"reflect"
//"strconv"
)

const (
	IntType = iota
	UintType
	Int64Type
	Uint64Type
	Int32Type
	Uint32Type
	Int16Type
	Uint16Type
	Int8Type
	Uint8Type
	Float64Type
	Float32Type
	StringType
	BoolType
)

var TypeStr []string = []string{
	"IntType",
	"UintType",
	"Int64Type",
	"Uint64Type",
	"Int32Type",
	"Uint32Type",
	"Int16Type",
	"Uint16Type",
	"Int8Type",
	"Uint8Type",
	"Float64Type",
	"Float32Type",
	"StringType",
	"BoolType",
}

// func getData(data interface{}, itype int) interface{} {
// 	value := reflect.ValueOf(data)
// 	switch f := value; value.Kind() {
// 	case reflect.Bool:
// 		rvalue := f.Bool()
// 		switch itype {
// 		case IntType:
// 			if rvalue {
// 				return int(1)
// 			} else {
// 				return int(0)
// 			}
// 		case UintType:
// 			if rvalue {
// 				return uint(1)
// 			} else {
// 				return uint(0)
// 			}
// 		case Int64Type:
// 			if rvalue {
// 				return int64(1)
// 			} else {
// 				return int64(0)
// 			}
// 		case Uint64Type:
// 			if rvalue {
// 				return uint64(1)
// 			} else {
// 				return uint64(0)
// 			}
// 		case Int32Type:
// 			if rvalue {
// 				return int32(1)
// 			} else {
// 				return int32(0)
// 			}
// 		case Uint32Type:
// 			if rvalue {
// 				return uint32(1)
// 			} else {
// 				return uint32(0)
// 			}
// 		case Int16Type:
// 			if rvalue {
// 				return int16(1)
// 			} else {
// 				return int16(0)
// 			}
// 		case Uint16Type:
// 			if rvalue {
// 				return uint16(1)
// 			} else {
// 				return uint16(0)
// 			}
// 		case Int8Type:
// 			if rvalue {
// 				return int8(1)
// 			} else {
// 				return int8(0)
// 			}
// 		case Uint8Type:
// 			if rvalue {
// 				return uint8(1)
// 			} else {
// 				return uint8(0)
// 			}
// 		case Float64Type:
// 			if rvalue {
// 				return float64(1)
// 			} else {
// 				return float64(0)
// 			}
// 		case Float32Type:
// 			if rvalue {
// 				return float32(1)
// 			} else {
// 				return float32(0)
// 			}
// 		case StringType:
// 			if rvalue {
// 				return "1"
// 			} else {
// 				return "0"
// 			}
// 		case BoolType:
// 			return rvalue
// 		default:
// 			panic("can't convert to " + TypeStr[itype] + " of type " + value.Kind().String())
// 		}
// 	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
// 		rvalue := f.Int()
// 		switch itype {
// 		case IntType:
// 			return int(rvalue)
// 		case UintType:
// 			return uint(rvalue)
// 		case Int64Type:
// 			return int64(rvalue)
// 		case Uint64Type:
// 			return uint64(rvalue)
// 		case Int32Type:
// 			return int32(rvalue)
// 		case Uint32Type:
// 			return uint32(rvalue)
// 		case Int16Type:
// 			return int16(rvalue)
// 		case Uint16Type:
// 			return uint16(rvalue)
// 		case Int8Type:
// 			return int8(rvalue)
// 		case Uint8Type:
// 			return uint8(rvalue)
// 		case Float64Type:
// 			return float64(rvalue)
// 		case Float32Type:
// 			return float32(rvalue)
// 		case StringType:
// 			return strconv.FormatInt(rvalue, 10)
// 		case BoolType:
// 			if rvalue == 0 {
// 				return false
// 			} else {
// 				return true
// 			}
// 		default:
// 			panic("can't convert to " + TypeStr[itype] + " of type " + value.Kind().String())
// 		}
// 	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
// 		rvalue := f.Uint()
// 		switch itype {
// 		case IntType:
// 			return int(rvalue)
// 		case UintType:
// 			return uint(rvalue)
// 		case Int64Type:
// 			return int64(rvalue)
// 		case Uint64Type:
// 			return uint64(rvalue)
// 		case Int32Type:
// 			return int32(rvalue)
// 		case Uint32Type:
// 			return uint32(rvalue)
// 		case Int16Type:
// 			return int16(rvalue)
// 		case Uint16Type:
// 			return uint16(rvalue)
// 		case Int8Type:
// 			return int8(rvalue)
// 		case Uint8Type:
// 			return uint8(rvalue)
// 		case Float64Type:
// 			return float64(rvalue)
// 		case Float32Type:
// 			return float32(rvalue)
// 		case StringType:
// 			return strconv.FormatUint(rvalue, 10)
// 		case BoolType:
// 			if rvalue == 0 {
// 				return false
// 			} else {
// 				return true
// 			}
// 		default:
// 			panic("can't convert to " + TypeStr[itype] + " of type " + value.Kind().String())
// 		}
// 	case reflect.Float32, reflect.Float64:
// 		rvalue := f.Float()
// 		switch itype {
// 		case IntType:
// 			return int(rvalue)
// 		case UintType:
// 			return uint(rvalue)
// 		case Int64Type:
// 			return int64(rvalue)
// 		case Uint64Type:
// 			return uint64(rvalue)
// 		case Int32Type:
// 			return int32(rvalue)
// 		case Uint32Type:
// 			return uint32(rvalue)
// 		case Int16Type:
// 			return int16(rvalue)
// 		case Uint16Type:
// 			return uint16(rvalue)
// 		case Int8Type:
// 			return int8(rvalue)
// 		case Uint8Type:
// 			return uint8(rvalue)
// 		case Float64Type:
// 			return float64(rvalue)
// 		case Float32Type:
// 			return float32(rvalue)
// 		case StringType:
// 			return strconv.FormatFloat(rvalue, 'f', -1, 64)
// 		case BoolType:
// 			if rvalue == 0.0 {
// 				return false
// 			} else {
// 				return true
// 			}
// 		default:
// 			panic("can't convert to " + TypeStr[itype] + " of type " + value.Kind().String())
// 		}
// 	case reflect.String:
// 		rvalue := f.String()
// 		switch itype {
// 		case IntType:
// 			d64, err := strconv.ParseInt(rvalue, 10, 0)
// 			if err != nil {
// 				panic("can't convert to " + TypeStr[itype] + " for string " + f.String())
// 			}
// 			return int(d64)
// 		case UintType:
// 			d64, err := strconv.ParseUint(rvalue, 10, 0)
// 			if err != nil {
// 				panic("can't convert to " + TypeStr[itype] + " for string " + f.String())
// 			}
// 			return uint(d64)
// 		case Int64Type:
// 			d64, err := strconv.ParseInt(rvalue, 10, 0)
// 			if err != nil {
// 				panic("can't convert to " + TypeStr[itype] + " for string " + f.String())
// 			}
// 			return int64(d64)
// 		case Uint64Type:
// 			d64, err := strconv.ParseUint(rvalue, 10, 0)
// 			if err != nil {
// 				panic("can't convert to " + TypeStr[itype] + " for string " + f.String())
// 			}
// 			return uint64(d64)
// 		case Int32Type:
// 			d64, err := strconv.ParseInt(rvalue, 10, 0)
// 			if err != nil {
// 				panic("can't convert to " + TypeStr[itype] + " for string " + f.String())
// 			}
// 			return int32(d64)
// 		case Uint32Type:
// 			d64, err := strconv.ParseUint(rvalue, 10, 0)
// 			if err != nil {
// 				panic("can't convert to " + TypeStr[itype] + " for string " + f.String())
// 			}
// 			return uint32(d64)
// 		case Int16Type:
// 			d64, err := strconv.ParseInt(rvalue, 10, 0)
// 			if err != nil {
// 				panic("can't convert to " + TypeStr[itype] + " for string " + f.String())
// 			}
// 			return int16(d64)
// 		case Uint16Type:
// 			d64, err := strconv.ParseUint(rvalue, 10, 0)
// 			if err != nil {
// 				panic("can't convert to " + TypeStr[itype] + " for string " + f.String())
// 			}
// 			return uint16(d64)
// 		case Int8Type:
// 			d64, err := strconv.ParseInt(rvalue, 10, 0)
// 			if err != nil {
// 				panic("can't convert to " + TypeStr[itype] + " for string " + f.String())
// 			}
// 			return int8(d64)
// 		case Uint8Type:
// 			d64, err := strconv.ParseUint(rvalue, 10, 0)
// 			if err != nil {
// 				panic("can't convert to " + TypeStr[itype] + " for string " + f.String())
// 			}
// 			return uint8(d64)
// 		case Float64Type:
// 			df, err := strconv.ParseFloat(rvalue, 64)
// 			if err != nil {
// 				panic("can't convert to " + TypeStr[itype] + " for string " + f.String())
// 			}
// 			return float64(df)
// 		case Float32Type:
// 			df, err := strconv.ParseFloat(rvalue, 32)
// 			if err != nil {
// 				panic("can't convert to " + TypeStr[itype] + " for string " + f.String())
// 			}
// 			return float32(df)
// 		case StringType:
// 			return rvalue
// 		case BoolType:
// 			if rvalue == "0" {
// 				return false
// 			} else {
// 				return true
// 			}
// 		default:
// 			panic("can't convert to " + TypeStr[itype] + " of type " + value.Kind().String())
// 		}
// 	case reflect.Interface:
// 		value := f.Elem()
// 		if !value.IsValid() {
// 			panic("can't convert to " + TypeStr[itype] + " of type " + value.Kind().String())
// 		} else {
// 			return getData(value, itype)
// 		}
// 	//case reflect.Chan, reflect.Func, reflect.UnsafePointer, reflect.Ptr, reflect.Array, reflect.Slice, reflect.Map, reflect.Struct, reflect.Complex64, reflect.Complex128, reflect.Invalid:
// 	default:
// 		panic("can't convert to " + TypeStr[itype] + " of type " + value.Kind().String())
// 	}
// }

// func Int(data interface{}) int {
// 	num, ok := getData(data, IntType).(int)
// 	if ok {
// 		return num
// 	} else {
// 		panic("can't convert to " + TypeStr[IntType] + " of type " + reflect.TypeOf(data).Kind().String())
// 	}
// }

// func Uint(data interface{}) uint {
// 	num, ok := getData(data, UintType).(uint)
// 	if ok {
// 		return num
// 	} else {
// 		panic("can't convert to " + TypeStr[UintType] + " of type " + reflect.TypeOf(data).Kind().String())
// 	}
// }

// func Int64(data interface{}) int64 {
// 	num, ok := getData(data, Int64Type).(int64)
// 	if ok {
// 		return num
// 	} else {
// 		panic("can't convert to " + TypeStr[Int64Type] + " of type " + reflect.TypeOf(data).Kind().String())
// 	}
// }

// func Uint64(data interface{}) uint64 {
// 	num, ok := getData(data, Uint64Type).(uint64)
// 	if ok {
// 		return num
// 	} else {
// 		panic("can't convert to " + TypeStr[Uint64Type] + " of type " + reflect.TypeOf(data).Kind().String())
// 	}
// }

// func Int32(data interface{}) int32 {
// 	num, ok := getData(data, Int64Type).(int32)
// 	if ok {
// 		return num
// 	} else {
// 		panic("can't convert to " + TypeStr[Int32Type] + " of type " + reflect.TypeOf(data).Kind().String())
// 	}
// }

// func Uint32(data interface{}) uint32 {
// 	num, ok := getData(data, Uint32Type).(uint32)
// 	if ok {
// 		return num
// 	} else {
// 		panic("can't convert to " + TypeStr[Uint32Type] + " of type " + reflect.TypeOf(data).Kind().String())
// 	}
// }

// func Int16(data interface{}) int16 {
// 	num, ok := getData(data, Int16Type).(int16)
// 	if ok {
// 		return num
// 	} else {
// 		panic("can't convert to " + TypeStr[Int16Type] + " of type " + reflect.TypeOf(data).Kind().String())
// 	}
// }

// func Uint16(data interface{}) uint16 {
// 	num, ok := getData(data, Uint16Type).(uint16)
// 	if ok {
// 		return num
// 	} else {
// 		panic("can't convert to " + TypeStr[Uint16Type] + " of type " + reflect.TypeOf(data).Kind().String())
// 	}
// }

// func Int8(data interface{}) int8 {
// 	num, ok := getData(data, Int8Type).(int8)
// 	if ok {
// 		return num
// 	} else {
// 		panic("can't convert to " + TypeStr[Int8Type] + " of type " + reflect.TypeOf(data).Kind().String())
// 	}
// }

// func Uint8(data interface{}) uint8 {
// 	num, ok := getData(data, Uint8Type).(uint8)
// 	if ok {
// 		return num
// 	} else {
// 		panic("can't convert to " + TypeStr[Uint8Type] + " of type " + reflect.TypeOf(data).Kind().String())
// 	}
// }

// func Float32(data interface{}) float32 {
// 	num, ok := getData(data, Float32Type).(float32)
// 	if ok {
// 		return num
// 	} else {
// 		panic("can't convert to " + TypeStr[Float32Type] + " of type " + reflect.TypeOf(data).Kind().String())
// 	}
// }

// func Float64(data interface{}) float64 {
// 	num, ok := getData(data, Float64Type).(float64)
// 	if ok {
// 		return num
// 	} else {
// 		panic("can't convert to " + TypeStr[Float64Type] + " of type " + reflect.TypeOf(data).Kind().String())
// 	}
// }

// func String(data interface{}) string {
// 	num, ok := getData(data, StringType).(string)
// 	if ok {
// 		return num
// 	} else {
// 		panic("can't convert to " + TypeStr[StringType] + " of type " + reflect.TypeOf(data).Kind().String())
// 	}
// }

// func Bool(data interface{}) bool {
// 	num, ok := getData(data, BoolType).(bool)
// 	if ok {
// 		return num
// 	} else {
// 		panic("can't convert to " + TypeStr[BoolType] + " of type " + reflect.TypeOf(data).Kind().String())
// 	}
// }
