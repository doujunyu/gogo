package utility

import (
	"bufio"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"strings"
)

// FileNew 文件地址成文件 (新文件路径,新文件名,文件地址)
func FileNew(NewFileName string, file multipart.File) (string, error) {
	//创建一个文件

	NewPath(NewFileName)
	ThisFile, err := os.OpenFile(NewFileName, os.O_CREATE|os.O_WRONLY, 0766)
	if err != nil {
		return "", err
	}
	defer ThisFile.Close()
	io.Copy(ThisFile, file) //把获取的文件写入到创建的文件内
	return NewFileName, nil
}

// CopyFile 复制文件 (新文件,被复制的文件)
func CopyFile(NewFileName string, OldFileName string) (written int64, err error) {
	srcFile, err := os.Open(OldFileName)
	if err != nil {
		fmt.Println("读取文件错误")
		return
	}
	defer srcFile.Close()
	//通过srcFile，获取到
	reader := bufio.NewReader(srcFile)
	//创建新文件的文件目录
	NewPath(NewFileName)
	//打开NewFileName
	dstFile, err := os.OpenFile(NewFileName, os.O_WRONLY|os.O_CREATE, 0766)
	if err != nil {
		fmt.Println("读取文件错误")
		return
	}
	defer dstFile.Close()
	writer := bufio.NewWriter(dstFile)
	return io.Copy(writer, reader)
}

// NewPath 传入完整路径包括文件名
func NewPath(File string) {
	number := strings.LastIndex(File, "/")
	NewPathIndex := File[:number]
	_ = os.MkdirAll(NewPathIndex, 0766)
}
