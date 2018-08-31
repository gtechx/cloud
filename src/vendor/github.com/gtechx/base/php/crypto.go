package php

import (
	//"bytes"
	"crypto/md5"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

// func Md5(s string) (string, error) {
// 	h := md5.New()
// 	if _, err := h.Write([]byte(s)); err != nil {
// 		return "", err
// 	}
// 	result := h.Sum(nil)
// 	return hex.EncodeToString(result), nil
// }

func Md5(str string) string {
	hash := md5.New()
	hash.Write([]byte(str))
	cipherText2 := hash.Sum(nil)
	hexText := make([]byte, 32)
	hex.Encode(hexText, cipherText2)
	return string(hexText)
}

// func Md5(s string) string {
// 	md5Ctx := md5.New()
// 	md5Ctx.Write([]byte(s))
// 	md5Hash := md5Ctx.Sum(nil)
// 	md5Str := string(hex.EncodeToString(md5Hash))
// 	return md5Str
// }

func Md5_file(filepath string) (string, error) {
	f, err := os.Open(filepath)
	if err != nil {
		return "", err
	}
	defer f.Close()
	contentbyte, err := ioutil.ReadAll(f)
	if err != nil {
		return "", err
	}
	h := md5.New()
	h.Write(contentbyte)
	return hex.EncodeToString(h.Sum(nil)), nil

}

// func Base64_decode(str string) (string, error) {
// 	bt, err := base64.StdEncoding.DecodeString(str)
// 	if err != nil {
// 		return "", err
// 	}
// 	return string(bt), nil
// }

// func Base64_encode(str string) string {
// 	return base64.StdEncoding.EncodeToString([]byte(str))
// }

func Base64_encode(str string) string {
	return base64.StdEncoding.EncodeToString([]byte(str))
}

func Base64_decode(str string) string {
	decoded, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		//fmt.Println("base64_decode error:", err)
		return ""
	}
	return string(decoded)
}

/* ============================================================================================ */
// func Base64Encode(data string) string {
//     var buf bytes.Buffer
//     encoder := base64.NewEncoder(base64.StdEncoding, &buf)
//     encoder.Write([]byte(data))
//     encoder.Close()
//     return buf.String()
// }
// /* ============================================================================================ */
// func Base64Decode(data string) string {
//    data = strings.Replace(data, "\r", "", -1)
//    data = strings.Replace(data, "\n", "", -1)
//    var buf = bytes.NewBufferString(data)
//    var res bytes.Buffer
//    decoder := base64.NewDecoder(base64.StdEncoding, buf)
//    res.ReadFrom(decoder)
//    return res.String()
// }
// /* ============================================================================================ */
// func Md5(data string) string {
//     var h = md5.New()
//     h.Write([]byte(data))
//     return fmt.Sprintf("%x", h.Sum(nil))
// }
/* ============================================================================================ */
func Sha1(data string) string {
	var h = sha1.New()
	h.Write([]byte(data))
	return fmt.Sprintf("%x", h.Sum(nil))
}

/* ============================================================================================ */
func QuotedPrintableDecode(data string) string {
	data = strings.Replace(data, "\r", "", -1)
	data = strings.Replace(data, "=\n", "", -1)
	var reg, _ = regexp.Compile(`[=][0-9A-F]{2}`)
	for _, line := range reg.FindAllString(data, -1) {
		dst, err := hex.DecodeString(strings.Trim(line, "="))
		if err == nil {
			data = strings.Replace(data, line, string(dst), 1)
		}
	}
	return data
}
