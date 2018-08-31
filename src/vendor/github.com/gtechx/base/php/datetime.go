/**
此包內是关于时间和日期的相关函数
*/
package php

import (
	"strings"
	"time"
)

func Time() int64 {
	return time.Now().Unix()
}

// func Time() int64 {
// 	return time.Now().Unix()
// }

func Microtime() int64 {
	return time.Now().UnixNano() / 1000000
}

func Sleep(s int) {
	time.Sleep(time.Duration(s) * time.Second)
}

func Checkdate(month, day, year int) bool {
	//check month
	if month > 12 || month < 1 {
		return false
	}
	//check day
	if day < 1 {
		return false
	}
	if month == 1 || month == 3 || month == 5 || month == 7 || month == 8 || month == 10 || month == 12 { //31 days
		if day > 31 {
			return false
		}
	} else if month == 4 || month == 6 || month == 9 || month == 11 { //30 days
		if day > 30 {
			return false
		}
	} else { // february
		if checkIfLeapYear(year) {
			if day > 29 {
				return false
			}
		} else {
			if day > 28 {
				return false
			}
		}
	}
	//check year
	if year < 1 || year > 32767 {
		return false
	}
	return true
}

func checkIfLeapYear(year int) bool {

	if year%100 == 0 {
		if year%400 == 0 {
			return true
		}
		return false
	}
	if year%4 == 0 {
		return true
	}
	return false
}

// func Date(str string, timestamp int64) (string, error) {

// 	return "", nil
// }

/* ============================================================================================ */
func prepare(format string) string {
	var formatArr = strings.Split(format, "")
	for i, _ := range formatArr {
		switch formatArr[i] {
		case "d":
			formatArr[i] = "02"
		case "j":
			formatArr[i] = "_2"
		case "m":
			formatArr[i] = "01"
		case "n":
			formatArr[i] = "_1"
		case "t":
			formatArr[i] = "31"
		case "Y":
			formatArr[i] = "2006"
		case "y":
			formatArr[i] = "06"
		case "a":
			formatArr[i] = "pm"
		case "A":
			formatArr[i] = "PM"
		case "g":
			formatArr[i] = "3"
		case "G":
			formatArr[i] = "15" // Что-то нет, 24 часового формата без ведущего нуля (((
		case "H":
			formatArr[i] = "15"
		case "h":
			formatArr[i] = "03"
		case "i":
			formatArr[i] = "04"
		case "s":
			formatArr[i] = "05"
		case "O":
			formatArr[i] = "-0700"
		case "P":
			formatArr[i] = "-07:00"
		case "T":
			formatArr[i] = "MST"
		case "Z":
			formatArr[i] = "Z0700"
		case "D":
			formatArr[i] = "Mon"
		case "l":
			formatArr[i] = "Monday"
		case "F":
			formatArr[i] = "January"
		case "M":
			formatArr[i] = "Jan"
		}
	}
	return strings.Join(formatArr, "")
}

/* ============================================================================================ */
func Date(format string, timestamp int64) string {
	if timestamp == -1 {
		timestamp = time.Now().Unix()
	}
	t := time.Unix(timestamp, 0)
	return t.Format(prepare(format))
}

/* ============================================================================================ */
// func Time() int64 {
// 	return time.Now().Unix()
// }

/* ============================================================================================ */
// func MicroTime(asFloat bool) string {
// 	var strTime string
// 	if asFloat == true {
// 		strTime = strconv.FormatInt(time.Now().Unix(), 10) + "." + strconv.Itoa(time.Now().Nanosecond())
// 	} else {
// 		strTime = strconv.Itoa(time.Now().Nanosecond()) + " " + strconv.FormatInt(time.Now().Unix(), 10)
// 	}
// 	return strTime
// }

/* ============================================================================================ */
func Strtotime(strTime string, args ...string) int64 {
	format := "2006-01-02 15:04:05"

	if len(args) > 0 && args[0] != "" {
		format = prepare(args[0])
	}
	// if format == "" {
	// 	format = "2006-01-02 15:04:05"
	// } else {
	// 	format = prepare(format)
	// }
	t, err := time.Parse(format, Trim(strTime))
	setErr(err)
	if err != nil {
		return 0
	}
	return t.UTC().Unix()
}
