package common

import (
	"bytes"
	"encoding/binary"
	"reflect"
	"strconv"
)

func getIntData(value reflect.Value, itype int, args ...binary.ByteOrder) interface{} {
	//value := reflect.ValueOf(data)
	var endian binary.ByteOrder = binary.LittleEndian
	if value.Kind() == reflect.Invalid {
		panic("can't convert to " + TypeStr[itype] + " of type " + value.Kind().String())
	}
	if len(args) > 0 {
		endian = args[0]
	}
	if value.Type() == nil {
		switch itype {
		case IntType:
			return int(0)
		case UintType:
			return uint(0)
		case Int64Type:
			return int64(0)
		case Uint64Type:
			return uint64(0)
		case Int32Type:
			return int32(0)
		case Uint32Type:
			return uint32(0)
		case Int16Type:
			return int16(0)
		case Uint16Type:
			return uint16(0)
		case Int8Type:
			return int8(0)
		case Uint8Type:
			return uint8(0)
		default:
			goto err //panic("can't convert to " + TypeStr[itype] + " of type " + value.Kind().String())
		}
	}
	switch f := value; value.Kind() {
	case reflect.Bool:
		rvalue := f.Bool()
		switch itype {
		case IntType:
			if rvalue {
				return int(1)
			} else {
				return int(0)
			}
		case UintType:
			if rvalue {
				return uint(1)
			} else {
				return uint(0)
			}
		case Int64Type:
			if rvalue {
				return int64(1)
			} else {
				return int64(0)
			}
		case Uint64Type:
			if rvalue {
				return uint64(1)
			} else {
				return uint64(0)
			}
		case Int32Type:
			if rvalue {
				return int32(1)
			} else {
				return int32(0)
			}
		case Uint32Type:
			if rvalue {
				return uint32(1)
			} else {
				return uint32(0)
			}
		case Int16Type:
			if rvalue {
				return int16(1)
			} else {
				return int16(0)
			}
		case Uint16Type:
			if rvalue {
				return uint16(1)
			} else {
				return uint16(0)
			}
		case Int8Type:
			if rvalue {
				return int8(1)
			} else {
				return int8(0)
			}
		case Uint8Type:
			if rvalue {
				return uint8(1)
			} else {
				return uint8(0)
			}
		default:
			//panic("can't convert to " + TypeStr[itype] + " of type " + value.Kind().String())
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		rvalue := f.Int()
		switch itype {
		case IntType:
			return int(rvalue)
		case UintType:
			return uint(rvalue)
		case Int64Type:
			return int64(rvalue)
		case Uint64Type:
			return uint64(rvalue)
		case Int32Type:
			return int32(rvalue)
		case Uint32Type:
			return uint32(rvalue)
		case Int16Type:
			return int16(rvalue)
		case Uint16Type:
			return uint16(rvalue)
		case Int8Type:
			return int8(rvalue)
		case Uint8Type:
			return uint8(rvalue)
		default:
			//panic("can't convert to " + TypeStr[itype] + " of type " + value.Kind().String())
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		rvalue := f.Uint()
		switch itype {
		case IntType:
			return int(rvalue)
		case UintType:
			return uint(rvalue)
		case Int64Type:
			return int64(rvalue)
		case Uint64Type:
			return uint64(rvalue)
		case Int32Type:
			return int32(rvalue)
		case Uint32Type:
			return uint32(rvalue)
		case Int16Type:
			return int16(rvalue)
		case Uint16Type:
			return uint16(rvalue)
		case Int8Type:
			return int8(rvalue)
		case Uint8Type:
			return uint8(rvalue)
		default:
			//panic("can't convert to " + TypeStr[itype] + " of type " + value.Kind().String())
		}
	case reflect.Float32, reflect.Float64:
		rvalue := f.Float()
		switch itype {
		case IntType:
			return int(rvalue)
		case UintType:
			return uint(rvalue)
		case Int64Type:
			return int64(rvalue)
		case Uint64Type:
			return uint64(rvalue)
		case Int32Type:
			return int32(rvalue)
		case Uint32Type:
			return uint32(rvalue)
		case Int16Type:
			return int16(rvalue)
		case Uint16Type:
			return uint16(rvalue)
		case Int8Type:
			return int8(rvalue)
		case Uint8Type:
			return uint8(rvalue)
		default:
			//panic("can't convert to " + TypeStr[itype] + " of type " + value.Kind().String())
		}
	case reflect.String:
		rvalue := f.String()
		if rvalue == "" {
			switch itype {
			case IntType:
				return int(0)
			case UintType:
				return uint(0)
			case Int64Type:
				return int64(0)
			case Uint64Type:
				return uint64(0)
			case Int32Type:
				return int32(0)
			case Uint32Type:
				return uint32(0)
			case Int16Type:
				return int16(0)
			case Uint16Type:
				return uint16(0)
			case Int8Type:
				return int8(0)
			case Uint8Type:
				return uint8(0)
			default:
				goto err //panic("can't convert to " + TypeStr[itype] + " of type " + value.Kind().String())
			}
		}
		switch itype {
		case IntType:
			d64, err := strconv.ParseInt(rvalue, 0, INT_SIZE*8)
			if err != nil {
				goto err //panic("can't convert to " + TypeStr[itype] + " for string " + f.String())
			}
			return int(d64)
		case UintType:
			d64, err := strconv.ParseUint(rvalue, 0, INT_SIZE*8)
			if err != nil {
				goto err //panic("can't convert to " + TypeStr[itype] + " for string " + f.String())
			}
			return uint(d64)
		case Int64Type:
			d64, err := strconv.ParseInt(rvalue, 0, 64)
			if err != nil {
				goto err //panic("can't convert to " + TypeStr[itype] + " for string " + f.String())
			}
			return int64(d64)
		case Uint64Type:
			d64, err := strconv.ParseUint(rvalue, 0, 64)
			if err != nil {
				goto err //panic("can't convert to " + TypeStr[itype] + " for string " + f.String())
			}
			return uint64(d64)
		case Int32Type:
			d64, err := strconv.ParseInt(rvalue, 0, 32)
			if err != nil {
				goto err //panic("can't convert to " + TypeStr[itype] + " for string " + f.String())
			}
			return int32(d64)
		case Uint32Type:
			d64, err := strconv.ParseUint(rvalue, 0, 32)
			if err != nil {
				goto err //panic("can't convert to " + TypeStr[itype] + " for string " + f.String())
			}
			return uint32(d64)
		case Int16Type:
			d64, err := strconv.ParseInt(rvalue, 0, 16)
			if err != nil {
				goto err //panic("can't convert to " + TypeStr[itype] + " for string " + f.String())
			}
			return int16(d64)
		case Uint16Type:
			d64, err := strconv.ParseUint(rvalue, 0, 16)
			if err != nil {
				goto err //panic("can't convert to " + TypeStr[itype] + " for string " + f.String())
			}
			return uint16(d64)
		case Int8Type:
			d64, err := strconv.ParseInt(rvalue, 0, 8)
			if err != nil {
				goto err //panic("can't convert to " + TypeStr[itype] + " for string " + f.String())
			}
			return int8(d64)
		case Uint8Type:
			d64, err := strconv.ParseUint(rvalue, 0, 8)
			if err != nil {
				goto err //panic("can't convert to " + TypeStr[itype] + " for string " + f.String())
			}
			return uint8(d64)
		default:
			//panic("can't convert to " + TypeStr[itype] + " of type " + value.Kind().String())
		}
	case reflect.Interface:
		value := f.Elem()
		if !value.IsValid() {
			goto err //panic("can't convert to " + TypeStr[itype] + " of type " + value.Kind().String())
		} else {
			return getIntData(value, itype)
		}
	case reflect.Array, reflect.Slice:
		//byte array to int type
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
			case IntType:
				datasize = INT_SIZE
			case UintType:
				datasize = INT_SIZE
			case Int64Type:
				datasize = 8
			case Uint64Type:
				datasize = 8
			case Int32Type:
				datasize = 4
			case Uint32Type:
				datasize = 4
			case Int16Type:
				datasize = 2
			case Uint16Type:
				datasize = 2
			case Int8Type:
				datasize = 1
			case Uint8Type:
				datasize = 1
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
			case IntType:
				if INT_SIZE == 4 {
					num := int32(0)
					err := binary.Read(reader, endian, &num)
					if err != nil {
						goto err //panic("can't convert to " + TypeStr[itype] + " of type " + value.Kind().String())
					}
					return int(num)
				} else if INT_SIZE == 8 {
					num := int64(0)
					err := binary.Read(reader, endian, &num)
					if err != nil {
						//println(err.Error())
						//println(buffer)
						goto err //panic("can't convert to " + TypeStr[itype] + " of type " + value.Kind().String())
					}
					return int(num)
				} else if INT_SIZE == 2 {
					num := int16(0)
					err := binary.Read(reader, endian, &num)
					if err != nil {
						//println(err.Error())
						//println(buffer)
						goto err //panic("can't convert to " + TypeStr[itype] + " of type " + value.Kind().String())
					}
					return int(num)
				}
			case UintType:
				if INT_SIZE == 4 {
					num := uint32(0)
					err := binary.Read(reader, endian, &num)
					if err != nil {
						goto err //panic("can't convert to " + TypeStr[itype] + " of type " + value.Kind().String())
					}
					return uint(num)
				} else if INT_SIZE == 8 {
					num := uint64(0)
					err := binary.Read(reader, endian, &num)
					if err != nil {
						//println(err.Error())
						//println(buffer)
						goto err //panic("can't convert to " + TypeStr[itype] + " of type " + value.Kind().String())
					}
					return uint(num)
				} else if INT_SIZE == 2 {
					num := uint16(0)
					err := binary.Read(reader, endian, &num)
					if err != nil {
						//println(err.Error())
						//println(buffer)
						goto err //panic("can't convert to " + TypeStr[itype] + " of type " + value.Kind().String())
					}
					return uint(num)
				}
				panic("can't convert to " + TypeStr[itype] + " of type " + value.Kind().String())
			case Int64Type:
				num := int64(0)
				err := binary.Read(reader, endian, &num)
				if err != nil {
					goto err //panic("can't convert to " + TypeStr[itype] + " of type " + value.Kind().String())
				}
				return num
			case Uint64Type:
				num := uint64(0)
				err := binary.Read(reader, endian, &num)
				if err != nil {
					goto err //panic("can't convert to " + TypeStr[itype] + " of type " + value.Kind().String())
				}
				return num
			case Int32Type:
				num := int32(0)
				err := binary.Read(reader, endian, &num)
				if err != nil {
					goto err //panic("can't convert to " + TypeStr[itype] + " of type " + value.Kind().String())
				}
				return num
			case Uint32Type:
				num := uint32(0)
				err := binary.Read(reader, endian, &num)
				if err != nil {
					goto err //panic("can't convert to " + TypeStr[itype] + " of type " + value.Kind().String())
				}
				return num
			case Int16Type:
				num := int16(0)
				err := binary.Read(reader, endian, &num)
				if err != nil {
					goto err //panic("can't convert to " + TypeStr[itype] + " of type " + value.Kind().String())
				}
				return num
			case Uint16Type:
				num := uint16(0)
				err := binary.Read(reader, endian, &num)
				if err != nil {
					goto err //panic("can't convert to " + TypeStr[itype] + " of type " + value.Kind().String())
				}
				return num
			case Int8Type:
				num := int8(0)
				err := binary.Read(reader, endian, &num)
				if err != nil {
					goto err //panic("can't convert to " + TypeStr[itype] + " of type " + value.Kind().String())
				}
				return num
			case Uint8Type:
				num := uint8(0)
				err := binary.Read(reader, endian, &num)
				if err != nil {
					goto err //panic("can't convert to " + TypeStr[itype] + " of type " + value.Kind().String())
				}
				return num
			default:
			}
		}
	case reflect.Ptr:
		// pointer to array or slice or struct?  ok at top level
		// but not embedded (avoid loops)
		if f.Pointer() != 0 {
			return getIntData(f.Elem(), itype)
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
	if value.IsValid() {
		panic("can't convert to " + TypeStr[itype] + " of type " + value.Kind().String() + " value:" + String(value.Interface()))
	} else {
		panic("can't convert to " + TypeStr[itype] + " of type " + value.Kind().String())
	}
}

func Int(data interface{}) int {
	num, ok := getIntData(reflect.ValueOf(data), IntType).(int)
	if ok {
		return num
	} else {
		panic("can't convert to " + TypeStr[IntType] + " of type " + reflect.TypeOf(data).Kind().String())
	}
}

func Uint(data interface{}) uint {
	num, ok := getIntData(reflect.ValueOf(data), UintType).(uint)
	if ok {
		return num
	} else {
		panic("can't convert to " + TypeStr[UintType] + " of type " + reflect.TypeOf(data).Kind().String())
	}
}

func Int64(data interface{}) int64 {
	num, ok := getIntData(reflect.ValueOf(data), Int64Type).(int64)
	if ok {
		return num
	} else {
		panic("can't convert to " + TypeStr[Int64Type] + " of type " + reflect.TypeOf(data).Kind().String())
	}
}

func Uint64(data interface{}) uint64 {
	num, ok := getIntData(reflect.ValueOf(data), Uint64Type).(uint64)
	if ok {
		return num
	} else {
		panic("can't convert to " + TypeStr[Uint64Type] + " of type " + reflect.TypeOf(data).Kind().String())
	}
}

func Int32(data interface{}) int32 {
	num, ok := getIntData(reflect.ValueOf(data), Int32Type).(int32)
	if ok {
		return num
	} else {
		panic("can't convert to " + TypeStr[Int32Type] + " of type " + reflect.TypeOf(data).Kind().String())
	}
}

func Uint32(data interface{}) uint32 {
	num, ok := getIntData(reflect.ValueOf(data), Uint32Type).(uint32)
	if ok {
		return num
	} else {
		panic("can't convert to " + TypeStr[Uint32Type] + " of type " + reflect.TypeOf(data).Kind().String())
	}
}

func Int16(data interface{}) int16 {
	num, ok := getIntData(reflect.ValueOf(data), Int16Type).(int16)
	if ok {
		return num
	} else {
		panic("can't convert to " + TypeStr[Int16Type] + " of type " + reflect.TypeOf(data).Kind().String())
	}
}

func Uint16(data interface{}) uint16 {
	num, ok := getIntData(reflect.ValueOf(data), Uint16Type).(uint16)
	if ok {
		return num
	} else {
		panic("can't convert to " + TypeStr[Uint16Type] + " of type " + reflect.TypeOf(data).Kind().String())
	}
}

func Int8(data interface{}) int8 {
	num, ok := getIntData(reflect.ValueOf(data), Int8Type).(int8)
	if ok {
		return num
	} else {
		panic("can't convert to " + TypeStr[Int8Type] + " of type " + reflect.TypeOf(data).Kind().String())
	}
}

func Uint8(data interface{}) uint8 {
	num, ok := getIntData(reflect.ValueOf(data), Uint8Type).(uint8)
	if ok {
		return num
	} else {
		panic("can't convert to " + TypeStr[Uint8Type] + " of type " + reflect.TypeOf(data).Kind().String())
	}
}

func Uintptr(data interface{}) uintptr {
	return uintptr(Uint(reflect.ValueOf(data)))
}

//big endian
func BInt(data interface{}) int {
	num, ok := getIntData(reflect.ValueOf(data), IntType, binary.BigEndian).(int)
	if ok {
		return num
	} else {
		panic("can't convert to " + TypeStr[IntType] + " of type " + reflect.TypeOf(data).Kind().String())
	}
}

func BUint(data interface{}) uint {
	num, ok := getIntData(reflect.ValueOf(data), UintType, binary.BigEndian).(uint)
	if ok {
		return num
	} else {
		panic("can't convert to " + TypeStr[UintType] + " of type " + reflect.TypeOf(data).Kind().String())
	}
}

func BInt64(data interface{}) int64 {
	num, ok := getIntData(reflect.ValueOf(data), Int64Type, binary.BigEndian).(int64)
	if ok {
		return num
	} else {
		panic("can't convert to " + TypeStr[Int64Type] + " of type " + reflect.TypeOf(data).Kind().String())
	}
}

func BUint64(data interface{}) uint64 {
	num, ok := getIntData(reflect.ValueOf(data), Uint64Type, binary.BigEndian).(uint64)
	if ok {
		return num
	} else {
		panic("can't convert to " + TypeStr[Uint64Type] + " of type " + reflect.TypeOf(data).Kind().String())
	}
}

func BInt32(data interface{}) int32 {
	num, ok := getIntData(reflect.ValueOf(data), Int32Type, binary.BigEndian).(int32)
	if ok {
		return num
	} else {
		panic("can't convert to " + TypeStr[Int32Type] + " of type " + reflect.TypeOf(data).Kind().String())
	}
}

func BUint32(data interface{}) uint32 {
	num, ok := getIntData(reflect.ValueOf(data), Uint32Type, binary.BigEndian).(uint32)
	if ok {
		return num
	} else {
		panic("can't convert to " + TypeStr[Uint32Type] + " of type " + reflect.TypeOf(data).Kind().String())
	}
}

func BInt16(data interface{}) int16 {
	num, ok := getIntData(reflect.ValueOf(data), Int16Type, binary.BigEndian).(int16)
	if ok {
		return num
	} else {
		panic("can't convert to " + TypeStr[Int16Type] + " of type " + reflect.TypeOf(data).Kind().String())
	}
}

func BUint16(data interface{}) uint16 {
	num, ok := getIntData(reflect.ValueOf(data), Uint16Type, binary.BigEndian).(uint16)
	if ok {
		return num
	} else {
		panic("can't convert to " + TypeStr[Uint16Type] + " of type " + reflect.TypeOf(data).Kind().String())
	}
}

func BInt8(data interface{}) int8 {
	num, ok := getIntData(reflect.ValueOf(data), Int8Type, binary.BigEndian).(int8)
	if ok {
		return num
	} else {
		panic("can't convert to " + TypeStr[Int8Type] + " of type " + reflect.TypeOf(data).Kind().String())
	}
}

func BUint8(data interface{}) uint8 {
	num, ok := getIntData(reflect.ValueOf(data), Uint8Type, binary.BigEndian).(uint8)
	if ok {
		return num
	} else {
		panic("can't convert to " + TypeStr[Uint8Type] + " of type " + reflect.TypeOf(data).Kind().String())
	}
}

func BUintptr(data interface{}) uintptr {
	return uintptr(BUint(reflect.ValueOf(data)))
}
