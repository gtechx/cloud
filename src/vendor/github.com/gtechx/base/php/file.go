package php

import (
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func Basename(file_path string) string {
	return filepath.Base(file_path)
}

func File_get_contents(url string) string {

	if strings.Index(url, "http") == 0 {
		return string(file_get_contents_url(url))
	}

	if strings.Index(url, "https") == 0 {
		return string(file_get_contents_url(url))
	}

	return string(file_get_contents_file(url))
	// (strings.Split("a,b,c,d,e,f,g", ",")) // [a b c d e f g]
}

func file_get_contents_url(url string) []byte {

	response, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}

	return body
}

func file_get_contents_file(url string) []byte {

	file, err := os.Open(url)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	body, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}

	return body
}

func File_put_contents(fileName string, write_data string) int {
	file, _ := os.Create(fileName)
	defer file.Close()

	wrote_byte, _ := file.Write([]byte(write_data))
	file.Sync()

	return wrote_byte
}

func File_exists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

// func Basename(path, suffix string) (string, error) {
// 	fi, err := os.Stat(path)
// 	if err != nil {
// 		return "", err
// 	}
// 	return fi.Name(), nil
// }

func Chdir(dir string) error {
	return os.Chdir(dir)
}

func Chgrp(filename string, group interface{}) error {
	//todo
	//uid := os.Geteuid()
	//switch group.(type) {
	//case string:
	//	grp := os.
	//case int:
	//	return os.Chown(filename,uid,group.(int))
	//default:
	//	return errors.New("unsupported group type")
	//
	//}
	//os.Chown()
	return nil
}

func Chmod(filename string, mode int) error {
	return os.Chmod(filename, os.FileMode(mode))
}

func Chroot(dir string) error {
	//todo
	return nil
}

func Chown(filename string, user interface{}) error {
	//todo
	//switch user.(type) {
	//case string:
	//	grp := os.
	//case int:
	//	return os.Chown(filename,uid,group.(int))
	//default:
	//	return errors.New("unsupported user type")
	//}
	//return os.Chown()
	return nil
}

func Is_dir(filepath string) bool {
	fileinfo, err := os.Stat(filepath)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return fileinfo.IsDir()
}

func Mkdir(filepath string, mode uint32) bool {
	err := os.MkdirAll(filepath, os.FileMode(mode))
	if err != nil {
		return false
	}
	return true
}

// func File_exists(filepath string) bool {
// 	if _, err := os.Stat(filepath); err != nil {
// 		if os.IsNotExist(err) {
// 			return false
// 		}
// 	}
// 	return true
// }

// func FileSize(size int64) string {
// 	s := float64(size)
// 	if s > 1024*1024 {
// 		return fmt.Sprintf("%.1f M", s/(1024*1024))
// 	}
// 	if s > 1024 {
// 		return fmt.Sprintf("%.1f K", s/1024)
// 	}
// 	return fmt.Sprintf("%f B", s)
// }

func IsFile(path string) bool {
	f, e := os.Stat(path)
	if e != nil {
		return false
	}
	if f.IsDir() {
		return false
	}
	return true
}

func IsDir(path string) bool {
	f, e := os.Stat(path)
	if e != nil {
		return false
	}
	return f.IsDir()
}

func Realpath(path string) string {
	rpath, err := filepath.Abs(path)

	if err != nil {
		return path
	}

	return rpath
}

/* ============================================================================================ */
// func PathInfo(path string) string {
// 	var result []string
// 	var arr = strings.Split(path, "/")
// 	for i := 0; i < len(arr)-1; i++ {
// 		result = append(result, arr[i])
// 	}
// 	return strings.Join(result, "/")
// }

// /* ============================================================================================ */
// func BaseName(path string, suffix string) string {
// 	stat, err := os.Stat(path)
// 	setErr(err)
// 	if err != nil {
// 		return ""
// 	}
// 	path = stat.Name()
// 	if suffix != "" {
// 		index := strings.LastIndex(path, suffix)
// 		if index+len(suffix) == len(path) {
// 			path = SubStr(path, 0, index)
// 		}
// 	}
// 	return path
// }

// /* ============================================================================================ */
// func ChGrp(filename string, group int) bool {
// 	err := os.Chown(filename, -1, group)
// 	setErr(err)
// 	if err == nil {
// 		return true
// 	}
// 	return false
// }

// /* ============================================================================================ */
// func ChMod(filename string, mode int) bool {
// 	err := syscall.Chmod(filename, uint32(mode))
// 	if err == nil {
// 		return true
// 	}
// 	return false
// }

// /* ============================================================================================ */
// func ChOwn(filename string, user int) bool {
// 	err := os.Chown(filename, user, -1)
// 	if err == nil {
// 		return true
// 	}
// 	return false
// }

// /* ============================================================================================ */
// func FileOwner(filename string) uint32 {
// 	fi, err := os.Stat(filename)
// 	setErr(err)
// 	return fi.Sys().(*syscall.Stat_t).Uid
// }

// /* ============================================================================================ */
// func FileGroup(filename string) uint32 {
// 	fi, err := os.Stat(filename)
// 	setErr(err)
// 	return fi.Sys().(*syscall.Stat_t).Gid
// }

