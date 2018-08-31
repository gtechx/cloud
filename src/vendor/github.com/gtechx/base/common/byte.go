package common

import (
	"bytes"
	"encoding/binary"
	"reflect"
	"unsafe"
)

func init() {
	if getEndian() {
		Endian = binary.BigEndian
		bigEndian = true
	} else {
		Endian = binary.LittleEndian
		bigEndian = false
	}
}

//保存机器大小端
var Endian binary.ByteOrder
var bigEndian bool

func IsBigEndian() bool {
	return bigEndian
}

func IsLittleEndian() bool {
	return !bigEndian
}

//以下代码判断机器大小端
const INT_SIZE int = int(unsafe.Sizeof(0))

//true = big endian, false = little endian
func getEndian() (ret bool) {
	var i int = 0x1
	bs := (*[INT_SIZE]byte)(unsafe.Pointer(&i))
	if bs[0] == 0 {
		return true
	} else {
		return false
	}
}

func Byte(data interface{}) byte {
	return byte(Uint8(data))
}

func Rune(data interface{}) rune {
	return rune(Int32(data))
}

func BRune(data interface{}) rune {
	return rune(BInt32(data))
}

func Bytes(args ...interface{}) []byte {
	buffer := []byte{}
	for _, arg := range args {
		buffer = append(buffer, getBytes(arg)...)
	}

	return buffer
}

func getBytes(data interface{}) []byte {
	if reflect.TypeOf(data) == nil {
		return []byte{}
	}
	switch data.(type) {
	case []byte:
		rdata := data.([]byte)
		newbuf := make([]byte, len(rdata))
		copy(newbuf, rdata)
		return newbuf
	case string:
		rdata := []byte(data.(string))
		newbuf := make([]byte, len(rdata))
		copy(newbuf, []byte(rdata))
		return newbuf
	}

	return getBytesData(reflect.ValueOf(data))
}

func BBytes(args ...interface{}) []byte {
	buffer := []byte{}
	for _, arg := range args {
		buffer = append(buffer, getBBytes(arg)...)
	}

	return buffer
}

func getBBytes(data interface{}) []byte {
	if reflect.TypeOf(data) == nil {
		return []byte{}
	}
	switch data.(type) {
	case []byte:
		rdata := data.([]byte)
		newbuf := make([]byte, len(rdata))
		copy(newbuf, rdata)
		return newbuf
	case string:
		rdata := []byte(data.(string))
		newbuf := make([]byte, len(rdata))
		copy(newbuf, []byte(rdata))
		return newbuf
	}

	return getBytesData(reflect.ValueOf(data), binary.BigEndian)
}

