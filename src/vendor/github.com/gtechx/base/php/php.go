package php

import (
	"fmt"
	//"net/http"
	//"bufio"
	"net"
	"net/url"
	"reflect"
	//"regexp"
	"strconv"
	"strings"
	"time"
)

// func Sleep(seconds time.Duration) {
// 	time.Sleep(seconds * time.Second)
// }

var errList []string
var err error

/* ============================================================================================ */
func setErr(err error) {
	if err != nil {
		errList = append(errList, err.Error())
	}
}

/* ============================================================================================ */
func Error(separator string) string {
	if separator == "" {
		separator = "\n"
	}
	return strings.Join(errList, separator)
}

/* ============================================================================================ */
func ShowError() {
	if len(errList) > 0 {
		fmt.Println(Error("\n"))
	}
}

/* ============================================================================================ */
func ClearError() {
	errList = make([]string, 0)
}

func Sprintf(format string, a ...interface{}) string {
	return fmt.Sprintf(format, a...)
}

func Range(start int, end int, args ...int) []int {
	step := 1
	if len(args) > 0 {
		step = args[0]
	}
	count := (end - start + 1*step) / step
	arr := make([]int, count)

	for i := 0; i < count; i++ {
		arr[i] = i * step //start + i*step
	}

	return arr
}

func Uniqid(str string) string {
	return Md5(str)
}

func Gethostbyname(hostname string) string {
	conn, err := net.Dial("udp", hostname+":80")
	if err != nil {
		return ""
	}
	defer conn.Close()
	return strings.Split(conn.LocalAddr().String(), ":")[0]
}

func Parse_url(strurl string) (map[string]string, bool) {
	urldata, err := url.Parse(strurl)

	if err != nil {
		return nil, false
	}

	// scheme - e.g. http
	// host
	// port
	// user
	// pass
	// path
	// query - after the question mark ?
	// fragment - after the hashmark #
	ret := make(map[string]string)
	ret["scheme"] = urldata.Scheme

	hostarr := strings.Split(urldata.Host, ":")
	ret["host"] = hostarr[0]
	if len(hostarr) > 1 {
		ret["port"] = hostarr[1]
	} else {
		ret["port"] = "80"
	}

	ret["user"] = ""
	ret["pass"] = ""
	if urldata.User != nil {
		userarr := strings.Split(urldata.User.String(), ":")

		if len(userarr) > 0 {
			ret["user"] = userarr[0]
			if len(userarr) > 1 {
				ret["pass"] = userarr[1]
			} else {
				ret["pass"] = ""
			}
		}
	}

	ret["path"] = urldata.Path
	ret["query"] = urldata.RawQuery
	ret["fragment"] = urldata.Fragment

	return ret, true
}

func Isset(args ...interface{}) bool {
	arrVal := reflect.ValueOf(args[0])
	if !arrVal.IsValid() {
		return false
	}

	if len(args) < 2 || args[1] == nil {
		return issetValue(args[0])
	}

	switch arrVal.Kind() {
	case reflect.Array, reflect.Slice:
		if index, ok := args[1].(int); ok {
			return arrVal.Index(index).IsValid()
		}

		return false
	case reflect.Map:
		indexVal := arrVal.MapIndex(reflect.ValueOf(args[1]))
		if indexVal.IsValid() == false {
			return false
		}

		return true
	case reflect.Struct:
		key := reflect.ValueOf(args[1]).String()
		if !arrVal.FieldByName(key).IsValid() {
			return false
		} else {
			return true
		}
	}

	return false
}

func issetValue(value interface{}) bool {
	if value == nil {
		return false
	}

	val := reflect.ValueOf(value)
	switch val.Kind() {
	case reflect.Bool:
		return val.Bool()

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return val.Int() != 0

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return val.Uint() != 0

	case reflect.Float32, reflect.Float64:
		return val.Float() != 0.00

	case reflect.Complex64, reflect.Complex128:
		return val.Complex() != 0+0i

	case reflect.String:
		return val.String() != ""

	}

	return val.IsValid()
}

func Empty(args ...interface{}) bool {
	arrVal := reflect.ValueOf(args[0])
	if !arrVal.IsValid() {
		return true
	}

	if len(args) < 2 || args[1] == nil {
		return isEmptyValue(args[0])
	}

	switch arrVal.Kind() {
	case reflect.Array, reflect.Slice:
		if index, ok := args[1].(int); ok {
			return isEmptyValue(arrVal.Index(index).Interface())
		}
		return true

	case reflect.Map:
		indexVal := arrVal.MapIndex(reflect.ValueOf(args[1]))
		if indexVal.IsValid() == false {
			return true
		}
		return isEmptyValue(indexVal.Interface())

	case reflect.Struct:
		key := reflect.ValueOf(args[1]).String()
		if !arrVal.FieldByName(key).IsValid() {
			return true
		} else {
			return isEmptyValue(arrVal.FieldByName(key).Interface())
		}
	}

	return false
}

