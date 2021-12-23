package utility

import (
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"time"
)
type LogWriteFunc func (string,string,string,string)
type LogWriteStrings struct {
	Url string
	FileName string
	Prefix string
	Content string
}
// LogWrite 写入日志 (日志路径，日志文件名，内容前缀，内容)
func LogWrite(url string, fileName string, prefix string, content string) {
	//日志内存
	path := url + "/" + fileName + "/" + time.Now().Format("2006-01-02")
	_ = os.MkdirAll(path, os.ModePerm) //生成文件
	//拼接文件绝对路径+文件名
	fileNamePath := path + "/servicing.log"
	fileInfo,err := os.Stat(fileNamePath)
	if err == nil{
		if fileInfo.Size() > int64(1024){
			files,_ :=ioutil.ReadDir(path)
			_ = os.Rename(fileNamePath,path + "/_"+strconv.Itoa(len(files))+".log")
		}
	}
	filePath, err := os.OpenFile(fileNamePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModePerm)
	if err != nil {
		return
	}
	defer func(filePath *os.File) {
		err := filePath.Close()
		if err != nil {

		}
	}(filePath)
	//写入文件
	logger := log.New(filePath, "["+prefix+"]:", log.Lmicroseconds)
	logger.Printf(content)
	return
}