func getBytesData(value reflect.Value, args ...binary.ByteOrder) []byte {
	var endian binary.ByteOrder = binary.LittleEndian
	if len(args) > 0 {
		endian = args[0]
	}

	switch f := value; value.Kind() {
	case reflect.String:
		return []byte(f.String())
	case reflect.Bool:
		if f.Bool() {
			return []byte{1}
		} else {
			return []byte{0}
		}
	case reflect.Int:
		if INT_SIZE == 4 {
			num := int32(f.Int())
			bytesBuffer := new(bytes.Buffer)
			err := binary.Write(bytesBuffer, endian, num)
			if err != nil {
				return []byte{}
			}

			return bytesBuffer.Bytes()
		} else if INT_SIZE == 8 {
			num := int64(f.Int())
			bytesBuffer := new(bytes.Buffer)
			err := binary.Write(bytesBuffer, endian, num)
			if err != nil {
				return []byte{}
			}

			return bytesBuffer.Bytes()
		} else if INT_SIZE == 2 {
			num := int16(f.Int())
			bytesBuffer := new(bytes.Buffer)
			err := binary.Write(bytesBuffer, endian, num)
			if err != nil {
				return []byte{}
			}

			return bytesBuffer.Bytes()
		}
	case reflect.Uint, reflect.Uintptr:
		if INT_SIZE == 4 {
			num := uint32(f.Uint())
			bytesBuffer := new(bytes.Buffer)
			err := binary.Write(bytesBuffer, endian, num)
			if err != nil {
				return []byte{}
			}

			return bytesBuffer.Bytes()
		} else if INT_SIZE == 8 {
			num := uint64(f.Uint())
			bytesBuffer := new(bytes.Buffer)
			err := binary.Write(bytesBuffer, endian, num)
			if err != nil {
				return []byte{}
			}

			return bytesBuffer.Bytes()
		} else if INT_SIZE == 2 {
			num := uint16(f.Uint())
			bytesBuffer := new(bytes.Buffer)
			err := binary.Write(bytesBuffer, endian, num)
			if err != nil {
				return []byte{}
			}

			return bytesBuffer.Bytes()
		}
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128:
		bytesBuffer := new(bytes.Buffer)

		err := binary.Write(bytesBuffer, endian, f.Interface())
		if err != nil {
			return []byte{}
		}

		return bytesBuffer.Bytes()
	case reflect.Array, reflect.Slice:
		l := f.Len()
		buffer := []byte{}
		for i := 0; i < l; i++ {
			buffer = append(buffer, getBytesData(f.Index(i), endian)...)
		}
		return buffer
	case reflect.Struct:
		t := f.Type()
		l := f.NumField()
		buffer := []byte{}
		for i := 0; i < l; i++ {
			// see comment for corresponding code in decoder.value()
			if v := f.Field(i); t.Field(i).Name != "_" {

				switch v.Kind() {
				case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
					switch v.Type().Kind() {
					case reflect.Int8:
						num := reflect.ValueOf(interface{}(int8(v.Int())))
						buffer = append(buffer, getBytesData(num, endian)...)
					case reflect.Int16:
						num := reflect.ValueOf(interface{}(int16(v.Int())))
						buffer = append(buffer, getBytesData(num, endian)...)
					case reflect.Int32:
						num := reflect.ValueOf(interface{}(int32(v.Int())))
						buffer = append(buffer, getBytesData(num, endian)...)
					case reflect.Int64:
						num := reflect.ValueOf(interface{}((v.Int())))
						buffer = append(buffer, getBytesData(num, endian)...)
					}

				case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
					switch v.Type().Kind() {
					case reflect.Uint8:
						num := reflect.ValueOf(interface{}(uint8(v.Uint())))
						buffer = append(buffer, getBytesData(num, endian)...)
					case reflect.Uint16:
						num := reflect.ValueOf(interface{}(uint16(v.Uint())))
						buffer = append(buffer, getBytesData(num, endian)...)
					case reflect.Uint32:
						num := reflect.ValueOf(interface{}(uint32(v.Uint())))
						buffer = append(buffer, getBytesData(num, endian)...)
					case reflect.Uint64:
						num := reflect.ValueOf(interface{}(uint64(v.Uint())))
						buffer = append(buffer, getBytesData(num, endian)...)
					}

				case reflect.Float32, reflect.Float64:
					switch v.Type().Kind() {
					case reflect.Float32:
						num := reflect.ValueOf(interface{}(float32(v.Float())))
						buffer = append(buffer, getBytesData(num, endian)...)
					case reflect.Float64:
						num := reflect.ValueOf(interface{}(float64(v.Float())))
						buffer = append(buffer, getBytesData(num, endian)...)
					}

				case reflect.Complex64, reflect.Complex128:
					switch v.Type().Kind() {
					case reflect.Complex64:
						x := v.Complex()
						num := reflect.ValueOf(interface{}(float32(real(x))))
						buffer = append(buffer, getBytesData(num, endian)...)
						num = reflect.ValueOf(interface{}(float32(imag(x))))
						buffer = append(buffer, getBytesData(num, endian)...)
					case reflect.Complex128:
						x := v.Complex()
						num := reflect.ValueOf(interface{}(float64(real(x))))
						buffer = append(buffer, getBytesData(num, endian)...)
						num = reflect.ValueOf(interface{}(float64(imag(x))))
						buffer = append(buffer, getBytesData(num, endian)...)
					}
				default:
					buffer = append(buffer, getBytesData(v, endian)...)
				}
			}
		}
		return buffer
	case reflect.Map:
		keys := f.MapKeys()
		buffer := []byte{}
		for _, key := range keys {
			buffer = append(buffer, getBytesData(f.MapIndex(key), endian)...)
		}
		return buffer
	case reflect.Interface:
		value := f.Elem()
		if !value.IsValid() {
			goto err //panic("can't convert to bool of type " + value.Kind().String())
		} else {
			return getBytesData(value, endian)
		}
	case reflect.Ptr:
		if f.Pointer() != 0 {
			return getBytesData(f.Elem(), endian)
		}
		//case reflect.Chan, reflect.Func, reflect.UnsafePointer, reflect.Ptr, reflect.Map, reflect.Invalid:
	}
err:
	return []byte{} //panic("can't convert to bool of type " + value.Kind().String())
}