func isEmptyValue(value interface{}) bool {
	if value == nil {
		return true
	}

	val := reflect.ValueOf(value)
	switch val.Kind() {
	case reflect.Bool:
		return !val.Bool()

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return val.Int() == 0

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return val.Uint() == 0

	case reflect.Float32, reflect.Float64:
		return val.Float() == 0.00

	case reflect.Complex64, reflect.Complex128:
		return val.Complex() == 0+0i

	case reflect.String:
		realVal := val.String()
		return realVal == "" || realVal == "0"

	case reflect.Array, reflect.Slice, reflect.Map, reflect.Chan:
		return val.Len() == 0

	case reflect.Struct:
		return val.NumField() == 0
	}

	return false
}

// func In_array(str interface{}, list []interface{}) bool {
// 	for _, v := range list {
// 		if v == str {
// 			return true
// 		}
// 	}
// 	return false
// }

func is_array(value reflect.Value) int {
	result := 0

	switch value.Kind() {
	case reflect.Invalid:
	case reflect.Bool:
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
	case reflect.Float32:
	case reflect.Float64:
	case reflect.Complex64:
	case reflect.Complex128:
	case reflect.String:
	case reflect.Struct:
	case reflect.Interface:
		if value.Elem().IsValid() {
			return is_array(value.Elem())
		}
	case reflect.Ptr:
	case reflect.Chan, reflect.Func, reflect.UnsafePointer:
	case reflect.Map:
		result = 2
	case reflect.Array, reflect.Slice:
		result = 1
	default:
	}

	return result
}

//return 1 is array, 2 is map
func Is_array(value interface{}) int {
	result := 0

	switch f := value.(type) {
	case bool:
	case float32:
	case float64:
	case complex64:
	case complex128:
	case int:
	case int8:
	case int16:
	case int32:
	case int64:
	case uint:
	case uint8:
	case uint16:
	case uint32:
	case uint64:
	case uintptr:
	case string:
	case []byte:
		result = 1
	case reflect.Value:
		result = is_array(f)
	default:
		// If the type is not simple, it might have methods.
		// if !p.handleMethods(verb) {
		// 	// Need to use reflection, since the type had no
		// 	// interface methods that could be used for formatting.
		// 	p.printValue(reflect.ValueOf(f), verb, 0)
		// }
		result = is_array(reflect.ValueOf(f))
	}

	return result
}

func IsNumeric(x interface{}) bool {

	// Initialize
	result := false

	// Figure out result
	switch x.(type) {
	default:
		result = false

	case uint8:
		result = true
	case uint16:
		result = true
	case uint32:
		result = true
	case uint64:
		result = true

	case int8:
		result = true
	case int16:
		result = true
	case int32:
		result = true
	case int64:
		result = true

	case float32:
		result = true
	case float64:
		result = true

	case complex64:
		result = true
	case complex128:
		result = true

	//case byte:
	//	result = true
	//case rune:
	//	result = true

	case uint:
		result = true
	case int:
		result = true

	case string:
		if xAsString, ok := x.(string); ok {
			result = isStringNumeric(xAsString)
		} else {
			result = false
		}
	} // switch

	// Return.
	return result
}

func isStringNumeric(x string) bool {

	result := true

	isFirstLoop := true

	hasPeriod := false

	for i, c := range x {

		switch c {
		default:
			result = false
			return result

		case '-':
			if isFirstLoop {
				// Nothing here.
			} else {
				result = false
				return result
			}

		case '.':
			if hasPeriod {
				result = false
				return result
			}

			if isFirstLoop {
				result = false
				return result
			}

			if len(x) == 1+i {
				result = false
				return result
			}

			hasPeriod = true

		case '0': // Nothing here.
		case '1': // Nothing here.
		case '2': // Nothing here.
		case '3': // Nothing here.
		case '4': // Nothing here.
		case '5': // Nothing here.
		case '6': // Nothing here.
		case '7': // Nothing here.
		case '8': // Nothing here.
		case '9': // Nothing here.
		} // switch

		isFirstLoop = false

	} // for

	return result
}

