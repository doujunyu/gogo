package ff

import (
	"errors"
	"gogo/utility"
	"os"
	"path"
	"strings"
)

type Files struct {
	Path string
}

//文件上传验证
type CheckConfig struct {
	Size   int64
	Suffix string
}

// JobNewFile 初始化日志
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
	if Check["suffix"] != nil && strings.Index(Check["suffix"].(string), handler.Filename) >= 0 {
		return "", errors.New("文件格式不正确")
	}
	//更新文件名
	if Check["name"] != nil {
		FileName = Check["name"].(string) + path.Ext(handler.Filename)
	}
	FilePathName := j.File.Path + "/" + FilePath
	return utility.FileNew(FilePathName+"/"+FileName, file)

}