// /* ============================================================================================ */
// func FilePerms(filename string) string {
// 	fi, err := os.Stat(filename)
// 	setErr(err)
// 	return fi.Mode().Perm().String()
// }

// /* ============================================================================================ */
// func Copy(source string, dest string) bool {
// 	if FileType(source) == "dir" {
// 		fi, err := os.Stat(source)
// 		setErr(err)
// 		_, err = os.Open(dest)
// 		MkDir(dest, int(fi.Mode()), true)
// 		entries, err := ioutil.ReadDir(source)
// 		for _, entry := range entries {
// 			sfp := source + "/" + entry.Name()
// 			dfp := dest + "/" + entry.Name()
// 			Copy(sfp, dfp)
// 		}
// 	} else {
// 		sf, err := os.Open(source)
// 		setErr(err)
// 		defer sf.Close()
// 		df, err := os.Create(dest)
// 		setErr(err)
// 		defer df.Close()
// 		_, err = io.Copy(df, sf)
// 		if err == nil {
// 			si, err := os.Stat(source)
// 			if err != nil {
// 				err = os.Chmod(dest, si.Mode())
// 			}
// 		}
// 	}
// 	return true
// }

/* ============================================================================================ */
func Dirname(path string) string {
	return filepath.Dir(path)
}

// /* ============================================================================================ */
// func FileExists(filename string) bool {
// 	stat, err := os.Stat(filename)
// 	setErr(err)
// 	if stat == nil && err != nil {
// 		return false
// 	}
// 	return true
// }

// /* ============================================================================================ */
// func File(filename string) []string {
// 	file := FileGetContents(filename)
// 	return strings.Split(file, "\n")
// }

//  ============================================================================================
func Filemtime(filename string) int64 {
	stat, err := os.Stat(filename)
	setErr(err)
	if err != nil {
		return 0
	}
	time := stat.ModTime()
	return time.Unix()
}

/* ============================================================================================ */
func Filesize(filename string) int64 {
	stat, err := os.Stat(filename)
	setErr(err)
	if err != nil {
		return 0
	}
	return stat.Size()
}

/* ============================================================================================ */
func Filetype(filename string) string {
	stat, err := os.Stat(filename)
	//setErr(err)
	if err != nil {
		return ""
	}
	if stat.IsDir() {
		return "dir"
	}
	return "file"
}

/* ============================================================================================ */
// func Is_dir(filename string) bool {
// 	stat, err := os.Stat(filename)
// 	return stat.IsDir()
// }

/* ============================================================================================ */
func Is_file(filename string) bool {
	stat, err := os.Stat(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return !stat.IsDir()
}

/* ============================================================================================ */
func Rename(oldName string, newName string) bool {
	err := os.Rename(oldName, newName)
	if err == nil {
		return true
	}
	return false
}

/* ============================================================================================ */
func Unlink(filename string) bool {
	err := os.Remove(filename)
	if err == nil {
		return true
	}
	return false
}

// /* ============================================================================================ */
// func MkDir(dirname string, mode int, recursive bool) bool {
// 	if mode == -1 {
// 		mode = 0777
// 	}
// 	if recursive == true {
// 		err = os.MkdirAll(dirname, os.FileMode(mode))
// 	} else {
// 		err = os.Mkdir(dirname, os.FileMode(mode))
// 	}
// 	setErr(err)
// 	if err == nil {
// 		return true
// 	}
// 	return false
// }

/* ============================================================================================ */
func Rmdir(dirname string, args ...bool) bool {
	all := true
	if len(args) > 0 {
		all = args[0]
	}
	if all == true {
		err = os.RemoveAll(dirname)
	} else {
		err = os.Remove(dirname)
	}
	//setErr(err)
	if err == nil {
		return true
	}
	return false
}

// /* ============================================================================================ */
// func FilePutContents(filename string, data string) bool {
// 	err := ioutil.WriteFile(filename, []byte(data), 0775)
// 	setErr(err)
// 	if err == nil {
// 		return true
// 	}
// 	return false
// }

// /* ============================================================================================ */
// func FileGetContents(filename string) string {
// 	var file []byte
// 	if !strings.HasPrefix(filename, "http://") && !strings.HasPrefix(filename, "https://") {
// 		file, err = ioutil.ReadFile(filename)
// 		setErr(err)
// 	} else {
// 		var timeout = time.Duration(10 * time.Second)
// 		var dialTimeout = func(network, addr string) (net.Conn, error) {
// 			return net.DialTimeout(network, addr, timeout)
// 		}
// 		transport := http.Transport{
// 			Dial: dialTimeout,
// 		}
// 		client := http.Client{
// 			Transport: &transport,
// 		}
// 		var res *http.Response
// 		res, err := client.Get(filename)
// 		setErr(err)
// 		defer func() {
// 			if res != nil && res.Body != nil {
// 				res.Body.Close()
// 			}
// 		}()
// 		if res != nil && res.Body != nil {
// 			file, err = ioutil.ReadAll(res.Body)
// 			setErr(err)
// 		}
// 	}
// 	return string(file)
// }