func Fsockopen(hostname string, port int, errno *int, errstr *string, timeout int64) *net.Conn {
	conn, err := net.DialTimeout("tcp", hostname+":"+strconv.Itoa(port), time.Duration(timeout))
	if err != nil {
		// handle error
		if errstr != nil {
			*errstr = err.Error()
		}
		return nil
	}
	//fmt.Fprintf(conn, "GET / HTTP/1.0\r\n\r\n")
	//status, err := bufio.NewReader(conn).ReadString('\n')

	return &conn //bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))
}

func Pfsockopen(hostname string, port int, errno *int, errstr *string, timeout int64) *net.Conn {
	return Fsockopen(hostname, port, errno, errstr, timeout)
}

func stream_set_blocking(fp *net.Conn, block bool) {

}

func stream_set_timeout(fp *net.Conn, timeout int64) {
	s := timeout / 1000000
	n := timeout % 1000000
	(*fp).SetDeadline(time.Unix(s, n))
}

func Fwrite(fp *net.Conn, data string) int {
	num, err := (*fp).Write([]byte(data))
	//num, err := fp.WriteString(data)

	if err != nil {
		return -1
	}

	return num
}

func Fread(fp *net.Conn, count int) string {
	buffer := make([]byte, count)
	num, err := (*fp).Read(buffer)

	if err != nil {
		return ""
	}

	return string(buffer[:num])
}

//just for copy paste
type in_array_type struct{}

// wait for generics........
// func in_array_type(single interface{}, array interface{}, Type interface{}) (in bool, err error) {
// 	a := single.(reflect.TypeOf(Type))
// 	b, ok := array.([]string)
// 	if !ok {
// 		return false, errors.New("Second parameter does not in law")
// 	}
// 	for _, v := range b {
// 		if v == a {
// 			return true, nil
// 		}
// 	}
// 	return false, nil
// }

// package aes

// import (
// 	"beanstalkd-queue/lib/signature"
// 	"bytes"
// 	"crypto/aes"
// 	"crypto/cipher"
// 	"errors"
// )

// // AES-128。key长度：16, 24, 32 bytes 对应 AES-128, AES-192, AES-256

// func AppEncrypt(origData []byte, AppID string) ([]byte, error) {
// 	key := signature.GetAppKey(AppID)
// 	if len(key) == 0 {
// 		return nil, errors.New(("key is not exist"))
// 	}
// 	return Encrypt(origData, []byte(key))
// }

// func AppDecrypt(origData []byte, AppID string) ([]byte, error) {
// 	key := signature.GetAppKey(AppID)
// 	if len(key) == 0 {
// 		return nil, errors.New(("key is not exist"))
// 	}
// 	return Decrypt(origData, []byte(key))
// }

// func Encrypt(origData, key []byte) ([]byte, error) {
// 	block, err := aes.NewCipher(key)
// 	if err != nil {
// 		return nil, err
// 	}
// 	blockSize := block.BlockSize()
// 	origData = PKCS5Padding(origData, blockSize)
// 	// origData = ZeroPadding(origData, block.BlockSize())
// 	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
// 	crypted := make([]byte, len(origData))
// 	// 根据CryptBlocks方法的说明，如下方式初始化crypted也可以
// 	// crypted := origData
// 	blockMode.CryptBlocks(crypted, origData)
// 	return crypted, nil
// }

// func Decrypt(crypted, key []byte) ([]byte, error) {
// 	block, err := aes.NewCipher(key)
// 	if err != nil {
// 		return nil, err
// 	}
// 	blockSize := block.BlockSize()
// 	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
// 	origData := make([]byte, len(crypted))
// 	// origData := crypted
// 	blockMode.CryptBlocks(origData, crypted)
// 	origData = PKCS5UnPadding(origData)
// 	// origData = ZeroUnPadding(origData)
// 	return origData, nil
// }

// func ZeroPadding(ciphertext []byte, blockSize int) []byte {
// 	padding := blockSize - len(ciphertext)%blockSize
// 	padtext := bytes.Repeat([]byte{0}, padding)
// 	return append(ciphertext, padtext...)
// }

// func ZeroUnPadding(origData []byte) []byte {
// 	length := len(origData)
// 	unpadding := int(origData[length-1])
// 	return origData[:(length - unpadding)]
// }

// func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
// 	padding := blockSize - len(ciphertext)%blockSize
// 	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
// 	return append(ciphertext, padtext...)
// }

// func PKCS5UnPadding(origData []byte) []byte {
// 	length := len(origData)
// 	// 去掉最后一个字节 unpadding 次
// 	unpadding := int(origData[length-1])
// 	return origData[:(length - unpadding)]
// }
