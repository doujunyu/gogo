package gogo

import (
	"os"
	"strings"
)

// UrlRootPath 获取当前项目的跟目录
func UrlRootPath() (url string, err error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return strings.Replace(dir, "\\", "/", -1), nil
}
