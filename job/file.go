package job

import (
	"errors"
	"github.com/doujunyu/gogo/utility"
	"os"
	"path"
	"strings"
)

// Files 文件配置
type Files struct {
	Path string `Testing:"文件地址(会从项目跟目录开始)"`
}

// CheckConfig 文件上传验证
type CheckConfig struct {
	Size   int64
	Suffix string
}

// JobNewFile 初始化文件
func JobNewFile() *Files {
	root, _ := utility.UrlRootPath()
	filesPath := os.Getenv("FILE_PATH")
	return &Files{
		Path: root + filesPath,
	}
}

// +----------------------------------------------------------------------
// | 获取文件
// +----------------------------------------------------------------------

// InputFile 接收文件 (接参数,文件放入位置，文件名)(文件位置，error)
func (j *Job) InputFile(FileName string, FilePath string, Check map[string]interface{}) (string, error) {
	//获取文件
	file, handler, err := j.R.FormFile(FileName)
	if err != nil {
		return "", err
	}
	defer file.Close()
	//验证文件大小
	if Check["size"] != nil && handler.Size > Check["size"].(int64) {
		Mb := string(rune(Check["size"].(int64) / 1024 / 1024))
		return "", errors.New("文件过大:上传的文件超过设置的" + Mb + "MB")
	}
	//验证文件格式
	ext := path.Ext(handler.Filename)
	extStr := ext[1:len(ext)]
	if Check["suffix"] != nil && strings.Index(Check["suffix"].(string), extStr) == -1 {
		return "", errors.New("文件格式不正确")
	}
	//更新文件名
	if Check["name"] != nil {
		FileName = Check["name"].(string) + path.Ext(handler.Filename)
	}
	FilePathName := j.File.Path + "/" + FilePath
	_,err = utility.FileNew(FilePathName+"/"+FileName, file)
	if err != nil {
		return "",err
	}
	return os.Getenv("FILE_PATH") + "/" + FilePath+"/"+FileName ,nil
}
