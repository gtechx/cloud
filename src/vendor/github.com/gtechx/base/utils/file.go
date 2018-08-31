package utils

import (
	"io"
	"os"
)

func CopyFile(src, des string) bool {
	srcFile, err := os.Open(src)
	if err != nil {
		return false
	}
	defer srcFile.Close()

	desFile, err := os.Create(des)
	if err != nil {
		return false
	}
	defer desFile.Close()

	_, err = io.Copy(desFile, srcFile)

	return err == nil
}
