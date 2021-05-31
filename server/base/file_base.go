package base

import (
	"fmt"
	"os"
)

func DeleteFile(dir string, file string) bool {
	path := fmt.Sprintf("%s/%s", dir, file)
	if Exists(path) {
		if IsDir(path) {
			err := os.RemoveAll(path)
			if err != nil {
				return false
			}
		} else {
			err := os.Remove(path)
			if err != nil {
				return false
			}
		}
	}
	return true
}
func Exists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}
