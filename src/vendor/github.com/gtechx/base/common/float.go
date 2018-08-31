package common

import (
	"bytes"
	"encoding/binary"
	"reflect"
	"strconv"
)

func getFloatData(value reflect.Value, itype int, args ...binary.ByteOrder) interface{} {
	var endian binary.ByteOrder = binary.LittleEndian
	if len(args) > 0 {
		endian = args[0]
	}
	//value := reflect.ValueOf(data)
	if value.Type() == nil {
		switch itype {
		case Float64Type:
			return float64(0)
		case Float32Type:
			return float32(0)
		default:
			panic("can't convert to " + TypeStr[itype] + " of type " + value.Kind().String())
		}
	}
	switch f := value; value.Kind() {
	case reflect.Bool:
		rvalue := f.Bool()
		switch itype {
		case Float64Type:
			if rvalue {
				return float64(1)
			} else {
				return float64(0)
			}
		case Float32Type:
			if rvalue {
				return float32(1)
			} else {
				return float32(0)
			}
		default:
			panic("can't convert to " + TypeStr[itype] + " of type " + value.Kind().String())
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		rvalue := f.Int()
		switch itype {
		case Float64Type:
			return float64(rvalue)
		case Float32Type:
			return float32(rvalue)
		default:
			panic("can't convert to " + TypeStr[itype] + " of type " + value.Kind().String())
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		rvalue := f.Uint()
		switch itype {
		case Float64Type:
			return float64(rvalue)
		case Float32Type:
			return float32(rvalue)
		default:
			panic("can't convert to " + TypeStr[itype] + " of type " + value.Kind().String())
		}
	case reflect.Float32, reflect.Float64:
		rvalue := f.Float()
		switch itype {
		case Float64Type:
			return float64(rvalue)
		case Float32Type:
			return float32(rvalue)
		default:
			panic("can't convert to " + TypeStr[itype] + " of type " + value.Kind().String())
		}
	case reflect.String:
		rvalue := f.String()
		if rvalue == "" {
			switch itype {
			case Float64Type:
				return float64(0)
			case Float32Type:
				return float32(0)
			default:
				panic("can't convert to " + TypeStr[itype] + " of type " + value.Kind().String())
			}
		}
		switch itype {
		case Float64Type:
			df, err := strconv.ParseFloat(rvalue, 64)
			if err != nil {
				panic("can't convert to " + TypeStr[itype] + " for string " + f.String())
			}
			return float64(df)
		case Float32Type:
			df, err := strconv.ParseFloat(rvalue, 32)
			if err != nil {
				panic("can't convert to " + TypeStr[itype] + " for string " + f.String())
			}
			return float32(df)
		default:
			panic("can't convert to " + TypeStr[itype] + " of type " + value.Kind().String())
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
			bytelen := f.Len()

			datasize := 8
			switch itype {
			case Float64Type:
				datasize = 8
			case Float32Type:
				datasize = 4
			default:
				goto err //panic("can't convert to " + TypeStr[itype] + " of type " + value.Kind().String())
			}
			if bytelen < datasize {
				tmpbuf := make([]byte, datasize-bytelen)
				if endian == binary.BigEndian {
					buffer = append(tmpbuf, buffer...)
				} else {
					buffer = append(buffer, tmpbuf...)
				}
			}
			reader := bytes.NewReader(buffer)
			switch itype {
			case Float64Type:
				num := float64(0)
				err := binary.Read(reader, endian, &num)
				if err != nil {
					goto err //panic("can't convert to " + TypeStr[itype] + " of type " + value.Kind().String())
				}
				return num
			case Float32Type:
				num := float32(0)
				err := binary.Read(reader, endian, &num)
				if err != nil {
					goto err //panic("can't convert to " + TypeStr[itype] + " of type " + value.Kind().String())
				}
				return num
			default:
			}
		}
	case reflect.Interface:
		value := f.Elem()
		if !value.IsValid() {
			panic("can't convert to " + TypeStr[itype] + " of type " + value.Kind().String())
		} else {
			return getFloatData(value, itype)
		}
	case reflect.Ptr:
		// pointer to array or slice or struct?  ok at top level
		// but not embedded (avoid loops)
		if f.Pointer() != 0 {
			return getFloatData(f.Elem(), itype)
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
	panic("can't convert to " + TypeStr[itype] + " of type " + value.Kind().String())
}

func Float32(data interface{}) float32 {
	num, ok := getFloatData(reflect.ValueOf(data), Float32Type).(float32)
	if ok {
		return num
	} else {
		panic("can't convert to " + TypeStr[Float32Type] + " of type " + reflect.TypeOf(data).Kind().String())
	}
}

func Float64(data interface{}) float64 {
	num, ok := getFloatData(reflect.ValueOf(data), Float64Type).(float64)
	if ok {
		return num
	} else {
		panic("can't convert to " + TypeStr[Float64Type] + " of type " + reflect.TypeOf(data).Kind().String())
	}
}

func BFloat32(data interface{}) float32 {
	num, ok := getFloatData(reflect.ValueOf(data), Float32Type, binary.BigEndian).(float32)
	if ok {
		return num
	} else {
		panic("can't convert to " + TypeStr[Float32Type] + " of type " + reflect.TypeOf(data).Kind().String())
	}
}

func BFloat64(data interface{}) float64 {
	num, ok := getFloatData(reflect.ValueOf(data), Float64Type, binary.BigEndian).(float64)
	if ok {
		return num
	} else {
		panic("can't convert to " + TypeStr[Float64Type] + " of type " + reflect.TypeOf(data).Kind().String())
	}
}
