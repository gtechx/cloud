/*
字符串操作相关函数
*/
package php

import (
	//"bytes"
	//"encoding/hex"
	"errors"
	"fmt"
	//"math"
	"net/url"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

// func Explode(str, sep string) ([]string, error) {
// 	if sep == "" {
// 		return []string{str}, nil
// 	}
// 	return strings.Split(str, sep), nil
// }

func Implode(sep string, strarr []string) string {
	return strings.Join(strarr, sep)
}

func Explode(sep, s string) []string {
	return strings.Split(s, sep)
}

func Htmlspecialchars_decode(s string) string {
	s = strings.Replace(s, "&lt;", "<", -1)
	s = strings.Replace(s, "&gt;", ">", -1)
	s = strings.Replace(s, "&amp;", "&", -1)
	s = strings.Replace(s, `&quot;`, `"`, -1)
	s = strings.Replace(s, `&#039;`, `'`, -1)
	return s
}

func Htmlspecialchars(s string) string {
	s = strings.Replace(s, "<", "&lt;", -1)
	s = strings.Replace(s, ">", "&gt;", -1)
	s = strings.Replace(s, "&", "&amp;", -1)
	s = strings.Replace(s, `"`, "&quot;", -1)
	s = strings.Replace(s, `'`, "&#039;", -1)
	return s
}

// func Trim(s1, s2 string) string {
// 	return strings.Trim(s1, s2)
// }

func Trim(s string) string {
	return strings.TrimSpace(s)
}

func Ltrim(s1, s2 string) string {
	return strings.TrimLeft(s1, s2)
}

func Rtrim(s1, s2 string) string {
	return strings.TrimRight(s1, s2)
}

func Strlen(str string) int {
	return len(str)
}

// func Str_replace(oldstr, newstr, s string, n int) string {
// 	return strings.Replace(s, oldstr, newstr, n)
// }

func Strpos(str string, sep string) int {
	return strings.Index(str, sep)
}

func Strtolower(str string) string {
	return strings.ToLower(str)
}

func Strtoupper(str string) string {
	return strings.ToUpper(str)
}

func Strcasecmp(str1, str2 string) bool {
	newstr1 := strings.ToLower(str1)
	newstr2 := strings.ToLower(str2)
	return !(newstr1 == newstr2)
}

func Substr(str string, args ...int) string {
	strlen := len(str)
	if strlen == 0 {
		return ""
	}
	spos := 0
	sublen := strlen
	argc := len(args)
	// if len(args) == 1 {
	// 	spos = args[0]

	// 	if spos < 0 {
	// 		sublen = -spos
	// 	} else {
	// 		sublen = strlen - spos
	// 	}
	// } else if len(args) == 2 {
	// 	spos = args[0]
	// 	sublen = args[1]
	// }

	if argc >= 2 {
		sublen = args[1]
	}

	if argc >= 1 {
		spos = args[0]
	}

	if spos < 0 {
		spos = strlen + spos
	}

	if sublen < 0 {
		sublen = strlen + sublen - spos
		if sublen < 0 {
			sublen = 0
		}
	}

	if spos+sublen > strlen {
		sublen = strlen - spos
	}

	if sublen <= 0 {
		return ""
	}

	if spos > strlen {
		return ""
	}

	//fmt.Println(str, spos, sublen)
	r := make([]byte, sublen)
	//fmt.Println(str, spos, sublen)
	copy(r, str[spos:spos+sublen])

	return string(r)
}

// func Strrev(s string) string {
// 	if s == "" {
// 		return s
// 	}
// 	defer func() {
// 		if err := recover(); err != nil {
// 			fmt.Println(err)
// 		}
// 	}()
// 	length := len([]rune(s))
// 	var re []rune = make([]rune, length)
// 	i := 0
// 	for _, v := range s {
// 		re[length-i-1] = v
// 		i++
// 	}
// 	return string(re)
// }

func Chr(ascii int) string {
	// if ascii > 127 || ascii < 32 {
	// 	return "" //, errors.New("invalid ascii code.")
	// }
	// if ascii > 255 {
	// 	ascii = ascii % 256
	// }
	if ascii < 0 {
		return ""
	}
	ascii = ascii % 256
	// var buf bytes.Buffer
	// buf.Write([]byte{byte(ascii)})
	return string([]byte{byte(ascii)})
}

func Ord(str string) int {
	if len(str) <= 0 {
		return -1
	}
	return int(str[0])
}

// func Chr(ch int) string {
// 	return string(rune(ch))
// }

/**
 *使用此函数将字符串分割成小块非常有用。
 *例如将 base64_encode() 的输出转换成符合 RFC 2045 语义的字符串。
 *它会在每 chunklen 个字符后边插入 end。(按byte计数，不是rune！)
 */
// func Chunk_split(body string, chunklen int, end string) (str string, err error) {
// 	count := int(math.Ceil(float64(len(body) / chunklen)))
// 	var buf bytes.Buffer
// 	for i := 0; i < count-1; i++ { //防止访问越界
// 		start := i * chunklen
// 		_, err = buf.Write([]byte(body)[start : start+chunklen])
// 		if err != nil {
// 			break
// 		}
// 		_, err = buf.WriteString(end)
// 		if err != nil {
// 			break
// 		}
// 	}
// 	_, err = buf.Write([]byte(body)[chunklen*(count-1):])
// 	_, err = buf.WriteString(end)
// 	if err != nil {
// 		return "", err
// 	}
// 	return buf.String(), nil
// }

/* ============================================================================================ */
// func Trim(str string) string {
// 	return LTrim(RTrim(str))
// }

// /* ============================================================================================ */
// func LTrim(str string) string {
// 	return strings.TrimLeft(str, " \n\r\t")
// }

// /* ============================================================================================ */
// func RTrim(str string) string {
// 	return strings.TrimRight(str, " \n\r\t")
// }

// /* ============================================================================================ */
// func HtmlSpecialChars(str string) string {
// 	str = html.EscapeString(str)
// 	str = strings.Replace(str, `&#34;`, `&quot;`, -1)
// 	return str
// }

// /* ============================================================================================ */
// func HtmlSpecialCharsDecode(str string) string {
// 	return html.UnescapeString(str)
// }

/*
rawurlencode遵守是94年国际标准备忘录RFC 1738，
urlencode实现的是传统做法，和上者的主要区别是对空格的转义是'+'而不是'%20'
javascript的encodeURL也是94年标准，
而javascript的escape是另一种用"%xxx"标记unicode编码的方法。
推荐在PHP中使用用rawurlencode。弃用urlencode
*/
/* ============================================================================================ */
func Urlencode(str string) string {
	//t := &url.URL{Path: str}
	return url.QueryEscape(str) //strings.Replace(t.String(), "+", "%2B", -1)
}

/* ============================================================================================ */
func Urldecode(str string) string {
	res, err := url.QueryUnescape(str)
	//setErr(err)
	if err == nil {
		return res
	}
	return str
}

func Rawurlencode(str string) string {
	//t := &url.URL{Path: str}
	//return strings.Replace(t.String(), "+", "%20", -1)
	str = url.QueryEscape(str)
	return strings.Replace(str, "+", "%20", -1)
}

/* ============================================================================================ */
func Rawurldecode(str string) string {
	newstr := strings.Replace(str, "%20", "+", -1)
	res, err := url.QueryUnescape(newstr)
	//setErr(err)
	if err == nil {
		return res
	}
	return str

	// re := regexp.MustCompile(`(?Ui)%[0-9A-F]{2}`)
	// str = re.ReplaceAllStringFunc(str, func(s string) string {
	// 	b, err := hex.DecodeString(s[1:])
	// 	if err == nil {
	// 		return string(b)
	// 	}
	// 	return s
	// })
	// return str
}

// /* ============================================================================================ */
// func Substr(str string, start int, length int) string {
// 	if start < 0 {
// 		start = 0
// 	}
// 	if length == 0 {
// 		length = len(str)
// 	}
// 	strArr := strings.Split(str, "")
// 	i := 0
// 	str = ""
// 	for i = start; i < start+length; i++ {
// 		if i >= len(strArr) {
// 			break
// 		}
// 		str += strArr[i]
// 	}
// 	return str
// }

// /* ============================================================================================ */
// func StrPos(haystack string, needle string, offset int) int {
// 	if offset > 0 {
// 		haystack = SubStr(haystack, offset, len(haystack))
// 	}
// 	return strings.Index(haystack, needle) + offset
// }

/* ============================================================================================ */
func Strrpos(haystack string, needle string) int {
	return strings.LastIndex(haystack, needle)
}

/* ============================================================================================ */
func Strstr(haystack string, needle string) string {
	start := Strpos(haystack, needle)
	if start == -1 {
		return ""
	}
	return Substr(haystack, start, len(haystack))
}

/* ============================================================================================ */
func Strchr(haystack string, needle string) string {
	return Strstr(haystack, needle)
}

/* ============================================================================================ */
func Stristr(haystack string, needle string) string {
	start := Strpos(strings.ToLower(haystack), strings.ToLower(needle))
	return Substr(haystack, start, len(haystack))
}

/* ============================================================================================ */
func Strrchr(haystack string, needle string) string {
	start := Strrpos(haystack, needle)
	if start == -1 {
		return ""
	}
	return Substr(haystack, start, len(haystack))
}

/* ============================================================================================ */
func Substr_count(haystack string, needle string) int {
	return strings.Count(haystack, needle)
}

/* ============================================================================================ */
func Strspn(str1 string, str2 string) int {
	strArr1 := strings.Split(str1, "")
	strArr2 := strings.Split(str2, "")
	var i int
	var j int
	var ok bool
	for i = 0; i < len(strArr1); i++ {
		ok = false
		for j = 0; j < len(strArr2); j++ {
			if strArr1[i] == strArr2[j] {
				ok = true
				break
			}
		}
		if ok == false {
			break
		}
	}
	return i
}

/* ============================================================================================ */
func Strcspn(str1 string, str2 string) int {
	strArr1 := strings.Split(str1, "")
	strArr2 := strings.Split(str2, "")
	var i int
	var j int
	var ok bool
	for i = 0; i < len(strArr1); i++ {
		ok = false
		for j = 0; j < len(strArr2); j++ {
			if strArr1[i] == strArr2[j] {
				ok = true
				break
			}
		}
		if ok == true {
			break
		}
	}
	return i
}

// /* ============================================================================================ */
// func StrLen(str string) int {
// 	return len(str)
// }

/* ============================================================================================ */
func Wordwrap(str string, width int, breakStr string, cut bool) string {
	var brRune = []rune(breakStr)
	var strArr = []rune(str)
	var res []rune
	var i, j int

	for i = 0; i < len(strArr); i++ {
		j++
		if j >= width && (unicode.IsSpace(strArr[i]) == true || cut == true) {
			j = 0
			res = append(res, brRune...)
		}
		res = append(res, strArr[i])
	}
	return string(res)
}

/* ============================================================================================ */
func Htmlwrap(str string, width int, breakStr string, cut bool) string {
	var brRune = []rune(breakStr)
	var last = str
	var lines []string
	var counter int

	re := regexp.MustCompile(`(?Ui)<[^>]*>`)
	for _, htmlTag := range re.FindAllString(str, -1) {
		arr := strings.SplitN(last, htmlTag, 2)
		if len(arr[0]) > 0 {
			lines = append(lines, arr[0])
		}
		lines = append(lines, htmlTag)
		last = arr[1]
	}
	if len(last) > 0 {
		lines = append(lines, last)
	}
	for i, line := range lines {
		if strings.HasPrefix(line, "<") == false && strings.HasSuffix(line, ">") == false {
			var strArr = []rune(line)
			var res []rune
			for j := 0; j < len(strArr); j++ {
				counter++
				if counter >= width && (unicode.IsSpace(strArr[j]) == true || cut == true) {
					counter = 0
					res = append(res, brRune...)
				}
				res = append(res, strArr[j])
			}
			lines[i] = string(res)
		}
	}
	return strings.Join(lines, "")
}

/* ============================================================================================ */
// func Str_replace(from string, to string, str string) string {
// 	return strings.Replace(str, from, to, -1)
// }

func Str_replace(from, to interface{}, str string, count ...int) string {
	repcount := -1
	if len(count) >= 1 {
		repcount = count[0]
	}
	// if repcount <= 0 {
	// 	return str
	// }

	isfromok := true
	switch from.(type) {
	case string:
	case []string:
	default:
		isfromok = false
	}

	istook := true
	switch to.(type) {
	case string:
	case []string:
	default:
		istook = false
	}

	if !isfromok || !istook {
		//fmt.Println("isfromok or istook is:", isfromok, istook)
		return str
	}

	switch from.(type) {
	case string:
		switch to.(type) {
		case string:
			//fmt.Println("1111isfromok or istook is:", isfromok, istook)
			return strings.Replace(str, from.(string), to.(string), repcount)
		case []string:
			todata := to.([]string)
			if len(todata) == 0 {
				return str
			}
			return strings.Replace(str, from.(string), to.([]string)[0], repcount)
		default:
			return str
		}
	case []string:
		fromdata := from.([]string)
		if len(fromdata) == 0 {
			return str
		}
		switch to.(type) {
		case string:
			str = strings.Replace(str, fromdata[0], to.(string), repcount)
			for i := 1; i < len(fromdata); i++ {
				str = strings.Replace(str, fromdata[i], "", repcount)
			}
			return str
		case []string:
			todata := to.([]string)
			if len(todata) == 0 {
				return str
			}
			for i := 0; i < len(fromdata); i++ {
				if i >= len(todata) {
					str = strings.Replace(str, fromdata[i], "", repcount)
				} else {
					str = strings.Replace(str, fromdata[i], todata[i], repcount)
				}
			}
			return str
		default:
		}
	default:
	}

	return str
}

/* ============================================================================================ */
func Substr_replace(str string, replacement string, start int, length int) string {
	if length < 0 {
		length = len(str)
	}
	strArr := strings.Split(str, "")
	str = ""
	var i int
	for i = 0; i < len(strArr); i++ {
		if i == start {
			str += replacement
		} else if i > start && i < start+length {
			continue
		} else {
			str += strArr[i]
		}
	}
	return str
}

/* ============================================================================================*/
func Strtr(str string, from string, to string) string {
	if len(from) > len(to) {
		from = Substr(from, 0, len(to))
	}
	strArr := strings.Split(str, "")
	fromArr := strings.Split(from, "")
	toArr := strings.Split(to, "")
	var i int
	var j int
	for i = 0; i < len(strArr); i++ {
		for j = 0; j < len(fromArr); j++ {
			if fromArr[j] == strArr[i] {
				strArr[i] = toArr[j]
				break
			}
		}
	}
	return strings.Join(strArr, "")
}

/* ============================================================================================ */
func Stripslashes(str string) string {
	twoSlashes := "###TWO_SLASHES###"
	for Strpos(str, twoSlashes) > -1 {
		twoSlashes += "#"
	}
	str = strings.Replace(str, `\\`, twoSlashes, -1)
	str = strings.Replace(str, `\`, "", -1)
	str = strings.Replace(str, twoSlashes, `\\`, -1)
	return str
}

/* ============================================================================================ */
func Addslashes(str string) string {
	str = strings.Replace(str, `\`, `\\`, -1)
	str = strings.Replace(str, `"`, `\"`, -1)
	str = strings.Replace(str, `'`, `\'`, -1)
	str = strings.Replace(str, "`", "\\`", -1)
	return str
}

/* ============================================================================================ */
func Quotemeta(str string) string {
	return regexp.QuoteMeta(str)
}

/* ============================================================================================ */
func Strrev(str string) string {
	strArr := strings.Split(str, "")
	str = ""
	var i int
	for i = len(strArr) - 1; i >= 0; i-- {
		str += strArr[i]
	}
	return str
}

/* ============================================================================================ */
func Str_repeat(str string, num int) string {
	return strings.Repeat(str, num)
}

/* ============================================================================================ */
func Str_pad(str string, padLength int, padStr string, padType int) string {
	if len(str) >= padLength {
		return str
	}
	if padType == 1 { // right
		for padLength > len(str) {
			if len(str)+len(padStr) > padLength {
				padStr = Substr(padStr, 0, padLength-len(str))
			}
			str += padStr
		}
	} else if padType == 2 { //left
		for padLength > len(str) {
			if len(str)+len(padStr) > padLength {
				padStr = Substr(padStr, 0, padLength-len(str))
			}
			str = padStr + str
		}
	} else if padType == 3 { // both
		for padLength > len(str) {
			if len(str)+len(padStr) > padLength {
				padStr = Substr(padStr, 0, padLength-len(str))
			}
			str += padStr
			if len(str)+len(padStr) > padLength {
				padStr = Substr(padStr, 0, padLength-len(str))
			}
			str = padStr + str
		}
	}
	return str
}

/* ============================================================================================ */
func Chunk_split(str string, chunkLen int, end string) string {
	var rStr = []rune(str + " ")
	var rEnd = []rune(end)
	var result []rune

	if chunkLen <= 0 {
		chunkLen = 76
	}
	if len(rEnd) == 0 {
		rEnd = []rune("\r\n")
	}

	for i := 0; i < len(rStr); i += chunkLen {
		from := i
		to := i + chunkLen
		if to >= len(rStr) {
			to = len(rStr) - 1
		}
		result = append(result, rStr[from:to]...)
		result = append(result, rEnd...)
	}

	return strings.TrimSpace(string(result))
}

// /* ============================================================================================ */
// func Explode(arg string, str string, maxlimit int) []string {
// 	if maxlimit <= 0 {
// 		return strings.Split(str, arg)
// 	}
// 	return strings.SplitN(str, arg, maxlimit)
// }

// /* ============================================================================================ */
// func Implode(arg string, array []string) string {
// 	return strings.Join(array, arg)
// }

/* ============================================================================================ */
func Join(arg string, array []string) string {
	return strings.Join(array, arg)
}

/* ============================================================================================ */
func Intjoin(a []int, sep string) string {
	var buf []string
	for _, val := range a {
		buf = append(buf, strconv.Itoa(val))
	}
	return strings.Join(buf, sep)
}

/* ============================================================================================ */
func Intsplit(s, sep string) []int {
	var result []int
	for _, val := range strings.Split(strings.Trim(s, "{}"), sep) {
		iVal, err := strconv.Atoi(strings.TrimSpace(val))
		if err == nil {
			result = append(result, iVal)
		}
	}
	return result
}

/* ============================================================================================ */
func Int64join(a []int64, sep string) string {
	var buf []string
	for _, val := range a {
		buf = append(buf, strconv.FormatInt(val, 10))
	}
	return strings.Join(buf, sep)
}

/* ============================================================================================ */
func Int64split(s, sep string) []int64 {
	var result []int64
	for _, val := range strings.Split(strings.Trim(s, "{}"), sep) {
		iVal, err := strconv.ParseInt(val, 10, 64)
		if err == nil {
			result = append(result, iVal)
		}
	}
	return result
}

// /* ============================================================================================ */
// func ParceURL(strUrl string) *url.URL {
// 	result, err := url.Parse(strUrl)
// 	setErr(err)
// 	return result
// }

// /* ============================================================================================ */
// func StrToLower(str string) string {
// 	return strings.ToLower(str)
// }

// /* ============================================================================================ */
// func StrToUpper(str string) string {
// 	return strings.ToUpper(str)
// }

/* ============================================================================================ */
func Ucfirst(str string) string {
	first := Substr(str, 0, 1)
	last := Substr(str, 1, 0)
	return strings.ToUpper(first) + strings.ToLower(last)
}

/* ============================================================================================ */
func Ucwords(str string) string {
	strArr := strings.Split(str, " ")
	for i := 0; i < len(strArr); i++ {
		strArr[i] = Ucfirst(strArr[i])
	}
	return strings.Join(strArr, " ")
}

/* ============================================================================================ */
func Strip_tags(s string) string {
	reg := regexp.MustCompile(`(?Uis)<script[^>]*>.*</script>`)
	s = string(reg.ReplaceAll([]byte(s), []byte("")))
	reg = regexp.MustCompile(`(?Uis)<.*>`)
	return string(reg.ReplaceAll([]byte(s), []byte("")))
}

/* ============================================================================================ */
func Parse_str(s string) map[string][]string {
	var arguments = make(map[string][]string)
	for _, argument := range strings.Split(s, "&") {
		arr := strings.SplitN(argument, "=", 2)
		if len(arr) > 1 {
			key := Urldecode(arr[0])
			value := Urldecode(arr[1])
			key = strings.Replace(key, " ", "_", -1)
			key = strings.Replace(key, ".", "_", -1)
			arguments[key] = append(arguments[key], value)
		}
	}
	return arguments
}

func Parse_str_s(s string) map[string]string {
	var arguments = make(map[string]string)
	for _, argument := range strings.Split(s, "&") {
		arr := strings.SplitN(argument, "=", 2)
		if len(arr) > 1 {
			key := Urldecode(arr[0])
			value := Urldecode(arr[1])
			key = strings.Replace(key, " ", "_", -1)
			key = strings.Replace(key, ".", "_", -1)
			arguments[key] = value
		}
	}
	return arguments
}

func Preg_quote(str, delimiter string) string {
	// . \ + * ? [ ^ ] $ ( ) { } = ! < > | : -
	//.\+*?[^]$(){}=!<>|:-
	//rstr := `.\+*?[^]$(){}=!<>|:-` + delimiter
	rstr := `\` + delimiter

	for i := 0; i < len(rstr); i += 1 {
		str = strings.Replace(str, string(rstr[i]), string(`\`+string(rstr[i])), -1)
	}
	return str
}

func Preg_replace(pattern, replacement interface{}, str string) string {
	//reg := MustCompile(pattern)
	//s := string(reg.ReplaceAll([]byte(subject), []byte(replacement)))

	//return s

	isfromok := true
	switch pattern.(type) {
	case string:
	case []string:
	default:
		isfromok = false
	}

	istook := true
	switch replacement.(type) {
	case string:
	case []string:
	default:
		istook = false
	}

	if !isfromok || !istook {
		return str
	}

	switch pattern.(type) {
	case string:
		switch replacement.(type) {
		case string:
			reg := MustCompile(pattern.(string))
			str = string(reg.ReplaceAll([]byte(str), []byte(replacement.(string))))
			return str
		case []string:
			todata := replacement.([]string)
			if len(todata) == 0 {
				return str
			}
			reg := MustCompile(pattern.(string))
			str = string(reg.ReplaceAll([]byte(str), []byte(replacement.([]string)[0])))
			return str
		default:
			return str
		}
	case []string:
		fromdata := pattern.([]string)
		if len(fromdata) == 0 {
			return str
		}
		switch replacement.(type) {
		case string:
			reg := MustCompile(fromdata[0])
			str = string(reg.ReplaceAll([]byte(str), []byte(replacement.(string))))
			for i := 1; i < len(fromdata); i++ {
				reg = MustCompile(fromdata[i])
				str = string(reg.ReplaceAll([]byte(str), []byte("")))
			}
			return str
		case []string:
			todata := replacement.([]string)
			if len(todata) == 0 {
				return str
			}
			for i := 0; i < len(fromdata); i++ {
				if i >= len(todata) {
					reg := MustCompile(fromdata[i])
					str = string(reg.ReplaceAll([]byte(str), []byte("")))
					//str = strings.Replace(str, fromdata[i], "", repcount)
				} else {
					reg := MustCompile(fromdata[i])
					str = string(reg.ReplaceAll([]byte(str), []byte(todata[i])))
					//str = strings.Replace(str, fromdata[i], todata[i], repcount)
				}
			}
			return str
		default:
		}
	default:
	}

	return str
}

func Preg_match(pattern, subject string) bool {
	reg, err := Compile(pattern)
	//ret, err := regexp.MatchString(pattern, subject)

	if err != nil {
		return false
	}

	return reg.MatchString(subject)
}

func Preg_match_r(pattern, subject string) ([]string, bool) {
	//ret, err := regexp.MatchString(pattern, subject)

	reg, err := Compile(pattern)

	if err == nil {
		match := reg.FindStringSubmatch(subject)

		return match, true
	}

	return nil, false
}

var (
	// ErrMalformedRegexp is returned when provided regexp string too short
	ErrMalformedRegexp = errors.New("Malformed Regexp")
	// ErrCantFindClosingSeparator is returned when provided regexp string
	// has wrong open/close separator configuration
	ErrCantFindClosingSeparator = errors.New("Cant find closing separator")
)

// MustCompile is like Compile but panics if the expression cannot be parsed.
// It simplifies safe initialization of global variables holding compiled regular
// expressions.
func MustCompile(str string) *regexp.Regexp {
	regexp, err := Compile(str)
	if err != nil {
		panic(`regexp: Compile(` + quote(str) + `): ` + err.Error())
	}

	return regexp
}

// Compile parses a regular expression in PHP format and returns,
// if successful, a Regexp object that can be used to match against text.
func Compile(str string) (*regexp.Regexp, error) {
	if len(str) < 3 {
		return nil, ErrMalformedRegexp
	}

	// Searching for opening and closing separator
	sep := str[0:1]
	last := strings.LastIndex(str, sep)
	if last == 0 || last <= 0 {
		return nil, ErrCantFindClosingSeparator
	}

	// Reading modifiers
	prefix := ""
	mods := str[last+1:]
	if len(mods) > 0 {
		for _, m := range mods {
			switch m {
			case 'i': // Case insensitive
				prefix += "i"
			case 'm': // Multiline
				prefix += "m"
			case 's': // Single line
				prefix += "s"
			case 'U': // Ungreedy
				prefix += "U"
			}
		}

		if prefix != "" {
			prefix = "(?" + prefix + ")"
		}
	}

	return regexp.Compile(prefix + str[1:last])
}

func quote(s string) string {
	if strconv.CanBackquote(s) {
		return "`" + s + "`"
	}
	return strconv.Quote(s)
}

func Bin2hex(str string) string {
	return fmt.Sprintf("%x", []byte(str))
}

func Hex2bin(str string) (ret string) {
	fmt.Sscanf(str, "%x", &ret)
	return
}

func Is_numeric(data interface{}) bool {
	value := reflect.ValueOf(data)
	if value.Kind() >= reflect.Int && value.Kind() <= reflect.Uint64 {
		return true
	}
	// if value.Kind() >= reflect.Int && value.Kind() <= reflect.Uint64{
	// 	return true
	// }
	switch f := value; value.Kind() {
	// case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
	// case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
	case reflect.Float32, reflect.Float64:
		return true
	case reflect.String:
		str := f.String()
		if strings.Index(str, ".") != -1 {
			_, err := strconv.ParseFloat(str, 64)
			if err == nil {
				return true
			}
		}
		if strings.Index(str, "0b") == 0 || strings.Index(str, "0B") == 0 {
			str = strings.Replace(str, "0b", "", 1)
			str = strings.Replace(str, "0B", "", 1)

			_, err = strconv.ParseInt(str, 2, 64)
			if err == nil {
				return true
			}
			_, err = strconv.ParseUint(str, 2, 64)
			if err == nil {
				return true
			}
		}
		if strings.Index(str, "0x") == 0 || strings.Index(str, "0X") == 0 {
			_, err = strconv.ParseInt(str, 16, 64)
			if err == nil {
				return true
			}
			_, err = strconv.ParseUint(str, 16, 64)
			if err == nil {
				return true
			}
		}
		if string(str[0]) == "0" {
			_, err = strconv.ParseInt(str, 8, 64)
			if err == nil {
				return true
			}
			_, err = strconv.ParseUint(str, 8, 64)
			if err == nil {
				return true
			}
		}
		_, err = strconv.ParseInt(str, 0, 64)
		if err == nil {
			return true
		}
		_, err = strconv.ParseUint(str, 0, 64)
		if err == nil {
			return true
		}
		_, err := strconv.ParseFloat(str, 64)
		if err == nil {
			return true
		}

		// case reflect.Complex64:
		// case reflect.Complex128:
		// case reflect.Invalid:
		// case reflect.Bool:
		// case reflect.Map:
		// case reflect.Struct:
		// case reflect.Interface:
		// case reflect.Array, reflect.Slice:
		// case reflect.Ptr, reflect.Uintptr:
		// case reflect.Chan, reflect.Func, reflect.UnsafePointer:
		//default:
	}
	return false
}

func Nl2br(str string, args ...bool) string {
	xhtml := true
	if len(args) > 0 {
		xhtml = args[0]
	}
	if xhtml {
		str = strings.Replace(str, "\n", "<br />", -1)
	} else {
		str = strings.Replace(str, "\n", "<br>", -1)
	}
	return str
}
